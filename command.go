package compare_msg

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: `Work with messages files`,
}

var compare = &cobra.Command{
	Use:   "compare",
	Short: "Compare two message files",
	Run: func(cmd *cobra.Command, args []string) {
		baseFile, _ := cmd.Flags().GetString("file-a")
		messageFile, _ := cmd.Flags().GetString("file-b")
		if err := Compare(baseFile, messageFile); err != nil {
			fmt.Println(err)
		}
	},
}

var update = &cobra.Command{
	Use:   "update",
	Short: "Update message file from translation",
	Run: func(cmd *cobra.Command, args []string) {
		baseFile, _ := cmd.Flags().GetString("file-a")
		translationFile, _ := cmd.Flags().GetString("translation-file")
		outFile, _ := cmd.Flags().GetString("out-file")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		force, _ := cmd.Flags().GetBool("force")
		if err := Update(baseFile, translationFile, outFile, overwrite, force); err != nil {
			fmt.Println(err)
		}
	},
}

var compareUpdate = &cobra.Command{
	Use:   "compare-update",
	Short: "Compare with second and update missing from translation",
	Run: func(cmd *cobra.Command, args []string) {
		baseFile, _ := cmd.Flags().GetString("file-a")
		messageFile, _ := cmd.Flags().GetString("file-b")
		translationFile, _ := cmd.Flags().GetString("translation-file")
		outFile, _ := cmd.Flags().GetString("out-file")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		force, _ := cmd.Flags().GetBool("force")
		if err := CompareUpdate(baseFile, messageFile, translationFile, outFile, overwrite, force); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(compare)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(compareUpdate)
	compare.Flags().String("file-a", "", "Message file")
	compare.Flags().String("file-b", "", "Base file")
	update.Flags().String("file-a", "", "Message file")
	update.Flags().String("translation-file", "", "Translation file")
	update.Flags().String("out-file", "", "Output file")
	update.Flags().Bool("force", false, "Force update translation even if already exits")
	update.Flags().Bool("overwrite", false, "Overwrite out file if exits")
	compareUpdate.Flags().String("file-a", "", "Message file")
	compareUpdate.Flags().String("file-b", "", "Base file")
	compareUpdate.Flags().String("translation-file", "", "Translation file")
	compareUpdate.Flags().String("out-file", "", "Output file")
	compareUpdate.Flags().Bool("overwrite", false, "Overwrite out file if exits")
	compareUpdate.Flags().Bool("force", false, "Force update translation even if already exits")
	compare.MarkFlagRequired("file-a")
	compare.MarkFlagRequired("file-b")
	update.MarkFlagRequired("file-a")
	update.MarkFlagRequired("translation-file")
	compareUpdate.MarkFlagRequired("file-a")
	compareUpdate.MarkFlagRequired("file-b")
	compareUpdate.MarkFlagRequired("translation-file")
}

func Execute() error {
	return rootCmd.Execute()
}
