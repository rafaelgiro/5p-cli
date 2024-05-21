package champion

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/5pots-com/cli/internal/common"
)

type Champion struct {
	Name string
	Data []byte
}

const (
	champURL = "https://raw.communitydragon.org/%s/game/data/characters/%s/%s.bin.json"
)

func (c *Champion) Download(patch common.Patch) ([]byte, error) {
	if !common.Validate(patch) {
		return nil, fmt.Errorf("invalid patch: %s", patch)
	}

	url := fmt.Sprintf(champURL, patch, c.Name, c.Name)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}

func RemoveNoise(data []byte) []byte {
	replacements := map[string]string{
		`"mFormat":"\{.*?\}",`:     "",
		`"mAllStartingItemIds":.*`: "}}",
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

	if err := os.WriteFile(filePath, data, 0666); err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}

	return nil
}
