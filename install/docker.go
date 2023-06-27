package install

import (
	"log"
)

func Docker(cfg *SSHConfiguration) error {
	log.Print("install docker")

	return runScript(cfg, []string{
		"apk update && sleep 5",
		"apk add docker && sleep 5",
		"rc-update add docker boot && sleep 5",
		"service docker start && sleep 5",
		"docker version",
	})
}
