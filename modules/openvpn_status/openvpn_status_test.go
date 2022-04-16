package openvpn_status

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	v1statusLog = `OpenVPN CLIENT LIST
Updated,2022-04-17 00:33:28
Common Name,Real Address,Bytes Received,Bytes Sent,Connected Since
gofle,192.168.122.135:58999,19265,261631,2022-04-17 00:33:19
client_bsd2,192.168.122.136:11445,2752,3545,2022-04-17 00:33:04
ROUTING TABLE
Virtual Address,Common Name,Real Address,Last Ref
10.8.0.2,gofle,192.168.122.135:58999,2022-04-17 00:33:27
10.8.0.3,client_bsd2,192.168.122.136:11445,2022-04-17 00:33:04
GLOBAL STATS
Max bcast/mcast queue length,2
END
`
	v2statusLog = `TITLE,OpenVPN 2.5.5 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [LZ4] [EPOLL] [PKCS11] [MH/PKTINFO] [AEAD] built on Mar 22 2022
TIME,2022-03-31 11:47:46,1648716466
HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID,Data Channel Cipher
CLIENT_LIST,gofle,192.168.122.135:59523,10.8.0.2,,19265,261631,2022-03-31 11:47:34,1648716454,UNDEF,10,0,AES-128-GCM
CLIENT_LIST,client_bsd2,192.168.122.136:20053,10.8.0.3,,2752,3545,2022-04-17 00:02:39,1650142959,UNDEF,0,0,AES-128-GCM
HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
ROUTING_TABLE,10.8.0.2,gofle,192.168.122.135:59523,2022-03-31 11:47:37,1648716457
ROUTING_TABLE,10.8.0.3,client_bsd2,192.168.122.136:20053,2022-04-17 00:02:39,1650142959
GLOBAL_STATS,Max bcast/mcast queue length,1
END`
	v3statusLog = `TITLE	OpenVPN 2.5.5 x86_64-pc-linux-gnu [SSL (OpenSSL)] [LZO] [LZ4] [EPOLL] [PKCS11] [MH/PKTINFO] [AEAD] built on Mar 22 2022
TIME	2022-03-31 11:53:29	1648716809
HEADER	CLIENT_LIST	Common Name	Real Address	Virtual Address	Virtual IPv6 Address	Bytes Received	Bytes Sent	Connected Since	Connected Since (time_t)	UsernameClient ID	Peer ID	Data Channel Cipher
CLIENT_LIST	gofle	192.168.122.135:49199	10.8.0.2		19265	261631	2022-03-31 11:53:19	1648716799	UNDEF	0	0	AES-128-GCM
CLIENT_LIST     client_bsd2     192.168.122.136:21748   10.8.0.3       2001:0db8:85a3:0000:0000:8a2e:0370:7334         2752   3545  2022-04-16 20:48:27     1650131307      UNDEF   624     0       AES-128-GCM
HEADER	ROUTING_TABLE	Virtual Address	Common Name	Real Address	Last Ref	Last Ref (time_t)
ROUTING_TABLE	10.8.0.2	gofle	192.168.122.135:49199	2022-03-31 11:53:24	1648716804
ROUTING_TABLE   10.8.0.3        client_bsd2     192.168.122.136:21748   2022-04-16 20:48:27     1650131307
GLOBAL_STATS	Max bcast/mcast queue length	0
END`
)

var statusLogList = []string{v1statusLog, v2statusLog, v3statusLog}

var expected = map[string]int64{
	"bytes_in":              22017,
	"bytes_out":             265176,
	"clients":               2,
	"gofle_bytes_in":        19265,
	"gofle_bytes_out":       261631,
	"client_bsd2_bytes_in":  2752,
	"client_bsd2_bytes_out": 3545,
}

func createTempFile(data string) (string, func()) {
	file, err := ioutil.TempFile("", "my_test_case_")
	if err != nil {
		fmt.Printf("%v", err)
	}
	cleanup := func() {
		_ = file.Close()
		_ = os.RemoveAll(file.Name())
	}
	if _, err := file.Write([]byte(data)); err != nil {
		fmt.Printf("%v", err)
	}
	return file.Name(), cleanup
}

func TestNew(t *testing.T) {
	job := New()

	assert.Implements(t, (*module.Module)(nil), job)
	assert.Equal(t, defaultFilePath, job.StatusPath)
	assert.NotNil(t, job.charts)
	assert.NotNil(t, job.collectedUsers)
}

func TestOpenVPN_Status_Init(t *testing.T) {
	assert.True(t, New().Init())
}

func TestOpenVPN_Status_Collect(t *testing.T) {
	for _, logData := range statusLogList {
		statusPath, cleanup := createTempFile(logData)
		defer cleanup()
		job := New()
		job.StatusPath = statusPath

		require.True(t, job.Init())
		job.perUserMatcher = matcher.TRUE()
		require.True(t, job.Check())

		mx := job.Collect()
		require.NotNil(t, mx)
		delete(mx, "gofle_connection_time")
		delete(mx, "client_bsd2_connection_time")
		assert.Equal(t, expected, mx)
	}
}
