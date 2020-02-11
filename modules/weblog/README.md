# Apache/NGINX logs monitoring with Netdata

This module parses [`Apache`](https://httpd.apache.org/) and [`NGINX`](https://nginx.org/en/) web servers logs.

## Charts

Module produces following charts:

-   Total Requests in `requests/s`
-   Excluded Requests in `requests/s`
-   Requests By Type in `requests/s`
-   Responses By Status Code Class in `responses/s`
-   Responses By Status Code in `responses/s`
-   Informational Responses By Status Code in `responses/s`
-   Successful Responses By Status Code in `responses/s`
-   Redirects Responses By Status Code in `responses/s`
-   Client Errors Responses By Status Code in `responses/s`
-   Server Errors Responses By Status Code in `responses/s`
-   Bandwidth in `kilobits/s`
-   Request Processing Time in `milliseconds`
-   Requests Processing Time Histogram in `requests/s`
-   Upstream Response Time in `requests/s`
-   Upstream Responses Time Histogram in `responses/s`
-   Current Poll Unique Clients in `clients`
-   Requests By Vhost in `requests/s`
-   Requests By Port in `requests/s`
-   Requests By Scheme in `requests/s`
-   Requests By HTTP Method in `requests/s`
-   Requests By HTTP Version in `requests/s`
-   Requests By IP Protocol in `requests/s`
-   Requests By SSL Connection Protocol in `requests/s`
-   Requests By SSL Connection Cipher Suite in `requests/s`
-   URL Field Requests By Pattern `requests/s`

For every Custom field:

-   Requests By Pattern in `requests/s`

For every URL pattern:

-   Responses By Status Code in `responses/s`
-   Bandwidth in `kilobits/s`
-   Request Processing Time in `milliseconds`

## Log Parsers

Weblog supports 3 different log parsers:

-   CSV
-   [LTSV](http://ltsv.org/)
-   RegExp

Try to avoid using RegExp because it's much slower than the other two parsers. RegExp should be used only if LTSV and CSV parsers dont work for you.

There is an example job for every log parser.

```yaml
jobs:
  - name: csv_parser_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: 'FORMAT'
      fields_per_record: -1
      delimiter: ' '
      trim_leading_space: no

  - name: ltsv_parser_example
    path: /path/to/file.log
    log_type: ltsv
    ltsv_config:
      field_delimiter: ' '
      value_delimiter: ':'
      mapping:
        label1: field1
        label2: field2

  - name: regexp_parser_example
    path: /path/to/file.log
    log_type: regexp
    regexp_config:
      pattern: 'PATTERN'
```

## Log Parser Auto-Detection

If `log_type` parameter is set to `auto` (which is default), weblog will try to auto-detect appropriate log parser and log format
using the last line of the log file.

To auto-detect parser type the module checks if the line is in LTSV format first. If it is not the case it assumes that the format is CSV.

To auto-detect CSV format weblog uses list of predefined csv formats. It tries to parse the line using each of them in the following order:

```
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time $upstream_response_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time $upstream_response_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time
$host:$server_port $remote_addr - - [$time_local] "$request" $status $body_bytes_sent
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time $upstream_response_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent - - $request_length $request_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time $upstream_response_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent     $request_length $request_time
                   $remote_addr - - [$time_local] "$request" $status $body_bytes_sent
```

The first one that matches will be used later. If you use default Apache/NGINX log format auto-detect will do for you.
If it doesnt work you need [to set format manually](#custom-log-format).

## Known Fields

These are [NGINX](http://nginx.org/en/docs/varindex.html) and [Apache](http://httpd.apache.org/docs/current/mod/mod_log_config.html) log format variables.

Weblog is aware how to parse and interpret the fields:

| nginx                   | apache    | description                                   |
|-------------------------|-----------|-----------------------------------------------|
| $host ($http_host)      | %v        | Name of the server which accepted a request.
| $server_port            | %p        | Port of the server which accepted a request.
| $scheme                 | -         | Request scheme. "http" or "https".
| $remote_addr            | %a (%h)   | Client address.
| $request                | %r        | Full original request line. The line is "$request_method $request_uri $server_protocol".
| $request_method         | %m        | Request method. Usually "GET" or "POST".
| $request_uri            | %U        | Full original request URI.
| $server_protocol        | %H        | Request protocol. Usually "HTTP/1.0", "HTTP/1.1", or "HTTP/2.0".
| $status                 | %s (%>s)  | Response status code.
| $request_length         | %I        | Bytes received from a client, including request and headers.
| $bytes_sent             | %O        | Bytes sent to a client, including request and headers.
| $body_bytes_sent        | %B (%b)   | Bytes sent to a client, not counting the response header.
| $request_time           | %D        | Request processing time.
| $upstream_response_time | -         | Time spent on receiving the response from the upstream server.
| $ssl_protocol           | -         | Protocol of an established SSL connection.
| $ssl_cipher             | -         | String of ciphers used for an established SSL connection.

In addition to that weblog understands [user defined fields](#custom-fields-feature).

Notes:

-   Apache `%h` logs the IP address if [HostnameLookups](https://httpd.apache.org/docs/2.4/mod/core.html#hostnamelookups) is Off.
    Weblog counts hostname as IPv4 address. We recommend either to disable HostnameLookups or use `%a` instead of `%h`. 
-   Since httpd 2.0, unlike 1.3, the `%b` and `%B` format strings do not represent the number of bytes sent to the client,
    but simply the size in bytes of the HTTP response. It will will differ, for instance, if the connection is aborted,
    or if SSL is used. The `%O` format provided by [`mod_logio`](https://httpd.apache.org/docs/2.4/mod/mod_logio.html)
    will log the actual number of bytes sent over the network.
-   To get `%I` and `%O` working you need to enable `mod_logio` on Apache.
-   NGINX logs URI with query parameters, Apache doesnt.
-   `$request` is parsed into `$request_method`, `$request_uri` and `$server_protocol`. If you have `$request` in your log format, 
    there is no sense to have others.
-   Don't use both `$bytes_sent` and `$body_bytes_sent` (`%O` and `%B` or `%b`). The module does not distinguish between these parameters.


## Custom Log Format

Custom log format is easy. Use [known fields](#known-fields) to construct your log format.

-   If using CSV parser

Since weblog understands 
 and Apache variables all you need is to copy your log format and... that is it!
If there is a field that is not known by the weblog it's not a problem. It will skip it during parsing.
But we suggest to replace all unknown fields with `-` for optimization purposes.

Let's take as an example some non default format.

```bash
# apache
LogFormat "\"%{Referer}i\" \"%{User-agent}i\" %h %l %u %t \"%r\" %>s %b" custom

# nginx
log_format custom '"$http_referer" "$http_user_agent" '
                  '$remote_addr - $remote_user [$time_local] '
                  '"$request" $status $body_bytes_sent'
```

To get it working we need to copy the format without any changes (make it a line for nginx). Replacing unknown fields
is optional but recommended.

Special case:

`%t` and `$time_local` represent time in [Common Log Format](https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format).
It is a special case because it's in fact 2 fields after csv parse (ex.: `[22/Mar/2009:09:30:31 +0100]`).
Weblog understands it and you don't need to replace it with `-` (if we want to do it we need to make it `- -`).

```yaml
jobs:
  - name: apache_csv_custom_format_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '- - %h - - %t \"%r\" %>s %b'

  - name: nginx_csv_custom_format_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '- - $remote_addr - - [$time_local] "$request" $status $body_bytes_sent'
```

-   If using LTSV parser

Provide fields mapping if needed. Dont use `$` and `%` prefixes for mapped field names. They are only needed in CSV format.

-   If using RegExp parser

Use pattern with subexpressions names. These names should be known by weblog.


## Custom Fields Feature

Weblog is able to extract user defined fields and count patterns matches against these fields.

This feature needs:
-   custom log format with user defined fields
-   list of patterns to match against appropriate fields

Pattern syntax: [matcher](https://github.com/netdata/go.d.plugin/tree/master/pkg/matcher#supported-format).
 
There is an example with 2 custom fields - `$http_referer` and `$http_user_agent`. Weblog is unaware of these fields, but
we still can get some info from them.

```yaml
  - name: nginx_csv_custom_fields_example
    path: /path/to/file.log
    log_type: csv
    csv_config:
      format: '- - $remote_addr - - [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"'
    custom_fields:
      - name:  http_referer     # same name as in 'format' without $
        patterns:
          - name:  cacti
            match: '~ cacti'
          - name:  observium
            match: '~ observium'
      - name:  http_user_agent  # same name as in 'format' without $
        patterns:
          - name:  android
            match: '~ Android'
          - name:  iphone
            match: '~ iPhone'
          - name:  other
            match: '* *'
```

## Configuration

Edit the `go.d/web_log.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/web_log.conf
```

This module needs only `path` to log file. If it fails to auto-detect your log format you need [to set it manually](#custom-log-format). 

```yaml
jobs:
  - name: nginx
    path: /var/log/nginx/access.log

  - name: apache
    path: /var/log/apache2/access.log
    log_type: csv
    csv_config
      format: '- - %h - - %t \"%r\" %>s %b'
```
 
For all available options, please see the module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/web_log.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m web_log
