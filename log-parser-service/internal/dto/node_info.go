package dto

type NodeInfo struct {
	ID     string `json:"node_info_id"`
	LogID  string `json:"log_id"`
	SWGUID string `json:"sw_guid"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
