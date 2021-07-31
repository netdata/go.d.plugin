package mongo

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
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
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

func (l *Local) valid() bool {
	return options.Client().ApplyURI(l.connectionString()).Validate() == nil
}

func (l *Local) connectionString() string {
	mongoURL := url.URL{
		Scheme: "mongodb",
		Host:   fmt.Sprintf("%s:%d", l.Host, l.Port),
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
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Authdb string `yaml:"authdb"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
}

func (a *Auth) valid() bool {
	return options.Client().ApplyURI(a.connectionString()).Validate() == nil
}

func (a *Auth) connectionString() string {
	mongoURL := url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(a.User, a.Pass),
		Host:   fmt.Sprintf("%s:%d", a.Host, a.Port),
	}
	query := mongoURL.Query()

	if a.Authdb != "" {
		query.Set("authSource", a.Authdb)
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
	Local         `yaml:"local"`
	Auth          `yaml:"auth"`
	ConnectionStr string        `yaml:"connectionStr"`
	Timeout       time.Duration `yaml:"timeout"`
}
