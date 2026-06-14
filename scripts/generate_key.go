package main

import (
	"bytes"
	"crypto"
	"encoding/pem"
	"log"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func main() {
	// Создаем новую сущность PGP
	entity, err := openpgp.NewEntity("Bank System", "bank@example.com", "", &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       2048,
	})
	if err != nil {
		log.Fatalf("Ошибка при создании PGP сущности: %v", err)
	}

	// Сохраняем публичный ключ
	pubKeyFile, err := os.Create("key.pem")
	if err != nil {
		log.Fatalf("Ошибка при создании файла публичного ключа: %v", err)
	}
	defer pubKeyFile.Close()

	var pubBuf bytes.Buffer
	if err := entity.Serialize(&pubBuf); err != nil {
		log.Fatalf("Ошибка при сериализации публичного ключа: %v", err)
	}

	pubKeyBlock := &pem.Block{
		Type:  "PGP PUBLIC KEY",
		Bytes: pubBuf.Bytes(),
	}
	if err := pem.Encode(pubKeyFile, pubKeyBlock); err != nil {
		log.Fatalf("Ошибка при сохранении публичного ключа: %v", err)
	}

	// Сохраняем приватный ключ
	privKeyFile, err := os.Create("key.priv.pem")
	if err != nil {
		log.Fatalf("Ошибка при создании файла приватного ключа: %v", err)
	}
	defer privKeyFile.Close()

	var privBuf bytes.Buffer
	if err := entity.SerializePrivate(&privBuf, nil); err != nil {
		log.Fatalf("Ошибка при сериализации приватного ключа: %v", err)
	}

	privKeyBlock := &pem.Block{
		Type:  "PGP PRIVATE KEY",
		Bytes: privBuf.Bytes(),
	}
	if err := pem.Encode(privKeyFile, privKeyBlock); err != nil {
		log.Fatalf("Ошибка при сохранении приватного ключа: %v", err)
	}

	log.Println("PGP ключи успешно сгенерированы")
}
