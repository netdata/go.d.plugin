package unity

import (
	"github.com/netdata/go-orchestrator/module"
	"os"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

func init() {
	creator := module.Creator{
		Create: func() module.Module { return New() },
	}
	module.Register("unity", creator)
}

func New() *Unity {
	return &Unity{
		charts: charts.Copy(),
		lunsNames: make(map[string]map[string]string),
	}
}

type Unity struct {
	module.Base // should be embedded by every module
	charts    *Charts
	lunsNames   map[string]map[string]string
	config Config
}

func (Unity) Cleanup() {}

func (u *Unity) Init() bool {
	var config Config //start load config
	jsonFile, err := os.Open("/etc/netdata/go.d/unity.json")
	if err != nil {
		u.Error(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	u.config = config //end load config

	u.lunsNames, err = u.getLuns() //get the lun id/name dictionary


	lunsIDs := make(map[string]map[string]string, len(u.lunsNames))
	for k, v := range u.lunsNames {
		lunsIDs[k] = make(map[string]string, len(u.lunsNames[k]))
		for kk, vv := range v{
			lunsIDs[k][vv] = kk
		}
    }
	
	for k,server := range u.config.Servers{
		for kk,lunName := range server.Targets.LUNs{
			u.config.Servers[k].Targets.LUNs[kk] = lunsIDs[server.Adress][lunName]
		}
	}
	return true
}

func (u *Unity) Check() bool {
	return len(u.Collect()) > 0
}

func (u *Unity) Charts() *module.Charts {
	charts := module.Charts{}
	toUnderscore,_ := regexp.Compile("\\.|-")
	r, _ := regexp.Compile("lun")
	for _,server := range u.config.Servers{
		for _,path := range getPaths(server.Targets,u.config.Metrics){
			if r.MatchString(path){
				id := getID(path)
				name := u.lunsNames[server.Adress][id]
				r2, _ := regexp.Compile(id)
				path = r2.ReplaceAllString(path,name)
			}
			path = toUnderscore.ReplaceAllString(server.Adress+"_"+path,"_")
			charts.Add(
				&module.Chart{
					ID:path,
					Title:path,
					Units:"u",
					Fam:toUnderscore.ReplaceAllString(server.Adress,"_"),
					Dims:module.Dims{
						&module.Dim{
							ID: path,
							Name: path,
						},
					},
				},
			)
		}
	}
	return &charts
}

func (u *Unity) Collect() map[string]int64 {
	mx, err := u.collect()

	if err != nil {
		u.Error(err)
		return nil
	}

	return mx
}