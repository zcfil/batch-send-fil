package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filecoin-project/go-address"
	crypto2 "github.com/filecoin-project/go-crypto"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/lotus/lib/sigs"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"strings"
)
type Key struct {
	KeyInfo

	PublicKey []byte
	Address   address.Address
}

type bufs []byte
func (b bufs)Read(p []byte) (n int, err error){
	n = copy(p,b)
	return n, nil
}
const (
	iter       = 2048
	keyLen     = 64
	MnemonicCount = 12
	fsKeystore      = "keystore"
	lotus_path		= "LOTUS_PATH"
	KNamePrefix  	= "wallet-"
)
const (
	KTBLS             KeyType = "bls"
	KTSecp256k1       KeyType = "secp256k1"
	KTSecp256k1Ledger KeyType = "secp256k1-ledger"
)
var English []string

func init(){
	f, _ := ioutil.ReadFile("conf/english.txt")
	English = strings.Split(string(f),"\n")
}


func CreateKey(mnemonic string) (*Key, error) {
	pr,err := GenerateKeyFromSeed(mnemonic)
	if err!=nil{
		return nil,err
	}
	keyinfo := KeyInfo{
		Type:       KTSecp256k1,
		PrivateKey: pr,
	}
	key, err :=  NewKey(keyinfo)
	if err!=nil{
		return nil,err
	}
	return key,nil
}
func NewKey(keyinfo KeyInfo) (*Key, error) {
	k := &Key{
		KeyInfo: keyinfo,
	}

	var err error
	k.PublicKey, err = sigs.ToPublic(ActSigType(k.Type), k.PrivateKey)
	if err != nil {
		return nil, err
	}
	fmt.Println("云构:",keyinfo.Type)
	fmt.Println("云构:",len(k.PublicKey))
	switch k.Type {
	case KTSecp256k1:
		k.Address, err = address.NewSecp256k1Address(k.PublicKey)
		if err != nil {
			return nil, errors.New("converting Secp256k1 to address: "+err.Error())
		}
	case KTBLS:
		//k.Address, err = address.NewBLSAddress(k.PublicKey)
		//if err != nil {
		//	return nil, errors.New("converting BLS to address: "+err.Error())
		//}
	default:
		return nil, errors.New("unsupported key type: "+ string(k.Type))
	}
	return k, nil

}

func ActSigType(typ KeyType) crypto.SigType {
	switch typ {
	case KTBLS:
		return crypto.SigTypeBLS
	case KTSecp256k1:
		return crypto.SigTypeSecp256k1
	default:
		return crypto.SigTypeUnknown
	}
}

func Str2DEC(s string) (num int) {
	l := len(s)
	for i := l - 1; i >= 0; i-- {
		num += (int(s[l-i-1]) & 0xf) << uint8(i)
	}
	return
}
func CreateMnemonic()string{
	mnemonic := ""
	p,_ := rand.Prime(rand.Reader,128)
	hash := sha256.New()
	ps := p.Bytes()

	//做256哈希取前4位
	hash.Write(ps)
	bytes := hash.Sum(nil)

	b1 := bytes[0]>>4

	str := ""
	for _,v := range ps{
		str += fmt.Sprintf("%08b",v)
	}
	str += fmt.Sprintf("%04b",b1)
	//index := make([]int,MnemonicCount)
	//获取12位助记词
	for i:=0;i< MnemonicCount;i++{
		//将str  132位12等分 成为12个助记词下标
		mnemonic += English[Str2DEC(str[i*11:11*(i+1)])]
		if i< MnemonicCount-1{
			mnemonic += " "
		}
	}
	return mnemonic
}
func GenerateKeyFromSeed(mnemonic string) ([]byte, error) {
	p := pbkdf2.Key([]byte(mnemonic),nil, iter, keyLen,sha512.New)

	seed := make(bufs,len(p)/2)
	copy(seed, p[:len(p)/2])

	return crypto2.GenerateKeyFromSeed(seed)
}


func ExportKey(name string) ([]byte, error) {
	var keyInfo KeyInfo
	keyinfo,err := keyInfo.Get(KNamePrefix +name)
	if err!=nil{
		return nil,err
	}
	keystr,_ := json.Marshal(keyinfo)
	return keystr, err
}
