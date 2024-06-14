package champion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/5pots-com/cli/internal/common"
)

func (c *Champion) DownSavePBE(dir string, dirty bool, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Printf("Downloading %s data on patch \"%s\"...\n", c.Name, common.PBE)

	pbe, err := c.Download(common.PBE, !dirty)
	if err != nil {
		return fmt.Errorf("failed to Download PBE data: %v", err)
	}

	if !dirty {
		pbe = RemoveNoise(pbe)
	}

	if err := c.SaveToFile(dir, common.PBEFile, pbe); err != nil {
		return fmt.Errorf("failed to save file %s: %v", common.PBEFile, err)
	}

	return nil
}

func (c *Champion) DownSaveLatest(dir string, dirty bool, wg *sync.WaitGroup) error {
	defer wg.Done()
	fmt.Printf("Downloading %s data on patch \"%s\"...\n", c.Name, common.Latest)

	live, err := c.Download(common.Latest, !dirty)
	if err != nil {
		return fmt.Errorf("failed to Download Live data: %v", err)
	}

	if !dirty {
		live = RemoveNoise(live)
	}

	if err := c.SaveToFile(dir, common.LatestFile, live); err != nil {
		return fmt.Errorf("failed to save file %s: %v", common.LatestFile, err)
	}

	return nil
}

func (c *Champion) DownSaveMetaData(dir string, wg *sync.WaitGroup) error {
	defer wg.Done()

	m, err := common.DownloadMetadata()
	if err != nil {
		return fmt.Errorf("failed to download metadata: %v", err)
	}

	if err := c.SaveToFile(dir, common.MetaFile, m); err != nil {
		return fmt.Errorf("failed to save file %s: %v", common.MetaFile, err)
	}

	return nil
}

func (c *Champion) Download(patch common.Patch, clean bool) ([]byte, error) {
	if !common.Validate(patch) {
		return nil, fmt.Errorf("invalid patch: %s", patch)
	}

	d, err := downChamp(c.Name, patch)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s champion data: %v", c.Name, err)
	}

	s, err := downStrings(patch, c.Name, d)

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

func downStrings(patch common.Patch, name string, c []byte) ([]byte, error) {
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
		return nil, fmt.Errorf("failed to parse strings to json: %s: %v", url, err)
	}

	var ch map[string]interface{}
	if err := json.Unmarshal(c, &ch); err != nil {
		return nil, fmt.Errorf("failed to parse champion to json: %s: %v", url, err)
	}

	final := make(map[string]string)

	for k := range ch {
		sp := strings.Split(k, "/")
		if sp[0] == "Characters" && sp[2] == "Spells" && len(sp) == 5 {
			ttp := fmt.Sprintf("generatedtip_spell_%s_tooltipextended", strings.ToLower(sp[4]))
			if len(data.Entries[ttp]) != 0 {
				final[ttp] = data.Entries[ttp]
			}
		}
	}

	d, err := json.Marshal(final)

	if err != nil {
		return nil, fmt.Errorf("failed convert json %s: %v", name, err)
	}

	return d, nil
}
