/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/5pots-com/cli/internal/champion"
	"github.com/5pots-com/cli/internal/common"
	"github.com/spf13/cobra"
)

var downloadChampionCmd = &cobra.Command{
	Use:   "champion [name]",
	Short: "Download data for a specific champion",
	Long: `Download the latest data for a specific champion from the PBE and Live server.

This command fetches the latest champion data and optionally removes noise from the data.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if dirty {
			fmt.Printf("WARNING: --dirty flag detected, this will leave diff noises on files...\n")
		}

		champName := strings.ToLower(args[0])
		c := &champion.Champion{Name: champName}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		dir := fmt.Sprintf("%s/data/champions/%s", wd, c.Name)

		hasNewLatest, err := common.HasNewVersion(common.Latest, dir)
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		hasNewPBE, err := common.HasNewVersion(common.PBE, dir)
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		var wg sync.WaitGroup

		if hasNewLatest || force {
			wg.Add(1)
			go c.DownSaveLatest(dir, dirty, &wg)
		} else {
			fmt.Printf("Live version already downloaded for %s\n", c.Name)
		}
		if hasNewPBE || force {
			wg.Add(1)
			go c.DownSavePBE(dir, dirty, &wg)
		} else {
			fmt.Printf("PBE version already downloaded for %s\n", c.Name)
		}
		if hasNewLatest || hasNewPBE {
			wg.Add(1)
			go c.DownSaveMetaData(dir, &wg)
		}

		wg.Wait()
		fmt.Printf("Success!\n")
	},
	Example: `  # Download data for a specific champion
  pots download champion Sion
  
  # Download data for a specific champion and don't remove noise
  pots download champion Ahri --dirty
  
  # Download data even when latest version is available on local
  pots download champion Jinx --force`,
}

func init() {
	downloadCmd.AddCommand(downloadChampionCmd)
	downloadChampionCmd.Flags().BoolVarP(&dirty, "dirty", "d", false, "leave noise in data for debugging")
	downloadChampionCmd.Flags().BoolVarP(&force, "force", "f", false, "ignores metadata matching version and downloads the champion anyway")
}
