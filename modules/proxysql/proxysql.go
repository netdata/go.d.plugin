package proxysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("proxysql", module.Creator{
		Create: func() module.Module { return New() },
	})
}

func New() *ProxySQL {
	return &ProxySQL{
		Config: Config{
			DSN:     "stats:stats@tcp(127.0.0.1:6032)/",
			Timeout: web.Duration{Duration: time.Second * 2},
		},

		charts: baseCharts.Copy(),

		commands: make(map[string]bool),
		users:    make(map[string]bool),
	}
}

type Config struct {
	DSN     string       `yaml:"dsn"`
	MyCNF   string       `yaml:"my.cnf"`
	Timeout web.Duration `yaml:"timeout"`
}

type ProxySQL struct {
	module.Base
	Config `yaml:",inline"`

	db *sql.DB

	charts *Charts

	commands map[string]bool
	users    map[string]bool
}

func (p *ProxySQL) Init() bool {
	if p.MyCNF != "" {
		dsn, err := dsnFromFile(p.MyCNF)
		if err != nil {
			p.Error(err)
			return false
		}
		p.DSN = dsn
	}

	if p.DSN == "" {
		p.Error("'dsn' not set")
		return false
	}

	p.Debugf("using DSN [%s]", p.DSN)
	return true
}

func (p *ProxySQL) Check() bool {
	return len(p.Collect()) > 0
}

func (p *ProxySQL) Charts() *Charts {
	return p.charts
}

func (p *ProxySQL) Collect() map[string]int64 {
	mx, err := p.collect()
	if err != nil {
		p.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (p *ProxySQL) Cleanup() {
	if p.db == nil {
		return
	}
	if err := p.db.Close(); err != nil {
		p.Errorf("cleanup: error on closing the ProxySQL instance [%s]: %v", p.DSN, err)
	}
	p.db = nil
}
