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
	pbeFileName    = "pbe.json"
	latestFileName = "live.json"
)

var clean bool

var downloadChampionCmd = &cobra.Command{
	Use:   "champion [name]",
	Short: "Download data for a specific champion",
	Long: `Download the latest data for a specific champion from the PBE and Live server.

This command fetches the latest champion data and optionally removes noise from the data.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		champName := strings.ToLower(args[0])

		c := &champion.Champion{Name: champName}

		pbe, err := c.Download(common.PBE)
		if err != nil {
			log.Fatalf("Failed to Download PBE data: %v", err)
		}

		live, err := c.Download(common.Latest)
		if err != nil {
			log.Fatalf("Failed to Download Live data: %v", err)
		}

		if clean {
			pbe = champion.RemoveNoise(pbe)
			live = champion.RemoveNoise(live)
		}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		dir := fmt.Sprintf("%s/data", wd)

		if err := c.SaveToFile(dir, latestFileName, live); err != nil {
			log.Fatalf("Failed to save file %s: %v", latestFileName, err)
		}

		if err := c.SaveToFile(dir, pbeFileName, pbe); err != nil {
			log.Fatalf("Failed to save file %s: %v", pbeFileName, err)
		}

		fmt.Printf("Success!\n")
	},
	Example: `  # Download data for a specific champion
  pots download champion Sion
  
  # Download data for a specific champion and remove noise
  pots download champion Ahri --clean
  
  # Download data for a specific champion using shorthand flag for clean
  pots download champion Jinx -c`,
}

func init() {
	downloadCmd.AddCommand(downloadChampionCmd)
	downloadChampionCmd.Flags().BoolVarP(&clean, "clean", "c", false, "automatically cleans noise from data")
}
