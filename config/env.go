package config

import (
	"log"

	"github.com/Netflix/go-env"
)

type Environment struct {
	Port     int  `env:"PORT,default=8000"`
	Debug    bool `env:"DEBUG,default=false"`
	DbConfig DbConfig
	Extras   env.EnvSet
}

type DbConfig struct {
	User     string `env:"DB_USER,default=postgres"`
	Port     string `env:"DB_PORT,default=5432"`
	Password string `env:"DB_PASSWORD,default=postgres"`
	Host     string `env:"DB_HOST,default=localhost"`
	Name     string `env:"DB_NAME,default=postgres"`
}

var Env Environment

func init() {
	_, err := env.UnmarshalFromEnviron(&Env)
	if err != nil {
		log.Fatal(err)
	}

}

func GetEnv() Environment {
	return Env
}
