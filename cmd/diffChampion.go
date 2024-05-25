/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/5pots-com/cli/internal/champion"
	"github.com/5pots-com/cli/internal/common"
	"github.com/spf13/cobra"
	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

var simple bool

// diffChampionCmd represents the diffChampion command
var diffChampionCmd = &cobra.Command{
	Use:   "champion [name]",
	Short: "Find differences from PBE and Live for Champions",
	Long:  `Find differences from PBE and Live for Champions and prints them on the screen`,
	Run: func(cmd *cobra.Command, args []string) {
		champName := strings.ToLower(args[0])

		c := &champion.Champion{Name: champName}

		fmt.Printf("Downloading %s data on patch %s...\n", c.Name, common.PBE)
		pbe, err := c.Download(common.PBE, !dirty)
		if err != nil {
			log.Fatalf("Failed to Download PBE data: %v", err)
		}

		fmt.Printf("Downloading %s data on patch %s...\n", c.Name, common.Latest)
		live, err := c.Download(common.Latest, !dirty)
		if err != nil {
			log.Fatalf("Failed to Download Live data: %v", err)
		}

		fmt.Printf("Finding differences...\n")

		pbe = champion.RemoveNoise(pbe)
		live = champion.RemoveNoise(live)

		fpbe, err := common.Format(pbe)
		flive, err := common.Format(live)

		jd := diff.New()

		diffs, err := jd.Compare(flive, fpbe)

		formatter := formatter.NewDeltaFormatter()
		diffString, err := formatter.Format(diffs)

		if simple {
			fmt.Println(diffs.Deltas())
		} else {

			fmt.Println(diffString)
		}

	},
}

func init() {
	diffCmd.AddCommand(diffChampionCmd)
	diffChampionCmd.Flags().BoolVarP(&simple, "simple", "s", false, "automatically cleans noise from data")
}
