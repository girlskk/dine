package db

import (
	"fmt"
	"net/url"
)

type Config struct {
	AutoMigrate bool `default:"false"`
	Debug       bool `default:"true"`

	Host     string `default:"localhost"`
	Port     string `default:"3306"`
	User     string
	Password string
	Name     string
	TimeZone string `default:"Asia/Shanghai"`

	MaxLifetime  int `default:"3600"`
	MaxOpenConns int `default:"150"`
	MaxIdleConns int `default:"50"`
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		url.QueryEscape(c.TimeZone),
	)
}
