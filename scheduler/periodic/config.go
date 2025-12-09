package periodic

type baseConfig struct {
	Cron    string `json:"cron"`
	Timeout int    `json:"timeout"`
}
