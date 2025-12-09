package util

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
)

type FlagStringArray []string

func (f *FlagStringArray) String() string {
	return fmt.Sprint(*f)
}

func (f *FlagStringArray) Set(value string) error {
	if len(*f) > 0 {
		return errors.New("flag already set")
	}
	for _, s := range SplitString(value, ",") {
		*f = append(*f, s)
	}
	return nil
}

func SplitString(s, sep string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, sep)
}

// NodeETH0Addr 获取运行环境eth0 IPv4
func NodeETH0Addr() (string, error) {
	inter, err := net.InterfaceByName("eth0")
	if err != nil {
		inter, err = net.InterfaceByName("en0")
		if err != nil {
			return "", err
		}
	}
	addrs, err := inter.Addrs()
	if err != nil {
		return "", err
	}
	for _, v := range addrs {
		ipNet, ok := v.(*net.IPNet)
		if !ok {
			continue
		}
		addr := ipNet.IP.String()
		if strings.Contains(addr, ".") {
			return addr, nil
		}
	}
	return "", errors.New("not found IPv4 addr")
}

func LastErrorMessage(err error) (msg string) {
	if err == nil {
		return
	}
	errs := strings.Split(err.Error(), ":")
	return errs[len(errs)-1]
}

func SHA1Hash(raw []byte) string {
	s := sha1.Sum(raw)
	return hex.EncodeToString(s[:])
}
