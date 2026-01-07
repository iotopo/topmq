package main

import (
	"github.com/iotopo/topmq/cmd"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
		return
	}
	start()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		logrus.Info("收到退出信号，正在关闭...")
		cancel()
	}()

	run()
	stop()
}
