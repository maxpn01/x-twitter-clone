package apperror

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func MapPostgresUniqueViolation(err error, constraintErrors map[string]*AppError) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != pgerrcode.UniqueViolation {
		return err
	}

	if appErr, ok := constraintErrors[pgErr.ConstraintName]; ok {
		return appErr
	}

	return err
}
