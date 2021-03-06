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
    myRouter.HandleFunc("/login", newLogin).Methods("POST") //Faz login
    myRouter.HandleFunc("/transfers", newTransfer).Methods("POST") //Realiza transferência
    myRouter.HandleFunc("/transfers", getAllTransfers) //Retorna todas transferências feitas pelo usuário logado
    myRouter.HandleFunc("/accounts", newAccount).Methods("POST") //Cria nova conta OK
    myRouter.HandleFunc("/accounts", getAllAccounts) //Retorna todas as contas cadastradas OK
    myRouter.HandleFunc("/accounts/{ID}/balance", getBalance) //Retorna o saldo da conta que pertence ao ID informado OK
    myRouter.HandleFunc("/deposits", newDeposit).Methods("POST") //Realiza um depósito em uma conta cadastrada OK
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
	var jsonAcc = []byte(`{ "ID" : 3, "name" : "Fulano", "cpf" : "5555", "secret" : "12345", "balance" : 0.0, "created_at" : ""}`)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(jsonAcc))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	//Pega os dados do response para realizar a validação
	var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    //Valida os dados passados com os dados do body

    if m["name"] != "Fulano" {
        t.Errorf("Expected name to be 'Fulano'. Got '%v'", m["name"])
    }

    if m["cpf"] != "5555"{
        t.Errorf("Expected cpf to be '5555'. Got '%v'", m["cpf"])
    }

    // O id é comparado com 1.0 porque o JSON unmarshaling converte números para
    // floats, quando usa o map[string]interface{}
    if m["ID"] != 3.0 {
        t.Errorf("Expected ID to be '3'. Got '%v'", m["ID"])
    }
}

//Testa rota GET /accounts, que retorna todas as contas cadastradas
func TestGetAllAccounts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/accounts", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

//Testa a rota /deposits que faz um depósito
func TestNewDeposit(t *testing.T) {

	//Simulamos a criação de um depósito
	var jsonDep = []byte(`{ "ID" : "", "cpf" : "5555", "account_destination_id" : 3, "amount" : 10.0, "created_at" : ""}`)
	req, _ := http.NewRequest("POST", "/deposits", bytes.NewBuffer(jsonDep))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m string
    json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Println(m)	
}

//Testa a rota de adquirir balance de uma conta
func TestGetBalance(t *testing.T) {

	//Chamamos a rota para o id 3
	req, _ := http.NewRequest("GET", "/accounts/3/balance", nil)

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

    var m float64
    json.Unmarshal(response.Body.Bytes(), &m)
    //printa o Balance da conta com o ID passado na solicitação
    fmt.Println("Balance da conta com ID 3 é:", m)
}

//Teste da rota POST '/transfers' que realiza transferências de uma conta logada para outra cadastrada
func TestNewTransfer(t *testing.T) {
	//Simulamos uma transferência
		//Lembrando que o dado de id_origem é pego no token
		//Usaremos o token da conta criada na função TestNewAccount()
	var jsonTrans = []byte(`{ "ID" : "", "account_origin_id" : 0, "account_destination_id" : 1, "amount" : 10.0, "created_at" : ""}`)
	req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonTrans))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiY3BmIjoiNTU1NSIsInNlY3JldCI6IjEyMzQ1In0.25BL0qmCrmpZKmPkUatLi5gfnLMtnZv2N-5aCXKHY1o")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m string
    json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Println(m)
	
}

//Teste da rota GET '/transfers' que retorna todas as transferencias de uma conta 
func TestGetAllTransfers(t *testing.T) {
	//Chamamos a rota com o mesmo token da conta criada em TestNewAccount()
	req, _ := http.NewRequest("GET", "/transfers", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiY3BmIjoiNTU1NSIsInNlY3JldCI6IjEyMzQ1In0.25BL0qmCrmpZKmPkUatLi5gfnLMtnZv2N-5aCXKHY1o")
	
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

    var m []Transfer
    json.Unmarshal(response.Body.Bytes(), &m)
    
    //Printa todas as transferencias realizadas pelo usuário autenticado
    for _,i := range(m) {
    	fmt.Println("----------------")
    	fmt.Println("ID:",i.ID)
    	fmt.Println("Account_Origin_Id:", i.Account_Origin_Id)
    	fmt.Println("Account_Destination_Id:", i.Account_Destination_Id)
    	fmt.Println("Amount:", i.Amount)
    	fmt.Println("Created_At:", i.Created_At)
    }
}

//Teste da '/login' que realiza um login e faz a autenticação do usuário logado
func TestNewLogin(t *testing.T){
	//Login com os dados da conta criada em TestNewAccount()
	var jsonLogin = []byte(`{ "cpf" : "5555", "secret" : "12345"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m string
    json.Unmarshal(response.Body.Bytes(), &m)
	fmt.Println(m)
}

