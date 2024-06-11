package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/multiformats/go-base32"
	"io/ioutil"
	"os"
	"path/filepath"
)
type KeyType string
type KeyInfo struct {
	Type       KeyType
	PrivateKey []byte
}
//name : KNamePrefix+Address.String()
func (kinfo *KeyInfo)Put(name string) error {
	encName := base32.RawStdEncoding.EncodeToString([]byte(name))
	path := os.Getenv(lotus_path)
	keyPath := filepath.Join(append([]string{path}, fsKeystore, encName)...)
	fmt.Println("keyPath:",keyPath)
	_, err := os.Stat(keyPath)
	if err == nil {
		return errors.New("key已经存在！ "+ name)
	} else if !os.IsNotExist(err) {
		return errors.New("checking key before put "+name+":"+err.Error())
	}

	keyData, err := json.Marshal(kinfo)
	if err != nil {
		return errors.New("encoding key "+name+":"+err.Error())
	}

	err = ioutil.WriteFile(keyPath, keyData, 0600)
	if err != nil {
		return errors.New("writing key "+name+":"+err.Error())
	}
	return nil
}

func (kinfo *KeyInfo)Get(name string) (KeyInfo, error) {

	encName := base32.RawStdEncoding.EncodeToString([]byte(name))
	path := os.Getenv(lotus_path)
	keyPath := filepath.Join(append([]string{path}, fsKeystore, encName)...)

	fstat, err := os.Stat(keyPath)
	if os.IsNotExist(err) {
		return KeyInfo{}, errors.New("key不存在！ "+ name)
	} else if err != nil {
		return KeyInfo{}, errors.New("checking key before put "+name+":"+err.Error())
	}

	if fstat.Mode()&0077 != 0 {
		return KeyInfo{},errors.New("kstrPermissionMsg" + name + fstat.Mode().String())
	}

	file, err := os.Open(keyPath)
	if err != nil {
		return KeyInfo{}, errors.New("opening key "+name+": "+err.Error())
	}
	defer file.Close() //nolint: errcheck // read only op

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return KeyInfo{}, errors.New("reading key "+name+": "+err.Error())
	}

	var res KeyInfo
	err = json.Unmarshal(data, &res)
	if err != nil {
		return KeyInfo{}, errors.New("decoding key "+name+": "+err.Error())
	}

	return res, nil
}

//func (kstore *KeyInfo)Import(keyinfo types.KeyInfo)error{
//	key,err := NewKey(keyinfo)
//	if err != nil {
//		return err
//	}
//	err = kstore.Put(KNamePrefix+key.Address.String(),key.KeyInfo)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
