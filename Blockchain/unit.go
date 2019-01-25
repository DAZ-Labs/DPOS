//inside blockchain package
//simple DPoS mechanism
package blockchain

import (
	"bytes"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
)

//interface with different functions to be implemented
type Signable interface {
	GetSigner() string //getters are not in conventional form
	GetPublicKey() []byte
	GetData() *bytes.Buffer
	GetSign() []byte
	GetTimestamp() int64
}

//Basic Unit that can be presented in blockchain
//All fields should be public
//with interface signable as a part
type Unit struct {
	Signable
	Hash      string
	Signer    string
	Sign      []byte
	PublicKey []byte
	TimeStamp int64
}

//getting the block signer
func (u *Unit) GetSigner() string {
	return u.Signer
}

//getting public key
func (u *Unit) GetPublicKey() []byte {
	return u.PublicKey
}

//getting data
func (u *Unit) GetData() *bytes.Buffer {
	data := new(bytes.Buffer)
	data.Write([]byte(u.Signer))
	data.Write(getBytes(u.TimeStamp))
	return data
}

//getting sign
func (u *Unit) GetSign() []byte {
	return u.Sign
}

//getting timestamp
func (u *Unit) GetTimestamp() int64 {
	return u.TimeStamp
}

//verifying
func verify(s Signable) (result bool, err error) {
	result = false
	id, err := peer.IDB58Decode(s.GetSigner())
	if err != nil {
		return
	}
	public, err := crypto.UnmarshalPublicKey(s.GetPublicKey())
	if err != nil {
		return
	}

	if id.MatchesPublicKey(public) {
		result, err = public.Verify(s.GetData().Bytes(), s.GetSign())
	}
	return
}
