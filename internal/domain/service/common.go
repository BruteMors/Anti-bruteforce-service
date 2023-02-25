package service

import (
	"strconv"
	"strings"
)

func getPrefix(inputIp string, inputMask string) (string, error) {
	ip := strings.Split(inputIp, ".")
	mask := strings.Split(inputMask, ".")
	var prefix string
	for index, ipOct := range ip {
		intIpOct, err := strconv.Atoi(ipOct)
		if err != nil {
			return "", err
		}
		intMaskOct, err := strconv.Atoi(mask[index])
		if err != nil {
			return "", err
		}
		prefix += strconv.Itoa(intIpOct & intMaskOct)
		if index != len(ip)-1 {
			prefix += "."
		}
	}
	return prefix, nil
}
