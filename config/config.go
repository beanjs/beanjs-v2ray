package config

import (
	"github.com/alexflint/go-arg"
	_ "github.com/joho/godotenv/autoload"
)

type ConfigurationDef struct {
	VultrKey   string `arg:"required,env:VULTR_KEY"`
	VultrLabel string `arg:"env:VULTR_LABEL" default:"v2ray-proxy"`
	VultrPort  int    `arg:"env:VULTR_PORT" default:"3000"`

	V2RayID   string `arg:"env:V2RAY_ID" default:""`
	V2RayPath string `arg:"env:V2RAY_PATH" default:"/world"`
	V2RayDir  string `arg:"env:V2RAY_DIR" default:"/home/v2ray"`
	Timeout   int    `default:"10"`
}

var Configuration = ConfigurationDef{}

func (ConfigurationDef) Version() string {
	return "beanjs-v2ray 1.0.0"
}

func init() {
	arg.MustParse(&Configuration)
}
