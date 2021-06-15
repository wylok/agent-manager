package modules

import (
	"agent-manager/config"
	"agent-manager/kits"
	_ "bytes"
	_ "github.com/CodyGuo/godaemon"
	"github.com/mitchellh/go-ps"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func DownloadAgent(agent, AgentFile string) bool {
	// 文件重新下载
	res, err := http.Get(config.AgentFileUrl + agent)
	kits.Log("重新下载"+agent, "info", "CheckAgent")
	if err != nil {
		kits.Log(err.Error(), "error", "CheckAgent")
	} else {
		f, err := os.Create(AgentFile)
		if err != nil {
			kits.Log(err.Error(), "error", "CheckAgent")
		} else {
			_, err = io.Copy(f, res.Body)
			if err == nil {
				if kits.CheckFile(AgentFile) {
					_ = f.Chmod(0755)
					f.Close()
					return true
				}
			}
		}
	}
	return false
}

func CheckAgent(agent, AgentPath string, force bool) {
	AgentPid := config.PidPath + agent + ".pid"
	AgentFile := AgentPath + agent
	cmd := exec.Command(AgentFile, "-d")
	if kits.CheckFile(AgentPid) {
		f, err := ioutil.ReadFile(AgentPid)
		if err != nil {
			kits.Log(err.Error(), "error", "CheckAgent")
		} else {
			p, err := ps.FindProcess(int(kits.BytesToInt64(f)))
			if err != nil {
				kits.Log(err.Error(), "error", "CheckAgent")
			}
			if p == nil {
				if kits.CheckFile(AgentFile) {
					_ = cmd.Run()
				} else {
					if DownloadAgent(agent, AgentFile) {
						_ = cmd.Run()
					}
				}
			}
		}
	} else {
		if force {
			if kits.CheckFile(AgentFile) {
				_ = cmd.Run()
			} else {
				if DownloadAgent(agent, AgentFile) {
					_ = cmd.Run()
				}
			}
		}
	}
}
