package godplugin

import (
	"sync"

	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
)

type jobSet struct {
	KeyFunc func(job job.Job) interface{}
	m       sync.Map
}

func keyFuncFullName(job job.Job) interface{} {
	return job.FullName()
}

func (s *jobSet) PutIfNotExist(job job.Job) bool {
	_, loaded := s.m.LoadOrStore(s.getKey(job), job)
	return !loaded
}

func (s *jobSet) Delete(job job.Job) {
	s.m.Delete(s.getKey(job))
}

func (s *jobSet) Exist(job job.Job) bool {
	_, exist := s.m.Load(s.getKey(job))
	return exist
}

func (s *jobSet) Range(f func(job job.Job) bool) {
	s.m.Range(func(k, v interface{}) bool {
		return f(v.(job.Job))
	})
}

func (s *jobSet) getKey(job job.Job) interface{} {
	if s.KeyFunc != nil {
		return s.KeyFunc(job)
	}
	return job
}
