package node

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipoluianov/cc_node/logger"
)

func GenerateKey() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	privateKeyBS := crypto.FromECDSA(privateKey)
	publicKeyBS := crypto.FromECDSAPub(&privateKey.PublicKey)
	fmt.Println("private_key:", hex.EncodeToString(privateKeyBS))
	fmt.Println("public_key:", hex.EncodeToString(publicKeyBS))
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	logger.Println("Address:", address.Hex())
}
