package httpserver

type Config struct {
	Port           string `default:"8080"`
	RequestTimeout int    `default:"120"`
}
