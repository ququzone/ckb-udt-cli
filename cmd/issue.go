package cmd

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ququzone/ckb-sdk-go/crypto/secp256k1"
	"github.com/ququzone/ckb-sdk-go/rpc"
	"github.com/ququzone/ckb-sdk-go/transaction"
	"github.com/ququzone/ckb-sdk-go/types"
	"github.com/ququzone/ckb-sdk-go/utils"
	"github.com/ququzone/ckb-udt-cli/config"
	"github.com/spf13/cobra"
)

var (
	key    *string
	amount *string
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue sUDT token",
	Long:  `Issue sUDT with secp256k1 cell.`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rpc.Dial(config.C.RPC)
		if err != nil {
			Fatalf("create rpc client error: %v", err)
		}

		key, err := secp256k1.HexToKey(*key)
		if err != nil {
			Fatalf("import private key error: %v", err)
		}

		scripts, err := utils.NewSystemScripts(client)
		if err != nil {
			Fatalf("load system script error: %v", err)
		}

		change, err := key.Script(scripts)

		capacity := uint64(14200001000)
		fee := uint64(1000)

		cellCollector := utils.NewCellCollector(client, change, capacity+fee)
		cells, total, err := cellCollector.Collect()
		if err != nil {
			Fatalf("collect cell error: %v", err)
		}

		if total < capacity+fee {
			Fatalf("insufficient capacity: %d < %d", total, capacity+fee)
		}

		tx := transaction.NewSecp256k1SingleSigTx(scripts)
		for _, dep := range config.C.UDT.Deps {
			tx.CellDeps = append(tx.CellDeps, &types.CellDep{
				OutPoint: &types.OutPoint{
					TxHash: types.HexToHash(dep.TxHash),
					Index:  dep.Index,
				},
				DepType: types.DepType(dep.DepType),
			})
		}
		uuid, _ := change.Hash()

		tx.Outputs = append(tx.Outputs, &types.CellOutput{
			Capacity: uint64(capacity),
			Lock: &types.Script{
				CodeHash: change.CodeHash,
				HashType: change.HashType,
				Args:     change.Args,
			},
			Type: &types.Script{
				CodeHash: types.HexToHash(config.C.UDT.Script.CodeHash),
				HashType: types.ScriptHashType(config.C.UDT.Script.HashType),
				Args:     uuid.Bytes(),
			},
		})
		a, _ := big.NewInt(0).SetString(*amount, 10)
		b := a.Bytes()
		for i := 0; i < len(b)/2; i++ {
			b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
		}
		tx.OutputsData = append(tx.OutputsData, b)
		if total-capacity+fee > 6100000000 {
			tx.Outputs = append(tx.Outputs, &types.CellOutput{
				Capacity: total - capacity - fee,
				Lock: &types.Script{
					CodeHash: change.CodeHash,
					HashType: change.HashType,
					Args:     change.Args,
				},
			})
			tx.OutputsData = append(tx.OutputsData, []byte{})
		} else {
			tx.Outputs[1].Capacity = tx.Outputs[1].Capacity + total - capacity + fee
		}

		group, witnessArgs, err := transaction.AddInputsForTransaction(tx, cells)
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

		fmt.Printf("issue udt hash: %s, uuid: %s\n", hash.String(), uuid.String())
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)

	key = issueCmd.Flags().StringP("key", "k", "", "Issue private key")
	amount = issueCmd.Flags().StringP("amount", "a", "", "Issue amount")
	_ = issueCmd.MarkFlagRequired("key")
	_ = issueCmd.MarkFlagRequired("amount")
}
