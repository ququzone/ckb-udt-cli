package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	transferConf   *string
	transferKey    *string
	transferAmount *string
	transferTo     *string
	transferUUID   *string
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

	transferConf = transferCmd.Flags().StringP("config", "c", "", "Config file")
	transferKey = transferCmd.Flags().StringP("key", "k", "", "From private key")
	transferUUID = transferCmd.Flags().StringP("uuid", "u", "", "UDT uuid")
	transferAmount = transferCmd.Flags().StringP("amount", "a", "", "Transfer amount")
	transferTo = transferCmd.Flags().StringP("to", "t", "", "Transfer recipient address")
	_ = transferCmd.MarkFlagRequired("key")
	_ = transferCmd.MarkFlagRequired("amount")
	_ = transferCmd.MarkFlagRequired("uuid")
	_ = transferCmd.MarkFlagRequired("to")
}
