package nginxvts

type vtsStatus struct {
	// HostName     string
	// NginxVersion string
	LoadMsec    int64
	NowMsec     int64
	Connections struct {
		Active   int64 `stm:"active"`
		Reading  int64 `stm:"reading"`
		Writing  int64 `stm:"writing"`
		Waiting  int64 `stm:"waiting"`
		Accepted int64 `stm:"accepted"`
		Handled  int64 `stm:"handled"`
		Request  int64 `stm:"requests"`
	} `stm:"connections"`
	SharedZones struct {
		// Name     string
		MaxSize  int64 `stm:"maxsize"`
		UsedSize int64 `stm:"usedsize"`
		UsedNode int64 `stm:"usednode"`
	}
	ServerZones   map[string]Server
	UpstreamZones map[string][]Upstream
	FilterZones   map[string]map[string]Server
	CacheZones    map[string]Cache
}

func (v vtsStatus) hasServerZones() bool   { return v.ServerZones != nil }
func (v vtsStatus) hasUpstreamZones() bool { return v.UpstreamZones != nil }
func (v vtsStatus) hasFilterZones() bool   { return v.FilterZones != nil }
func (v vtsStatus) hasCacheZones() bool    { return v.CacheZones != nil }

// Server is Nginx virtual host name
type Server struct {
	RequestCounter int64 `stm:"requestcounter"`
	InBytes        int64 `stm:"inbytes"`
	OutBytes       int64 `stm:"outbytes"`
	Responses      struct {
		OneXx       int64 `stm:"1xx" json:"1xx"`
		TwoXx       int64 `stm:"2xx" json:"2xx"`
		ThreeXx     int64 `stm:"3xx" json:"3xx"`
		FourXx      int64 `stm:"4xx" json:"4xx"`
		FiveXx      int64 `stm:"5xx" json:"5xx"`
		Miss        int64 `stm:"miss"`
		Bypass      int64 `stm:"bypass"`
		Expired     int64 `stm:"expired"`
		Stale       int64 `stm:"stale"`
		Updating    int64 `stm:"updating"`
		Revalidated int64 `stm:"revalidated"`
		Hit         int64 `stm:"hit"`
		Scarce      int64 `stm:"scarce"`
	} `stm:"responses"`
}

// Upstream is Nginx proxy upstream
type Upstream struct {
	Server         string
	RequestCounter int64 `stm:"requestcounter"`
	InBytes        int64 `stm:"inbytes"`
	OutBytes       int64 `stm:"outbytes"`
	Responses      struct {
		OneXx   int64 `stm:"1xx" json:"1xx"`
		TwoXx   int64 `stm:"2xx" json:"2xx"`
		ThreeXx int64 `stm:"3xx" json:"3xx"`
		FourXx  int64 `stm:"4xx" json:"4xx"`
		FiveXx  int64 `stm:"5xx" json:"5xx"`
	} `stm:"responses"`
}

// Cache is Nginx proxy cache
type Cache struct {
	MaxSize   int64 `stm:"maxsize"`
	UsedSize  int64 `stm:"usedsize"`
	InBytes   int64 `stm:"inbytes"`
	OutBytes  int64 `stm:"outbytes"`
	Responses struct {
		Miss        int64 `stm:"miss"`
		Bypass      int64 `stm:"bypass"`
		Expired     int64 `stm:"expired"`
		Stale       int64 `stm:"stale"`
		Updating    int64 `stm:"updating"`
		Revalidated int64 `stm:"revalidated"`
		Hit         int64 `stm:"hit"`
		Scarce      int64 `stm:"scarce"`
	} `stm:"responses"`
}
