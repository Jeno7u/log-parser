-- +goose Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS logs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name text NOT NULL,
    source_path text NOT NULL,
    status text NOT NULL,
    error_message text NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS nodes (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id uuid NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_desc text NOT NULL,
    num_ports integer NOT NULL,
    node_type integer NOT NULL,
    class_version integer NOT NULL,
    base_version integer NOT NULL,
    system_image_guid text NOT NULL,
    node_guid text NOT NULL,
    port_guid text NOT NULL,
    UNIQUE (log_id, node_guid)
);

CREATE TABLE IF NOT EXISTS ports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id uuid NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
	node_id uuid NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    node_guid text NOT NULL,
    port_guid text NOT NULL,
    port_num integer NOT NULL,
    local_port_num integer NOT NULL,
    port_state text NOT NULL,
    port_phy_state text NOT NULL,
    link_speed_actv text NOT NULL,
    link_width_actv text NOT NULL,
    raw_line text NOT NULL,
    UNIQUE (log_id, node_guid, port_guid, port_num)
);

CREATE TABLE IF NOT EXISTS nodes_info (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id uuid NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    sw_guid text NOT NULL,
    key text NOT NULL,
    value text NOT NULL,
    UNIQUE (log_id, sw_guid, key)
);

CREATE TABLE IF NOT EXISTS topology_links (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    log_id uuid NOT NULL REFERENCES logs(id) ON DELETE CASCADE,
    node_id uuid NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    port_id uuid NOT NULL REFERENCES ports(id) ON DELETE CASCADE,
    relation_type text NOT NULL DEFAULT 'node_port',
    UNIQUE (log_id, node_id, port_id, relation_type)
);



-- +goose Down
DROP TABLE IF EXISTS nodes_info;

DROP TABLE IF EXISTS topology_links;

DROP TABLE IF EXISTS ports;

DROP TABLE IF EXISTS nodes;

DROP TABLE IF EXISTS logs;
