package config

import (
	"github.com/verystar/ini"
	"log"
)

var cfg *ini.Ini
var env string

func Setup(environment string) {
	var err error
	cfg, err = ini.Load("conf.ini")
	if err != nil {
		log.Fatalln(err)
	}
	env = environment
}

func Get(key string) string {
	return cfg.Read(env, key)
}

func Getenv() string {
	return env
}
