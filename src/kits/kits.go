package kits

import (
	"agent-manager/config"
	"log"
	"os"
	"syscall"
)

func Log(Msg, MsgType string, FuncName string) {

	Prefix := map[string]string{"info": "[Info]", "error": "[Error]", "debug": "[Debug]"}
	_, err := os.Stat(config.LogFile)
	if err == nil {
		logFile, err := os.OpenFile(config.LogFile, syscall.O_RDWR|syscall.O_APPEND, 0666)
		if err == nil {
			defer logFile.Close()
			debugLog := log.New(logFile, FuncName+Prefix[MsgType], log.LstdFlags)
			debugLog.Println(Msg)
		}
	} else {
		logFile, err := os.Create(config.LogFile)
		if err == nil {
			defer logFile.Close()
			debugLog := log.New(logFile, FuncName+Prefix[MsgType], log.LstdFlags)
			debugLog.Println(Msg)
		}
	}
}
