package repository

import (
	"context"
	"fmt"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LogRepository interface {
	CreateLog(ctx context.Context, log dto.Log) (string, error)
	UpdateLogStatus(ctx context.Context, logID string, status string, message string) error
	GetLogByID(ctx context.Context, logID string) (dto.Log, error)
	CreateLinks(ctx context.Context, logID string, links []dto.TopologyLink) error
	ListLinksByLogID(ctx context.Context, logID string) ([]dto.TopologyLink, error)
}

type logRepository struct {
	pool *pgxpool.Pool
}

func NewLogRepository(pool *pgxpool.Pool) LogRepository {
	return &logRepository{pool}
}

func (r *logRepository) CreateLog(ctx context.Context, log dto.Log) (string, error) {
	query := `
        INSERT INTO logs (file_name, source_path, status, created_at, updated_at) 
        VALUES ($1, $2, $3, now(), now()) 
        RETURNING id
    `

	var id string
	err := r.pool.QueryRow(ctx, query, log.FileName, log.SourcePath, log.Status).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *logRepository) UpdateLogStatus(ctx context.Context, logID string, status string, message string) error {
	query := `
		UPDATE logs SET status = $2, error_message = NULLIF($3, ''), updated_at = now() WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, logID, status, message)
	if err != nil {
		return fmt.Errorf("update log status: %w", err)
	}
	return nil
}

func (r *logRepository) GetLogByID(ctx context.Context, logID string) (dto.Log, error) {
	query := `
		SELECT id::text, file_name, source_path, status, COALESCE(error_message, ''), created_at, updated_at
		FROM logs
		WHERE id = $1
	`

	var log dto.Log
	err := r.pool.QueryRow(ctx, query, logID).Scan(&log.ID, &log.FileName, &log.SourcePath, &log.Status, &log.ErrorMessage, &log.CreatedAt, &log.UpdatedAt)
	if err != nil {
		return dto.Log{}, err
	}

	return log, nil
}

func (r *logRepository) CreateLinks(ctx context.Context, logID string, links []dto.TopologyLink) error {
	query := `
		INSERT INTO topology_links (log_id, node_id, port_id, relation_type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (log_id, node_id, port_id, relation_type) DO NOTHING
	`

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, link := range links {
		if _, err := tx.Exec(ctx, query, logID, link.NodeID, link.PortID, link.RelationType); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *logRepository) ListLinksByLogID(ctx context.Context, logID string) ([]dto.TopologyLink, error) {
	query := `
		SELECT id::text, log_id::text, node_id::text, port_id::text, relation_type
		FROM topology_links
		WHERE log_id = $1
		ORDER BY relation_type, node_id
	`

	rows, err := r.pool.Query(ctx, query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []dto.TopologyLink
	for rows.Next() {
		var link dto.TopologyLink
		if err := rows.Scan(&link.ID, &link.LogID, &link.NodeID, &link.PortID, &link.RelationType); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, rows.Err()
}
