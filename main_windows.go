// https://github.com/golang/sys/blob/master/windows/svc/example/service.go
package main

import (
	"github.com/iotopo/topmq/cmd"
	"github.com/iotopo/topmq/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"os"
	"time"
)

type daemoService struct {
}

func (m *daemoService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	start()
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
	go run()

	defer stop()

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				// Testing deadlock from https://code.google.com/p/winsvc/issues/detail?id=4
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// golang.org/x/sys/windows/svc.TestExample is verifying this output.
				changes <- svc.Status{State: svc.StopPending}
				return
			default:
				logrus.Errorf("unexpected control request #%d", c)
			}
		}
	}
}

func main() {
	inService, err := svc.IsWindowsService()
	if err != nil {
		logrus.Fatalf("failed to determine if we are running in service: %v", err)
	}
	if !inService {
		if len(os.Args) > 1 {
			cmd.Execute()
			return
		}
	}

	name := config.AppName

	logrus.Infof("starting %s service", name)
	if inService {
		err = svc.Run(name, &daemoService{})
	} else {
		err = debug.Run(name, &daemoService{})
	}
	if err != nil {
		logrus.Errorf("%s service failed: %v", name, err)
		return
	}
	logrus.Infof("%s service stopped", name)
}
