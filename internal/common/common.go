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
