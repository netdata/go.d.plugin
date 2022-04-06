package k8s_inventory

import (
	"errors"
	"time"
)

func (ki *KubernetesInventory) collect() (map[string]int64, error) {
	if ki.discoverer == nil {
		return nil, errors.New("nil discoverer")
	}

	ki.once.Do(func() {
		in := make(chan resource)

		ki.wg.Add(1)
		go func() { defer ki.wg.Done(); ki.runCollectKubernetes(ki.ctx, in) }()

		ki.wg.Add(1)
		go func() { defer ki.wg.Done(); ki.discoverer.run(ki.ctx, in) }()
	})

	time.Sleep(time.Second * 20)
	ki.ctxCancel()
	ki.wg.Wait()
	<-ki.ctx.Done()

	return nil, nil
}
