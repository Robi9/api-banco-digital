package main

import(
    "github.com/gorilla/mux"
    "net/http"
    "net/http/httptest"
    "fmt"
    "encoding/json"
    "bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)
//Realiza a solicitação
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    Router().ServeHTTP(rr, req)

    return rr
}
//Faz a validação da resposta da solicitação
func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}
//Função de rotas igual a de main.go, para auxiliar os testes.
func Router() *mux.Router {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/login", authLogin).Methods("POST") //Faz login
    myRouter.HandleFunc("/transfers", newTransfer).Methods("POST") //Realiza transferência
    myRouter.HandleFunc("/transfers", getAllTransfers) //Retorna todas transferências feitas pelo usuário logado
    myRouter.HandleFunc("/accounts", newAccount).Methods("POST") //Cria nova conta
    myRouter.HandleFunc("/accounts", getAllAccounts) //Retorna todas as contas cadastradas
    myRouter.HandleFunc("/accounts/{ID}/balance", getBalance) //Retorna o saldo da conta que pertence ao ID informado
    myRouter.HandleFunc("/deposits", newDeposit).Methods("POST") //Realiza um depósito em uma conta cadastrada
    return myRouter
}

//Testa função getMongoClient
func TestGetMongoClient(t *testing.T) {
	//Chama a função e valida se retorna erro
	_,err := getMongoClient()
	assert.NoError(t, err)
}

//Testa função rota /accounts (POST). (Rota de criação de conta)
func TestNewAccount(t *testing.T) {
	//Simulamos a criação de uma conta
		//Mude os dados caso teste mais de uma vez!
	var jsonAcc = []byte(`{ "ID" : 1, "name" : "Fulano", "cpf" : "7777", "secret" : "12345", "balance" : 99.0, "created_at" : ""}`)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonAcc))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	//Pega os dados do response para realizar a validação
	var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)
    fmt.Println(m["name"])

    //Valida os dados passados com os dados do body

    if m["name"] != "Fulano" {
        t.Errorf("Expected name to be 'Fulano'. Got '%v'", m["name"])
    }

    if m["cpf"] != "7777"{
        t.Errorf("Expected cpf to be '7777'. Got '%v'", m["cpf"])
    }

    // O id é comparado com 1.0 porque o JSON unmarshaling converte números para
    // floats, quando usa o map[string]interface{}
    if m["ID"] != 1.0 {
        t.Errorf("Expected product ID to be '1'. Got '%v'", m["ID"])
    }
}