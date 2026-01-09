package mqttserver

import (
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/packets"
	"slices"
	"strings"
	"time"
)

type OverviewInfo struct {
	Uptime           int64 `json:"uptime"`           // the number of seconds the server has been online
	ClientsConnected int64 `json:"clientsConnected"` // number of currently connected clients
	ClientsMaximum   int64 `json:"clientsMaximum"`   // maximum number of active clients that have been connected
	ClientsTotal     int64 `json:"clientsTotal"`     // total number of connected and disconnected clients with a persistent session currently connected and registered
	Retained         int64 `json:"retained"`         // total number of retained messages active on the broker
	Subscriptions    int64 `json:"subscriptions"`    // total number of subscriptions active on the broker
	Inflight         int64 `json:"inflight"`         // the number of messages currently in-flight
	InflightDropped  int64 `json:"inflightDropped"`  // the number of inflight messages which were dropped
	MessagesSent     int   `json:"messagesSent"`     // 每秒发送的消息数量
	MessagesReceived int   `json:"messagesReceived"` // 每秒接收的消息数量
	MessagesDropped  int   `json:"messagesDropped"`  // 每秒丢弃的消息数量
}

type Client struct {
	Listener              string     `json:"listener"` // listener id of the client: ws/tcp
	ClientID              string     `json:"clientID"`
	Username              string     `json:"username"`
	Remote                string     `json:"remote"`
	Keepalive             uint16     `json:"keepalive"`
	Clean                 bool       `json:"clean"`
	ProtocolVersion       byte       `json:"protocolVersion"`
	Subscriptions         int        `json:"subscriptions"`
	Closed                bool       `json:"closed"`
	DisconnectedAt        *time.Time `json:"disconnectedAt,omitempty"`
	SessionExpiryInterval uint32     `json:"sessionExpiryInterval"`
}

type ClientRequest struct {
	ClientID string `json:"clientID" form:"clientID"`
	Username string `json:"username" form:"username"`
	Remote   string `json:"remote" form:"remote"`
}

type RetainedMessage struct {
	Topic     string     `json:"topic"`
	ClientID  string     `json:"clientID"`
	Qos       byte       `json:"qos"`
	CreatedAt time.Time  `json:"createdAt"`
	ExpiredAt *time.Time `json:"expiredAt"`
}

type Subscription struct {
	ClientID       string `json:"clientID"`
	Remote         string `json:"remote"`
	Topic          string `json:"topic"`
	Qos            byte   `json:"qos"`
	Retain         bool   `json:"retain"`         // 发布时状态保留
	NoLocal        bool   `json:"noLocal"`        // 禁止本地转发
	RetainHandling byte   `json:"retainHandling"` // 保留消息处理
}

type SubscriptionRequest struct {
	ClientID string `json:"clientID" form:"clientID"`
	Topic    string `json:"topic" form:"topic"`
}

func CloseClient(clientID string) error {
	if server == nil {
		return nil
	}
	if clientID == "" {
		return nil
	}
	client, ok := server.Clients.Get(clientID)
	if ok {
		return server.DisconnectClient(client, packets.CodeDisconnect)
	}
	return nil
}

func GetClientDetails(clientID string) *Client {
	client, ok := server.Clients.Get(clientID)
	if ok {
		clientInfo := &Client{
			Listener:              client.Net.Listener,
			Closed:                client.Closed(),
			ClientID:              clientID,
			Username:              string(client.Properties.Username),
			Remote:                client.Net.Remote,
			Clean:                 client.Properties.Clean,
			Keepalive:             client.State.Keepalive,
			Subscriptions:         client.State.Subscriptions.Len(),
			SessionExpiryInterval: client.Properties.Props.SessionExpiryInterval,
			ProtocolVersion:       client.Properties.ProtocolVersion,
		}
		disconnectedTime := client.StopTime()
		if disconnectedTime > 0 {
			t := time.Unix(disconnectedTime, 0)
			clientInfo.DisconnectedAt = &t
		}
		return clientInfo
	}
	return nil
}

