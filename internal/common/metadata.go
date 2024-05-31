package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	MetadataURL string = "https://raw.communitydragon.org/%s/content-metadata.json"
	MetaFile    string = "metadata.json"
)

func DownloadMetadata() ([]byte, error) {
	pbe, err := down(PBE)
	if err != nil {
		return nil, fmt.Errorf("failed to download PBE metadata: %v", err)
	}

	live, err := down(Latest)
	if err != nil {
		return nil, fmt.Errorf("failed to download latest metadata: %v", err)
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
	curr, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, MetaFile))

	if err != nil {
		return true, nil
	}

	m, err := down(p)

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

func SaveMeta(m []byte, dir string) error {
	if err := os.MkdirAll(DataFolder, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", DataFolder, err)
	}

	filePath := fmt.Sprintf("%s/%s/%s", dir, DataFolder, "metadata.json")

	fdata, err := Format(m)

	if err != nil {
		return fmt.Errorf("failed to format file %s: %v", DataFolder, err)
	}

	if err := os.WriteFile(filePath, fdata, 0666); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}

func down(p Patch) ([]byte, error) {
	url := fmt.Sprintf(MetadataURL, p)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata from URL %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body from URL %s: %v", url, err)
	}

	return body, nil
}
