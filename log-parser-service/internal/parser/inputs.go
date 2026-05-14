package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ArchiveInputs struct {
	ArchiveName string
	DBCSV       []byte
	SharpInfo   []byte
}

func ReadInputs(path string) (ArchiveInputs, error) {
	return readDirectory(path)
}

func readDirectory(path string) (ArchiveInputs, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return ArchiveInputs{}, err
	}

	inputs := ArchiveInputs{ArchiveName: filepath.Base(path)}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, readErr := os.ReadFile(filepath.Join(path, entry.Name()))
		if readErr != nil {
			return ArchiveInputs{}, readErr
		}
		assignArchiveFile(&inputs, entry.Name(), data)
	}

	return ensureArchiveInputs(inputs)
}

func assignArchiveFile(inputs *ArchiveInputs, name string, data []byte) {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "db_csv"):
		inputs.DBCSV = data
	case strings.Contains(lower, "sharp_an_info"):
		inputs.SharpInfo = data
	}
}

func ensureArchiveInputs(inputs ArchiveInputs) (ArchiveInputs, error) {
	if len(inputs.DBCSV) == 0 {
		return ArchiveInputs{}, fmt.Errorf("archive is missing ibdiagnet2.db_csv")
	}
	if len(inputs.SharpInfo) == 0 {
		return ArchiveInputs{}, fmt.Errorf("archive is missing ibdiagnet2.sharp_an_info")
	}

	return inputs, nil
}
