package phpfpm

import (
	"encoding/json"
	"github.com/netdata/go.d.plugin/pkg/stm"
	fcgiclient "github.com/tomasen/fcgi_client"
	"io/ioutil"
	"log"
)



func (p *Phpfpm) initSocket() error {

	env := make(map[string]string)

	env["SCRIPT_NAME"] = "/status"
	env["SCRIPT_NAME"] = "/status"
	env["SCRIPT_FILENAME"] = "/status"
	env["SERVER_SOFTWARE"] = "go / fcgiclient "
	env["REMOTE_ADDR"] = "127.0.0.1"
	env["QUERY_STRING"] = "json&full"
	env["REQUEST_METHOD"] = "GET"
	env["CONTENT_TYPE"] = "application/json"

	p.env = env

	return nil

}
func (p *Phpfpm) isSocket() bool {
	if len(p.Socket) > 0 {
		return true
	}
	return false
}

func (p *Phpfpm) collectSocket() map[string]int64  {

	socket, err := fcgiclient.Dial("unix", p.Socket)
	if err != nil {
		log.Println("err:", err)
	}
	resp, err := socket.Get(p.env)
	if err != nil {
		log.Println("err:", err)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err:", err)
	}

	st := &status{}

	err2 := json.Unmarshal(content, st)
	if err2 != nil {
		log.Println(err2)
	}

	socket.Close()
	metrics := stm.ToMap(st)

	calcIdleProcessesRequestsDuration(metrics, st.Processes)
	calcIdleProcessesLastRequestCPU(metrics, st.Processes)
	calcIdleProcessesLastRequestMemory(metrics, st.Processes)

	return metrics
}
