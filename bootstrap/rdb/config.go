package rdb

type Config struct {
	Addr         string `default:"localhost:6379"`
	Username     string
	Password     string
	DB           int `default:"0"`
	PoolSize     int
	MinIdleConns int
	MaxIdleConns int
}
