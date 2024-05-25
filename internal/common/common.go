package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Patch string

const (
	PBE         Patch  = "pbe"
	Latest      Patch  = "latest"
	MetadataURL string = "https://raw.communitydragon.org/%s/content-metadata.json"
	dataFolder  string = "data"
)

const (
	// Colors and font options via ANSI escape codes
	Reset     = "\033[0m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Invert    = "\033[7m"
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

func DownloadMetadata() ([]byte, error) {
	pbe, err := downMeta(PBE)
	if err != nil {
		return nil, fmt.Errorf("failed to pbe metadata: %v", err)
	}

	live, err := downMeta(Latest)
	if err != nil {
		return nil, fmt.Errorf("failed to latest metadata: %v", err)
	}

	var map1, map2 map[string]string
	if err := json.Unmarshal(pbe, &map1); err != nil {
		return nil, fmt.Errorf("error unmarshaling pbe json: %v", err)
	}
	if err := json.Unmarshal(live, &map2); err != nil {
		return nil, fmt.Errorf("error unmarshaling live json: %v", err)
	}

	mergedMap := make(map[string]string)
	mergedMap["pbe"] = map1["version"]
	mergedMap["latest"] = map2["version"]

	mergedJSON, err := json.Marshal(mergedMap)
	if err != nil {
		return nil, fmt.Errorf("error marshaling merged map: %v", err)
	}

	return mergedJSON, nil
}

func HasNewVersion(p Patch, dir string) (bool, error) {
	curr, err := os.ReadFile(fmt.Sprintf("%s/metadata.json", dir))

	if err != nil {
		return true, nil
	}

	m, err := downMeta(p)

	if err != nil {
		return false, err
	}

	var map1 map[Patch]string
	var map2 map[string]string
	if err := json.Unmarshal(curr, &map1); err != nil {
		return false, fmt.Errorf("error unmarshaling pbe json: %v", err)
	}

	if err := json.Unmarshal(m, &map2); err != nil {
		return false, fmt.Errorf("error unmarshaling pbe json: %v", err)
	}

	if map2["version"] != map1[p] {
		return true, nil
	} else {
		return false, nil
	}
}

func downMeta(p Patch) ([]byte, error) {
	url := fmt.Sprintf(MetadataURL, p)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadatadata from %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func SaveMeta(m []byte, dir string) error {
	if err := os.MkdirAll(dataFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dataFolder, err)
	}

	filePath := fmt.Sprintf("%s/%s/%s", dir, dataFolder, "metadata.json")

	fdata, err := Format(m)

	if err != nil {
		return fmt.Errorf("failed to format file %s: %v", dataFolder, err)
	}

	if err := os.WriteFile(filePath, fdata, 0666); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}
