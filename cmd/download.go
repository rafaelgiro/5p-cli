/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

type Patch string

const (
	PBE    Patch = "pbe"
	Latest Patch = "latest"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads current live and PBE data.",
	Long:  `Downloads the latest data from live and PBE versions into two files.`,
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
