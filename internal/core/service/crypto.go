package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash"

	"os"
)

func NewPublicKey(file string) (*rsa.PublicKey, error) {
	if file == "" {
		return nil, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading public key error: %w", err)
	}

	block, _ := pem.Decode(data)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing public key error: %w", err)
	}

	return publicKey.(*rsa.PublicKey), err
}

func NewPrivateKey(file string) (*rsa.PrivateKey, error) {
	if file == "" {
		return nil, nil
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("reading private key error: %w", err)
	}

	block, _ := pem.Decode(data)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing private key error: %w", err)
	}

	return privateKey, nil
}

func encryptChunk(data []byte, hash hash.Hash, pupKey *rsa.PublicKey) ([]byte, error) {
	b, err := rsa.EncryptOAEP(hash, rand.Reader, pupKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa OAEP encrypt error:%w", err)
	}

	return b, nil
}

func Encrypt(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	var encData []byte

	dataLen := len(data)
	hash := sha512.New()
	step := chunkSize(pubKey.Size(), hash.Size())
	for begin := 0; begin < dataLen; begin += step {
		end := begin + step
		if end > dataLen {
			end = dataLen
		}

		encChunk, err := encryptChunk(data[begin:end], hash, pubKey)
		if err != nil {
			return nil, fmt.Errorf("encrypt chunk error:%w", err)
		}

		encData = append(encData, encChunk...)
	}

	return encData, nil
}

func Decrypt(privKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	var decData []byte

	dataLen := len(data)
	hash := sha512.New()
	step := 512
	for begin := 0; begin < dataLen; begin += step {
		end := begin + step
		if end > dataLen {
			end = dataLen
		}

		decChunk, err := decryptChunk(data[begin:end], hash, privKey)
		if err != nil {
			return nil, fmt.Errorf("decrypt chunk error:%w", err)
		}

		decData = append(decData, decChunk...)
	}

	return decData, nil
}

func decryptChunk(data []byte, hash hash.Hash, privKey *rsa.PrivateKey) ([]byte, error) {
	b, err := rsa.DecryptOAEP(hash, rand.Reader, privKey, data, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa OAEP decrypt error:%w", err)
	}

	return b, nil
}

// The message must be no longer than the length of the public modulus minus
// twice the hash length, minus a further 2.
func chunkSize(keySize int, hashSize int) int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.23.1:src/crypto/rsa/rsa.go;l=527
	return keySize - 2*hashSize - 2
}
