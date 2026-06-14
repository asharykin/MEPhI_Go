package main

import (
	"crypto/rand"
	"encoding/pem"
	"flag"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"time"
)

func main() {
	name := flag.String("name", "Bank System", "Имя владельца ключа")
	email := flag.String("email", "bank@example.com", "Email владельца ключа")
	output := flag.String("output", "key.pem", "Путь для сохранения ключа")
	flag.Parse()

	// Создаем конфигурацию для генерации ключа
	config := &packet.Config{
		DefaultHash:            packet.SHA256,
		DefaultCipher:          packet.CipherAES256,
		DefaultCompressionAlgo: packet.CompressionZLIB,
		RSABits:                4096,
	}

	// Создаем сущность
	entity, err := openpgp.NewEntity(*name, "", *email, config)
	if err != nil {
		fmt.Printf("Ошибка создания сущности: %v\n", err)
		os.Exit(1)
	}

	// Генерируем подпись
	for _, id := range entity.Identities {
		err := id.SelfSignature.SignUserId(id.UserId.Id, entity.PrimaryKey, entity.PrivateKey, config)
		if err != nil {
			fmt.Printf("Ошибка подписи: %v\n", err)
			os.Exit(1)
		}
	}

	// Сохраняем ключ в файл
	file, err := os.Create(*output)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Кодируем ключ в PEM формат
	block := &pem.Block{
		Type:  "PGP PRIVATE KEY",
		Bytes: entity.PrivateKey.SerializePrivateKey(nil, config),
	}

	err = pem.Encode(file, block)
	if err != nil {
		fmt.Printf("Ошибка кодирования ключа: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ключ успешно создан и сохранен в %s\n", *output)
}
