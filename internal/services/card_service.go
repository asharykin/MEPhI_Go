package services

import (
	"banksystem/internal/models"
	"banksystem/internal/repositories"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"math/rand"
	"strconv"
	"time"
)

type CardService struct {
	cardRepo *repositories.CardRepository
	key      *openpgp.Entity
}

func NewCardService(cardRepo *repositories.CardRepository, key *openpgp.Entity) *CardService {
	return &CardService{
		cardRepo: cardRepo,
		key:      key,
	}
}

func (s *CardService) CreateCard(accountID int64) (*models.Card, error) {
	// Генерация номера карты по алгоритму Луна
	cardNumber := generateLuhnCardNumber()

	// Генерация CVV
	cvv := strconv.Itoa(rand.Intn(900) + 100)

	// Хеширование CVV
	hashedCVV, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Шифрование данных карты
	encryptedData, err := s.encryptCardData(cardNumber, cvv)
	if err != nil {
		return nil, err
	}

	// Создание HMAC
	hmac := s.createHMAC(encryptedData)

	card := &models.Card{
		AccountID:     accountID,
		EncryptedData: encryptedData,
		HashedCVV:     string(hashedCVV),
		HMAC:          hmac,
		CreatedAt:     time.Now(),
	}

	err = s.cardRepo.Create(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (s *CardService) GetCard(id int64) (*models.Card, error) {
	return s.cardRepo.GetByID(id)
}

func (s *CardService) GetUserCards(userID int64) ([]*models.Card, error) {
	return s.cardRepo.GetByUserID(userID)
}

func (s *CardService) encryptCardData(cardNumber, cvv string) (string, error) {
	// Создаем буфер для зашифрованных данных
	buf := new(bytes.Buffer)

	// Создаем armor-кодировщик
	armorWriter, err := armor.Encode(buf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}

	// Создаем зашифрованный writer
	encryptedWriter, err := openpgp.Encrypt(armorWriter, []*openpgp.Entity{s.key}, nil, nil, nil)
	if err != nil {
		return "", err
	}

	// Записываем данные карты
	_, err = encryptedWriter.Write([]byte(cardNumber + "|" + cvv))
	if err != nil {
		return "", err
	}

	// Закрываем writers
	err = encryptedWriter.Close()
	if err != nil {
		return "", err
	}

	err = armorWriter.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *CardService) createHMAC(data string) string {
	h := hmac.New(sha256.New, []byte("secret-key"))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func generateLuhnCardNumber() string {
	// Генерация 15 цифр
	digits := make([]int, 15)
	for i := 0; i < 15; i++ {
		digits[i] = rand.Intn(10)
	}

	// Вычисление контрольной суммы по алгоритму Луна
	sum := 0
	for i := 0; i < 15; i++ {
		if i%2 == 0 {
			digits[i] *= 2
			if digits[i] > 9 {
				digits[i] -= 9
			}
		}
		sum += digits[i]
	}

	// Добавление контрольной цифры
	checkDigit := (10 - (sum % 10)) % 10
	digits = append(digits, checkDigit)

	// Преобразование в строку
	result := ""
	for _, d := range digits {
		result += strconv.Itoa(d)
	}

	return result
}

func (s *CardService) VerifyCard(cardID int64, cvv string) error {
	card, err := s.cardRepo.GetByID(cardID)
	if err != nil {
		return err
	}

	if !card.VerifyCVV(cvv) {
		return errors.New("неверный CVV")
	}

	return nil
}
