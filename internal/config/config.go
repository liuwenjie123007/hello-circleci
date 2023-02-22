package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Env struct {
	MySQL MySQL
}

type MySQL struct {
	Host     string `envconfig:"MYSQL_HOST" required:"true"`
	Port     string `envconfig:"MYSQL_PORT" required:"true"`
	User     string `envconfig:"MYSQL_USER" required:"true"`
	Password string `envconfig:"MYSQL_PASSWORD" required:"true"`
}

func GetAPIConfig() (*Env, error) {
	var n Env
	if err := envconfig.Process("", &n); err != nil {
		return nil, fmt.Errorf("get api config error: %w", err)
	}

	n.MySQL = MySQL{
		Host:     n.MySQL.Host,
		Port:     n.MySQL.Port,
		User:     n.MySQL.User,
		Password: n.MySQL.Password,
	}

	return &n, nil
}
