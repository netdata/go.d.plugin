package mongo

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultName          = "admin"
	defaultHost          = "localhost"
	defaultPort          = 27017
	defaultTimeout       = 20
	defaultAuthDb        = "admin"
	defaultUser          = ""
	defaultPass          = ""
	defaultConnectionStr = ""
)

type Validator interface {
	valid() bool
	connectionString() string
}

type Local struct {
	host string `yaml:"host"`
	port uint   `yaml:"port"`
}

func (l *Local) valid() bool {
	return options.Client().ApplyURI(l.connectionString()).Validate() == nil
}

func (l *Local) connectionString() string {
	mongoURL := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%d", l.host, l.port),
	}
	return mongoURL.String()
}

type SSL struct {
	Ssl              bool   `yaml:"dsl"`
	SslCertReqs      bool   `yaml:"ssl_cert_reqs"`
	SslCaCerts       string `yaml:"ssl_ca_certs"`
	SslCrlFile       string `yaml:"ssl_crlfile"`
	SslCertfile      string `yaml:"ssl_certfile"`
	SslKeyfile       string `yaml:"ssl_keyfile"`
	SslPemPassphrase string `yaml:"ssl_pem_passphrase"`
}

type Auth struct {
	SSL
	host   string `yaml:"host"`
	port   int    `yaml:"port"`
	authdb string `yaml:"authdb"`
	user   string `yaml:"user"`
	pass   string `yaml:"pass"`
}

func (a *Auth) valid() bool {
	return options.Client().ApplyURI(a.connectionString()).Validate() == nil
}

func (a *Auth) connectionString() string {
	mongoURL := url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(a.user, a.pass),
		Host:   fmt.Sprintf("%s:%d", a.host, a.port),
	}
	query := mongoURL.Query()

	if a.authdb != "" {
		query.Set("authSource", a.authdb)
	}
	if a.SSL.Ssl {
		query.Set("ssl", "true")
	}
	if a.SSL.Ssl {
		query.Set("ssl", "true")
	}
	if a.SSL.SslCaCerts != "" {
		query.Set("tlsCAFile", a.SSL.SslCaCerts)
	}
	if a.SSL.SslCertfile != "" {
		query.Set("tlsCertificateKeyFile", a.SSL.SslCertfile)
	}
	if a.SSL.SslKeyfile != "" {
		query.Set("tlsCertificateKeyFile", a.SSL.SslKeyfile)
	}
	if a.SSL.SslPemPassphrase != "" {
		query.Set("tlsCertificateKeyFilePassword", a.SSL.SslPemPassphrase)
	}
	return strings.Join([]string{mongoURL.String(), query.Encode()}, "/")
}

type Config struct {
	Local
	Auth
	name          string        `yaml:"name"`
	ConnectionStr string        `yaml:"connectionStr"`
	Timeout       time.Duration `yaml:"timeout"`
}
