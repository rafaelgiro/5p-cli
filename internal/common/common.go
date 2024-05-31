package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Patch string

type Character struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Mtime string `json:"mtime"`
}

const (
	PBE          Patch  = "pbe"
	Latest       Patch  = "latest"
	PBEFile      string = "pbe.json"
	LatestFile   string = "live.json"
	DataFolder   string = "data"
	OutputFolder string = "results"
)

func Validate(patch Patch) bool {
	switch patch {
	case PBE, Latest:
		return true
	}
	return false
}

func Format(jsonData []byte) ([]byte, error) {
	var out bytes.Buffer

	if err := json.Indent(&out, jsonData, "", "  "); err != nil {
		return nil, fmt.Errorf("failed to format json file: %v ", err)
	}

	return out.Bytes(), nil
}

func SaveToFile(dir, fileName string, data []byte) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	filePath := fmt.Sprintf("%s/%s", dir, fileName)
	fdata, err := Format(data)

	if err != nil {
		return fmt.Errorf("failed to format file %s: %v", dir, err)
	}

	if err := os.WriteFile(filePath, fdata, 0666); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}
