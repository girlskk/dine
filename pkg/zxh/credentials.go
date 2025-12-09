package zxh

type Credentials struct {
	MchID     string
	AccessKey string
}

func NewCredentials(mchID, accessKey string) *Credentials {
	return &Credentials{
		MchID:     mchID,
		AccessKey: accessKey,
	}
}

func (c *Credentials) auth(params map[string]any) {
	params["mchId"] = c.MchID
	params["sign"] = sign(params, c.AccessKey)
}
