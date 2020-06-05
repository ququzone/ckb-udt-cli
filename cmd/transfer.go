package cmd

import (
	"context"
	"fmt"
	"math/big"

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
		c, err := config.Init(*transferConf)
		if err != nil {
			Fatalf("load config error: %v", err)
		}

		client, err := rpc.Dial(c.RPC)
		if err != nil {
			Fatalf("create rpc client error: %v", err)
		}

		key, err := secp256k1.HexToKey(*transferKey)
		if err != nil {
			Fatalf("import private key error: %v", err)
		}

		amount := big.NewInt(0)
		amount, _ = amount.SetString(*transferAmount, 10)
		if amount == nil || amount.Uint64() == 0 {
			Fatalf("transfer amount error: %s", *transferAmount)
		}

		uuid := types.HexToHash(*transferUUID).Bytes()

		capacity := uint64(28400000000)
		fee := uint64(1000)
		recipientAcp := false
		recipientAddr, err := address.Parse(*transferTo)
		var recipientCell *types.Cell
		if err != nil {
			Fatalf("parse to address error: %v", err)
		}
		if recipientAddr.Script.CodeHash.String() == c.ACP.Script.CodeHash {
			capacity -= 14200000000
			recipientAcp = true
			cells, err := CollectUDT(client, c, recipientAddr.Script, uuid, big.NewInt(0))
			if err != nil {
				Fatalf("collect cell error: %v", err)
			}
			if len(cells.Cells) == 0 {
				Fatalf("can't find anyone can pay cell for %s", *transferTo)
			}
			recipientCell = cells.Cells[0]
		}

		scripts, err := utils.NewSystemScripts(client)
		if err != nil {
			Fatalf("load system script error: %v", err)
		}

		fromSecp256k1Script, err := key.Script(scripts)
		fromScript := &types.Script{
			CodeHash: types.HexToHash(c.ACP.Script.CodeHash),
			HashType: types.ScriptHashType(c.ACP.Script.HashType),
			Args:     fromSecp256k1Script.Args,
		}
		fromAcp := true

		cells, err := CollectUDT(client, c, fromScript, uuid, amount)
		if err != nil {
			Fatalf("collect cell error: %v", err)
		}

		if cells.Options["total"].(*big.Int).Cmp(amount) < 0 {
			fromAcp = false
			fromScript = fromSecp256k1Script
			cells, err = CollectUDT(client, c, fromScript, uuid, amount)
			if err != nil {
				Fatalf("collect cell error: %v", err)
			}
			if cells.Options["total"].(*big.Int).Cmp(amount) < 0 {
				Fatalf("insufficient UDT balance")
			}
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
		if fromAcp || recipientAcp {
			for _, dep := range c.ACP.Deps {
				tx.CellDeps = append(tx.CellDeps, &types.CellDep{
					OutPoint: &types.OutPoint{
						TxHash: types.HexToHash(dep.TxHash),
						Index:  dep.Index,
					},
					DepType: types.DepType(dep.DepType),
				})
			}
		}

		var feeCells *utils.CollectResult
		if cells.Capacity < capacity+fee {
			cellCollector := utils.NewCellCollector(client, fromScript, utils.NewCapacityCellProcessor(capacity+fee-cells.Capacity))
			feeCells, err = cellCollector.Collect()
			if err != nil {
				Fatalf("collect cell error: %v", err)
			}

			if feeCells.Capacity < capacity+fee-cells.Capacity {
				Fatalf("insufficient capacity: %d < %d", cells.Capacity, capacity+fee)
			}
		}

		if recipientAcp {
			input := &types.CellInput{
				Since: 0,
				PreviousOutput: &types.OutPoint{
					TxHash: recipientCell.OutPoint.TxHash,
					Index:  recipientCell.OutPoint.Index,
				},
			}
			tx.Inputs = append(tx.Inputs, input)
			tx.Witnesses = append(tx.Witnesses, []byte{})

			tx.Outputs = append(tx.Outputs, &types.CellOutput{
				Capacity: recipientCell.Capacity,
				Lock:     recipientCell.Lock,
				Type:     recipientCell.Type,
			})
			b := amount.Bytes()
			for i := 0; i < len(b)/2; i++ {
				b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
			}
			if len(b) < 16 {
				for i := len(b); i < 16; i++ {
					b = append(b, 0)
				}
			}
			tx.OutputsData = append(tx.OutputsData, b)
		} else {
			tx.Outputs = append(tx.Outputs, &types.CellOutput{
				Capacity: uint64(14200000000),
				Lock:     recipientAddr.Script,
				Type: &types.Script{
					CodeHash: types.HexToHash(c.UDT.Script.CodeHash),
					HashType: types.ScriptHashType(c.UDT.Script.HashType),
					Args:     uuid,
				},
			})
			b := amount.Bytes()
			for i := 0; i < len(b)/2; i++ {
				b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
			}
			if len(b) < 16 {
				for i := len(b); i < 16; i++ {
					b = append(b, 0)
				}
			}
			tx.OutputsData = append(tx.OutputsData, b)
		}

		changeCapacity := cells.Capacity - fee
		if feeCells != nil {
			changeCapacity += capacity + fee - cells.Capacity
		}
		if !recipientAcp {
			changeCapacity -= 14200000000
		}

		tx.Outputs = append(tx.Outputs, &types.CellOutput{
			Capacity: changeCapacity,
			Lock:     fromScript,
			Type: &types.Script{
				CodeHash: types.HexToHash(c.UDT.Script.CodeHash),
				HashType: types.ScriptHashType(c.UDT.Script.HashType),
				Args:     uuid,
			},
		})
		if cells.Options["total"].(*big.Int).Cmp(amount) == 0 {
			tx.OutputsData = append(tx.OutputsData, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		} else {
			b := big.NewInt(0).Sub(cells.Options["total"].(*big.Int), amount).Bytes()
			for i := 0; i < len(b)/2; i++ {
				b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
			}
			if len(b) < 16 {
				for i := len(b); i < 16; i++ {
					b = append(b, 0)
				}
			}
			tx.OutputsData = append(tx.OutputsData, b)
		}

		if feeCells != nil && feeCells.Capacity-capacity+fee-cells.Capacity > 6100000000 {
			tx.Outputs = append(tx.Outputs, &types.CellOutput{
				Capacity: feeCells.Capacity - capacity + fee - cells.Capacity,
				Lock:     fromScript,
			})
			tx.OutputsData = append(tx.OutputsData, []byte{})
		}

		var inputs []*types.Cell
		inputs = append(inputs, cells.Cells...)
		if feeCells != nil {
			inputs = append(inputs, feeCells.Cells...)
		}

		group, witnessArgs, err := transaction.AddInputsForTransaction(tx, inputs)
		if err != nil {
			Fatalf("add inputs to transaction error: %v", err)
		}

		err = transaction.SingleSignTransaction(tx, group, witnessArgs, key)
		if err != nil {
			Fatalf("sign transaction error: %v", err)
		}

		fmt.Println(rpc.TransactionString(tx))

		hash, err := client.SendTransaction(context.Background(), tx)
		if err != nil {
			Fatalf("send transaction error: %v", err)
		}

		fmt.Printf("transfer transaction hash: %s\n", hash.String())
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)

	transferConf = transferCmd.Flags().StringP("config", "c", "config.yaml", "Config file")
	transferKey = transferCmd.Flags().StringP("key", "k", "", "From private key")
	transferUUID = transferCmd.Flags().StringP("uuid", "u", "", "UDT uuid")
	transferAmount = transferCmd.Flags().StringP("amount", "a", "", "Transfer amount")
	transferTo = transferCmd.Flags().StringP("to", "t", "", "Transfer recipient address")
	_ = transferCmd.MarkFlagRequired("key")
	_ = transferCmd.MarkFlagRequired("amount")
	_ = transferCmd.MarkFlagRequired("uuid")
	_ = transferCmd.MarkFlagRequired("to")
}
