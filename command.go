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
		baseFile, _ := cmd.Flags().GetString("base-file")
		messageFile, _ := cmd.Flags().GetString("message-file")
		translationFile := ""
		outFile := ""
		overwrite := false
		if err := CompareUpdate(baseFile, messageFile, translationFile, outFile, overwrite); err != nil {
			fmt.Println(err)
		}
	},
}

var update = &cobra.Command{
	Use:   "update",
	Short: "Update message file from translation",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var compareUpdate = &cobra.Command{
	Use:   "compare-update",
	Short: "Compare with second and update missing from translation",
	Run: func(cmd *cobra.Command, args []string) {
		baseFile, _ := cmd.Flags().GetString("base-file")
		messageFile, _ := cmd.Flags().GetString("message-file")
		translationFile, _ := cmd.Flags().GetString("translation-file")
		outFile, _ := cmd.Flags().GetString("out-file")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		if err := CompareUpdate(baseFile, messageFile, translationFile, outFile, overwrite); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(compare)
	rootCmd.AddCommand(update)
	rootCmd.AddCommand(compareUpdate)
	rootCmd.PersistentFlags().String("message-file", "", "Message file")
	compare.Flags().String("base-file", "", "Base file")
	update.Flags().String("translation-file", "", "Translation file")
	compareUpdate.Flags().String("translation-file", "", "Translation file")
	compareUpdate.Flags().String("base-file", "", "Base file")
	compareUpdate.Flags().String("out-file", "", "Output file")
	compareUpdate.Flags().Bool("overwrite", false, "Overwrite out file")
	rootCmd.MarkFlagRequired("message-file")
	compare.MarkFlagRequired("base-file")
	update.MarkFlagRequired("translation-file")
	compareUpdate.MarkFlagRequired("base-file")
	compareUpdate.MarkFlagRequired("translation-file")
}

func Execute() error {
	return rootCmd.Execute()
}
