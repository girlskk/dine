package alert

type Config struct {
	AppName     string
	Disabled    bool `default:"true"`
	AccessToken string
	Secret      string
	Ats         []string
}
