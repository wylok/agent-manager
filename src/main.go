package main

import (
	"agent-manager/config"
	"agent-manager/kits"
	"agent-manager/modules"
	_ "bytes"
	_ "github.com/CodyGuo/godaemon"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	// 检查进程是否重复启动
	if kits.CheckFile(config.PidFile) {
		f, err := ioutil.ReadFile(config.PidFile)
		if err != nil {
			kits.Log(err.Error(), "error", "main")
		}
		p, err := ps.FindProcess(int(kits.BytesToInt64(f)))
		if p != nil {
			println("manager进程已运行!")
			os.Exit(1)
		}
	}
	for {
		pid := os.Getpid()
		err := ioutil.WriteFile(config.PidFile, kits.Int64ToBytes(int64(pid)), 0666)
		if err != nil {
			kits.Log(err.Error(), "error", "main")
		}
		// 监视agent进程状况
		go modules.CheckAgent("cmdb-agent", "/usr/bin/", true)
		go modules.CheckAgent("monitor-agent", "/usr/bin/", false)
		// 监视proxy进程状况
		go modules.CheckAgent("agent-proxy", "/usr/bin/", false)
		time.Sleep(time.Duration(config.SleepInterval) * time.Second)
	}
}
