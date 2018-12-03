package godplugin

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoopQueue_add(t *testing.T) {
	var l loopQueue
	var wg sync.WaitGroup

	workers := 10
	addNum := 1000

	f := func() {
		for i := 0; i < addNum; i++ {
			l.add(nil)
		}
		wg.Done()
	}

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go f()
	}

	wg.Wait()

	assert.Equal(t, workers*addNum, len(l.queue))
}

func TestLoopQueue_remove(t *testing.T) {

}
