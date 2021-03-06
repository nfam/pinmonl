package main

import (
	"time"

	"github.com/spf13/viper"
)

type config struct {
	Address string
	Verbose int

	DefaultUser bool

	Web struct {
		DevServer string
	}

	JWT struct {
		Secret string
		Issuer string
		Expire time.Duration
	}

	DB struct {
		Driver string
		DSN    string
	}

	Exchange struct {
		Enabled bool
		Address string
	}

	Queue struct {
		Job    int
		Worker int
	}
}

func unmarshalConfig() (*config, error) {
	var c config
	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
