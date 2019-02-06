package bind

type xml3Stats struct {
	Server xml3Server `xml:"server"`
	Views  []xml3View `xml:"views>view"`
}

type xml3Server struct {
	CounterGroups []xml3CounterGroup `xml:"counters"`
}

type xml3CounterGroup struct {
	Type     string `xml:"type,attr"`
	Counters []struct {
		Name  string `xml:"name,attr"`
		Value int    `xml:",chardata"`
	} `xml:"counter"`
}

type xml3View struct {
	// Omitted branches: zones
	Name          string             `xml:"name,attr"`
	CounterGroups []xml3CounterGroup `xml:"counters"`
	Caches        []struct {
		Name   string `xml:"name,attr"`
		RRSets []struct {
			Name  string `xml:"name"`
			Value int    `xml:"counter"`
		} `xml:"rrset"`
	} `xml:"cache"`
}

type xml3Client struct{}
