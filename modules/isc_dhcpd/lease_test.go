package isc_dhcpd

import (
	"testing"
	"os"
	"bufio"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseFile(t *testing.T) {
	tests := map[string]struct {
		File string
		numberOfHosts int
		waitFail bool
	} {
		"no file" : {
			File : "testdata/nothing_here.lease",
			numberOfHosts : 0,
			waitFail : true,
		},
		"one host" : {
			File : "testdata/ipv4_dhcpd1.leases",
			numberOfHosts : 1,
			waitFail : false,
		},
		"four hosts" : {
			File : "testdata/ipv4_dhcpd4.leases",
			numberOfHosts : 4,
			waitFail : false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fp, err := os.Open(test.File)
			if test.waitFail {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				defer fp.Close()

				r := bufio.NewReader(fp)

				list := ParseDHCPd(r)
				assert.Equal(t, test.numberOfHosts, len(list))
			}

		})
	}
}