package cmd

import (
	"github.com/spf13/cobra"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue sUDT token",
	Long: `Issue sUDT with secp256k1 cell.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(issueCmd)

	issueCmd.Flags().StringP("key", "k", "", "Issue private key")
	issueCmd.Flags().StringP("amount", "a", "", "Issue amount")
}
