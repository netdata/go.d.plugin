package unity

import (
	"fmt"
	"regexp"
)

func (u *Unity) collect() (map[string]int64, error) {
	temp, err := u.getStatus()

	if err != nil {
		return nil, err
	}
	if len(temp) == 0 {
		return nil, fmt.Errorf("nothing was collected")
	}

	mx := make(map[string]int64)
	toUnderscore,_ := regexp.Compile("\\.|-")
	r, _ := regexp.Compile("lun")
	for k,v := range temp{
		if r.MatchString(k) {
			id := getID(k)
			name := u.lunsNames[getIP(k)][id]
			r2, _ := regexp.Compile(id)
			mx[toUnderscore.ReplaceAllString(r2.ReplaceAllString(k,name),"_")]=v
		}else{
			mx[toUnderscore.ReplaceAllString(k,"_")]=v
		}
	}
	u.Infof("%d metrics collected",len(mx))
	return mx, nil
}

func getIP(path string) string{
	r, _ := regexp.Compile("(.*)?.kpi")
	path = r.FindString(path)
	match := r.FindStringSubmatch(path)
	if len(match)==2{
		return match[1]
	}else{
		return ""
	}
}

func getID(path string) string{
	r, _ := regexp.Compile("(.*lun.)(.*?)(.sp.*)")
	match := r.FindStringSubmatch(path)
	if len(match)==4{
		return match[2]
	}else{
		return ""
	}
}