package main

import(

    "fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)
//Testa função getAccount
func TestGetAccount(t *testing.T) {
	fmt.Println("Teste de getAccount()")
	//Chama a função e valida se retorna erro
	_,err := getAccount("5555")
	assert.NoError(t, err)
}
//Testa função SecretToHash
func TestSecretToHash(t *testing.T) {
	fmt.Println("Teste de SecretToHash()")
	//Chama a função e valida se retorna erro
	hashSecret,err := SecretToHash("12345")
	fmt.Println(hashSecret)
	assert.NoError(t, err)
}
//Testa função checkSecret
func TestCheckSecret(t *testing.T) {
	fmt.Println("Teste de CheckSecret()")
	//Chama a função e valida se retorna erro
	hashSecret := "$2a$10$gQo1NUM8YJzONzgROUW9JuAAMwP.vVtEYZCERDrlzRFDV1Aks36dG"
	secret := "12345"
	_,err := checkSecret(hashSecret,secret)
	assert.NoError(t, err)
}
//Testa função updateBalanceAccount
func TestUpdateBalanceAccount(t *testing.T) {
	fmt.Println("Teste de updateBalanceAccount()")
	//Chama a função e valida se retorna erro
	err := updateBalanceAccount(1,2.0)
	assert.NoError(t, err)
}
//Testa função storeTransfer
func TestStoreTransfer(t *testing.T) {
	fmt.Println("Teste de storeTransfer()")
	//Chama a função e valida se retorna erro
	var transfer Transfer
	transfer = Transfer{ID : "2", Account_Origin_Id : 3, Account_Destination_Id : 2, Amount : 10.1, Created_At : "01/06/2021 22:37:00"}
	err := storeTransfer(transfer)
	assert.NoError(t, err)
}

//Testa função storeDeposit
func TestStoreDeposit(t *testing.T) {
	fmt.Println("Teste de storeDeposit()")
	//Chama a função e valida se retorna erro
	var deposit Deposit
	deposit = Deposit{ ID : "2", CPF : "5555", Account_Destination_Id : 3, Amount : 10.0, Created_At : "01/06/2021 22:37:00"}
	err := storeDeposit(deposit)
	assert.NoError(t, err)
}

//Testa função verifyAccountID
func TestVerifyAccountID(t *testing.T) {
	fmt.Println("Teste de verifyAccountID()")
	//Chama a função e valida se retorna erro
	_,err := verifyAccountID(2)
	assert.NoError(t, err)
}

