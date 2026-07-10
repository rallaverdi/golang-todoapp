package core_postgres_pool

import "errors"

var (
	ErrNoRows                       = errors.New("no rows")
	ErrViolatesForeignKeyConstraint = errors.New("violates foreign key")
	ErrUnknown                      = errors.New("unknown error")
)
