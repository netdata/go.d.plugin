package pihole

import (
	"github.com/netdata/go.d.plugin/modules/pihole/client"

	"sync"
)

func (p *Pihole) collect() (map[string]int64, error) {
	var (
		summary             *client.SummaryRaw
		queryTypes          *client.QueryTypes
		forwardDestinations *client.ForwardDestinations
		topClients          *client.TopClients
		topItems            *client.TopItems
	)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		var err error
		summary, err = p.client.SummaryRaw()
		if err != nil {
			p.Error(err)
		}
		wg.Done()
	}()
	go func() {
		var err error
		queryTypes, err = p.client.QueryTypes()
		if err != nil {
			p.Error(err)
		}
	}()
	go func() {
		var err error
		forwardDestinations, err = p.client.ForwardDestinations()
		if err != nil {
			p.Error(err)
		}
	}()
	go func() {
		var err error
		topClients, err = p.client.TopClients(5)
		if err != nil {
			p.Error(err)
		}
	}()
	go func() {
		var err error
		topItems, err = p.client.TopItems(5)
		if err != nil {
			p.Error(err)
		}
	}()

	wg.Wait()

	return nil, nil
}
