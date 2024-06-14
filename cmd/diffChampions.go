/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/5pots-com/cli/internal/champion"
	"github.com/5pots-com/cli/internal/common"
	"github.com/spf13/cobra"
)

// diffChampionCmd represents the diffChampion command
var diffChampionsCmd = &cobra.Command{
	Use:   "champions",
	Short: "Find differences from PBE and Live for Champions",
	Long:  `Find differences from PBE and Live for Champions and prints them on the screen`,
	Run: func(cmd *cobra.Command, args []string) {

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		dir := fmt.Sprintf("%s/data/champions", wd)

		if err := common.CheckDownload([]string{dir}); err != nil {
			log.Fatalf("Champions folder not found. Please run the download champions command first: %v", err)
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			log.Fatalf("Failed to read champions folder: %v", err)
		}

		for _, e := range entries {
			c := &champion.Champion{Name: e.Name()}

			res, err := c.PrepareDiff(dir, common.OutputFolder)
			if err != nil {
				log.Fatalf("Failed to prepare diff file for %s on folder %s: %v", e.Name(), common.OutputFolder, err)
			}

			if len(res.Keys) != 0 {
				ch, err := json.Marshal(res)
				if err != nil {
					log.Fatalf("failed to convert %s champion json: %v", c.Name, err)
				}

				d := fmt.Sprintf("%s/results", wd)
				common.SaveToFile(d, fmt.Sprintf("%s.json", e.Name()), ch)
			}

		}

		// if err := c.SaveDiff(dir, common.OutputFolder); err != nil {
		// 	log.Fatalf("Failed to save diff file for %s on folder %s: %v", champName, common.OutputFolder, err)
		// }

		fmt.Println("Success!")
	},
}

func init() {
	diffCmd.AddCommand(diffChampionsCmd)
}
