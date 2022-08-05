package pgbouncer

import "errors"

func (p *PgBouncer) validateConfig() error {
	if p.DSN == "" {
		return errors.New("DSN not set")
	}
	return nil
}
