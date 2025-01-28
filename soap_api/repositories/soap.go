package repositories

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"watcharis/go-poc-protocal/pkg/httpclient"
	"watcharis/go-poc-protocal/soap_api/models"
)

const (
	SOAP_ENDPOINT = "http://localhost:8889/verify/user"
)

type SoapRepository interface {
	SoapVerifyUser(ctx context.Context, req models.VerifyUser) (models.ResponseVerifyUser, error)
}

type soapRepository struct{}

func NewSoapRepository() SoapRepository {
	return &soapRepository{}
}

func (r *soapRepository) SoapVerifyUser(ctx context.Context, data models.VerifyUser) (models.ResponseVerifyUser, error) {

	soapBody, err := r.generateSoapVerifyUserTemplate(data)
	if err != nil {
		return models.ResponseVerifyUser{}, err
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, SOAP_ENDPOINT, bytes.NewBuffer([]byte(soapBody)))
	if err != nil {
		return models.ResponseVerifyUser{}, err
	}

	// Set the required headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	// req.Header.Set("SOAPAction", soapAction)

	// Send the HTTP request
	httpClient := httpclient.CreateHttpClient()

	resp, err := httpClient.Do(req)
	if err != nil {
		return models.ResponseVerifyUser{}, err
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.ResponseVerifyUser{}, err
	}

	// Print the response
	fmt.Println("response status_code:", resp.Status)
	fmt.Println("response body:", string(respBody))

	var verifyUserResponse models.ResponseVerifyUser
	if err := xml.Unmarshal(respBody, &verifyUserResponse); err != nil {
		return models.ResponseVerifyUser{}, err
	}
	// fmt.Printf("verifyUserResponse : %+v\n", verifyUserResponse)

	return verifyUserResponse, nil
}

func (r *soapRepository) generateSoapVerifyUserTemplate(data models.VerifyUser) (string, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("VerifyUser").Parse(models.SoapVerifyUserSchemaTemplate)
	if err != nil {
		return "", err
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, data)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		return "", err
	}

	return doc.String(), nil
}
