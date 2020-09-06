package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	Debug  bool `env:"DEBUG,required"`
	Server serverConf
	Dns    dnsConf
}

type serverConf struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Login         string           `env:"SERVER_LOGIN,required"`
	Password         string           `env:"SERVER_PASSWORD,required"`
}

type dnsConf struct {
	Host          string `env:"RFC2136_HOST,required"`
	Port          int    `env:"RFC2136_PORT,required"`
	Zone          string `env:"RFC2136_ZONE,required"`
	TsigSecret    string `env:"RFC2136_TSIG_SECRET"`
	TsigSecretAlg string `env:"RFC2136_TSIG_SECRET_ALG"`
	TsigKeyname   string `env:"RFC2136_TSIG_KEYNAME"`
}

func AppConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
