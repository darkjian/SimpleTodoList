package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/darkjian/simpletodolist/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func format(cause error, err error) error {
	return fmt.Errorf("%w: %v", cause, err)
}

func normalizePGError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return format(domain.ErrNotFound, err)
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return format(domain.ErrContextCancelled, err)
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return format(domain.ErrNotFound, err)
		case "23503":
			return format(domain.ErrAlreadyExists, err)
		}
	}

	return fmt.Errorf("%w: %v", domain.ErrInternal, err)
}
