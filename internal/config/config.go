package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Services map[string]string `json:"Services" yaml:"Services"`
}
