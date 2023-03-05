package handlers

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateIP(t *testing.T) {
	cases := []struct {
		name    string
		network entity.IpNetwork
		expRes  bool
	}{
		{name: "valid ip and mask", network: entity.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.255.0",
		}, expRes: true},
		{name: "invalid ip", network: entity.IpNetwork{
			Ip:   "192.12.256.1",
			Mask: "255.255.255.0",
		}, expRes: false},
		{name: "invalid mask", network: entity.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.257.1",
		}, expRes: false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			isValid := ValidateIP(testCase.network)
			require.Equal(t, testCase.expRes, isValid)
		})
	}
}

func TestValidateRequest(t *testing.T) {
	cases := []struct {
		name    string
		request entity.Request
		expRes  bool
	}{
		{name: "valid request", request: entity.Request{
			Login:    "user",
			Password: "1234",
			Ip:       "192.168.1.1",
		}, expRes: true},

		{name: "invalid Login", request: entity.Request{
			Login:    "",
			Password: "1234",
			Ip:       "192.168.1.1",
		}, expRes: false},

		{name: "invalid pass", request: entity.Request{
			Login:    "user",
			Password: "",
			Ip:       "192.168.1.1",
		}, expRes: false},

		{name: "invalid ip", request: entity.Request{
			Login:    "user",
			Password: "1234",
			Ip:       "192.168.256.1",
		}, expRes: false},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			isValid := ValidateRequest(testCase.request)
			require.Equal(t, testCase.expRes, isValid)
		})
	}
}
