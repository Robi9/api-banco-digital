package main

import (
    "net/http"
    "io/ioutil"
    "encoding/json"
    "fmt"
    "github.com/dgrijalva/jwt-go"
)
//Estrutura de Login
type Login struct {
    CPF string    `json:"cpf"`
    Secret string `json:"secret"`
}
//Gera o Token ao realizar um  login
func authLogin(w http.ResponseWriter, r *http.Request) {

    reqBody,_ := ioutil.ReadAll(r.Body)

    var result Login
    err := json.Unmarshal(reqBody, &result)
    if err != nil {
        fmt.Println(err)
    }
    //Busca a conta com o CPF informado no login
    account := getAccount(result.CPF)

    //Valida se o secret informado no login é igual ao cadastrado, se sim inicia a geração do token
    if checkSecret(account.Secret, result.Secret) {
        //Usa o CPF e o secret para a geração
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "cpf":  result.CPF,
            "secret": result.Secret,
        })
        //Faz a assinatura com a key "secret" e valida
        tokenString, err := token.SignedString([]byte("secret")) //VER SE EU CONSIGO ESCONDER

        if err != nil {
            json.NewEncoder(w).Encode(err.Error())
            return
        }
        //Se ocorrer tudo certo print o Token, se não mensagem de erro.
        json.NewEncoder(w).Encode(tokenString)
    }else{
        json.NewEncoder(w).Encode("Secret Incorreto.")
    }    
}
