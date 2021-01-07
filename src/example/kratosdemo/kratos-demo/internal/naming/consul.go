package consul

// 参考自 https://github.com/cocktail18/kratos-consul-discovery/blob/master/consul.go
// https://github.com/go-kratos/kratos/blob/master/pkg/naming/etcd/etcd.go
import (
	"context"
	"log"
	"sync"
	"sync/atomic"

	"github.com/go-kratos/kratos/pkg/naming"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

// Resolver resolve naming service
type Resolver struct {
	appID   string
	c       chan struct{}
	client  *api.Client
	agent   *api.Agent
	plan    *watch.Plan
	builder *Builder
	ins     atomic.Value
}

// Fetch fetch resolver instance.
func (r *Resolver) Fetch(ctx context.Context) (ins *naming.InstancesInfo, ok bool) {
	v := r.ins.Load()
	ins, ok = v.(*naming.InstancesInfo)
	return
}

func (r *Resolver) Watch() <-chan struct{} {
	return r.c
}
func (r *Resolver) watch() error {
	// TODO
	return nil
}
func (r *Resolver) Close() error {
	if r.plan != nil && !r.plan.IsStopped() {
		r.plan.Stop()
	}
	return nil
}

type Builder struct {
	client *api.Client
	agent  *api.Agent
	r      map[string]*Resolver
	locker sync.RWMutex
	c      *Config
}

// Config discovery configures.
type Config struct {
	Nodes  []string
	Region string
	Zone   string
	Env    string
	Host   string
}

func (builder Builder) Build(id string, options ...naming.BuildOpt) naming.Resolver {
	builder.locker.RLock()
	if r, ok := builder.r[id]; ok {
		builder.locker.RUnlock()
		return r
	}
	builder.locker.RUnlock()
	builder.locker.Lock()
	r := &Resolver{
		appID:   id,
		client:  builder.client,
		agent:   builder.agent,
		builder: &builder,
	}
	r.c = make(chan struct{}, 10)
	builder.r[id] = r
	builder.locker.Unlock()
	err := r.watch()
	if err != nil {
		log.Fatalf("watch error %s", err.Error())
	}
	return r
}

func (builder Builder) Scheme() string {
	return "consul"
}

func NewConsulDiscovery(c Config) (builder Builder, err error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return
	}
	builder.client = client
	builder.agent = client.Agent()
	builder.r = make(map[string]*Resolver)
	builder.c = &c
	return
}
