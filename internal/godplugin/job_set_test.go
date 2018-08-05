package godplugin

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/l2isbad/go.d.plugin/internal/godplugin/job"
	"github.com/l2isbad/go.d.plugin/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestJobSet_PutIfNotExist_NoKeyFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	job1 := mock.NewMockJob(ctrl)
	job2 := mock.NewMockJob(ctrl)

	set := jobSet{}
	assert.True(t, set.PutIfNotExist(job1))
	assert.False(t, set.PutIfNotExist(job1))

	assert.True(t, set.PutIfNotExist(job2))
	assert.False(t, set.PutIfNotExist(job2))

	var list []job.Job
	set.Range(func(job job.Job) bool {
		list = append(list, job)
		return true
	})
	assert.Len(t, list, 2)
	assert.Contains(t, list, job1)
	assert.Contains(t, list, job2)
}

func TestJobSet_PutIfNotExist_KeyFunc(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	job11 := mock.NewMockJob(ctrl)
	job11.EXPECT().FullName().Return("job1").AnyTimes()
	job12 := mock.NewMockJob(ctrl)
	job12.EXPECT().FullName().Return("job1").AnyTimes()
	job2 := mock.NewMockJob(ctrl)
	job2.EXPECT().FullName().Return("job2").AnyTimes()

	set := jobSet{KeyFunc: keyFuncFullName}
	assert.True(t, set.PutIfNotExist(job11))
	assert.False(t, set.PutIfNotExist(job11))
	assert.False(t, set.PutIfNotExist(job12))

	assert.True(t, set.PutIfNotExist(job2))
	assert.False(t, set.PutIfNotExist(job2))

	var list []job.Job
	set.Range(func(job job.Job) bool {
		list = append(list, job)
		return true
	})
	assert.Len(t, list, 2)
	assert.Contains(t, list, job11)
	assert.Contains(t, list, job2)
}

func TestJobSet_Exist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	job1 := mock.NewMockJob(ctrl)
	job2 := mock.NewMockJob(ctrl)

	set := jobSet{}
	assert.False(t, set.Exist(job1))
	assert.False(t, set.Exist(job2))

	assert.True(t, set.PutIfNotExist(job1))
	assert.True(t, set.Exist(job1))
	assert.False(t, set.Exist(job2))

	assert.True(t, set.PutIfNotExist(job2))
	assert.True(t, set.Exist(job1))
	assert.True(t, set.Exist(job2))

	set.Delete(job1)
	assert.False(t, set.Exist(job1))
	assert.True(t, set.Exist(job2))
}
