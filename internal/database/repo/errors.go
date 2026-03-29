package repo

import (
	"errors"
	"parus-test/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// 23505 - уникальность

// 23503 - нарушение внешнего ключа

func mapDBError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return domain.ErrAlreadyExists
		case "23503": // foreign_key_violation
			return domain.ErrNotFound
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}

		return err
	}
	return err
}
