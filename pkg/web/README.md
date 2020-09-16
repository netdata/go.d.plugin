# web

This package contains HTTP related configurations (`Request`, `Client` and `HTTP` structs) and functions to
create `http.Request` and `http.Client` from them.

`HTTP` embeds both `Request` and `Client`. 

Every module that collects metrics doing HTTP requests should use `HTTP`.
It allows to have same set of user configurable options across all modules.

## Usage

Just make `HTTP` part of your module configuration.

```go
package example

import "github.com/netdata/go.d.plugin/pkg/web"

type Config struct {
	web.HTTP `yaml:",inline"`
}

type Example struct {
	Config `yaml:",inline"`
}
```

## Configuration options

HTTP request related options:

-   `url`: the URL to access.
-   `username`: the username for basic HTTP authentication.
-   `password`: the password for basic HTTP authentication.
-   `proxy_username`: the username for basic HTTP authentication of a user agent to a proxy server.
-   `proxy_password`: the password for basic HTTP authentication of a user agent to a proxy server.
-   `body`: the HTTP request body to be sent by the client.
-   `method`: the HTTP method (GET, POST, PUT, etc.).
-   `headers`: the HTTP request header fields to be sent by the client.

HTTP client related options:

-   `timeout`: the HTTP request time limit.
-   `not_follow_redirects`: the policy for handling redirects.
-   `proxy_url`: the URL of the proxy to use.
-   `tls_skip_verify`: controls whether a client verifies the server's certificate chain and host name.
-   `tls_ca`: certificate authority to use when verifying server certificates.
-   `tls_cert`: tls certificate to use.
-   `tls_key`: tls key to use.

