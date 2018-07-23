package springboot2

import "github.com/l2isbad/go.d.plugin/internal/modules"

type Springboot2 struct {
	modules.Charts
	modules.Logger
	modules.NoConfiger

	data map[string]int64
}

func (w *Springboot2) Check() bool {
	return false
}

func (w *Springboot2) GetData() map[string]int64 {
	return nil
}

func init() {
	f := func() modules.Module {
		return &Springboot2{
			data: make(map[string]int64),
		}
	}
	modules.Add(f)
}
