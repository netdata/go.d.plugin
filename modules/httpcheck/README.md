<!--
title: "HTTP endpoint monitoring with Netdata"
description: "Monitor the health and performance of any HTTP endpoint with zero configuration, per-second metric granularity, and interactive visualizations."
custom_edit_url: "https://github.com/netdata/go.d.plugin/edit/master/modules/httpcheck/README.md"
sidebar_label: "HTTP endpoints"
learn_status: "Published"
learn_topic_type: "References"
learn_rel_path: "Integrations/Monitor/Remotes"
-->

# HTTP endpoint collector

This module monitors one or more http servers availability and response time.

## Metrics

All metrics have "httpcheck." prefix.

Labels per scope:

- global: url.

| Metric          | Scope  |                        Dimensions                        |   Units    |
|-----------------|:------:|:--------------------------------------------------------:|:----------:|
| response_time   | global |                           time                           |     ms     |
| response_length | global |                          length                          | characters |
| status          | global | success, no_connection, timeout, bad_content, bad_status |  boolean   |

## Check statuses

| Status        | Description                                                                              |
|---------------|------------------------------------------------------------------------------------------|
| success       | No error on HTTP request, body reading and body content checking                         |
| timeout       | Timeout error on HTTP request                                                            |
| bad content   | The body of the response didn't match the regex (only if `response_match` option is set) |
| bad status    | Response status code not in `status_accepted`                                            |
| no connection | Any other network error not specifically handled by the module                           |

## Configuration

Edit the `go.d/httpcheck.conf` configuration file using `edit-config` from the
Netdata [config directory](https://github.com/netdata/netdata/blob/master/docs/configure/nodes.md), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/httpcheck.conf
```

Here is an example for 2 servers:

```yaml
jobs:
  - name: cool_website1
    url: http://cool.website1:8080/home

  - name: cool_website2
    url: http://cool.website2:8080/home
    status_accepted:
      - 200
      - 201
      - 202
    response_match: <title>My cool website!<\/title>

  - name: cool_website3
    url: http://cool.website3:8080/home
    headers:
      User-Agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"
      body: 'j_username=user&j_password=pass'
    cookie_file: '/tmp/cookie.txt'
```

For all available options please see
module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/httpcheck.conf).

## Troubleshooting

To troubleshoot issues with the `httpcheck` collector, run the `go.d.plugin` with the debug option enabled. The output
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
  ./go.d.plugin -d -m httpcheck
  ```

