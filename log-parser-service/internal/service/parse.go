package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jeno7u/log-parser/internal/dto"
	"github.com/Jeno7u/log-parser/internal/parser"
	"github.com/Jeno7u/log-parser/internal/repository"
)

type ParseService interface {
	Parse(ctx context.Context, sourcePath string, fileName string) (string, error)
}

type Parse struct {
	dataDir  string
	logRepo  repository.LogRepository
	nodeRepo repository.NodeRepository
	portRepo repository.PortRepository
}

func NewParse(dataDir string, logRepo repository.LogRepository, nodeRepo repository.NodeRepository, portRepo repository.PortRepository) ParseService {
	return &Parse{dataDir: dataDir, logRepo: logRepo, nodeRepo: nodeRepo, portRepo: portRepo}
}

func (s *Parse) Parse(ctx context.Context, sourcePath string, fileName string) (string, error) {
	resolvedPath := s.resolveSourcePath(sourcePath)

	inputs, err := parser.ReadInputs(resolvedPath)
	if err != nil {
		return "", err
	}

	parsed, err := parser.Parse(inputs)
	if err != nil {
		return "", err
	}

	if fileName == "" {
		fileName = filepath.Base(resolvedPath)
	}

	logID, err := s.logRepo.CreateLog(ctx, dto.Log{FileName: fileName, SourcePath: resolvedPath, Status: "parsing", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err != nil {
		return "", err
	}

	if err := s.nodeRepo.CreateNodes(ctx, logID, parsed.Nodes); err != nil {
		_ = s.logRepo.UpdateLogStatus(ctx, logID, "failed", err.Error())
		return "", err
	}

	nodes, err := s.nodeRepo.ListNodesByLogID(ctx, logID)
	if err != nil {
		_ = s.logRepo.UpdateLogStatus(ctx, logID, "failed", err.Error())
		return "", err
	}

	nodeIDs := make(map[string]string, len(nodes))
	for _, node := range nodes {
		nodeIDs[node.NodeGUID] = node.ID
	}

	if err := s.portRepo.CreatePorts(ctx, logID, nodeIDs, parsed.Ports); err != nil {
		_ = s.logRepo.UpdateLogStatus(ctx, logID, "failed", err.Error())
		return "", err
	}

	ports, err := s.portRepo.ListPortsByLogID(ctx, logID)
	if err != nil {
		_ = s.logRepo.UpdateLogStatus(ctx, logID, "failed", err.Error())
		return "", err
	}

	links := make([]dto.TopologyLink, 0, len(ports))
	for _, port := range ports {
		if port.NodeID == "" || port.ID == "" {
			continue
		}
		links = append(links, dto.TopologyLink{LogID: logID, NodeID: port.NodeID, PortID: port.ID, RelationType: "node_port"})
	}

	if err := s.logRepo.CreateLinks(ctx, logID, links); err != nil {
		_ = s.logRepo.UpdateLogStatus(ctx, logID, "failed", err.Error())
		return "", err
	}

	if err := s.nodeRepo.CreateNodeInfo(ctx, logID, parsed.NodeInfo); err != nil {
		_ = s.logRepo.UpdateLogStatus(ctx, logID, "failed", err.Error())
		return "", err
	}

	if err := s.logRepo.UpdateLogStatus(ctx, logID, "parsed", ""); err != nil {
		return "", fmt.Errorf("update parsed log status: %w", err)
	}

	return logID, nil
}

func (s *Parse) resolveSourcePath(sourcePath string) string {
	if filepath.IsAbs(sourcePath) {
		return sourcePath
	}

	cleanPath := filepath.Clean(sourcePath)
	dataPrefix := "data" + string(filepath.Separator)
	if cleanPath == "data" {
		return s.dataDir
	}
	cleanPath = strings.TrimPrefix(cleanPath, dataPrefix)

	return filepath.Join(s.dataDir, cleanPath)
}
