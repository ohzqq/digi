package cmd

import (
	"fmt"

	"github.com/ohzqq/digi/db"
	"github.com/ohzqq/digi/tui"
	"github.com/spf13/cobra"
)

// collectionsCmd represents the collections command
var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		ids := tui.ListCollections()
		col := db.GetAlbumsByRoot(ids...)
		fmt.Printf("%+V\n", len(col.Albums))
	},
}

func init() {
	lsCmd.AddCommand(collectionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
