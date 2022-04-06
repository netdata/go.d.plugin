package k8s_inventory

type metrics struct {
}

type nodeMetrics struct {
	status struct {
		allocatable struct {
			cpu    int64
			memory int64
		}
	}
}
