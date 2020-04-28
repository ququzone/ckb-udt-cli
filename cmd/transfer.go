package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	to *string
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer sUDT token",
	Long:  `Transfer sUDT from secp256k1 lock cell.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("transfer cmd unimplemented")
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	key = transferCmd.Flags().StringP("key", "k", "", "From private key")
	amount = transferCmd.Flags().StringP("amount", "a", "", "Transfer amount")
	to = transferCmd.Flags().StringP("to", "t", "", "Transfer recipient address")
	_ = transferCmd.MarkFlagRequired("key")
	_ = transferCmd.MarkFlagRequired("amount")
}
