package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func CollectAll[Row, T any](rows pgx.Rows, toDomain func(Row) T) ([]T, error) {
	found, err := pgx.CollectRows(rows, pgx.RowToStructByName[Row])
	if err != nil {
		return nil, err
	}
	items := make([]T, 0, len(found))
	for _, row := range found {
		items = append(items, toDomain(row))
	}
	return items, nil
}

func CollectOne[Row, T any](rows pgx.Rows, toDomain func(Row) T) (T, bool, error) {
	var zero T
	row, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Row])
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return zero, false, nil
	case err != nil:
		return zero, false, err
	}
	return toDomain(row), true, nil
}
