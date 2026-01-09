package main

import (
	"context"
	"github.com/iotopo/topmq/accounts"
	"github.com/iotopo/topmq/acls"
	"github.com/iotopo/topmq/auth"
	"github.com/iotopo/topmq/black_list"
	"github.com/iotopo/topmq/cache"
	"github.com/iotopo/topmq/db"
	"github.com/iotopo/topmq/logger"
	"github.com/iotopo/topmq/metrics"
	"github.com/iotopo/topmq/migrator"
	"github.com/iotopo/topmq/mqttserver"
	"github.com/iotopo/topmq/web"
	"github.com/sirupsen/logrus"
	"sync"
)

var wg sync.WaitGroup
var ctx context.Context
var cancel context.CancelFunc

func start() {
	ctx, cancel = context.WithCancel(context.Background())
	logger.Init()
	cache.Init()
	db.Init()
	migrator.DatabaseMigrate()
	err := auth.InitDefaultUser()
	if err != nil {
		logrus.WithError(err).Fatal("failed to create default user")
	}

	web.Init()
	mqttserver.Init()
	auth.Init()
	metrics.Init()
	accounts.Init()
	acls.Init()
	black_list.Init()
}

func run() {
	wg.Go(func() {
		metrics.Run(ctx)
	})

	wg.Go(func() {
		err := mqttserver.Run(ctx)
		if err != nil {
			logrus.WithError(err).Fatal("mqtt server start failed")
		}
	})

	wg.Go(func() {
		web.Run(ctx)
	})
	wg.Wait()
}

func stop() {
	mqttserver.Close()
	cancel()
	db.Close()
	cache.Close()
	logger.Close()
}
