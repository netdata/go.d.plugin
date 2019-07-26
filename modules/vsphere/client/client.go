package client

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	datacenter      = "Datacenter"
	folder          = "Folder"
	computeResource = "ComputeResource"
	hostSystem      = "HostSystem"
	virtualMachine  = "VirtualMachine"
)

type Config struct {
	URL      string
	User     string
	Password string
	Timeout  time.Duration
}

func newContainerView(client *vim25.Client, root types.ManagedObjectReference, timeout time.Duration) (*view.ContainerView, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	viewManager := view.NewManager(client)
	return viewManager.CreateContainerView(ctx, root, []string{}, true)
}

func newPerformance(client *vim25.Client, timeout time.Duration) (*performance.Manager, error) {
	perfManager := performance.NewManager(client)
	perfManager.Sort = true

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// populate cache
	_, err := perfManager.CounterInfoByName(ctx)
	if err != nil {
		return nil, err
	}
	return perfManager, nil
}

func New(config Config) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	u, err := soap.ParseURL(config.URL)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("empty URL")
	}
	u.User = url.UserPassword(config.User, config.Password)

	vmomiClient, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return nil, err
	}
	vmomiClient.Timeout = config.Timeout

	containerView, err := newContainerView(vmomiClient.Client, vmomiClient.ServiceContent.RootFolder, config.Timeout)
	if err != nil {
		return nil, err
	}

	perfManager, err := newPerformance(vmomiClient.Client, config.Timeout)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Client:  vmomiClient,
		Perf:    perfManager,
		Root:    containerView,
		Timeout: config.Timeout,
		Lock:    new(sync.RWMutex),
		config:  config,
	}

	return client, nil
}

type Client struct {
	Client  *govmomi.Client
	Root    *view.ContainerView
	Perf    *performance.Manager
	Timeout time.Duration
	Lock    *sync.RWMutex
	config  Config
}

func (c *Client) Reconnect() error {
	cl, err := New(c.config)
	if err != nil {
		return err
	}

	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Client = cl.Client
	c.Root = cl.Root
	c.Perf = cl.Perf
	return nil
}

func (c *Client) IsSessionActive() (bool, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	return c.Client.SessionManager.SessionIsActive(ctx)
}

func (c *Client) Version() string {
	return c.Client.ServiceContent.About.Version
}

func (c *Client) Login() error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return c.Client.SessionManager.Login(ctx, url.UserPassword(c.config.User, c.config.Password))
}

func (c *Client) Logout() error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return c.Client.SessionManager.Logout(ctx)
}

func (c *Client) PerformanceMetrics(pqs []types.PerfQuerySpec) ([]performance.EntityMetric, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	metrics, err := c.Perf.Query(ctx, pqs)
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return c.Perf.ToMetricSeries(ctx, metrics)
}

func (c *Client) Datacenters(pathSet []string) (dcs []mo.Datacenter, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	err = c.Root.Retrieve(ctx, []string{datacenter}, pathSet, &dcs)
	return
}

func (c *Client) Folders(pathSet []string) (folders []mo.Folder, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	err = c.Root.Retrieve(ctx, []string{folder}, pathSet, &folders)
	return
}

func (c *Client) ComputeResources(pathSet []string) (computeResources []mo.ComputeResource, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	err = c.Root.Retrieve(ctx, []string{computeResource}, pathSet, &computeResources)
	return
}

func (c *Client) Hosts(pathSet []string) (hosts []mo.HostSystem, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	err = c.Root.Retrieve(ctx, []string{hostSystem}, pathSet, &hosts)
	return
}

func (c *Client) VirtualMachines(pathSet []string) (vms []mo.VirtualMachine, err error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	err = c.Root.Retrieve(ctx, []string{virtualMachine}, pathSet, &vms)
	return
}

func (c *Client) CounterInfoByName() (map[string]*types.PerfCounterInfo, error) {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	return c.Perf.CounterInfoByName(ctx)
}

func IsNotAuthenticatedError(err error) bool {
	_, ok := soap.ToVimFault(err).(*types.NotAuthenticated)
	return ok
}

//func (c *Client) GetServerTime() (time.Time, error) {
//	ctx, cancel := c.contextWithTimeout()
//	defer cancel()
//	t, err := methods.GetCurrentTime(ctx, c.Client)
//	if err != nil {
//		return time.Time{}, err
//	}
//	return *t, nil
//}
