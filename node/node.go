package node

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ipoluianov/cc_node/logger"
)

type Node struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address

	remotePublicKeys map[string]*ecdsa.PublicKey

	active bool
}

func NewNode() *Node {
	var c Node
	c.remotePublicKeys = make(map[string]*ecdsa.PublicKey)
	return &c
}

func (c *Node) Start(path string) error {
	privateKeyBSHex, err := os.ReadFile(path + "/private_key.txt")
	if err != nil {
		return err
	}

	privateKeyBS, err := hex.DecodeString(string(privateKeyBSHex))
	if err != nil {
		return err
	}

	c.privateKey, err = crypto.ToECDSA(privateKeyBS)
	if err != nil {
		return err
	}

	c.address = crypto.PubkeyToAddress(c.privateKey.PublicKey)

	com.RegisterReceiver(c.address.Hex(), c)

	_, err = os.ReadFile(path + "/active.txt")
	if err == nil {
		c.active = true
		go c.ThAction()
	}
	return nil
}

func (c *Node) Stop() {
}

func (c *Node) ThAction() {
	for {
		logger.Println("--------------------------------")
		logger.Println("node", c.address.Hex(), "action")
		err := com.Send("", NewFrame(c.address.Hex(), "0xAa4BF73B6e1Cd966db0EF8c9Da008BE8Fe344a01", "get_public_key", ""))
		if err != nil {
			logger.Println("Send Frame Error:", err)
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Node) Receive(frame Frame) {
	logger.Println("Node::Receive", frame)

	if frame.Type == "get_public_key" {
		publicKeyBS := crypto.FromECDSAPub(&c.privateKey.PublicKey)
		response := NewFrame(c.address.Hex(), frame.SrcAddress, "get_public_key_result", hex.EncodeToString(publicKeyBS))
		response.Content1 = c.address.Hex()
		err := com.Send("", response)
		if err != nil {
			logger.Println("Send response error:", err)
		}
	}

	if frame.Type == "get_public_key_result" {
		logger.Println("received public key of host", frame.Content1, " = ", frame.Content)
		remotePublicKeyBS, err := hex.DecodeString(frame.Content)
		if err != nil {
			logger.Println("Decode Public Key Error:", err)
			return
		}
		remotePublicKey, err := crypto.UnmarshalPubkey(remotePublicKeyBS)
		if err != nil {
			logger.Println("Unmarshal Public Key Error:", err)
			return
		}

		c.remotePublicKeys[frame.SrcAddress] = remotePublicKey

		eciesPublicKeyB := ecies.ImportECDSAPublic(remotePublicKey)

		message := "Hello, Node B!"
		ciphertext, err := ecies.Encrypt(rand.Reader, eciesPublicKeyB, []byte(message), nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Encrypted message: %x\n", ciphertext)

		com.Send("", NewFrame(c.address.Hex(), frame.SrcAddress, "message", hex.EncodeToString(ciphertext)))
	}

	if frame.Type == "message" {
		logger.Println("RECEIVED MESSAGE", frame.Content)
		messageBS, err := hex.DecodeString(frame.Content)
		if err != nil {
			fmt.Println("RECEIVED MESSAGE DecodeString error:", err)
		}
		eciesPrivateKeyB := ecies.ImportECDSA(c.privateKey)
		decryptedMessage, err := eciesPrivateKeyB.Decrypt(messageBS, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Decrypted message: %s\n", string(decryptedMessage))
	}
}
