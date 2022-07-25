package postgres

import "errors"

func (p Postgres) validateConfig() error {
	if p.DSN == "" {
		return errors.New("DSN not set")
	}
	return nil
}
