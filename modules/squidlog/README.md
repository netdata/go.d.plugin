<!--
title: "Web server log (Squid) monitoring with Netdata"
description: "Monitor the health and performance of Squid web server logs with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/squidlog/README.md"
sidebar_label: "Web server logs (Squid)"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Quantify logs to metrics"
-->

# Squid log monitoring with Netdata

[`Squid`](http://www.squid-cache.org/) is a caching and forwarding HTTP web proxy.

This module parses `Squid` access logs.

## Metrics

All metrics have "squidlog." prefix.

| Metric                                   | Scope  |                         Dimensions                         |    Units     |
|------------------------------------------|:------:|:----------------------------------------------------------:|:------------:|
| requests                                 | global |                          requests                          |  requests/s  |
| excluded_requests                        | global |                         unmatched                          |  requests/s  |
| type_requests                            | global |               success, bad, redirect, error                |  requests/s  |
| http_status_code_class_responses         | global |                  1xx, 2xx, 3xx, 4xx, 5xx                   | responses/s  |
| http_status_code_responses               | global |         <i>a dimension per HTTP response code</i>          | responses/s  |
| bandwidth                                | global |                            sent                            |  kilobits/s  |
| response_time                            | global |                       min, max, avg                        | milliseconds |
| uniq_clients                             | global |                          clients                           |   clients    |
| cache_result_code_requests               | global |          <i>a dimension per cache result code</i>          |  requests/s  |
| cache_result_code_transport_tag_requests | global | <i>a dimension per cache result delivery transport tag</i> |  requests/s  |
| cache_result_code_handling_tag_requests  | global |      <i>a dimension per cache result handling tag</i>      |  requests/s  |
| cache_code_object_tag_requests           | global |  <i>a dimension per cache result produced object tag</i>   |  requests/s  |
| cache_code_load_source_tag_requests      | global |    <i>a dimension per cache result load source tag</i>     |  requests/s  |
| cache_code_error_tag_requests            | global |       <i>a dimension per cache result error tag</i>        |  requests/s  |
| http_method_requests                     | global |             <i>a dimension per HTTP method</i>             |  requests/s  |
| mime_type_requests                       | global |              <i>a dimension per MIME type</i>              |  requests/s  |
| hier_code_requests                       | global |           <i>a dimension per hierarchy code</i>            |  requests/s  |
| server_address_forwarded_requests        | global |           <i>a dimension per server address</i>            |  requests/s  |

## Log Parsers

Squidlog supports 3 log parsers:

- CSV
- LTSV
- RegExp

RegExp is the slowest among them, but it is very likely you will need to use it if your log format is not default.

## Known Fields

These are `Squid` [log format codes](http://www.squid-cache.org/Doc/config/logformat/).

Squidlog is aware how to parse and interpret following codes:

| field          | squid format code | description                                                   |
|----------------|-------------------|---------------------------------------------------------------|
| resp_time      | %tr               | Response time (milliseconds).                                 |
| client_address | %>a               | Client source IP address.                                     |
| client_address | %>A               | Client FQDN.                                                  |
| cache_code     | %Ss               | Squid request status (TCP_MISS etc).                          |
| http_code      | %>Hs              | The HTTP response status code from Content Gateway to client. |
| resp_size      | %<st              | Total size of reply sent to client (after adaptation).        |
| req_method     | %rm               | Request method (GET/POST etc).                                |
| hier_code      | %Sh               | Squid hierarchy status (DEFAULT_PARENT etc).                  |
| server_address | %<a               | Server IP address of the last server or peer connection.      |
| server_address | %<A               | Server FQDN or peer name.                                     |
| mime_type      | %mt               | MIME content type.                                            |

In addition, to
make `Squid` [native log format](https://wiki.squid-cache.org/Features/LogFormat#Squid_native_access.log_format_in_detail)
csv parsable, squidlog understands these groups of codes:

| field       | squid format code | description                        |
|-------------|-------------------|------------------------------------|
| result_code | %Ss/%>Hs          | Cache code and http code.          |
| hierarchy   | %Sh/%<a           | Hierarchy code and server address. |

## Custom Log Format

Custom log format is easy. Use [known fields](#known-fields) to construct your log format.

- If using CSV parser

**Note**: can be used only if all known squid format codes are separated by csv delimiter. For example, if you
have `%Ss:%Sh`, csv parser cant extract `%Ss` and `%Sh` from it, and you need to use RegExp parser.

Copy your current log format. Replace all known squid format codes with corresponding [known](#known-fields) fields.
Replaces others with "-".

```yaml
jobs:
  - name: squid_log_custom_csv_exampla
    path: /var/log/squid/access.log
    log_type: csv
    csv_config:
      format: '- resp_time client_address result_code resp_size req_method - - hierarchy mime_type'
```

- If using LTSV parser

Provide fields mapping. You need to map your label names to [known](#known-fields) fields.

```yaml
  - name: squid_log_custom_ltsv_exampla
    path: /var/log/squid/access.log
    log_type: ltsv
    ltsv_config:
      mapping:
        label1: resp_time
        label2: client_address
        ...
```

- If using RegExp parser

Use pattern with subexpressions names. These names should be [known](#known-fields) by squidlog. We recommend to
use https://regex101.com/ to test your regular expression.

```yaml
jobs:
  - name: squid_log_custom_regexp_exampla
    path: /var/log/squid/access.log
    log_type: regexp
    regexp_config:
      format: '^[0-9.]+\s+(?P<resp_time>[0-9]+) (?P<client_address>[\da-f.:]+) (?P<cache_code>[A-Z_]+)\/(?P<http_code>[0-9]+) (?P<resp_size>[0-9]+) (?P<req_method>[A-Z]+) [^ ]+ [^ ]+ (?P<hier_code>[A-Z_]+)\/[\da-z.:-]+ (?P<mime_type>[A-Za-z-]+)'
```

## Configuration

Edit the `go.d/squidlog.conf` configuration file using `edit-config` from the
Netdata [config directory](https://learn.netdata.cloud/docs/configure/nodes), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/squidlog.conf
```

This module needs only `path` to log file if you
use [native log format](https://wiki.squid-cache.org/Features/LogFormat#Squid_native_access.log_format_in_detail). If
you use custom log format you need [to set it manually](#custom-log-format).

```yaml
jobs:
  - name: squid
    path: /var/log/squid/access.log
```

For all available options, please see the
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/squidlog.conf).

## Troubleshooting

To troubleshoot issues with the `squid_log` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m squid_log
  ```
