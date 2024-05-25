package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/5pots-com/cli/internal/champion"
	"github.com/5pots-com/cli/internal/common"
	"github.com/spf13/cobra"
)

type Character struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Mtime string `json:"mtime"`
}

const (
	charactersURL = "https://raw.communitydragon.org/json/%s/game/data/characters/"
)

var filteredCharacters = []string{"tft", "tutorial", "trinket", "bw_", "cherry_", "durian_", "ha_", "hexgate", "item_", "kingporo", "nexus", "pet", "slime_", "sru_", "ultbook", "sruap_", "srx", "test", "practicetool_", "preseason_", "spellbook", "sr_infernal", "summonerbeacon"}

var downloadChampionsCmd = &cobra.Command{
	Use:   "champions [name]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if dirty {
			fmt.Printf("WARNING: --dirty flag detected, this will leave diff noises on files...\n")
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		dir := fmt.Sprintf("%s/data/champions", wd)

		hasNewLatest, err := common.HasNewVersion(common.Latest, dir)
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		hasNewPBE, err := common.HasNewVersion(common.PBE, dir)
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		var wg sync.WaitGroup

		chs, err := downCharacters(common.PBE)
		if err != nil {
			log.Fatalf("Failed to download list of characters: %v", err)
		}

		m, err := common.DownloadMetadata()
		if err != nil {
			log.Fatalf("Failed to download patch metadata: %v", err)
		}

		for _, ch := range chs {
			c := &champion.Champion{Name: ch}

			if hasNewLatest || force {
				wg.Add(1)
				go DownSaveLatest(c, dir, &wg)
			} else {
				fmt.Printf("Live version already downloaded for %s\n", c.Name)
			}
			if hasNewPBE || force {
				wg.Add(1)
				go downSavePBE(c, dir, &wg)
			} else {
				fmt.Printf("PBE version already downloaded for %s\n", c.Name)
			}
			if hasNewLatest || hasNewPBE {
				wg.Add(1)
				go downSaveMetaData(c, dir, &wg)
			}
			wg.Wait()
			time.Sleep(time.Second / 2)
		}

		common.SaveMeta(m, wd)

		fmt.Println("Success!")
	},
}

func init() {
	downloadCmd.AddCommand(downloadChampionsCmd)
	downloadChampionsCmd.Flags().BoolVarP(&dirty, "dirty", "d", false, "leave noise in data for debugging")
	downloadChampionsCmd.Flags().BoolVarP(&force, "force", "f", false, "ignores metadata matching version and downloads the champion anyway")
}

func downCharacters(p common.Patch) ([]string, error) {
	if !common.Validate(p) {
		return nil, fmt.Errorf("invalid p: %s", p)
	}

	url := fmt.Sprintf(charactersURL, p)
	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var characters []Character
	if err := json.Unmarshal(body, &characters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	var list []string

	for _, c := range characters {
		if !contains(filteredCharacters, c.Name) {
			list = append(list, c.Name)
		}
	}

	return list, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(e, a) {
			return true
		}
	}
	return false
}
