package repository

import (
	"context"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeRepository interface {
	CreateNodes(ctx context.Context, logID string, nodes []dto.Node) error
	ListNodesByLogID(ctx context.Context, logID string) ([]dto.Node, error)
	GetNodeByID(ctx context.Context, nodeID string) (dto.Node, error)
	CreateNodeInfo(ctx context.Context, logID string, infos []dto.NodeInfo) error
}

type nodeRepository struct {
	pool *pgxpool.Pool
}

func NewNodeRepository(pool *pgxpool.Pool) NodeRepository {
	return &nodeRepository{pool}
}

func (r *nodeRepository) CreateNodes(ctx context.Context, logID string, nodes []dto.Node) error {
	query := `
		INSERT INTO nodes (log_id, node_desc, num_ports, node_type, class_version, base_version, system_image_guid, node_guid, port_guid)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, node := range nodes {
		if _, err := tx.Exec(ctx, query, logID, node.NodeDesc, node.NumPorts, node.NodeType, node.ClassVersion, node.BaseVersion, node.SystemImageGUID, node.NodeGUID, node.PortGUID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *nodeRepository) ListNodesByLogID(ctx context.Context, logID string) ([]dto.Node, error) {
	query := `
		SELECT id::text, log_id::text, node_desc, num_ports, node_type, class_version, base_version, system_image_guid, node_guid, port_guid
		FROM nodes
		WHERE log_id = $1
		ORDER BY node_desc
	`

	rows, err := r.pool.Query(ctx, query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []dto.Node
	for rows.Next() {
		var node dto.Node
		if err := rows.Scan(&node.ID, &node.LogID, &node.NodeDesc, &node.NumPorts, &node.NodeType, &node.ClassVersion, &node.BaseVersion, &node.SystemImageGUID, &node.NodeGUID, &node.PortGUID); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, rows.Err()
}

func (r *nodeRepository) GetNodeByID(ctx context.Context, nodeID string) (dto.Node, error) {
	query := `
		SELECT id::text, log_id::text, node_desc, num_ports, node_type, class_version, base_version, system_image_guid, node_guid, port_guid
		FROM nodes
		WHERE id = $1
	`

	var node dto.Node
	err := r.pool.QueryRow(ctx, query, nodeID).Scan(&node.ID, &node.LogID, &node.NodeDesc, &node.NumPorts, &node.NodeType, &node.ClassVersion, &node.BaseVersion, &node.SystemImageGUID, &node.NodeGUID, &node.PortGUID)
	if err != nil {
		return dto.Node{}, err
	}

	return node, nil
}

func (r *nodeRepository) CreateNodeInfo(ctx context.Context, logID string, infos []dto.NodeInfo) error {
	query := `
		INSERT INTO nodes_info (log_id, sw_guid, key, value)
		VALUES ($1, $2, $3, $4)
	`

	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, info := range infos {
		if _, err := tx.Exec(ctx, query, logID, info.SWGUID, info.Key, info.Value); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
