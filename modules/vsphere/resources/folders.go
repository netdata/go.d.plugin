package resources

type (
	Folder struct {
		Name     string
		ID       string
		ParentID string
	}

	Folders map[string]*Folder
)

func (fs Folders) Put(folder *Folder) {
	fs[folder.ID] = folder
}

func (fs Folders) Get(id string) *Folder {
	return fs[id]
}
