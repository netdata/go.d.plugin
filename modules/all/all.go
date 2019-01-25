package all

import (
	_ "github.com/netdata/go.d.plugin/modules/activemq"
	_ "github.com/netdata/go.d.plugin/modules/apache"
	_ "github.com/netdata/go.d.plugin/modules/consul"
	_ "github.com/netdata/go.d.plugin/modules/dnsquery"
	_ "github.com/netdata/go.d.plugin/modules/example"
	_ "github.com/netdata/go.d.plugin/modules/freeradius"
	_ "github.com/netdata/go.d.plugin/modules/httpcheck"
	_ "github.com/netdata/go.d.plugin/modules/lighttpd"
	_ "github.com/netdata/go.d.plugin/modules/lighttpd2"
	_ "github.com/netdata/go.d.plugin/modules/nginx"
	// _ "github.com/netdata/go.d.plugin/modules/nvidia_nvml" # wil be imported by cgo enabled version
	_ "github.com/netdata/go.d.plugin/modules/portcheck"
	_ "github.com/netdata/go.d.plugin/modules/rabbitmq"
	_ "github.com/netdata/go.d.plugin/modules/solr"
	_ "github.com/netdata/go.d.plugin/modules/springboot2"
	_ "github.com/netdata/go.d.plugin/modules/weblog"
)
