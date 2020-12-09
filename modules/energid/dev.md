## Installation 

Download  [generation 2](https://docs.energi.software/en/downloads/gen2-core-wallet) to your platform.

### Working with `energid` client

Start server:

```
./energid -testnet -port=9796
```

Requests using `energi-cli`:

```
./energi-cli -testnet getblockchaininfo
./energi-cli -testnet getmempoolinfo
./energi-cli -testnet getnetworkinfo
./energi-cli -testnet gettxoutsetinfo
```

### Working with `curl` client

Start server:

```
$ ./energid -testnet -rest -rpcallowip=192.168.0.0/24 -rpcport=9796 -server -rpcuser=netdata -rpcpassword=netdata
```

Requests using `curl`:

```
curl -v  --user netdata --data-binary '[{"jsonrpc": "1.0", "id":"1", "method": "getblockchaininfo", "params": [] }, {"jsonrpc": "1.0", "id":"2", "method": "getmempoolinfo", "params": [] }, {"jsonrpc": "1.0", "id":"3", "method": "getnetworkinfo", "params": [] }, {"jsonrpc": "1.0", "id":"4", "method": "gettxoutsetinfo", "params": [] }]' -H 'content-type: application/json;' http://127.0.0.1:9796/
```
