package application

import (
	"fmt"
	"stori-challenge-v1/domain/entities"
	"stori-challenge-v1/domain/repositories"
	"stori-challenge-v1/infrastructure/reosurces/email"
	"stori-challenge-v1/infrastructure/utils"
	"strconv"
	"strings"
)

type SummaryResponse struct {
	TotalBalance                float64        `json:"totalBalance"`
	AverageDebitAmount          float64        `json:"averageDebitAmount"`
	AverageCreditAmount         float64        `json:"averageCreditAmount"`
	TransactionsQuantityByMonth map[string]int `json:"transactionsQuantityByMonth"`
}

type SummaryService struct {
	TransactionRepository repositories.TransactionRepository
	EmailResource         email.EmailResource
}

func (s *SummaryService) ProcessTransactionsCsv(transactions [][]string, csvName string, email string) (SummaryResponse, error) {
	averageCreditAmount := 0.0
	averageDebitAmount := 0.0
	transactionsQuantityByMonth := make(map[string]int)
	totalBalance := 0.0
	process := entities.Process{
		OriginName:   csvName,
		Origin:       CSV_ORIGIN,
		Transactions: []entities.Transaction{},
	}
	for _, record := range transactions {
		id := record[0]
		dateSplit := strings.Split(record[1], "/")
		amount, err := strconv.ParseFloat(record[2], 64)
		month, err := strconv.ParseUint(dateSplit[0], 10, 8)
		day, err := strconv.ParseUint(dateSplit[1], 10, 8)
		if err != nil {
			return SummaryResponse{}, err
		}
		process.Transactions = append(process.Transactions, entities.Transaction{
			Id:     id,
			Amount: amount,
			Day:    uint8(day),
			Month:  uint8(month),
		})
		totalBalance += amount
		if amount > 0 {
			averageCreditAmount += amount
		} else {
			averageDebitAmount += amount
		}
		monthName := utils.GetMonthByUint8(uint8(month))
		transactionsQuantityByMonth[monthName]++
	}
	fmt.Println("Saving data")
	summary := SummaryResponse{
		TotalBalance:                totalBalance,
		AverageCreditAmount:         averageCreditAmount,
		AverageDebitAmount:          averageDebitAmount,
		TransactionsQuantityByMonth: transactionsQuantityByMonth,
	}
	err := s.sendSummaryByEmail(summary, email)
	if err != nil {
		return SummaryResponse{}, err
	}
	err = s.TransactionRepository.Save(process)
	if err != nil {
		return SummaryResponse{}, err
	}
	return summary, nil
}

func (s *SummaryService) sendSummaryByEmail(summary SummaryResponse, email string) error {
	emailBody := fmt.Sprintf(
		"<b>Hello,</b><br><br>Here is your account summary:<br><br>"+
			"Total Balance: $%.2f<br>"+
			"Average Debit Amount: $%.2f<br>"+
			"Average Credit Amount: $%.2f<br><br>"+
			"Transactions Quantity By Month:<br>",
		summary.TotalBalance,
		summary.AverageDebitAmount,
		summary.AverageCreditAmount,
	)

	for month, quantity := range summary.TransactionsQuantityByMonth {
		emailBody += fmt.Sprintf("- %s: %d transactions<br>", month, quantity)
	}
	emailBody += "<br>Best regards,<br>"
	emailBody += "<br><br>" +
		"<div style=\"display: flex; align-items: center;\">" +
		"<img " +
		"width=\"150\" " +
		"height=\"150\" " +
		"src=\"https://play-lh.googleusercontent.com/oXTAgpljdbV5LuAOt1NP9_JafUZe9BNl7pwQ01ndl4blYL4N4IQh4-n456P5l_hc1A\" " +
		"alt=\"Stori logo\"" +
		">" +
		"<div style=\"display: flex; flex-direction: column; margin-left: 10px;\">      " +
		"<span style=\"margin: 10px;\"><b>Stori Challenge</b></span>  " +
		"<span style=\"margin: 10px;\"><b>David Saldarriaga</b></span> " +
		"<span style=\"margin: 10px;\"><b>2024</b></span>" +
		"</div>" +
		"</div>"
	subject := "Stori: Transactions summary"
	s.EmailResource.SendEmail(email, subject, emailBody)
	return nil
}
