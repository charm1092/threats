package threat

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repository {
	return &Repository{
		pool: conn,
	}
}

func (r *Repository) CreateThreat(
	ctx context.Context,
	md5, name, description, severity string,
) (Threat, error) {
	if !IsMD5Allowed(md5) {
		return Threat{}, ErrMD5IsNotAllowed
	}

	if !IsSeverityAllowed(severity) {
		return Threat{}, ErrSeverityIsNotAllowed
	}

	threat := NewThreat(md5, name, description, severity)

	sqlQuery := `
		INSERT INTO threats (md5, name, description, severity)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (md5) DO NOTHING
	`

	commandTag, err := r.pool.Exec(ctx, sqlQuery,
		threat.MD5,
		threat.Name,
		threat.Description,
		threat.Severity,
	)
	if err != nil {
		return Threat{}, err
	}

	if commandTag.RowsAffected() == 0 {
		return Threat{}, ErrMD5Exists
	}

	return *threat, nil
}

func (r *Repository) GetByMD5(
	ctx context.Context,
	md5 string,
) (Threat, error) {
	if !IsMD5Allowed(md5) {
		return Threat{}, ErrMD5IsNotAllowed
	}

	sqlQuery := `
		SELECT md5, name, COALESCE(description, ''), severity
		FROM threats
		WHERE md5 = $1
	`

	var threat Threat

	err := r.pool.QueryRow(ctx, sqlQuery, md5).Scan(
		&threat.MD5,
		&threat.Name,
		&threat.Description,
		&threat.Severity,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Threat{}, ErrMD5NotFound
		}
		return Threat{}, err
	}

	return threat, nil
}

func (r *Repository) UpdateThreat(
	ctx context.Context,
	md5, name, description, severity string,
) error {
	if !IsMD5Allowed(md5) {
		return ErrMD5IsNotAllowed
	}

	if !IsSeverityAllowed(severity) {
		return ErrSeverityIsNotAllowed
	}

	sqlQuery := `
		UPDATE threats
		SET name = $1,
		    description = $2,
		    severity = $3
		WHERE md5 = $4
	`

	commandTag, err := r.pool.Exec(ctx, sqlQuery,
		name,
		description,
		severity,
		md5,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrMD5NotFound
	}

	return nil
}