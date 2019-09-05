package config

import (
	"encoding/hex"
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

//SignKeyValidate validate the keypair is valid
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

//PortValidate validate the port is valid
func PortValidate(port int) (bool, error) {
	if port < 1024 || port > 65535 {
		return false, errors.New("port number invalidation")
	}
	return true, nil
}

//MongoURLValidate validate the mongoDB Url is valid
func MongoURLValidate(url string) (bool, error) {
	if url == "" {
		return true, nil
	}
	match, _ := regexp.MatchString("mongodb://[\\w+.+:\\w+.+@]+[\\d+.+:\\d+.+,]+[/\\w]*[?\\w+;=]+", url)

	if false == match {
		return false, errors.New("mongodb url error")
	}

	return true, nil
}

//IPValidate validate the IP is valid
func IPValidate(ip string) (bool, error) {
	if ip == "" {
		return true, nil
	}
	if ip == "localhost" {
		return true, nil
	}
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false, errors.New("not an validate IPv4 or IPv6 address")
	}

	// specialIP := ipAddr.IsInterfaceLocalMulticast() || ipAddr.IsMulticast() || ipAddr.IsUnspecified() ||
	// 	(ipAddr.To4() != nil && ipAddr.Equal(net.IPv4(255, 255, 255, 255))) ||
	// 	ipAddr[15] == 0 || ipAddr[15] == 255
	if ipAddr.To4() == nil && ipAddr.To16() == nil {
		return false, errors.New("not an validate IPv4 or IPv6 address")
	}
	return true, nil
}

//AddressValidateAll validate the IP:port is valid
func AddressValidateAll(ipAddrs []string) (bool, error) {
	if len(ipAddrs) > 0 {
		for _, address := range ipAddrs {
			addressPair := strings.Split(address, ":")
			if len(addressPair) < 2 {
				return false, errors.New("Address should be like IP:Port")
			}
			boolFlag, err := IPValidate(addressPair[0])
			if err != nil {
				return boolFlag, err
			}
			port, err := strconv.Atoi(addressPair[1])
			if err != nil {
				return false, errors.New("Address should be like IP:Port, make sure the port is valid")
			}
			boolFlag, err = PortValidate(port)
			if err != nil {
				return boolFlag, err
			}
		}
	}
	return true, nil
}
