package cmd

import (
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		//r := db.Tags(16760, 17761)
		//r := db.GetAlbumsById(393, 394)
		//for _, a := range r.Tags() {
		//  fmt.Printf("%+V\n", a.Name)
		//  for _, img := range a.Images() {
		//    println(img.Name)
		//  }
		//}

	},
}

func init() {
	rootCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
