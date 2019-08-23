package zookeeper

func (z *Zookeeper) collect() (map[string]int64, error) {
	mx := make(map[string]int64)

	err := z.collectMntr(mx)
	if err != nil {
		return nil, err
	}
	return mx, nil
}

func (z *Zookeeper) collectMntr(mx map[string]int64) error {
	return nil
}
