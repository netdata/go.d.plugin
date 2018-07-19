package cooked

type flagsChart struct {
	push      bool
	created   bool
	updated   bool
	obsoleted bool
}

func (f *flagsChart) setPush(b bool) {
	f.push = b
}

func (f *flagsChart) setCreated(b bool) {
	f.created = b
}

func (f *flagsChart) setUpdated(b bool) {
	f.updated = b
}

func (f *flagsChart) setObsoleted(b bool) {
	f.obsoleted = b
}

// isPush returns whether we need to send a chart to netdata.
func (f flagsChart) isPush() bool {
	return f.push
}

// isCreated returns whether a chart was created.
func (f flagsChart) isCreated() bool {
	return f.created
}

// isUpdated returns whether a chart was updated on previous update.
func (f flagsChart) isUpdated() bool {
	return f.updated
}

// IsObsoleted returns whether a chart was obsoleted.
func (f flagsChart) IsObsoleted() bool {
	return f.obsoleted
}

type flagsDim struct {
	push       bool
	retries    int
	retriesMax int
}

func (f flagsDim) alive() bool {
	return f.retries < f.retriesMax
}
