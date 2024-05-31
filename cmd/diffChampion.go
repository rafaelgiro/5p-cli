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

// diffChampionCmd represents the diffChampion command
var diffChampionCmd = &cobra.Command{
	Use:   "champion [name]",
	Short: "Find differences from PBE and Live for Champions",
	Long:  `Find differences from PBE and Live for Champions and prints them on the screen`,
	Run: func(cmd *cobra.Command, args []string) {
		champName := strings.ToLower(args[0])
		c := &champion.Champion{Name: champName}

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		dir := fmt.Sprintf("%s/data/champions", wd)
		if err := c.CheckDownload(dir); err != nil {
			log.Fatalf("Files not found for \"%s\". Please run the download champion command first: %v", c.Name, err)
		}

		if err := c.SaveDiff(dir, common.OutputFolder); err != nil {
			log.Fatalf("Failed to save diff file for %s on folder %s: %v", champName, common.OutputFolder, err)
		}

		fmt.Println("Success!")
	},
}

func init() {
	diffCmd.AddCommand(diffChampionCmd)
}
