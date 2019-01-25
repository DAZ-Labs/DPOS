//inside blockchain package
//simple DPoS mechanism
package blockchain

import (
	"bytes"
	"encoding/binary"
	"fmt"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

//transaction structure
type Transaction struct {
	Unit
	Recipient string
	Amount    uint64
}

//getting transaction data
func (t *Transaction) GetData() *bytes.Buffer {
	// Gather data to Sign.
	data := t.Unit.GetData()
	data.Write([]byte(t.Recipient))
	amountBytes := make([]byte, 16)
	binary.PutUvarint(amountBytes, t.Amount)
	data.Write(amountBytes)
	return data

}

func (t *Transaction) String() string {
	return fmt.Sprintf("%d from %s to %s", t.Amount, t.Signer, t.Recipient)
}

//verifying transaction
func (t *Transaction) Verify() (result bool, err error) {
	if t.Signer == t.Recipient {
		err = fmt.Errorf("self payment")
		return
	}
	return verify(t)
}

//paying to recpient
func Pay(private crypto.PrivKey, recipient peer.ID, amount uint64) *Transaction {

	sender, _ := peer.IDFromPrivateKey(private)
	public, _ := private.GetPublic().Bytes()

	t := Transaction{Unit: Unit{Signer: sender.Pretty(), PublicKey: public, TimeStamp: GetTimeStamp()}, Recipient: recipient.Pretty(), Amount: amount}
	sign, _ := private.Sign(t.GetData().Bytes())
	t.Sign = sign

	return &t
}
