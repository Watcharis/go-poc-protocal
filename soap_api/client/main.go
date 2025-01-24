package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"watcharis/go-poc-protocal/soap_api/models"
	"watcharis/go-poc-protocal/soap_api/repositories"
)

func main() {
	fmt.Println("start request soap api")
	ctx := context.Background()

	start := time.Now()
	log.Printf("start time : %+v", start)

	soapRepository := repositories.NewSoapRepository()

	requestVerifyUser := models.VerifyUser{
		CitizenID: "1234560098897",
		Name:      "Mock",
		Email:     "mock@test.com",
		MobileNo:  "0998765432",
		Age:       "78",
	}

	result, err := soapRepository.SoapVerifyUser(ctx, requestVerifyUser)
	if err != nil {
		fmt.Println("[ERROR] soap verify user failed  err:", err)
		panic(err)
	}

	fmt.Printf("verifyUserResponse Resp: %+v\n", result.SoapBody.ResponseBody)
	fmt.Printf("verifyUserResponse Resp: %+v\n", result.SoapBody.ResponseBody.Data)

	log.Printf("end time : %+v", time.Since(start))
}
