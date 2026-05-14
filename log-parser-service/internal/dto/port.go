package dto

type Port struct {
	ID            string `json:"port_id"`
	LogID         string `json:"log_id"`
	NodeID        string `json:"node_id"`
	NodeGUID      string `json:"node_guid"`
	PortGUID      string `json:"port_guid"`
	PortNum       int    `json:"port_num"`
	LocalPortNum  int    `json:"local_port_num"`
	PortState     string `json:"port_state"`
	PortPhyState  string `json:"port_phy_state"`
	LinkSpeedActv string `json:"link_speed_actv"`
	LinkWidthActv string `json:"link_width_actv"`
	RawLine       string `json:"raw_line"`
}
