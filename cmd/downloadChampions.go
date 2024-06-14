package cmd

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/5pots-com/cli/internal/champion"
	"github.com/5pots-com/cli/internal/common"
	"github.com/spf13/cobra"
)

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

		dir := fmt.Sprintf("%s/data", wd)

		hasNewLatest, err := common.HasNewVersion(common.Latest, dir)
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		hasNewPBE, err := common.HasNewVersion(common.PBE, dir)
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		var wg sync.WaitGroup

		chs, err := common.DownCharacters(common.PBE, champion.FilteredCharacters)
		if err != nil {
			log.Fatalf("Failed to download list of characters: %v", err)
		}

		m, err := common.DownloadMetadata()
		if err != nil {
			log.Fatalf("Failed to download patch metadata: %v", err)
		}

		for _, ch := range chs {
			c := &champion.Champion{Name: ch}
			champDir := fmt.Sprintf("/%s/champions/%s", dir, c.Name)

			if hasNewLatest || force {
				wg.Add(1)
				go c.DownSaveLatest(champDir, dirty, &wg)
			} else {
				fmt.Printf("Live version already downloaded for %s\n", c.Name)
			}
			if hasNewPBE || force {
				wg.Add(1)
				go c.DownSavePBE(champDir, dirty, &wg)
			} else {
				fmt.Printf("PBE version already downloaded for %s\n", c.Name)
			}
			if hasNewLatest || hasNewPBE {
				wg.Add(1)
				go c.DownSaveMetaData(champDir, &wg)
			}
			wg.Wait()
			time.Sleep(time.Second / 4)
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
