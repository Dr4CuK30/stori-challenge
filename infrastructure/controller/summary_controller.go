package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"io"
	"mime"
	"mime/multipart"
	"stori-challenge-v1/application"
	"stori-challenge-v1/infrastructure/utils"
	"strings"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SummaryController struct {
	csvReader      utils.CsvReader
	summaryService *application.SummaryService
}

func (s *SummaryController) HandlePostRequest(ctx context.Context, req *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	boundary, err := s.validateBody(req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Content-Type is not multipart/form-data",
		}, nil
	}
	multipartReader := multipart.NewReader(strings.NewReader(req.Body), boundary)

	data, filename, email, err := s.readPostMultipartValues(multipartReader)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	response, err := s.summaryService.ProcessTransactionsCsv(data, filename, email)
	responseBody, err := json.Marshal(Response{Message: "Success", Data: response})
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func (s *SummaryController) validateBody(req *events.APIGatewayProxyRequest) (string, error) {
	contentType := req.Headers["Content-Type"]
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(contentType, "multipart/form-data") {
		return "", errors.New("Content-Type no es multipart/form-data")
	}
	boundary := params["boundary"]
	return boundary, nil
}

func (s *SummaryController) readPostMultipartValues(reader *multipart.Reader) ([][]string, string, string, error) {
	var email string
	var data [][]string
	var filename string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, "", "", err
		}

		if part.FormName() == "transactions" && part.FileName() != "" && strings.HasSuffix(part.FileName(), ".csv") {
			csvRes, err := s.csvReader.ReadCSV(part)
			if err != nil {
				return nil, "", "", err
			}
			data = csvRes
			filename = part.FileName()
		} else if part.FormName() == "destination" {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, part)
			if err != nil {
				return nil, "", "", err
			}
			email = buf.String()
		}
	}
	return data, filename, email, nil
}
