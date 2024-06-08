package champion

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/5pots-com/cli/internal/common"
	"github.com/mitchellh/mapstructure"
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
	diffs, ld, pd, err := c.LoadAndDiff(dir)
	if err != nil {
		return fmt.Errorf("failed to diff %s PBE and live files: %v", dir, err)
	}

	live, err := decode(ld)
	if err != nil {
		return fmt.Errorf("failed to decode %s Live files: %v", dir, err)
	}

	pbe, err := decode(pd)
	if err != nil {
		return fmt.Errorf("failed to decode %s PBE files: %v", dir, err)
	}

	mount(live, diffs)

	// ====================================================

	result := make(map[string]interface{})
	result["live"] = live
	result["pbe"] = pbe

	ch, err := json.Marshal(result)

	if err != nil {
		return fmt.Errorf("failed to convert %s champion json: %v", c.Name, err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	d := fmt.Sprintf("%s/data/test", wd)

	common.SaveToFile(d, "test.json", ch)

	// fmt.Println(live)

	// filteredLive := make(map[string]interface{})

	// for key := range diffs {
	// 	filteredLive[key] = diffs[key]
	// }

	// res := make(map[string]interface{})

	// res["latest"] = filteredLive
	// res["diff"] = diffs
	// res["result"] = make(map[string]interface{})

	// js, err := json.Marshal(res)
	// if err != nil {
	// 	return fmt.Errorf("failed to convert result json of %s: %v", c.Name, err)
	// }

	// fmt.Println("Saving to file...")
	// fileName := fmt.Sprintf("%s.json", c.Name)
	// if err := common.SaveToFile(outputDir, fileName, js); err != nil {
	// 	return fmt.Errorf("failed to save file %s to %s: %v", fileName, outputDir, err)
	// }

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

func mount(d JSONData, di map[string]interface{}) error {
	for key := range di {
		if strings.Contains(key, "/Spells/") {
			s := strings.Split(key, "/")
			an := strings.ToLower(s[len(s)-1])
			ttp := d.Tooltips[fmt.Sprintf("generatedtip_spell_%s_tooltipextended", an)]
			spl := d.Character[key].Spell

			f, err := handleTooltip(ttp, spl)

			if err != nil {
				return fmt.Errorf("failed to convert champion ability to tooltip: %s; %v", key, err)
			}

			fmt.Println(f)

		}
	}

	return nil
}

func handleTooltip(ttp string, spl SpellDataResource) (string, error) {
	c := removeHTMLTags(ttp)

	// Regex to find all variables
	re := regexp.MustCompile(`@(.*?)@`)
	matches := re.FindAllStringSubmatch(c, -1)

	// Regex to grab the variable name and index
	vre := regexp.MustCompile(`(\D+?)(\d+)`)

	for _, match := range matches {
		if len(match) > 1 {
			// fmt.Println(match[1])
			matches := vre.FindStringSubmatch(match[1])

			if len(matches) == 3 {
				w := matches[1]
				i, err := strconv.Atoi(matches[2])
				if err != nil {
					return "", fmt.Errorf("failed to find index to ability variable. %s; %v", matches[1], err)
				}

				for _, val := range spl.DataValues {
					if val.Name == w {
						str := fmt.Sprint(val.Values[i])
						c = strings.Replace(c, match[1], str, -1)
					}
				}

			}
		}
	}

	fmt.Println("-----------------")
	fmt.Println(ttp)
	fmt.Println("==")

	return c, nil
}

func removeHTMLTags(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	result := re.ReplaceAllString(input, "")

	reSpecial := regexp.MustCompile(`@[^@]*?(?:Postfix|Prefix)@`)
	result = reSpecial.ReplaceAllString(result, "")
	return result
}
