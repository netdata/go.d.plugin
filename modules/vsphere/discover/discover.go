package discover

import (
	"fmt"
	"strings"
	"time"

	rs "github.com/netdata/go.d.plugin/modules/vsphere/resources"

	"github.com/netdata/go-orchestrator/logger"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type APIClient interface {
	Datacenters(pathSet ...string) ([]mo.Datacenter, error)
	Folders(pathSet ...string) ([]mo.Folder, error)
	ComputeResources(pathSet ...string) ([]mo.ComputeResource, error)
	Hosts(pathSet ...string) ([]mo.HostSystem, error)
	VirtualMachines(pathSet ...string) ([]mo.VirtualMachine, error)

	CounterInfoByName() (map[string]*types.PerfCounterInfo, error)
}

func NewVSphereDiscoverer(client APIClient) *vSphereDiscoverer {
	return &vSphereDiscoverer{
		APIClient: client,
	}
}

type HostMatcher interface {
	Match(*rs.Host) bool
}

type VMMatcher interface {
	Match(*rs.VM) bool
}

type vSphereDiscoverer struct {
	*logger.Logger
	APIClient
	HostMatcher
	VMMatcher
}

type resources struct {
	dcs      []mo.Datacenter
	folders  []mo.Folder
	clusters []mo.ComputeResource
	hosts    []mo.HostSystem
	vms      []mo.VirtualMachine
}

func (d vSphereDiscoverer) Discover() (*rs.Resources, error) {
	startTime := time.Now()
	raw, err := d.discover()
	if err != nil {
		return nil, fmt.Errorf("discovering resources : %v", err)
	}

	res := d.build(raw)

	err = d.setHierarchy(res)
	if err != nil {
		// TODO: handle objects w/o hier?
		d.Error(err)
	}

	numH := len(res.Hosts)
	numV := len(res.VMs)
	removed := d.removeUnmatched(res)
	if removed == (numH + numV) {
		return nil, fmt.Errorf("all resoursces were filtered (%d hosts, %d vms)", numH, numV)
	}

	err = d.collectMetricLists(res)
	if err != nil {
		return nil, fmt.Errorf("collecting metric lists : %v", err)
	}

	d.Infof("discovering : %d hosts, %d vms, the whole process took %s",
		len(res.Hosts),
		len(res.VMs),
		time.Since(startTime))

	return res, nil
}

var (
	// properties to set
	datacenterPathSet = []string{"name", "parent"}
	folderPathSet     = []string{"name", "parent"}
	clusterPathSet    = []string{"name", "parent"}
	hostPathSet       = []string{"name", "parent", "runtime.powerState"}
	vmPathSet         = []string{"name", "runtime.host", "runtime.powerState"}
)

func (d vSphereDiscoverer) discover() (*resources, error) {
	d.Debug("discovering : starting resource discovering")

	start := time.Now()
	t := start
	datacenters, err := d.Datacenters(datacenterPathSet...)
	if err != nil {
		return nil, err
	}
	d.Debugf("discovering : found %d datacenters, discovering took %s", len(datacenters), time.Since(t))

	t = time.Now()
	folders, err := d.Folders(folderPathSet...)
	if err != nil {
		return nil, err
	}
	d.Debugf("discovering : found %d folders, discovering took %s", len(folders), time.Since(t))

	t = time.Now()
	clusters, err := d.ComputeResources(clusterPathSet...)
	if err != nil {
		return nil, err
	}
	d.Debugf("discovering : found %d clusters, discovering took %s", len(clusters), time.Since(t))

	t = time.Now()
	hosts, err := d.Hosts(hostPathSet...)
	if err != nil {
		return nil, err
	}
	d.Debugf("discovering : found %d hosts, discovering took %s", len(hosts), time.Since(t))

	t = time.Now()
	vms, err := d.VirtualMachines(vmPathSet...)
	if err != nil {
		return nil, err
	}
	d.Debugf("discovering : found %d vms, discovering took %s", len(hosts), time.Since(t))

	raw := resources{
		dcs:      datacenters,
		folders:  folders,
		clusters: clusters,
		hosts:    hosts,
		vms:      vms,
	}

	d.Infof("discovering : found %d datacenters, %d folders, %d clusters (%d dummy), %d hosts, %d vms, discovering took %s",
		len(raw.dcs),
		len(raw.folders),
		len(clusters),
		numOfDummyClusters(clusters),
		len(raw.hosts),
		len(raw.vms),
		time.Since(start),
	)

	return &raw, nil
}

func numOfDummyClusters(clusters []mo.ComputeResource) (num int) {
	for _, c := range clusters {
		// domain-s61 | domain-c52
		if strings.HasPrefix(c.Reference().Value, "domain-s") {
			num++
		}
	}
	return num
}
