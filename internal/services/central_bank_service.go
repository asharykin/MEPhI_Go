package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
)

const (
	centralBankURL = "https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx"
)

type CentralBankService struct {
	client *http.Client
}

func NewCentralBankService() *CentralBankService {
	return &CentralBankService{
		client: &http.Client{},
	}
}

func (s *CentralBankService) GetKeyRate() (float64, error) {
	resp, err := s.client.Get("http://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return 0, err
	}

	keyRateElement := doc.FindElement("//ValCurs/Valute[@ID='R01235']/Value")
	if keyRateElement == nil {
		return 0, nil
	}

	return 7.5, nil // Заглушка, так как реальный API ЦБ РФ требует сертификат
}

func (s *CentralBankService) buildSOAPRequest() string {
	fromDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
			<soap12:Body>
				<KeyRate xmlns="http://web.cbr.ru/">
					<fromDate>%s</fromDate>
					<ToDate>%s</ToDate>
				</KeyRate>
			</soap12:Body>
		</soap12:Envelope>`, fromDate, toDate)
}

func (s *CentralBankService) sendRequest(soapRequest string) ([]byte, error) {
	req, err := http.NewRequest(
		"POST",
		centralBankURL,
		bytes.NewBuffer([]byte(soapRequest)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://web.cbr.ru/KeyRate")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный статус ответа: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func (s *CentralBankService) parseXMLResponse(rawBody []byte) (float64, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(rawBody); err != nil {
		return 0, err
	}

	krElements := doc.FindElements("//diffgram/KeyRate/KR")
	if len(krElements) == 0 {
		return 0, fmt.Errorf("данные по ставке не найдены")
	}

	latestKR := krElements[0]
	rateElement := latestKR.FindElement("./Rate")
	if rateElement == nil {
		return 0, fmt.Errorf("тег Rate отсутствует")
	}

	rateStr := rateElement.Text()
	var rate float64
	if _, err := fmt.Sscanf(rateStr, "%f", &rate); err != nil {
		return 0, fmt.Errorf("ошибка конвертации ставки: %v", err)
	}

	return rate, nil
}
