/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Find differences on data between PBE and live",
	Long:  "Find and returns the differences on data between PBE and live patch",
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
