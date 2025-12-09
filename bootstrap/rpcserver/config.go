package rpcserver

type Config struct {
	Port           string `default:"50051"`
	RequestTimeout int    `default:"15"`
}
