package blockchain

import (
	"fmt"
	"bytes"
  "crypto/sha512"
  "encoding/gob"
	"encoding/hex"
)

type Transaction struct {
	ID []byte // hash
	Inputs []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value int // the worthy tokens
	PubKey string // value needed to unlock the tocken
}

type TxInput struct {
	ID []byte
	Out int
	Sig string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [64]byte
  encoder := gob.NewEncoder(&encoded)
  err := encoder.Encode(tx)
  if err != nil { fmt.Println(err) }
	hash = sha512.Sum512(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {data = fmt.Sprintf("Coins to %s", to)}
	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{100, to}
	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) IsCoinbase() bool { return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1 }
func (in *TxInput) CanUnlock(data string) bool { return in.Sig == data }
func (out *TxOutput) CanBeUnlocked(data string) bool { return out.PubKey == data }

func (b *block) HashTransactions() []byte {
  var txHashes [][]byte 
  var txHash [64]byte
  for _, tx := range b.Transactions { txHashes = append(txHashes, tx.ID) }
  txHash = sha512.Sum512(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

func (chain *BlockChain) FindUnspendtTxns(address string) []Transaction {
	var unspentTxns []Transaction
	spentTXOs := make(map[string][]int)
	iter := iterator(chain, "/Players")
	for {
		block := deeper(iter, false)
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, out := range tx.Outputs { // for each output inside each TX 
				if spentTXOs[txID] != nil {         // that is already in buffer 
					for _, spentOut := range spentTXOs[txID] { // if key of output == one of values in bufer 
						if spentOut == outIdx {                  // goto VVV
							continue Outputs
						}
					}
				}
				// goto here
				if out.CanBeUnlocked(address) { unspentTxns = append(unspentTxns, *tx) }
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(tx.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}
		if len(block.Prev) == 0 { break }
		return unspentTxns
	} 
	return unspentTxns
}

func (chain *BlockChain) FindUTxO(address string) []TxOutput {
	var UTxOs []TxOutput
	unspentTxns := chain.FindUnspendtTxns(address)
	for _, tx := range unspentTxns {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTxOs = append(UTxOs, out)
			}
		}
	}
	return UTxOs
} 

func (chain *BlockChain) FindSpendableOuts(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspendtTxns(address)
	accumulated := 0
Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				if accumulated >= amount { break Work }
			}
		}
	}
	return accumulated, unspentOuts
}

func NewTX(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	acc, validOuts := chain.FindSpendableOuts(from, amount)
	if acc < amount { fmt.Println("\u001b[7m\u001b[38;5;124m NOT ENOUGH FUNDS!!! \u001b[0m") }
	for txid, outs := range validOuts {
		txID, err := hex.DecodeString(txid)
		if err != nil { fmt.Println(err) }
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount { outputs = append(outputs, TxOutput{acc - amount, from}) }
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}

func GetBalance(chain *BlockChain, address string) int {
	balance := 0
	UTxOs := chain.FindUTxO(address)
	for _, out := range UTxOs { balance += out.Value }
	return balance
}