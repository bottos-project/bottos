package config

import (
	"errors"
	"regexp"
)

const (
	ipReg = "(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)"
)

func SignKeyValidate(privateKey, publicKey string) (bool, error) {
	if privateKey == "" && publicKey == "" {
		return true, nil
	}
	if len(privateKey) != 64 {
		return false, errors.New("privateKey length error")
	}

	if len(publicKey) != 130 {
		return false, errors.New("publicKey length error")
	}
	return true, nil
}

func PortValidate(port int) (bool, error) {
	if port < 1 && port > 65535 {
		return false, errors.New("port number invalidation")
	}
	return true, nil
}

func MongoUrlValidate(url string) (bool, error) {
	if url == "" {
		return true, nil
	}
	match, _ := regexp.MatchString("mongodb://[\\w+.+:\\w+.+@]+[\\w+.+:\\w+.+,]+[/\\w]*[?\\w+;=]+", url)

	if false == match {
		return false, errors.New("mongodb url error")
	}

	return true, nil
}

func IpValidate(ipAddr string) (bool, error) {
	if ipAddr == "" {
		return true, nil
	}
	match, _ := regexp.MatchString(ipReg, ipAddr)
	if false == match {
		return false, errors.New("ip address invalidation")
	}
	return true, nil
}

func IpValidateAll(ipAddrs []string) (bool, error) {
	if len(ipAddrs) > 0 {
		for _, ip := range ipAddrs {
			boolFlag, err := IpValidate(ip)
			if err != nil {
				return boolFlag, err
			}
		}
	}
	return true, nil
}
