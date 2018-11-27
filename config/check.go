package config

import (
	"encoding/hex"
	"errors"
	"net"
	"regexp"
)

func SignKeyValidate(privateKey, publicKey string) (bool, error) {
	if privateKey == "" && publicKey == "" {
		return true, nil
	}
	if len(privateKey) != 64 {
		return false, errors.New("privateKey length error")
	} else if _, err := isHex(privateKey); err != nil {
		return false, errors.New("privateKey string error")
	}

	if len(publicKey) != 130 {
		return false, errors.New("publicKey length error")
	} else if _, err := isHex(publicKey); err != nil {
		return false, errors.New("publicKey string error")
	}
	return true, nil
}

func isHex(key string) ([]byte, error) {
	return hex.DecodeString(key)
}

func PortValidate(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("port number invalidation")
	}
	return true, nil
}

func MongoUrlValidate(url string) (bool, error) {
	if url == "" {
		return true, nil
	}
	match, _ := regexp.MatchString("mongodb://[\\w+.+:\\w+.+@]+[\\d+.+:\\d+.+,]+[/\\w]*[?\\w+;=]+", url)

	if false == match {
		return false, errors.New("mongodb url error")
	}

	return true, nil
}

func IpValidate(ip string) (bool, error) {
	if ip == "" {
		return true, nil
	}
	ipAddr := net.ParseIP(ip)
	if ipAddr.To4() == nil && ipAddr.To16() == nil {
		return false, errors.New("not an IPv4 or IPv6 address")
	} else {
		return true, nil
	}
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
