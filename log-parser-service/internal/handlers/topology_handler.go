package handlers

import (
	"log/slog"
	"net/http"
	"sort"
	"strings"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/Jeno7u/log-parser/internal/repository"
)

type TopologyHandler struct {
	logRepo  repository.LogRepository
	nodeRepo repository.NodeRepository
	portRepo repository.PortRepository
	log      *slog.Logger
}

func NewTopology(logRepo repository.LogRepository, nodeRepo repository.NodeRepository, portRepo repository.PortRepository, log *slog.Logger) *TopologyHandler {
	return &TopologyHandler{logRepo: logRepo, nodeRepo: nodeRepo, portRepo: portRepo, log: log}
}

func (h *TopologyHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, dto.ErrorResponse{Error: "method not allowed"})
		return
	}

	logID := strings.TrimPrefix(r.URL.Path, "/api/v1/topology/")
	if logID == "" {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "log_id is required"})
		return
	}

	nodes, err := h.nodeRepo.ListNodesByLogID(r.Context(), logID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ports, err := h.portRepo.ListPortsByLogID(r.Context(), logID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	links, err := h.logRepo.ListLinksByLogID(r.Context(), logID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}

	result := buildTopology(logID, nodes, ports, links)
	writeJSON(w, http.StatusOK, result)
}

func buildTopology(logID string, nodes []dto.Node, ports []dto.Port, links []dto.TopologyLink) dto.Topology {
	groups := map[string][]dto.Node{
		"switches": []dto.Node{},
		"hosts":    []dto.Node{},
	}

	for _, node := range nodes {
		if isHostNode(node) {
			groups["hosts"] = append(groups["hosts"], node)
			continue
		}
		groups["switches"] = append(groups["switches"], node)
	}

	groupNames := make([]string, 0, len(groups))
	for name, groupNodes := range groups {
		if len(groupNodes) == 0 {
			continue
		}
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	topologyGroups := make([]dto.TopologyGroup, 0, len(groupNames))
	for _, name := range groupNames {
		topologyGroups = append(topologyGroups, dto.TopologyGroup{Name: name, Nodes: groups[name]})
	}

	return dto.Topology{LogID: logID, Groups: topologyGroups, Ports: ports, Links: links}
}

func isHostNode(node dto.Node) bool {
	if node.NodeType == 1 {
		return true
	}

	name := strings.TrimSpace(strings.ToLower(node.NodeDesc))
	return strings.HasPrefix(name, "host")
}
