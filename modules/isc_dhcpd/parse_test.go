package isc_dhcpd

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseFile(t *testing.T) {
	tests := map[string]struct {
		fileName string
		wantLeases []leaseEntry
	} {
		"backup" : {
			fileName : "testdata/backup.leases",
			wantLeases : []leaseEntry{
				{
					ip : "192.168.0.100",
					ends : "1 2020/09/28 19:58:46",
					bindingState: "active",
				},
			},
		},
		"ipv4" : {
			fileName : "testdata/ipv4_dhcpd.leases",
			wantLeases : []leaseEntry{
				{
					ip : "10.220.252.2",
					ends : "3 2020/09/15 09:12:16",
					bindingState: "active",
				},
				{
					ip : "10.220.252.3",
					ends : "3 2020/09/15 07:29:01",
					bindingState: "free",
				},
				{
					ip : "10.220.252.4",
					ends : "epoch 1600137200",
					bindingState: "active",
				},
				{
					ip : "10.220.252.5",
					ends : "3 2020/09/15 01:33:19",
					bindingState: "free",
				},
			},
		},
		"ipv6" : {
			fileName : "testdata/ipv6_dhcpd.leases",
			wantLeases : []leaseEntry{
				{
					ip : "1985:470:1f0b:c9a::000",
					ends : "2 2020/09/30 10:53:29",
					bindingState: "active",
				},
				{
					ip : "1985:470:1f0b:c9a::001",
					ends : "2 2020/09/30 23:59:58",
					bindingState: "active",
				},
				{
					ip : "1985:470:1f0b:c9a::002",
					ends : "2 2020/09/30 02:11:08",
					bindingState: "active",
				},
				{
					ip : "1985:470:1f0b:c9a::003",
					ends : "2 2020/09/30 18:48:39",
					bindingState: "active",
				},
				{
					ip : "1985:470:1f0b:c9a::004",
					ends : "2 2020/09/30 14:53:15",
					bindingState: "active",
				},
				{
					ip : "1985:470:1f0b:c9a::005",
					ends : "2 2020/09/30 11:33:17",
					bindingState: "active",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fp, err := os.Open(test.fileName)
			var leases []leaseEntry
			require.NoError(t, err)
			defer fp.Close()

			r := bufio.NewReader(fp)

			list := parseDHCPdLeases(leases, r)
			assert.Equal(t, test.wantLeases, list)
		})
	}
}