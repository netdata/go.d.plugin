package unity

import (
	"strconv"
	"reflect"
	"time"
	"regexp"

	gounity "github.com/equelin/gounity"
	
	types "github.com/equelin/gounity/types/v1"
)

const (
	cpuUtilization	= "cpuUtilization"
	readBandwith = "readBandwith"
	writeBandwith = "writeBandwith"
)

type Config struct{
	Servers []Server `json:"servers"`
	Insecure bool `json:"insecure"`
	Username string `json:"username"`
	Password string `json:"password"`
	Interval uint32 `json:"interval"`
	Metrics Metrics `json:"metrics"`
}
type Server struct{
	Name string `json:"name"`
	Adress string `json:"adress"`
	Targets Targets `json:"targets"`
}
type Targets struct{
	LUNs []string `json:"lun"`
	SPfc []string `json:"fc"`
}
type Metrics struct{
	General []string `json:"general"`
	FC []string `json:"fc"`
	LUN []string `json:"lun"`
}

func getPaths(targets Targets, metrics Metrics) ([]string){
	var paths []string

	for _,m := range metrics.General{
		paths = append(paths,"kpi.sp.spa."+m)
		paths = append(paths,"kpi.sp.spb."+m)
	}
  
	for _,fc := range targets.SPfc{
	  for _,fc_m := range metrics.FC{
		paths = append(paths,"kpi.fibreChannel."+fc+"."+fc_m)
	  }
	}
  
	for _,sv := range targets.LUNs{
	  for _,sv_m := range metrics.LUN{
		paths = append(paths,"kpi.lun."+sv+".sp.spa."+sv_m)
		paths = append(paths,"kpi.lun."+sv+".sp.spb."+sv_m)
	  }
	}

	return paths
}

func getFromNested(t interface{}, path string) (int64,error){
	tt := reflect.ValueOf(t).MapRange()
	r, _ := regexp.Compile(".*lun.*(B|b)andwidth$")
	for tt.Next() {
		v := tt.Value().Interface()
		m,b := v.(string)
		if(b){
			val,_ := strconv.ParseFloat(m, 64)
			if r.MatchString(path){ //MB => KB
				return int64(val*1000),nil
			}
			return int64(val),nil
		}
		return getFromNested(v, path)
	}
	return -1,nil
}

func parseStatus(result types.MetricRealTimeQueryResult, prefix string) (map[string]int64, error){
	var metrics = make(map[string]int64)
	for _, entry := range result.Entries{
		val, err := getFromNested(entry.Content.Values, entry.Content.Path)
		if err != nil{
			val = int64(-1)
		}
		metrics[prefix+"."+entry.Content.Path] = val
	}
	return metrics, nil
}

func (u *Unity) getStatus() (map[string]int64, error) {
	var allMetrics = make(map[string]int64)

	for _,server := range u.config.Servers{
		session, err := gounity.NewSession(server.Adress, u.config.Insecure, u.config.Username, u.config.Password)
		if err != nil {
			u.Error("error on creating session : %v", err)
		}

		Metric, err := session.NewMetricRealTimeQuery(getPaths(server.Targets,u.config.Metrics), u.config.Interval)
		if err != nil {
			u.Error("error on querying metrics : %v", err)
		}
	
		time.Sleep(time.Duration(Metric.Content.Interval) * time.Second)
	
		Result, err := session.GetMetricRealTimeQueryResult(Metric.Content.ID)
		if err != nil {
			u.Error("error on getting metrics from the query : %v", err)
		}
		metrics, err := parseStatus(*Result,server.Adress)
		if err != nil {
			u.Error("cannot parse status : %v", err)
		}	
		for k,v := range metrics{
			allMetrics[k]=v
		}
	}

	return allMetrics, nil
}
