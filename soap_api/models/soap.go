package models

import "encoding/xml"

type VerifyUser struct {
	CitizenID string
	Name      string
	Email     string
	MobileNo  string
	Age       string
}

const SoapVerifyUserSchemaTemplate string = `<?xml version="1.0" encoding="UTF-8"?>
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:web="http://example.com/webservice">
		<soapenv:Header/>
		<soapenv:Body>
			<web:VerifyUserRequest>
				<CitizenID>{{.CitizenID}}</CitizenID>
				<Name>{{.Name}}</Name>
				<Email>{{.Email}}</Email>
				<MobileNo>{{.MobileNo}}</MobileNo>
				<Age>{{.Age}}</Age>
			</web:VerifyUserRequest>
		</soapenv:Body>
	</soapenv:Envelope>`

type RequestVerifyUser struct {
	XMLName xml.Name               `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    *RequestVerifyUserBody `xml:"Body"`
}

type RequestVerifyUserBody struct {
	VerifyUserRequest *SoapVerifyUserRequest `xml:"VerifyUserRequest"`
}

type SoapVerifyUserRequest struct {
	// XMLName   xml.Name `xml:"VerifyUserRequest"`
	CitizenID string `xml:"CitizenID"`
	Name      string `xml:"Name"`
	Email     string `xml:"Email"`
	MobileNo  string `xml:"MobileNo"`
	Age       string `xml:"Age"`
}

const SoapVerifyUserResponseTemplate = `<?xml version="1.0" encoding="UTF-8"?>
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:web="http://example.com/webservice">
		<soapenv:Body>
		<web:VerifyUserResponse>
			<code>1000</code>
			<message>success</message>
			<data>
				<cid></cid>
			</data>
		</web:VerifyUserResponse>
		</soapenv:Body>
	</soapenv:Envelope>`

type ResponseVerifyUser struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *SOAPBodyVerifyUserResponse
}

type SOAPBodyVerifyUserResponse struct {
	XMLName      xml.Name `xml:"Body"`
	ResponseBody *ResponseVerifyUserBody
}

type ResponseVerifyUserBody struct {
	XMLName xml.Name                    `xml:"VerifyUserResponse"`
	Code    string                      `xml:"code"`
	Message string                      `xml:"message"`
	Data    *DataResponseVerifyUserBody `xml:"data"`
}

type DataResponseVerifyUserBody struct {
	Cid string `xml:"cid"`
}

// // Define the struct for the SOAP Envelope
// type Envelope struct {
// 	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
// 	Body    Body     `xml:"Body"`
// 	// Xmlns   string   `xml:"xmlns:soapenv,attr"`
// }

// // Define the struct for the Body
// type Body struct {
// 	VerifyUserResponse VerifyUserResponse `xml:"VerifyUserResponse"`
// }

// // Define the struct for the VerifyUserResponse
// type VerifyUserResponse struct {
// 	Code    int    `xml:"code"`
// 	Message string `xml:"message"`
// }
