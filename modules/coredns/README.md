<!--
title: "CoreDNS monitoring with Netdata"
description: "Monitor the health and performance of CoreDNS servers with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/coredns/README.md"
sidebar_label: "CoreDNS"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Networking"
-->

# CoreDNS monitoring with Netdata

[`CoreDNS`](https://coredns.io/) is a fast and flexible DNS server.

This module monitor one or more `CoreDNS` instances depending on configuration.

## Metrics

All metrics have "coredns." prefix.

| Metric                                    | Scope  |                                                                                     Dimensions                                                                                     |    Units    |
|-------------------------------------------|:------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-----------:|
| dns_request_count_total                   | global |                                                                                      requests                                                                                      | requests/s  |
| dns_responses_count_total                 | global |                                                                                     responses                                                                                      | responses/s |
| dns_request_count_total_per_status        | global |                                                                                 processed, dropped                                                                                 | requests/s  |
| dns_no_matching_zone_dropped_total        | global |                                                                                      dropped                                                                                       | requests/s  |
| dns_panic_count_total                     | global |                                                                                       panics                                                                                       |  panics/s   |
| dns_requests_count_total_per_proto        | global |                                                                                      udp, tcp                                                                                      | requests/s  |
| dns_requests_count_total_per_ip_family    | global |                                                                                       v4, v6                                                                                       | requests/s  |
| dns_requests_count_total_per_per_type     | global |                                              a, aaaa, mx, soa, cname, ptr, txt, ns, ds, dnskey, rrsig, nsec, nsec3, ixfr, any, other                                               | requests/s  |
| dns_responses_count_total_per_rcode       | global | noerror, formerr, servfail, nxdomain, notimp, refused, yxdomain, yxrrset, nxrrset, notauth, notzone, badsig, badkey, badtime, badmode, badname, badalg, badtrunc, badcookie, other | responses/s |
| server_dns_request_count_total            | server |                                                                                      requests                                                                                      | requests/s  |
| server_dns_responses_count_total          | server |                                                                                     responses                                                                                      | responses/s |
| server_dns_responses_count_total          | server |                                                                                     responses                                                                                      | responses/s |
| server_request_count_total_per_status     | server |                                                                                 processed, dropped                                                                                 | requests/s  |
| server_requests_count_total_per_proto     | server |                                                                                      udp, tcp                                                                                      | requests/s  |
| server_requests_count_total_per_ip_family | server |                                                                                       v4, v6                                                                                       | requests/s  |
| server_requests_count_total_per_per_type  | server |                                              a, aaaa, mx, soa, cname, ptr, txt, ns, ds, dnskey, rrsig, nsec, nsec3, ixfr, any, other                                               | requests/s  |
| server_responses_count_total_per_rcode    | server | noerror, formerr, servfail, nxdomain, notimp, refused, yxdomain, yxrrset, nxrrset, notauth, notzone, badsig, badkey, badtime, badmode, badname, badalg, badtrunc, badcookie, other | responses/s |
| zone_dns_request_count_total              | server |                                                                                      requests                                                                                      | requests/s  |
| zone_dns_responses_count_total            | server |                                                                                     responses                                                                                      | responses/s |
| zone_dns_responses_count_total            | server |                                                                                     responses                                                                                      | responses/s |
| zone_requests_count_total_per_proto       | server |                                                                                      udp, tcp                                                                                      | requests/s  |
| zone_requests_count_total_per_ip_family   | server |                                                                                       v4, v6                                                                                       | requests/s  |
| zone_requests_count_total_per_per_type    | server |                                              a, aaaa, mx, soa, cname, ptr, txt, ns, ds, dnskey, rrsig, nsec, nsec3, ixfr, any, other                                               | requests/s  |
| zone_responses_count_total_per_rcode      | server | noerror, formerr, servfail, nxdomain, notimp, refused, yxdomain, yxrrset, nxrrset, notauth, notzone, badsig, badkey, badtime, badmode, badname, badalg, badtrunc, badcookie, other | responses/s |

## Configuration

Edit the `go.d/coredns.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/coredns.conf
```

The module needs only the `url` to a CoreDNS `metrics-address`. Here is an example for several instances:

```yaml
jobs:
  - name: local
    url: http://127.0.0.1:9153/metrics

  - name: remote
    url: http://203.0.113.10:9153/metrics
```

For all available options, please see the
module's [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/coredns.conf).

## Troubleshooting

To troubleshoot issues with the `coredns` collector, run the `go.d.plugin` with the debug option enabled. The output
should give you clues as to why the collector isn't working.

- Navigate to the `plugins.d` directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on
  your system, open `netdata.conf` and look for the `plugins` setting under `[directories]`.

  ```bash
  cd /usr/libexec/netdata/plugins.d/
  ```

- Switch to the `netdata` user.

  ```bash
  sudo -u netdata -s
  ```

- Run the `go.d.plugin` to debug the collector:

  ```bash
  ./go.d.plugin -d -m coredns
  ```
