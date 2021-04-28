# Unity module

A collection module for Unity Dell EMC.

## Charts

The charts wished for are listed in the `unity.json` file, in `/etc/netdata/go.d/`

## Configuration

Disabled by default. Should be explicitly enabled in [go.d.conf](https://github.com/netdata/go.d.plugin/blob/master/config/go.d.conf).

```yaml
# go.d.conf
modules:
  example: yes
```

The file `unity.conf` should not be edited to change the configuration. Rather, editing `unity.json` which contain targets and wanted metrics would be right.

unity.conf :

```yaml
jobs:
  - name: unity
```

unity.json :

```
{
    "username":"...",
    "password":"...",
    "servers":[
        {
            "name":"local",
            "adress":"178.0.0.1",
            "targets":{
                "lun":[
                    "nss-X_Y"
                ],
                "fc":[
                    "spa_fcX",
                    "spb_fcX"
                ]
            }
        },
    ],
    "interval":20,
    "insecure":true,
    "metrics":{
        "general":[
            "utilization"
        ],
        "fc":[
            "readBandwidth",
            "writeBandwidth"
        ],
        "lun":[
            "responseTime"
        ]
    }
}

```


## Troubleshooting

To troubleshoot issues with the Unity collector, run the `go.d.plugin` orchestrator with the debug option enabled.
The output should give you clues as to why the collector isn't working.

First, navigate to your plugins directory, usually at `/usr/libexec/netdata/plugins.d/`. If that's not the case on your
system, open `netdata.conf` and look for the setting `plugins directory`. Once you're in the plugins directory, switch
to the `netdata` user.

```bash
cd /usr/libexec/netdata/plugins.d/
sudo -u netdata -s
```

You can now run the `go.d.plugin` orchestrator to debug the collector:

```bash
./go.d.plugin -d -m unity
```
