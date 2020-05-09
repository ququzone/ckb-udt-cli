package cmd

import (
	"fmt"
	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
	"github.com/ququzone/ckb-udt-cli/config"
	"github.com/spf13/cobra"
	"math/big"
)

var (
	balanceConf *string
	balanceUUID *string
	balanceAddr *string
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Query sUDT balance",
	Long:  `Query sUDT balance by address.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Init(*balanceConf)
		if err != nil {
			Fatalf("load config error: %v", err)
		}

		client, err := rpc.Dial(c.RPC)
		if err != nil {
			Fatalf("create rpc client error: %v", err)
		}

		addr, err := address.Parse(*balanceAddr)
		if err != nil {
			Fatalf("parse address error: %v", err)
		}

		cellCollector := utils.NewCellCollector(client, addr.Script, NewUDTCellProcessor(client, nil))
		cellCollector.EmptyData = false
		cellCollector.TypeScript = &types.Script{
			CodeHash: types.HexToHash(c.UDT.Script.CodeHash),
			HashType: types.ScriptHashType(c.UDT.Script.HashType),
			Args:     types.HexToHash(*balanceUUID).Bytes(),
		}
		cells, err := cellCollector.Collect()
		if err != nil {
			Fatalf("collect cell error: %v", err)
		}
		total, ok := cells.Options["total"]
		if !ok {
			total = big.NewInt(0)
		}

		fmt.Printf("Address %s amount: %s\n", *balanceAddr, total.(*big.Int).String())
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)

	balanceConf = balanceCmd.Flags().StringP("config", "c", "config.yaml", "Config file")
	balanceUUID = balanceCmd.Flags().StringP("uuid", "u", "", "UDT uuid")
	balanceAddr = balanceCmd.Flags().StringP("address", "a", "", "Address")
	_ = balanceCmd.MarkFlagRequired("uuid")
	_ = balanceCmd.MarkFlagRequired("address")
}
