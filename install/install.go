package install

import (
	"log"
	"strings"
	"time"

	"git.beanjs.com/wangjun/beanjs-v2ray/config"
	"github.com/appleboy/easyssh-proxy"
)

type SSHConfiguration struct {
	Host string
	User string
	Pass string
}

func c2s(cfg *SSHConfiguration) *easyssh.MakeConfig {
	return &easyssh.MakeConfig{
		Server:   cfg.Host,
		Port:     "22",
		User:     cfg.User,
		Password: cfg.Pass,
	}
}

func runScript(cfg *SSHConfiguration, scripts []string) error {
	ssh := c2s(cfg)

	stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream(strings.Join(scripts, "\n"), 12*time.Duration(config.Configuration.Timeout)*time.Second)
	if err != nil {
		return err
	}

	for {
		select {
		case <-doneChan:
			return nil
		case outline := <-stdoutChan:
			if outline != "" {
				log.Printf("ssh out: %v", outline)
			}
		case errline := <-stderrChan:
			if errline != "" {
				log.Printf("ssh err: %v", errline)
			}
		case err = <-errChan:
			return err
		}
	}
}
