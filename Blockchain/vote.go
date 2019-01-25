//inside blockchain package
//simple DPoS mechanism
package blockchain

import (
	"fmt"

	"bytes"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

//structure for vote
type Vote struct {
	Unit
	Candidate string
}

// getting data of voter
func (v *Vote) GetData() *bytes.Buffer {
	data := v.Unit.GetData()
	data.Write([]byte(v.Candidate))
	return data
}

//verifying voter
func (v *Vote) Verify() (result bool, err error) {
	if v.Signer == v.Candidate {
		err = fmt.Errorf("self voting")
		return
	}
	return verify(v)
}

func (v *Vote) String() string {
	return fmt.Sprintf("%s voted for %s", v.Signer, v.Candidate)
}

//new vote
func NewVote(private crypto.PrivKey, candidate peer.ID) *Vote {

	sender, _ := peer.IDFromPrivateKey(private)
	public, _ := private.GetPublic().Bytes()

	v := Vote{Unit: Unit{Signer: sender.Pretty(), PublicKey: public, TimeStamp: GetTimeStamp()}, Candidate: candidate.Pretty()}
	sign, _ := private.Sign(v.GetData().Bytes())
	v.Sign = sign

	return &v
}
