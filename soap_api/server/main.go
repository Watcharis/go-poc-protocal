package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/pkg/logger"
	"watcharis/go-poc-protocal/soap_api/models"

	"go.uber.org/zap"
)

func main() {

	startServer := make(chan bool, 1)

	http.Handle("/health", http.HandlerFunc(pkg.HealthCheck))

	http.HandleFunc("/verify/user", VerifyUserHandler)

	go func() {
		log.Println("Server is running on http://localhost:8889")
		if err := http.ListenAndServe(":8889", nil); err != nil {
			panic(err)
		}
	}()

	<-startServer
}

func VerifyUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Check if the request is POST
	if r.Method != http.MethodPost {
		logger.Error(ctx, "method not allow, want http method POST", zap.String("request_method", r.Method))
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(ctx, "Unable to read request body", zap.Error(err))
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Log the received request
	logger.Info(ctx, "Received SOAP Request", zap.String("request_body", string(body)))

	var result models.RequestVerifyUser
	if err := xml.Unmarshal(body, &result); err != nil {
		http.Error(w, "Unable to xml.Unmarshal body", http.StatusInternalServerError)
		return
	}
	logger.Info(ctx, "result xml.Unmarshal request_body", zap.Any("request", result.Body.VerifyUserRequest))

	// Send a mock SOAP response
	responseTemplate := `<?xml version="1.0" encoding="UTF-8"?>
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:web="http://example.com/webservice">
			<soapenv:Body>
			<web:VerifyUserResponse>
				<code>1000</code>
				<message>success</message>
				<data>
					<cid>%s</cid>
				</data>
			</web:VerifyUserResponse>
			</soapenv:Body>
		</soapenv:Envelope>`

	response := fmt.Sprintf(responseTemplate, result.Body.VerifyUserRequest.CitizenID)

	// Set response headers
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
