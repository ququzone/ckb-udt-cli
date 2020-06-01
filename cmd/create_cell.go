package cmd

import (
	"context"
	"fmt"

	"github.com/ququzone/ckb-sdk-go/address"
	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
	"github.com/ququzone/ckb-udt-cli/config"
	"github.com/spf13/cobra"
)

var (
	createCellConf *string
	createCellKey  *string
	createCellUUID *string
)

var createCellCmd = &cobra.Command{
	Use:   "create-cell",
	Short: "create anyone can pay cell for sUDT token",
	Long:  `create anyone can pay cell for sUDT token.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.Init(*createCellConf)
		if err != nil {
			Fatalf("load config error: %v", err)
		}

		client, err := rpc.Dial(c.RPC)
		if err != nil {
			Fatalf("create rpc client error: %v", err)
		}

		key, err := secp256k1.HexToKey(*createCellKey)
		if err != nil {
			Fatalf("import private key error: %v", err)
		}

		scripts, err := utils.NewSystemScripts(client)
		if err != nil {
			Fatalf("load system script error: %v", err)
		}

		change, err := key.Script(scripts)
		capacity := uint64(14200000000)
		fee := uint64(1000)

		cellCollector := utils.NewCellCollector(client, change, utils.NewCapacityCellProcessor(capacity+fee))
		cells, err := cellCollector.Collect()
		if err != nil {
			Fatalf("collect cell error: %v", err)
		}

		if cells.Capacity < capacity+fee {
			Fatalf("insufficient capacity: %d < %d", cells.Capacity, capacity+fee)
		}

		tx := transaction.NewSecp256k1SingleSigTx(scripts)
		for _, dep := range c.UDT.Deps {
			tx.CellDeps = append(tx.CellDeps, &types.CellDep{
				OutPoint: &types.OutPoint{
					TxHash: types.HexToHash(dep.TxHash),
					Index:  dep.Index,
				},
				DepType: types.DepType(dep.DepType),
			})
		}

		// cell
		lock := &types.Script{
			CodeHash: types.HexToHash(c.ACP.Script.CodeHash),
			HashType: types.ScriptHashType(c.ACP.Script.HashType),
			Args:     change.Args,
		}
		tx.Outputs = append(tx.Outputs, &types.CellOutput{
			Capacity: uint64(capacity),
			Lock:     lock,
			Type: &types.Script{
				CodeHash: types.HexToHash(c.UDT.Script.CodeHash),
				HashType: types.ScriptHashType(c.UDT.Script.HashType),
				Args:     types.HexToHash(*createCellUUID).Bytes(),
			},
		})
		tx.OutputsData = append(tx.OutputsData, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,})

		if cells.Capacity-capacity+fee > 6100000000 {
			tx.Outputs = append(tx.Outputs, &types.CellOutput{
				Capacity: cells.Capacity - capacity - fee,
				Lock:     change,
			})
			tx.OutputsData = append(tx.OutputsData, []byte{})
		} else {
			tx.Outputs[1].Capacity = tx.Outputs[1].Capacity + cells.Capacity - capacity + fee
		}

		group, witnessArgs, err := transaction.AddInputsForTransaction(tx, cells.Cells)
		if err != nil {
			Fatalf("add inputs to transaction error: %v", err)
		}

		err = transaction.SingleSignTransaction(tx, group, witnessArgs, key)
		if err != nil {
			Fatalf("sign transaction error: %v", err)
		}

		hash, err := client.SendTransaction(context.Background(), tx)
		if err != nil {
			Fatalf("send transaction error: %v", err)
		}
		addr, _ := address.Generate(address.Testnet, lock)

		fmt.Printf("create anyone can pay cell transaction hash: %s, address: %s\n", hash.String(), addr)
	},
}

func init() {
	rootCmd.AddCommand(createCellCmd)

	createCellConf = createCellCmd.Flags().StringP("config", "c", "config.yaml", "Config file")
	createCellKey = createCellCmd.Flags().StringP("key", "k", "", "Private key")
	createCellUUID = createCellCmd.Flags().StringP("uuid", "u", "", "UDT uuid")
	_ = createCellCmd.MarkFlagRequired("key")
}
