package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/raulsilva-tech/Address-Finder/internal/dto"
)

/*
	Neste desafio você terá que usar o que aprendemos com Multithreading e APIs para buscar o resultado mais rápido entre duas APIs distintas.

As duas requisições serão feitas simultaneamente para as seguintes APIs:
https://brasilapi.com.br/api/cep/v1/01153000 + cep
http://viacep.com.br/ws/" + cep + "/json/
Os requisitos para este desafio são:
- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.
*/

type AddressHandler struct{}

type Error struct {
	Message string
}

func NewAddressHandler() *AddressHandler {
	return &AddressHandler{}
}

func (ah *AddressHandler) GetFastestAddressAnswer(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	cepParam := chi.URLParam(r, "cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{"CEP is required"})
		return
	}

	chViaCEP := make(chan dto.ViaCEPDTO)
	chBrasilAPI := make(chan dto.BrasilAPIDTO)

	go viaCEP(cepParam, chViaCEP)
	go brasilAPI(cepParam, chBrasilAPI)

	select {
	case response := <-chViaCEP:
		json.NewEncoder(w).Encode(response)
		fmt.Printf("ViaCEP: \n %v \n", response)
	case response := <-chBrasilAPI:
		json.NewEncoder(w).Encode(response)
		fmt.Printf("BrasilAPI: \n %v \n", response)
	case <-time.After(time.Second):
		json.NewEncoder(w).Encode(Error{"Timeout"})
		fmt.Println(Error{"Timeout"})
	}

	w.WriteHeader(http.StatusOK)
}

func viaCEP(cep string, ch chan dto.ViaCEPDTO) {

	// time.Sleep(time.Second)

	resp, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Not Found")
		return
	}

	//transformando corpo na struct ViaCEP
	var data dto.ViaCEPDTO
	err = json.Unmarshal(bodyBytes, &data)

	ch <- data

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}

func brasilAPI(cep string, ch chan dto.BrasilAPIDTO) {

	// time.Sleep(time.Second)

	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Fprintf(os.Stderr, "Not Found")
		return
	}

	//transformando corpo na struct ViaCEP
	var data dto.BrasilAPIDTO
	err = json.Unmarshal(bodyBytes, &data)

	ch <- data

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
}
