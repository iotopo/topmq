package mqttserver

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/iotopo/topmq/accounts"
	"github.com/iotopo/topmq/acls"
	"github.com/iotopo/topmq/black_list"
	"github.com/iotopo/topmq/config"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// AuthHook is an authentication hook which implements an auth ledger.
type AuthHook struct {
	mqtt.HookBase
	ctx       context.Context
	logger    *logrus.Entry
	authCache *expirable.LRU[string, bool]
	aclCache  *expirable.LRU[string, bool]
}

// ID returns the ID of the hook.
func (h *AuthHook) ID() string {
	return "auth-hook"
}

// Provides indicates which hook methods this hook provides.
func (h *AuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
		mqtt.OnPublishDropped,
	}, []byte{b})
}

// Init configures the hook with the auth ledger to be used for checking.
func (h *AuthHook) Init(config any) error {
	//h.logger.Info("loaded auth config")
	return nil
}

func (h *AuthHook) authOK(cl *mqtt.Client, pk packets.Packet) bool {
	username := string(pk.Connect.Username)
	password := string(pk.Connect.Password)
	clientID := cl.ID
	remote := cl.Net.Remote

	// 检查黑名单和白名单
	if authItems, err := black_list.FindAll(h.ctx); err == nil {
		for _, item := range authItems {
			if item.ExpiredAt != nil {
				if item.ExpiredAt.Before(time.Now()) {
					continue
				}
			}
			if item.ClientID != "" && item.ClientID == cl.ID ||
				item.Username != "" && item.Username == username ||
				(item.Remote != "" && strings.HasPrefix(cl.Net.Remote, item.Remote)) {
				return false
			}
		}
	}

	// 检查 mqtt server 账号
	if account, err := accounts.FindByUsername(username); err == nil {
		if account != nil &&
			!account.Disabled &&
			(account.Remote == "" || strings.HasPrefix(remote, account.Remote)) &&
			(auth.RString(account.ClientID).Matches(clientID)) &&
			account.Password == password {
			return true
		}
	} else {
		h.logger.WithError(err).Errorf("failed to find account: %s", username)
	}

	// Redis 密码扩展认证
	if config.Conf.Redis.PasswordHash != "" {
		redisAuth, err := GetRedisAuth(h.ctx, username)
		if err == nil {
			if redisAuth != nil {
				if redisAuth.Verify(password) {
					return true
				} else {
					return false
				}
			}
		} else {
			h.logger.WithError(err).Warnf("failed to get mqtt redis user: %s", username)
		}
	}

	return false
}

func (h *AuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	username := string(pk.Connect.Username)
	password := string(pk.Connect.Password)
	clientID := cl.ID
	remote := cl.Net.Remote

	cacheKey := fmt.Sprintf("%s:%s:%s:%s", username, password, clientID, remote)
	if v, ok := h.authCache.Get(cacheKey); ok {
		return v
	}

	pass := h.authOK(cl, pk)
	h.authCache.Add(cacheKey, pass)

	if !pass {
		h.logger.Infof("client failed authentication check: username=%s client_id=%s remote=%s", username, clientID, remote)
	}
	return pass
}

func (h *AuthHook) aclOK(cl *mqtt.Client, topic string, write bool) bool {
	aclRules, err := acls.FindAll(h.ctx)
	if err == nil && aclRules != nil {
		clientID := cl.ID
		username := string(cl.Properties.Username)
		varMap := map[string]string{
			"clientid": clientID,
			"username": username,
		}
		for _, rule := range aclRules {
			ruleTopic := ReplaceVariables(rule.Topic, varMap)
			if rule.ClientID != "" && rule.ClientID == cl.ID ||
				rule.Username != "" && rule.Username == username ||
				(rule.Remote != "" && strings.HasPrefix(cl.Net.Remote, rule.Remote)) {
				if auth.RString(ruleTopic).FilterMatches(topic) {
					access := rule.Access
					if !write && (access == "r" || access == "rw") {
						return true
					} else if write && (access == "w" || access == "rw") {
						return true
					} else {
						return false
					}
				}
			}
		}
	}

	return false
}

// OnACLCheck returns true if the connecting client has matching read or write access to subscribe
// or publish to a given topic.
func (h *AuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	username := string(cl.Properties.Username)

	cacheKey := fmt.Sprintf("%s:%s:%v", username, topic, write)
	if v, ok := h.aclCache.Get(cacheKey); ok {
		return v
	}

	pass := h.aclOK(cl, topic, write)
	h.aclCache.Add(cacheKey, pass)
	if !pass {
		h.logger.Debugf("client failed allowed ACL check: client_id=%s username=%s topic=%s", cl.ID, string(cl.Properties.Username), topic)
	}
	return pass
}

func newAuthHook(ctx context.Context) *AuthHook {
	return &AuthHook{
		ctx: ctx,
		//logger:    log.GetLogger("mqtt_server", log.WithTopic()),
		logger:    logrus.WithField("module", "mqtt_server"),
		authCache: expirable.NewLRU[string, bool](1000, nil, time.Minute*1),
		aclCache:  expirable.NewLRU[string, bool](5000, nil, time.Minute*1),
	}
}
