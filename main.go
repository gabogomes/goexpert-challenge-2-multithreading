package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CEPResponseFromBrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type CEPResponseFromViaCEPApi struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApiResponse struct {
	Response interface{}
	Error    error
}

func main() {
	cep := "01153000"
	URLForBrasilApi := "https://brasilapi.com.br/api/cep/v1/" + cep
	URLForViaCEPApi := "http://viacep.com.br/ws/" + cep + "/json/"

	chBrasilApi := make(chan ApiResponse)
	chViaCEPApi := make(chan ApiResponse)

	go MakeCEPApiRequestToBrasilApi(URLForBrasilApi, chBrasilApi)
	go MakeCEPApiRequestToViaCEPApi(URLForViaCEPApi, chViaCEPApi)

	select {
	case BrasilApiResponse := <-chBrasilApi:
		fmt.Printf("Response from BrasilAPI: %+v\n", BrasilApiResponse.Response)
		fmt.Printf("URL of Request: %s\n", URLForBrasilApi)
	case ViaCEPApiResponse := <-chViaCEPApi:
		fmt.Printf("Response from ViaCEPAPI: %+v\n", ViaCEPApiResponse.Response)
		fmt.Printf("URL of Request: %s\n", URLForViaCEPApi)
	case <-time.After(1 * time.Second):
		err := fmt.Errorf("Timeout")
		fmt.Println(err)
	}
}

func MakeCEPApiRequestToBrasilApi(url string, ch chan ApiResponse) {
	response, err := http.Get(url)
	if err != nil {
		ch <- ApiResponse{nil, err}
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		ch <- ApiResponse{nil, err}
		return
	}

	var formattedResponse CEPResponseFromBrasilApi
	err = json.Unmarshal(body, &formattedResponse)
	if err != nil {
		ch <- ApiResponse{nil, err}
		return
	}

	ch <- ApiResponse{formattedResponse, nil}
}

func MakeCEPApiRequestToViaCEPApi(url string, ch chan ApiResponse) {
	response, err := http.Get(url)
	if err != nil {
		ch <- ApiResponse{nil, err}
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		ch <- ApiResponse{nil, err}
		return
	}

	var formattedResponse CEPResponseFromViaCEPApi
	err = json.Unmarshal(body, &formattedResponse)
	if err != nil {
		ch <- ApiResponse{nil, err}
		return
	}

	ch <- ApiResponse{formattedResponse, nil}
}
