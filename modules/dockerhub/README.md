# Docker Hub repositories monitoring with Netdata

[`Docker Hub`](https://docs.docker.com/docker-hub/) is a service provided by Docker for finding and sharing container images with your team.
 
This module will collect `Docker Hub` repositories statistics.

## Charts

It produces the following charts:

-   Pulls Summary in `pulls`
-   Pulls in `pulls`
-   Pulls Rate in `pulls/s`
-   Stars in `stars/s`
-   Current Status in `status`
-   Time Since Last Update in `seconds`

## Configuration

Edit the `go.d/dockerhub.conf` configuration file using `edit-config` from the your agent's [config
directory](../../../../docs/step-by-step/step-04.md#find-your-netdataconf-file), which is typically at `/etc/netdata`.

```bash
cd /etc/netdata # Replace this path with your Netdata config directory
sudo ./edit-config go.d/dockerhub.conf
```

Needs only list of `repositories`. Here is an example:

```yaml
jobs:
  - name: me
    repositories:
      - 'me/repo1'
      - 'me/repo2'
      - 'me/repo3' 
```

For all available options please see module [configuration file](https://github.com/netdata/go.d.plugin/blob/master/config/go.d/dockerhub.conf).

## Troubleshooting

Check the module debug output. Run the following command as `netdata` user:

> ./go.d.plugin -d -m dockerhub
