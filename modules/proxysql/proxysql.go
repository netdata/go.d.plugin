package proxysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/netdata/go.d.plugin/agent/module"
)

func init() {
	module.Register("proxysql", module.Creator{
		Create: func() module.Module { return New() },
	})
}

type (
	Config struct {
		DSN         string `yaml:"dsn"`
		MyCNF       string `yaml:"my.cnf"`
		UpdateEvery int    `yaml:"update_every"`
	}

	ProxySQL struct {
		module.Base
		Config `yaml:",inline"`

		db *sql.DB

		charts *Charts
	}
)

func New() *ProxySQL {
	return &ProxySQL{
		Config: Config{
			DSN: "stats:stats@tcp(127.0.0.1:6032)/",
		},

		charts: charts.Copy(),
	}
}

func (p *ProxySQL) Cleanup() {
	if p.db == nil {
		return
	}
	if err := p.db.Close(); err != nil {
		p.Errorf("cleanup: error on closing the proxysql instance [%s]: %v", p.DSN, err)
	}
	p.db = nil
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
		p.Error("DSN not set")
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
