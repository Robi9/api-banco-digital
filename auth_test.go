package main

import(
	"net/http"
    "github.com/dgrijalva/jwt-go"
    "fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"bytes"
)

//Testa função auth()
func TestAuth(t *testing.T) {
	fmt.Println("Teste de auth()")
	//Chama a função e valida se retorna erro
	account := Account{ ID : 3, Name : "Fulano", CPF : "5555", Secret : "12345", Balance : 0.0, Created_At : ""}
	login := Login{CPF : "5555", Secret : "12345"}
	token,err := auth(account, login)
	fmt.Println("Token:", token)
	assert.NoError(t, err)
}

func TestVerifyToken(t *testing.T) {
	fmt.Println("Teste de verifyToken()")
	
	//Simula uma transferencia
	var jsonTrans = []byte(`{ "ID" : "", "account_origin_id" : 0, "account_destination_id" : 1, "amount" : 10.0, "created_at" : ""}`)
	req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonTrans))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MywiY3BmIjoiNTU1NSIsInNlY3JldCI6IjEyMzQ1In0.25BL0qmCrmpZKmPkUatLi5gfnLMtnZv2N-5aCXKHY1o")

	response := executeRequest(req)

	//Chama a função e valida se retorna erro
	token,err := verifyToken(response, req)
	assert.NoError(t, err)
	claims,_ := token.Claims.(jwt.MapClaims);
	//Os dados recuperados são da conta criada para teste na main_test
	fmt.Println("Dados recuperados do Token:", claims)
	
}