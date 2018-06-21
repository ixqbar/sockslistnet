package spider

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/robfig/cron"
)

var GlobalVars = NewTGlobalVars()

func Run() {
	GlobalVars.Init();

	proxy := NewProxy()
	if GTSpiderConfig.PullOnStartup {
		proxy.Pull()
	}

	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)

	c := cron.New()
	c.AddFunc(GTSpiderConfig.Schedule, func() {
		index := GlobalVars.Index()
		Logger.Printf("schedule[%d] start execute", index)
		defer Logger.Printf("schedule[%d] stop execute", index)

		proxy.Pull()
	})
	c.Start()

	select {
		case <-stopSignal:
			Logger.Print("catch exit signal")
	}

	c.Stop()
	GlobalVars.Clear()

	Logger.Print("proxy spider exit")
}