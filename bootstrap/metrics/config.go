package metrics

type Config struct {
	Enable   bool
	Addr     string
	NodeName string
	Interval int `default:"15"`
}
