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