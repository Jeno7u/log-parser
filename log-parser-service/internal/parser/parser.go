package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jeno7u/log-parser/internal/dto"
)

type ParsedLog struct {
	ArchiveName string
	Nodes       []dto.Node
	Ports       []dto.Port
	NodeInfo    []dto.NodeInfo
}

func Parse(inputs ArchiveInputs) (ParsedLog, error) {
	nodes, err := parseNodes(inputs.DBCSV)
	if err != nil {
		return ParsedLog{}, err
	}

	ports, err := parsePorts(inputs.DBCSV)
	if err != nil {
		return ParsedLog{}, err
	}

	nodeInfo, err := parseNodeInfo(inputs.SharpInfo)
	if err != nil {
		return ParsedLog{}, err
	}

	return ParsedLog{ArchiveName: inputs.ArchiveName, Nodes: nodes, Ports: ports, NodeInfo: nodeInfo}, nil
}

func parseNodes(content []byte) ([]dto.Node, error) {
	lines, err := readSection(content, "START_NODES", "END_NODES")
	if err != nil {
		return nil, err
	}
	if len(lines) < 2 {
		return nil, fmt.Errorf("nodes section is incomplete")
	}

	var nodes []dto.Node
	for _, line := range lines[1:] {
		parts := csvSplit(line)
		if len(parts) < 8 {
			return nil, fmt.Errorf("invalid node row: %s", line)
		}

		numPorts, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid num ports in node row: %w", err)
		}

		nodeType, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid node type in node row: %w", err)
		}

		classVersion, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid class version in node row: %w", err)
		}

		baseVersion, err := strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("invalid base version in node row: %w", err)
		}

		nodes = append(nodes, dto.Node{NodeDesc: trimQuotes(parts[0]), NumPorts: numPorts, NodeType: nodeType, ClassVersion: classVersion, BaseVersion: baseVersion, SystemImageGUID: parts[5], NodeGUID: parts[6], PortGUID: parts[7]})
	}

	return nodes, nil
}

func parsePorts(content []byte) ([]dto.Port, error) {
	lines, err := readSection(content, "START_PORTS", "END_PORTS")
	if err != nil {
		return nil, err
	}
	if len(lines) < 2 {
		return nil, fmt.Errorf("ports section is incomplete")
	}

	var ports []dto.Port
	for _, line := range lines[1:] {
		parts := csvSplit(line)
		if len(parts) < 21 {
			return nil, fmt.Errorf("invalid port row: %s", line)
		}

		portNum, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid port number in port row: %w", err)
		}

		localPortNum, err := strconv.Atoi(parts[13])
		if err != nil {
			return nil, fmt.Errorf("invalid local port number in port row: %w", err)
		}

		ports = append(ports, dto.Port{NodeGUID: parts[0], PortGUID: parts[1], PortNum: portNum, LocalPortNum: localPortNum, LinkWidthActv: parts[10], LinkSpeedActv: parts[15], PortPhyState: parts[19], PortState: parts[20], RawLine: line})
	}

	return ports, nil
}

func parseNodeInfo(content []byte) ([]dto.NodeInfo, error) {
	var infos []dto.NodeInfo
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var currentSWGUID string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "-") {
			continue
		}
		if strings.HasPrefix(line, "SW_GUID=") {
			currentSWGUID = strings.TrimSpace(strings.TrimPrefix(line, "SW_GUID="))
			continue
		}
		if currentSWGUID == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid info line: %s", line)
		}
		infos = append(infos, dto.NodeInfo{SWGUID: currentSWGUID, Key: strings.TrimSpace(parts[0]), Value: strings.TrimSpace(parts[1])})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, fmt.Errorf("nodes info section is incomplete")
	}

	return infos, nil
}

func readSection(content []byte, start string, end string) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(content))
	inside := false
	foundEnd := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == start {
			inside = true
			lines = []string{}
			continue
		}
		if line == end {
			if !inside {
				continue
			}
			foundEnd = true
			break
		}
		if inside && line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("section %s not found", start)
	}
	if !foundEnd {
		return nil, fmt.Errorf("section %s is not terminated", start)
	}

	return lines, nil
}

func csvSplit(line string) []string {
	parts := strings.Split(line, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(trimQuotes(parts[i]))
	}
	return parts
}

func trimQuotes(value string) string {
	return strings.Trim(value, `"`)
}
