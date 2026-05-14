package dto

type Node struct {
	ID              string `json:"node_id"`
	LogID           string `json:"log_id"`
	NodeGUID        string `json:"node_guid"`
	PortGUID        string `json:"port_guid"`
	NodeDesc        string `json:"node_desc"`
	NumPorts        int    `json:"num_ports"`
	NodeType        int    `json:"node_type"`
	ClassVersion    int    `json:"class_version"`
	BaseVersion     int    `json:"base_version"`
	SystemImageGUID string `json:"system_image_guid"`
}
