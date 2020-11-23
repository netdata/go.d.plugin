`DNSdist` setup:

[Install DNSdist](https://dnsdist.org/install.html) on your environment, if you 
decide to compile the source code, it will be necessary to install the following
packages before:

- [colm](https://www.colm.net/open-source/colm/) programming language
- [kelbt](freecode.com/projects/kelbt) parser
- [ragel](https://www.colm.net/open-source/ragel/) state machine  compiler

Create the configuration file, it is necessary to enable the webserver inside it:

```
newServer("8.8.8.8")
webserver("127.0.0.1:8083", "netdata", "netdata")
setServerPolicy(firstAvailable)
```

Start the server running:

```cmd
dnsdist -C dnsdist.conf --local=0.0.0.0:5300
```

Do requests for the server:

```cmd
for a in {0..1000}; do dig netdata.cloud @127.0.0.1 -p 5300 +noall +nocookie > /dev/null; done
```

Finally verify the statistics:

```cmd
curl -H"X-API-Key: netdata" "http://127.0.0.1:8083/jsonstat?command=stats"
```