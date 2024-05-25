package champion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/5pots-com/cli/internal/common"
)

type Champion struct {
	Name string
}

type Strings struct {
	Entries map[string]string `json:"entries"`
}

const (
	champURL   = "https://raw.communitydragon.org/%s/game/data/characters/%s/%s.bin.json"
	stringsURL = "https://raw.communitydragon.org/%s/game/en_us/data/menu/en_us/main.stringtable.json"
)

func (c *Champion) Download(patch common.Patch, clean bool) ([]byte, error) {
	if !common.Validate(patch) {
		return nil, fmt.Errorf("invalid patch: %s", patch)
	}

	d, err := downChamp(c.Name, patch)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s champion data: %v", c.Name, err)
	}

	s, err := downStrings(patch, c.Name, clean)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tooltips data: %v", c.Name, err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(d, &data); err != nil {
		fmt.Printf("WARNING: Data not found for %s. Saving as blank.\n", c.Name)
		return []byte("{}"), nil
	}

	var strs map[string]interface{}

	if err := json.Unmarshal(s, &strs); err != nil {
		fmt.Printf("WARNING: Data not found for tooltips on %s. Saving as blank.\n", c.Name)
		return []byte("{}"), nil
	}

	data["tooltips"] = strs

	ch, err := json.Marshal(data)

	if err != nil {
		return nil, fmt.Errorf("failed to convert %s champion json: %v", c.Name, err)
	}

	return ch, nil
}

func downChamp(name string, patch common.Patch) ([]byte, error) {
	url := fmt.Sprintf(champURL, patch, name, name)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body from %s: %v", url, err)
	}

	return body, nil
}

func downStrings(patch common.Patch, name string, clean bool) ([]byte, error) {
	url := fmt.Sprintf(stringsURL, patch)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var data Strings

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse to json: %s: %v", url, err)
	}

	targetKey := fmt.Sprintf("generatedtip_spell_%s", name)
	final := make(map[string]string)

	for key, value := range data.Entries {
		if !clean && strings.Contains(key, targetKey) {
			final[key] = value
		} else if clean && strings.Contains(key, targetKey) && strings.Contains(key, "tooltipextended") {
			final[key] = value
		}
	}

	d, err := json.Marshal(final)

	if err != nil {
		return nil, fmt.Errorf("failed convert json %s: %v", name, err)
	}

	return d, nil
}

func RemoveNoise(data []byte) []byte {
	replacements := map[string]string{
		`"mFormat":"\{.*?}`:         `"mFormat":"{loveusion}`,
		`,"mAllStartingItemIds":.*`: "}}",
		`"EventToTrack":.*?,`:       `"EventToTrack": 0,`,
	}

	for pat, rep := range replacements {
		re := regexp.MustCompile(pat)
		data = re.ReplaceAll(data, []byte(rep))
	}

	return data
}

func (c Champion) SaveToFile(dir, fileName string, data []byte) error {
	finalDir := fmt.Sprintf("%s/%s", dir, c.Name)

	if err := os.MkdirAll(finalDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", finalDir, err)
	}

	filePath := fmt.Sprintf("%s/%s", finalDir, fileName)
	fdata, err := common.Format(data)

	if err != nil {
		return fmt.Errorf("failed to format file %s: %v", finalDir, err)
	}

	if err := os.WriteFile(filePath, fdata, 0666); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}
