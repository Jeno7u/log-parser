package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/jackc/pgx/v5"
)

func (r *logRepository) CreatePorts(ctx context.Context, logID string, nodeIDs map[string]string, ports []dto.Port) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, port := range ports {
		nodeID, ok := nodeIDs[port.NodeGUID]
		if !ok {
			return fmt.Errorf("node %s not found for port %s", port.NodeGUID, port.PortGUID)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO ports (log_id, node_id, node_guid, port_guid, port_num, local_port_num, port_state, port_phy_state, link_speed_actv, link_width_actv, raw_line) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, logID, nodeID, port.NodeGUID, port.PortGUID, port.PortNum, port.LocalPortNum, port.PortState, port.PortPhyState, port.LinkSpeedActv, port.LinkWidthActv, port.RawLine); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *logRepository) ListPortsByNodeID(ctx context.Context, nodeID string) ([]dto.Port, error) {
	rows, err := r.pool.Query(ctx, `SELECT id::text, log_id::text, node_id::text, node_guid, port_guid, port_num, local_port_num, port_state, port_phy_state, link_speed_actv, link_width_actv, raw_line FROM ports WHERE node_id = $1 ORDER BY port_num`, nodeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []dto.Port
	for rows.Next() {
		var port dto.Port
		if err := rows.Scan(&port.ID, &port.LogID, &port.NodeID, &port.NodeGUID, &port.PortGUID, &port.PortNum, &port.LocalPortNum, &port.PortState, &port.PortPhyState, &port.LinkSpeedActv, &port.LinkWidthActv, &port.RawLine); err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}

	return ports, rows.Err()
}

func (r *logRepository) ListPortsByLogID(ctx context.Context, logID string) ([]dto.Port, error) {
	rows, err := r.pool.Query(ctx, `SELECT id::text, log_id::text, node_id::text, node_guid, port_guid, port_num, local_port_num, port_state, port_phy_state, link_speed_actv, link_width_actv, raw_line FROM ports WHERE log_id = $1 ORDER BY node_guid, port_num`, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []dto.Port
	for rows.Next() {
		var port dto.Port
		if err := rows.Scan(&port.ID, &port.LogID, &port.NodeID, &port.NodeGUID, &port.PortGUID, &port.PortNum, &port.LocalPortNum, &port.PortState, &port.PortPhyState, &port.LinkSpeedActv, &port.LinkWidthActv, &port.RawLine); err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}

	return ports, rows.Err()
}

var _ = sql.ErrNoRows
