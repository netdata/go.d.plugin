package resources

type (
	ClusterHierarchy struct {
		Dc HierarchyValue
	}

	// When two or more physical machines are grouped to work and be managed as a whole,
	// the aggregate computing and memory resources form a cluster.
	Cluster struct {
		Name     string
		ID       string
		ParentID string
		Hier     ClusterHierarchy
	}

	Clusters map[string]*Cluster
)

func (h ClusterHierarchy) IsSet() bool {
	return h.Dc.IsSet()
}

func (cs Clusters) Put(cluster *Cluster) {
	cs[cluster.ID] = cluster
}

func (cs Clusters) Get(id string) *Cluster {
	return cs[id]
}
