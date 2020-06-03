module github.com/netdata/go.d.plugin

go 1.14

require (
	github.com/Wing924/ltsv v0.3.1
	github.com/axiomhq/hyperloglog v0.0.0-20191112132149-a4c4c47bc57f
	github.com/cloudflare/cfssl v1.4.1
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/likexian/whois-go v1.7.1
	github.com/likexian/whois-parser-go v1.14.5
	github.com/miekg/dns v1.1.29
	github.com/netdata/go-orchestrator v0.0.0-20200603131224-cb4c839115c3
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/prometheus v2.5.0+incompatible
	github.com/stretchr/testify v1.6.0
	github.com/vmware/govmomi v0.22.2
	golang.org/x/net v0.0.0-20191027093000-83d349e8ac1a // indirect; needed for freebsd/arm64
	gopkg.in/yaml.v2 v2.3.0
	layeh.com/radius v0.0.0-20190322222518-890bc1058917
)
