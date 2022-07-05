package proxysql

import (
	"fmt"

	"gopkg.in/ini.v1"
)

func dsnFromFile(filename string) (string, error) {
	f, err := ini.Load(filename)
	if err != nil {
		return "", err
	}

	section, err := f.GetSection("client")
	if err != nil {
		return "", err
	}

	defaultUser := "stats"
	defaultPassword := "stats"
	defaultHost := "127.0.0.1"
	defaultPort := "6033"

	user := section.Key("user").String()
	password := section.Key("password").String()
	host := section.Key("host").String()
	port := section.Key("port").String()

	var dsn string

	if user != "" {
		dsn = user
	} else {
		dsn = defaultUser
	}

	if password != "" {
		dsn += ":" + password
	} else {
		dsn += ":" + defaultPassword
	}

	switch {
	case host != "" && port != "":
		dsn += fmt.Sprintf("@tcp(%s:%s)/", host, port)
	case host != "":
		dsn += fmt.Sprintf("@tcp(%s:%s)/", host, defaultPort)
	case port != "":
		dsn += fmt.Sprintf("@tcp(%s:%s)/", defaultHost, port)
	default:
		dsn += "@/"
	}

	return dsn, nil
}
