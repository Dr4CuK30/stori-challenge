package controller

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"stori-challenge-v1/application"
	"stori-challenge-v1/domain/repositories"
	"stori-challenge-v1/infrastructure/reosurces/email"
	"stori-challenge-v1/infrastructure/utils"
)

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var response events.APIGatewayProxyResponse
	var final_err error
	controller := SummaryController{
		csvReader: &utils.CsvReaderImp{},
		summaryService: &application.SummaryService{
			TransactionRepository: &repositories.PgTransactionRepository{},
			EmailResource:         &email.GmailResource{},
		},
	}
	if req.HTTPMethod == http.MethodPost {
		request, err := controller.HandlePostRequest(ctx, &req)
		final_err = err
		response = request
	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "Method not allowed",
		}, nil
	}
	if final_err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error processing request",
		}, nil
	}
	return response, nil
}
