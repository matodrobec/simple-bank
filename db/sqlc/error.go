package db

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows
var ErrUniqueViolation = &pgconn.PgError{Code: UniqueViolation}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		fmt.Println(">> constraint name", pgErr.ConstraintName)
		return pgErr.Code
	}
	return ""
}