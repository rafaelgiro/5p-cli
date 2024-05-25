package champion

import (
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

func (c *Champion) LoadAndDiff(dir string) (map[string]interface{}, error) {
	live, err := c.Load(dir, common.LatestFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open \"%s\" live data", c.Name)
	}

	pbe, err := c.Load(dir, common.PBEFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open \"%s\" PBE data", c.Name)
	}

	fmt.Printf("Finding differences...\n")

	jd := diff.New()
	diffs, err := jd.Compare(live, pbe)
	if err != nil {
		return nil, fmt.Errorf("failed to compare /%s PBE and live files: %v", dir, err)
	}

	formatter := formatter.NewDeltaFormatter()

	diff, err := formatter.FormatAsJson(diffs)
	if err != nil {
		return nil, fmt.Errorf("failed to format diff result: %v", err)
	}

	return diff, nil

}
