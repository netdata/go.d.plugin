#### Setup Dnsmasq

Just run the docker command ([image page](https://hub.docker.com/r/andyshinn/dnsmasq)):

```cmd
docker run \
	--name dnsmasq \
	-d \
	-p 53:53/udp \
	--cap-add=NET_ADMIN \
	andyshinn/dnsmasq
```

#### Gather metrics.

See [cache statistics section](https://manpages.debian.org/stretch/dnsmasq-base/dnsmasq.8.en.html#NOTES).

The cache statistics are also available in the DNS as answers to queries of class CHAOS and type TXT in domain bind. The
domain names are `cachesize.bind`, `insertions.bind`, `evictions.bind`,
`misses.bind`, `hits.bind`, `auth.bind` and `servers.bind`. An example command to query this, using the dig utility
would be:

> dig +short chaos txt cachesize.bind @127.0.0.1
