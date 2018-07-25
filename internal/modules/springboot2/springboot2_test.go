package springboot2

import (
	"io/ioutil"
	"testing"
)

var testdata, _ = ioutil.ReadFile("tests/testdata.txt")

func TestSpringboot2(t *testing.T) {
	// TODO: Since modules.Charts and modules.Logger are injected by job_set.go, it's hard to do unit test.
	// https://github.com/l2isbad/go.d.plugin/issues/12

	// ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	if r.URL.Path == "/actuator/prometheus" {
	// 		w.Write(testdata)
	// 		return
	// 	}
	// }))
	// defer ts.Close()
	// plugin := &Springboot2{
	// 	Request: web.Request{
	// 		URL: ts.URL + "/actuator/prometheus",
	// 	},
	// }

	// assert.True(t, plugin.Check())

	// data := plugin.GetData()
	// assert.EqualValues(t, map[string]int64{}, data)
}
