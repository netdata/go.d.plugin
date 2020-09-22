package filecheck

func (fc *Filecheck) collect() (map[string]int64, error) {
	mx := make(map[string]int64)

	fc.collectFiles(mx)
	fc.collectDirs(mx)

	return mx, nil
}
