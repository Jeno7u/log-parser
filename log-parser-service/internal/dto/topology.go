package dto

type TopologyGroup struct {
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

type TopologyLink struct {
	ID           string `json:"link_id,omitempty"`
	LogID        string `json:"log_id"`
	NodeID       string `json:"node_id"`
	PortID       string `json:"port_id"`
	RelationType string `json:"relation_type"`
}

type Topology struct {
	LogID  string          `json:"log_id"`
	Groups []TopologyGroup `json:"groups"`
	Ports  []Port          `json:"ports"`
	Links  []TopologyLink  `json:"links"`
}
