package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type cliAccountKeys struct {
	PubKey  string `json:"pubkey"`
	PrivKey string `json:"privkey"`
}

type cryptoRecord struct {
	Cipher     string
	CipherText string
}

type smeshAccountKeystruct struct {
	Crypto cryptoRecord
}

type account struct {
	DisplayName string
	PublicKey   string
	SecretKey   string
}

type secretStuff struct {
	Accounts []account
}

func decrypt(block cipher.Block, ciphertext []byte, iv []byte) []byte {
	stream := cipher.NewCTR(block, iv)
	plain := make([]byte, len(ciphertext))
	stream.XORKeyStream(plain, ciphertext)
	return plain
}

func main() {
	fmt.Println("SMESH 2 CLI")
	fmt.Println()
	fmt.Println("Convert SMESH APP keystores to json files for CLIWallet")
	fmt.Println()
	fmt.Println("Warning : CLIWallet files are not currently encrypted.")
	fmt.Println("only use this if you really know what you are doing!")

	var (
		password string
		keystore string
		outfile  string
		secret   secretStuff
	)
	flag.StringVar(&password, "password", "", "password used to encrypt the keystore")
	flag.StringVar(&keystore, "keystore", "", "SMESH APP keystore file")
	flag.StringVar(&outfile, "output", "", "file base to store keys in (xxx => xxx01, xxx02 etc)")
	flag.Parse()
	if len(password) == 0 || len(keystore) == 0 || len(outfile) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	var smData smeshAccountKeystruct
	f, err := os.Open(keystore)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}
	err = json.NewDecoder(f).Decode(&smData)
	if err != nil {
		log.Fatal(err)
	}
	keyBytes := []byte(password)
	key := pbkdf2.Key(keyBytes, []byte("Spacemesh blockmesh"), 1000000, 32, sha512.New)
	ciphertext, err := hex.DecodeString(smData.Crypto.CipherText)
	if err != nil {
		fmt.Println("Keystore file not valid", err)
		os.Exit(1)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("AES error", err)
		os.Exit(1)
	}
	iv := make([]byte, c.BlockSize())
	iv[15] = byte(5)
	plaintextBytes := decrypt(c, ciphertext, iv)
	err = json.Unmarshal(plaintextBytes, &secret)
	if err != nil {
		fmt.Println("data error", err)
		os.Exit(1)
	}
	for n, acc := range secret.Accounts {
		cak := cliAccountKeys{PrivKey: acc.SecretKey, PubKey: acc.PublicKey}
		path := fmt.Sprintf("%s%d.json", outfile, n)
		f, err := os.Create(path)
		if err != nil {
			fmt.Println("could not create ", path, err)
			continue
		}
		err = json.NewEncoder(f).Encode(cak)
		if err != nil {
			fmt.Println("could not save to ", path, err)
			f.Close()
			continue
		}
		fmt.Printf("Saved [%s] to %s\n", acc.DisplayName, path)
		f.Close()
	}
}
