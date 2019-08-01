package resources

type (
	Datacenter struct {
		Name string
		ID   string
	}

	Dcs map[string]*Datacenter
)

func (dcs Dcs) Put(datacenter *Datacenter) {
	dcs[datacenter.ID] = datacenter
}

func (dcs Dcs) Get(id string) *Datacenter {
	return dcs[id]
}
