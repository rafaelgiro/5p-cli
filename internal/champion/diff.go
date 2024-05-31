package champion

import (
	"encoding/json"
	"fmt"

	"github.com/5pots-com/cli/internal/common"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

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

func (c *Champion) SaveDiff(dir, outputDir string) error {
	diffs, live, _, err := c.LoadAndDiff(dir)
	if err != nil {
		return fmt.Errorf("failed to diff %s PBE and live files: %v", dir, err)
	}

	var liveMap map[string]interface{}

	if err := json.Unmarshal(live, &liveMap); err != nil {
		return fmt.Errorf("failed to covnert live file %s to json: %v", dir, err)
	}

	filteredLive := make(map[string]interface{})

	for key := range diffs {
		filteredLive[key] = diffs[key]
	}

	res := make(map[string]interface{})

	res["latest"] = filteredLive
	res["diff"] = diffs
	res["result"] = make(map[string]interface{})

	js, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to convert result json of %s: %v", c.Name, err)
	}

	fmt.Println("Saving to file...")
	fileName := fmt.Sprintf("%s.json", c.Name)
	if err := common.SaveToFile(outputDir, fileName, js); err != nil {
		return fmt.Errorf("failed to save file %s to %s: %v", fileName, outputDir, err)
	}

	return nil

}

func doDiff(dir string, live, pbe []byte) (diff.Diff, error) {
	fmt.Printf("Finding differences...\n")

	jd := diff.New()
	diffs, err := jd.Compare(live, pbe)
	if err != nil {
		return nil, fmt.Errorf("failed to compare /%s PBE and live files: %v", dir, err)
	}

	return diffs, nil
}
