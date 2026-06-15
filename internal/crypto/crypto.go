package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

const (
	hmacKey = "your-secret-hmac-key"
)

func HashCVV(cvv string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("ошибка хеширования CVV: %v", err)
	}
	return string(hashed), nil
}

func VerifyCVV(hashedCVV, cvv string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedCVV), []byte(cvv)) == nil
}

func EncryptCardData(cardNumber, expiryDate string) (string, error) {
	entity, err := openpgp.NewEntity("Card Encryption", "", "", nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания ключа: %v", err)
	}

	data := fmt.Sprintf("%s|%s", cardNumber, expiryDate)
	buf := new(bytes.Buffer)

	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{entity}, nil, nil, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка шифрования: %v", err)
	}

	_, err = w.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("ошибка записи данных: %v", err)
	}
	w.Close()

	armored := new(bytes.Buffer)
	armorWriter, err := armor.Encode(armored, "PGP MESSAGE", nil)
	if err != nil {
		return "", fmt.Errorf("ошибка кодирования: %v", err)
	}

	_, err = armored.WriteTo(armorWriter)
	if err != nil {
		return "", fmt.Errorf("ошибка записи armored данных: %v", err)
	}
	armorWriter.Close()

	return armored.String(), nil
}

func DecryptCardData(encryptedData string) (string, string, error) {
	block, err := armor.Decode(strings.NewReader(encryptedData))
	if err != nil {
		return "", "", fmt.Errorf("ошибка декодирования: %v", err)
	}

	md, err := openpgp.ReadMessage(block.Body, nil, nil, nil)
	if err != nil {
		return "", "", fmt.Errorf("ошибка чтения сообщения: %v", err)
	}

	plaintext, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", "", fmt.Errorf("ошибка чтения данных: %v", err)
	}

	parts := strings.Split(string(plaintext), "|")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("неверный формат данных")
	}

	return parts[0], parts[1], nil
}

func ComputeHMAC(data string) string {
	h := hmac.New(sha256.New, []byte(hmacKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data, expectedHMAC string) bool {
	actualHMAC := ComputeHMAC(data)
	return hmac.Equal([]byte(actualHMAC), []byte(expectedHMAC))
}
