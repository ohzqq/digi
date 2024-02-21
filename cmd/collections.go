package cmd

import (
	"github.com/ohzqq/digi/tui"
	"github.com/spf13/cobra"
)

// collectionsCmd represents the collections command
var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		//ids := tui.ListCollections()
		tui.Start()

		//fmt.Println(sel)
		//tui.ListAlbums(col)
		//child := col.Albums()
		//node := col.OpenNode(child[0])
		//fmt.Printf("%+v\n", col.OpenNode(node[3]))
		//other := col.GetAlbums(child.Nodes()...)
		//fmt.Printf("%+v\n", other.Nodes())
		//otherO := col.GetAlbums(other.Nodes()...)
		//al := db.GetAlbumsById(1, 4, 416, 349)
		//for _, a := range al.Names {
		//fmt.Printf("%v\n", a)
		//fmt.Printf("dir %v f %s\n is root %v\n first %v\n\n", d, f, root, first)
		//}
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
