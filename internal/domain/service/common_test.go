package service

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPrefix(t *testing.T) {
	cases := []struct {
		name    string
		network entity.IpNetwork
		expRes  struct {
			prefix string
			err    error
		}
	}{
		{name: "1", network: entity.IpNetwork{
			Ip:   "192.168.1.1",
			Mask: "255.255.255.0",
		}, expRes: struct {
			prefix string
			err    error
		}{prefix: "192.168.1.0", err: nil}},
		{name: "2", network: entity.IpNetwork{
			Ip:   "88.147.254.238",
			Mask: "255.255.255.240",
		}, expRes: struct {
			prefix string
			err    error
		}{prefix: "88.147.254.224", err: nil}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			prefix, err := GetPrefix(testCase.network.Ip, testCase.network.Mask)

			require.Equal(t, testCase.expRes.err, err)

			require.Equal(t, testCase.expRes.prefix, prefix)

		})
	}

}
