package install

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"git.beanjs.com/wangjun/beanjs-v2ray/config"
	"github.com/google/uuid"
)

type logDef struct {
	LogLevel string `json:"loglevel"`
}

type inboundSettingsClientsDef struct {
	ID      string `json:"id"`
	AlterId int    `json:"alterId"`
}

type inboundSettingsDef struct {
	Clients []inboundSettingsClientsDef `json:"clients"`
}

type inboundStreamSettingsWsSettingsDef struct {
	Path string `json:"path"`
}

type inboundStreamSettingsDef struct {
	Network    string                             `json:"network"`
	WsSettings inboundStreamSettingsWsSettingsDef `json:"wsSettings"`
}

type inboundDef struct {
	Port           int                      `json:"port"`
	Listen         string                   `json:"listen"`
	Protocol       string                   `json:"protocol"`
	Settings       inboundSettingsDef       `json:"settings"`
	StreamSettings inboundStreamSettingsDef `json:"streamSettings"`
}

type outboundSettingsDef struct {
}

type outboundDef struct {
	Protocol string              `json:"protocol"`
	Settings outboundSettingsDef `json:"settings"`
}

type v2rayConfigDef struct {
	Log       logDef        `json:"log"`
	Inbounds  []inboundDef  `json:"inbounds"`
	Outbounds []outboundDef `json:"outbounds"`
}

var v2rayConfig = &v2rayConfigDef{
	Log: logDef{
		LogLevel: "debug",
	},
	Inbounds: []inboundDef{
		{
			Port:     80,
			Listen:   "0.0.0.0",
			Protocol: "vmess",
			Settings: inboundSettingsDef{
				Clients: []inboundSettingsClientsDef{
					{
						ID:      "",
						AlterId: 0,
					},
				},
			},
			StreamSettings: inboundStreamSettingsDef{
				Network: "ws",
				WsSettings: inboundStreamSettingsWsSettingsDef{
					Path: "",
				},
			},
		},
	},
	Outbounds: []outboundDef{
		{
			Protocol: "freedom",
			Settings: outboundSettingsDef{},
		},
	},
}

func init() {
	v2rayId := uuid.New().String()
	if config.Configuration.V2RayID != "" {
		v2rayId = config.Configuration.V2RayID
	}

	v2rayConfig.Inbounds[0].Settings.Clients[0].ID = v2rayId
	v2rayConfig.Inbounds[0].StreamSettings.WsSettings.Path = config.Configuration.V2RayPath
}

func V2Ray(cfg *SSHConfiguration) (string, error) {
	var err error

	log.Print("install v2ray")

	err = runScript(cfg, []string{
		fmt.Sprintf("mkdir -p %v", config.Configuration.V2RayDir),
	})
	if err != nil {
		return "", err
	}

	v2rayConfigBytes, err := json.MarshalIndent(v2rayConfig, "", "  ")
	if err != nil {
		return "", err
	}

	ssh := c2s(cfg)
	err = ssh.WriteFile(bytes.NewReader(v2rayConfigBytes), int64(len(v2rayConfigBytes)), fmt.Sprintf("%v/config.json", config.Configuration.V2RayDir))
	if err != nil {
		return "", err
	}

	err = runScript(cfg, []string{
		fmt.Sprintf("docker run -d -e V2RAY_VMESS_AEAD_FORCED=false -e xray.vmess.aead.forced=false -v %v/config.json:/etc/xray/config.json:ro -p %v:80 --name v2rayd --restart always teddysun/xray:latest", config.Configuration.V2RayDir, config.Configuration.VultrPort),
	})
	if err != nil {
		return "", err
	}

	return v2rayConfig.Inbounds[0].Settings.Clients[0].ID, nil
}
