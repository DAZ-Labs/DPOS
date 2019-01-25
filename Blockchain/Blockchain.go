//inside blockchain package
//simple DPoS mechanism
package blockchain

import (
	"encoding/binary"
	"fmt"
	"time"

	"bytes"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/spf13/viper"
)

// structure for blockdata
type BlockData struct {
	Unit
	ParentHash   string
	Transactions Transactions
	Votes        Votes
	Reward       uint
}

//structure for block
type Block struct {
	BlockData
	Parent *Block
	Next   *Block
}

//transaction structure in different file
type Transactions []Transaction

//voting structure in different file
type Votes []Vote

//getting data
func (block *BlockData) GetData() *bytes.Buffer {
	// Gather data to Sign.
	data := block.Unit.GetData()
	data.Write([]byte(block.ParentHash))
	for _, t := range block.Transactions {
		data.Write(t.Sign)
	}
	for _, v := range block.Votes {
		data.Write(v.Sign)
	}
	data.Write(getBytes(int64(block.Reward)))

	return data
}

func getBytes(n int64) []byte {
	bytes := make([]byte, 16)
	binary.PutVarint(bytes, n)
	return bytes
}

func (block *BlockData) String() string {
	return fmt.Sprintf("%s", block.Hash)
}

//create a genesis block
func CreateGenesis() *Block {
	return &Block{Parent: nil, Next: nil, BlockData: BlockData{ParentHash: "", Unit: Unit{Hash: HashStr("genesis"),
		TimeStamp: time.Date(2019, 01, 25, 06, 00, 00, 00, time.UTC).UnixNano()}}}
}

//crating new block
func CreateNewBlock(privateKey crypto.PrivKey, parent *Block, transactions []Transaction, votes []Vote) *Block {
	block := &Block{
		BlockData{
			Transactions: transactions,
			Votes:        votes,
			ParentHash:   parent.Hash,
		},
		parent,
		nil,
	}
	block.Hash = Hash(block.GetData().Bytes())
	parent.Next = block
	block.TimeStamp = GetTimeStamp()
	block.Reward = uint(viper.GetInt("reward"))

	//TODO error handling
	id, _ := peer.IDFromPrivateKey(privateKey)
	block.Signer = id.Pretty()
	block.PublicKey, _ = privateKey.GetPublic().Bytes()
	block.Sign, _ = privateKey.Sign(block.GetData().Bytes())

	return block
}

//Block contents verification. Checks basic cryptography and transactions timing
func (block *BlockData) Verify(parent *BlockData) (valid bool, err error) {
	if block.ParentHash != parent.Hash {
		err = fmt.Errorf("block %s has %s ParentHash, %s required", block.Hash, block.ParentHash, parent.Hash)
		return
	}

	requiredReward := uint(viper.GetInt("reward"))
	if block.Reward != requiredReward {
		err = fmt.Errorf("block %s has incorrect reward = %d, required: %d", block.Hash, block.Reward, requiredReward)
		return
	}

	valid, err = verify(block)
	if !valid {
		err = fmt.Errorf("invalid block: %s: %s", block.Hash, err)
		return
	}

	lastTS := parent.GetTimestamp()
	for _, t := range block.Transactions {
		transaction := t
		//Transactions must be crated only within block production time frame
		if t.GetTimestamp() < lastTS || t.GetTimestamp() > block.TimeStamp {
			err = fmt.Errorf("block %s transaction: %s has wrong timestamp", block.Hash, &t)
			return
		}
		lastTS = t.GetTimestamp()

		valid, err = t.Verify()
		if !valid {
			err = fmt.Errorf("invalid transaction in block: %s sent %d to %s: %s", transaction.GetSigner(),
				transaction.Amount, transaction.Recipient, err)
			return
		}
	}

	for _, v := range block.Votes {
		valid, err = v.Verify()
		if !valid {
			err = fmt.Errorf("invalid vote in block: %s voted for %s: %s", v.GetSigner(), v.Candidate, err)
			return
		}
	}

	return
}
