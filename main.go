package main

import (
	"log"
	"time"

	"git.beanjs.com/wangjun/beanjs-v2ray/config"
	"git.beanjs.com/wangjun/beanjs-v2ray/install"
	"git.beanjs.com/wangjun/beanjs-v2ray/vultr"
)

func main() {
	if err := vultr.ShowInfo(); err != nil {
		log.Fatal(err)
	}

	fwg, err := vultr.Firewall()
	if err != nil {
		log.Fatal(err)
	}

	ins, err := vultr.Instance(fwg.ID)
	if err != nil {
		log.Fatal(err)
	}

	waitTime := 6 * config.Configuration.Timeout
	log.Printf("wait %v second", waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)

	sshCfg := &install.SSHConfiguration{
		Host: ins.MainIP,
		User: "root",
		Pass: ins.DefaultPassword,
	}

	if err := install.Docker(sshCfg); err != nil {
		log.Fatal(err)
	}

	v2rayId, err := install.V2Ray(sshCfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("!!!!v2ray deploy done!!!!")
	log.Printf("v2ray id: %v", v2rayId)
	log.Printf("v2ray ip: %v", ins.MainIP)
	log.Printf("v2ray port: %v", config.Configuration.VultrPort)
	log.Printf("v2ray path: %v", config.Configuration.V2RayPath)
}
