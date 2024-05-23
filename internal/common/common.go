package common

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Patch string

const (
	PBE    Patch = "pbe"
	Latest Patch = "latest"
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
