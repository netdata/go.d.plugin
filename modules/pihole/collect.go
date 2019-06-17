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
	wg.Add(1)

	go func() {
		p.collectSummary(&summary)
		wg.Done()
	}()

	if p.Password != "" {
		wg.Add(4)
		go func() {
			p.collectQueryTypes(&queryTypes)
			wg.Done()
		}()

		go func() {
			p.collectForwardDestination(&forwardDestinations)
			wg.Done()
		}()

		go func() {
			p.collectTopClients(&topClients)
			wg.Done()
		}()

		go func() {
			p.collectTopItems(&topItems)
			wg.Done()
		}()
	}

	wg.Wait()

	return nil, nil
}

func (p *Pihole) collectSummary(s **client.SummaryRaw) {
	var err error
	if *s, err = p.client.SummaryRaw(); err != nil {
		p.Error(err)
	}
}

func (p *Pihole) collectQueryTypes(qt **client.QueryTypes) {
	var err error
	if *qt, err = p.client.QueryTypes(); err != nil {
		p.Error(err)
	}
}

func (p *Pihole) collectForwardDestination(fd **client.ForwardDestinations) {
	var err error
	if *fd, err = p.client.ForwardDestinations(); err != nil {
		p.Error(err)
	}
}

func (p *Pihole) collectTopClients(tc **client.TopClients) {
	var err error
	if *tc, err = p.client.TopClients(5); err != nil {
		p.Error(err)
	}
}

func (p *Pihole) collectTopItems(ti **client.TopItems) {
	var err error
	if *ti, err = p.client.TopItems(5); err != nil {
		p.Error(err)
	}
}
