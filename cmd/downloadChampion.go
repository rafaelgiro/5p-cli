/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/5pots-com/cli/internal/champion"
	"github.com/5pots-com/cli/internal/common"
	"github.com/spf13/cobra"
)

const (
	pbeFileName      = "pbe.json"
	latestFileName   = "live.json"
	metadataFileName = "metadata.json"
)

var dirty bool
var force bool

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

		dir := fmt.Sprintf("%s/data/champions", wd)

		hasNewLatest, err := common.HasNewVersion(common.Latest, fmt.Sprintf("%s/%s", dir, c.Name))
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		hasNewPBE, err := common.HasNewVersion(common.PBE, fmt.Sprintf("%s/%s", dir, c.Name))
		if err != nil {
			log.Fatalf("Failed to read local data: %v", err)
		}

		if hasNewLatest {
			downSaveLatest(c, dir)
		} else {
			fmt.Printf("Live version already downloaded for %s\n", c.Name)
		}
		if hasNewPBE {
			downSavePBE(c, dir)
		} else {
			fmt.Printf("PBE version already downloaded for %s\n", c.Name)
		}
		if hasNewLatest || hasNewPBE {
			downSaveMetaData(c, dir)
		}

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

func downSavePBE(c *champion.Champion, dir string) {
	fmt.Printf("Downloading %s data on patch \"%s\"...\n", c.Name, common.PBE)
	pbe, err := c.Download(common.PBE, !dirty)
	if err != nil {
		log.Fatalf("Failed to Download PBE data: %v", err)
	}

	if !dirty {
		pbe = champion.RemoveNoise(pbe)
	}

	if err := c.SaveToFile(dir, pbeFileName, pbe); err != nil {
		log.Fatalf("Failed to save file %s: %v", pbeFileName, err)
	}
}

func downSaveLatest(c *champion.Champion, dir string) {

	fmt.Printf("Downloading %s data on patch \"%s\"...\n", c.Name, common.Latest)
	live, err := c.Download(common.Latest, !dirty)
	if err != nil {
		log.Fatalf("Failed to Download Live data: %v", err)
	}

	if !dirty {
		live = champion.RemoveNoise(live)
	}

	if err := c.SaveToFile(dir, latestFileName, live); err != nil {
		log.Fatalf("Failed to save file %s: %v", latestFileName, err)
	}
}

func downSaveMetaData(c *champion.Champion, dir string) {
	m, err := common.DownloadMetadata()
	if err != nil {
		log.Fatalf("Failed to download metadata: %v", err)
	}

	if err := c.SaveToFile(dir, metadataFileName, m); err != nil {
		log.Fatalf("Failed to save file %s: %v", pbeFileName, err)
	}
}