func GetClients(req ClientRequest) []Client {
	if server == nil {
		return nil
	}

	limit := 3000

	var clients []Client
	clientMap := server.Clients.GetAll()
	var ids []string
	for clientID, client := range clientMap {
		if client.Net.Inline {
			continue
		}
		if req.ClientID != "" && !strings.Contains(clientID, req.ClientID) {
			continue
		}
		if req.Username != "" && !strings.Contains(string(client.Properties.Username), req.Username) {
			continue
		}
		if req.Remote != "" && !strings.Contains(client.Net.Remote, req.Remote) {
			continue
		}
		ids = append(ids, clientID)
	}
	slices.Sort(ids)
	for _, clientID := range ids {
		client := clientMap[clientID]
		clientInfo := Client{
			Listener:              client.Net.Listener,
			Closed:                client.Closed(),
			ClientID:              clientID,
			Username:              string(client.Properties.Username),
			Remote:                client.Net.Remote,
			Clean:                 client.Properties.Clean,
			Keepalive:             client.State.Keepalive,
			Subscriptions:         client.State.Subscriptions.Len(),
			SessionExpiryInterval: client.Properties.Props.SessionExpiryInterval,
			ProtocolVersion:       client.Properties.ProtocolVersion,
		}
		disconnectedTime := client.StopTime()
		if disconnectedTime > 0 {
			t := time.Unix(disconnectedTime, 0)
			clientInfo.DisconnectedAt = &t
		}
		clients = append(clients, clientInfo)
		if len(clients) >= limit {
			return clients
		}
	}
	return clients
}

func GetSubscriptions(req SubscriptionRequest) []Subscription {
	if server == nil {
		return nil
	}
	limit := 3000

	clientMap := server.Clients.GetAll()
	var ids []string
	for clientID, client := range clientMap {
		if client.Net.Inline {
			continue
		}
		if req.ClientID != "" {
			if !strings.Contains(clientID, req.ClientID) {
				continue
			}
		}
		ids = append(ids, clientID)
	}
	slices.Sort(ids)

	var subscriptions []Subscription
	for _, clientID := range ids {
		client, ok := server.Clients.Get(clientID)
		if !ok {
			continue
		}
		subscriptionMap := client.State.Subscriptions.GetAll()

		var keys []string
		for _, item := range subscriptionMap {
			if req.Topic != "" && !strings.Contains(item.Filter, req.Topic) && !auth.RString(req.Topic).FilterMatches(item.Filter) {
				continue
			}
			keys = append(keys, item.Filter)
		}
		slices.Sort(keys)

		for _, key := range keys {
			item := subscriptionMap[key]
			sub := Subscription{
				ClientID:       clientID,
				Remote:         client.Net.Remote,
				Topic:          item.Filter,
				Qos:            item.Qos,
				Retain:         item.RetainAsPublished,
				RetainHandling: item.RetainHandling,
				NoLocal:        item.NoLocal,
			}
			subscriptions = append(subscriptions, sub)
			if len(subscriptionMap) >= limit {
				return subscriptions
			}
		}
	}
	return subscriptions
}

func GetRetained(filter string) []RetainedMessage {
	if server == nil {
		return nil
	}
	topicMap := server.Topics.Retained.GetAll()
	var items []RetainedMessage
	var topics []string
	for topic, msg := range topicMap {
		if msg.Ignore {
			continue
		}
		if filter != "" && !strings.Contains(topic, filter) && !auth.RString(filter).FilterMatches(topic) {
			continue
		}
		topics = append(topics, topic)
	}
	slices.Sort(topics)

	limit := 3000

	for _, topic := range topics {
		msg := topicMap[topic]
		if msg.Ignore || msg.Origin == "" {
			continue
		}
		item := RetainedMessage{
			Topic:     msg.TopicName,
			ClientID:  msg.Origin,
			Qos:       msg.Properties.MaximumQos,
			CreatedAt: time.Unix(msg.Created, 0),
		}
		if msg.Expiry > 0 {
			t := time.Unix(msg.Expiry, 0)
			item.ExpiredAt = &t
		}
		items = append(items, item)
		if len(items) >= limit {
			return items
		}
	}
	return items
}

func DeleteRetained(topic string) {
	if server != nil {
		server.Topics.Retained.Delete(topic)
	}
}

func GetRetainedPayload(topic string) []byte {
	if server == nil {
		return nil
	}

	pkt, ok := server.Topics.Retained.Get(topic)
	if !ok {
		return nil
	}
	return pkt.Payload
}

func GetOverviewInfo() OverviewInfo {
	if server == nil {
		return OverviewInfo{}
	}
	info := server.Info.Clone()
	return OverviewInfo{
		ClientsConnected: info.ClientsConnected,
		ClientsMaximum:   info.ClientsMaximum,
		ClientsTotal:     info.ClientsTotal,
		Retained:         info.Retained,
		Subscriptions:    info.Subscriptions,
		Inflight:         info.Inflight,
		InflightDropped:  info.InflightDropped,
		MessagesReceived: int(CurrentMessagesReceived),
		MessagesSent:     int(CurrentMessagesSent),
		MessagesDropped:  int(CurrentMessagesDropped),
	}
}
