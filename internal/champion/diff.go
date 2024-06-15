package champion

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/5pots-com/cli/internal/common"
	"github.com/mitchellh/mapstructure"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

type diffResult struct {
	Keys []string               `json:"keys"`
	Live map[string]string      `json:"live"`
	PBE  map[string]string      `json:"pbe"`
	Diff map[string]interface{} `json:"diff"`
}

var blacklist = []string{"yuumi"}

func (c *Champion) CheckDownload(dir string) error {
	files := []string{
		fmt.Sprintf("%s/%s/%s", dir, c.Name, common.PBEFile),
		fmt.Sprintf("%s/%s/%s", dir, c.Name, common.LatestFile),
	}

	if err := common.CheckDownload(files); err != nil {
		return fmt.Errorf("champion \"%s\" not yet downloaded", c.Name)
	}

	return nil
}

func (c *Champion) LoadAndDiff(dir string) (map[string]interface{}, []byte, []byte, error) {
	live, err := c.Load(dir, common.LatestFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open \"%s\" live data", c.Name)
	}

	pbe, err := c.Load(dir, common.PBEFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open \"%s\" PBE data", c.Name)
	}

	fmt.Printf("Finding differences for: %s...\n", c.Name)
	diffs, err := doDiff(dir, live, pbe)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to diff \"%s\" PBE and live data", c.Name)
	}

	formatter := formatter.NewDeltaFormatter()

	diff, err := formatter.FormatAsJson(diffs)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to format diff result: %v", err)
	}

	return diff, live, pbe, nil

}

func (c *Champion) PrepareDiff(dir, outputDir string) (diffResult, error) {
	for _, bc := range blacklist {
		if bc == c.Name {
			fmt.Printf("Ignoring champion for now: %s\n", c.Name)
			return diffResult{}, nil
		}
	}

	diffs, ld, pd, err := c.LoadAndDiff(dir)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to diff %s PBE and live files: %v", dir, err)
	}

	keys := []string{}
	for key := range diffs {
		if key[0] != '{' {
			keys = append(keys, key)
		}
	}

	if len(keys) == 0 {
		fmt.Printf("No changes found for %s\n", c.Name)
		return diffResult{}, nil
	}

	live, err := decode(ld)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to decode %s Live files: %v", dir, err)
	}

	pbe, err := decode(pd)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to decode %s PBE files: %v", dir, err)
	}

	lttps, err := mount(live, diffs)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to read live tooltips: %v", err)
	}

	pttps, err := mount(pbe, diffs)
	if err != nil {
		return diffResult{}, fmt.Errorf("failed to read PBE tooltips: %v", err)
	}

	result := diffResult{}
	result.Live = lttps
	result.PBE = pttps
	result.Keys = keys
	result.Diff = diffs

	return result, nil
}

func doDiff(dir string, live, pbe []byte) (diff.Diff, error) {
	jd := diff.New()
	diffs, err := jd.Compare(live, pbe)
	if err != nil {
		return nil, fmt.Errorf("failed to compare /%s PBE and live files: %v", dir, err)
	}

	return diffs, nil
}

func decode(d []byte) (JSONData, error) {
	m := DownloadedData{}
	if err := json.Unmarshal(d, &m.Character); err != nil {
		return JSONData{}, fmt.Errorf("failed to convert live file to json: %v", err)
	}
	if n, ok := m.Character["tooltips"].(map[string]interface{}); ok {
		m.Tooltips = n
	}
	delete(m.Character, "tooltips")

	var result JSONData
	if err := mapstructure.Decode(m, &result); err != nil {
		return JSONData{}, fmt.Errorf("failed to decode data to json: %v", err)
	}

	return result, nil
}

func mount(d JSONData, diffs map[string]interface{}) (map[string]string, error) {
	ts := map[string]string{}

	for k := range diffs {
		sp := strings.Split(k, "/")

		// Handle spell tooltip
		if sp[0] == "Characters" && sp[2] == "Spells" && len(sp) == 5 {
			k := sp[4]
			tk := fmt.Sprintf("generatedtip_spell_%s_tooltipextended", strings.ToLower(k))
			tp := d.Tooltips[tk]
			spl := d.Character[k].Spell

			v, err := HandleTooltip(tp, spl)
			if err != nil {
				return map[string]string{}, fmt.Errorf("failed to convert champion ability to tooltip: %s; %v", k, err)
			}

			if len(v) != 0 {
				ts[k] = v
			}
		}
	}

	return ts, nil
}
