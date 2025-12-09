package oss

import "fmt"

type Config struct {
	Region          string `default:"cn-shanghai"`
	Bucket          string
	Dir             string
	AccessKeyID     string
	AccessKeySecret string
	RoleARN         string
	RoleSessionName string
	Expiration      int `default:"3600"`
	Domain          string
}

func (c Config) host() string {
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com", c.Bucket, c.Region)
}
