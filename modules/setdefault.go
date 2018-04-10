package modules

type globalVar int

const (
	UpdateEvery globalVar = iota
	ChartCleanup
)

var moduleDefault = make(map[string]map[globalVar]int)

type setter struct {
	name     string
	variable globalVar
}

func (u *setter) Set(v int) {
	if _, ok := moduleDefault[u.name]; !ok {
		moduleDefault[u.name] = make(map[globalVar]int)
	}
	moduleDefault[u.name][u.variable] = v
}

func SetDefault(gv globalVar) *setter {
	name := getFileName(2)
	switch gv {
	case UpdateEvery:
		return &setter{name, UpdateEvery}
	case ChartCleanup:
		return &setter{name, ChartCleanup}
	default:
		return &setter{name, globalVar(-1)}
	}
}

type getter struct {
	name string
}

func (g *getter) Get(gv globalVar) (int, bool) {
	v, ok := moduleDefault[g.name][gv]
	return v, ok
}

func GetDefault(module string) *getter {
	return &getter{name: module}
}
