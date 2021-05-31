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

    w.Header().Set("Content-Type", "application/json")

    reqBody,_ := ioutil.ReadAll(r.Body)

    var result Login
    err := json.Unmarshal(reqBody, &result)
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println(result.CPF)
    //Busca a conta com o CPF informado no login
    account := getAccount(result.CPF)

    //Valida se o secret informado no login é igual ao cadastrado, se sim inicia a geração do token
    if checkSecret(account.Secret, result.Secret) {
        //Usa o CPF e o secret para a geração
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "ID"    : account.ID,
            "cpf"   : result.CPF,
            "secret": result.Secret,
        })
        //Faz a assinatura com a key "secret" e valida
        tokenString, err := token.SignedString([]byte("secret")) //VER SE EU CONSIGO ESCONDER
        //fmt.Println(err)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("Erro ao gerar JWT token: " + err.Error()))
            return
        }

        w.Header().Set("Authorization", tokenString)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Token: " + tokenString))

    }else{
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("Name e secret não conferem, tente novamente!"))
        return
    }    
}

//Verifica e retorna o token da autenticação
func verifyToken(w http.ResponseWriter, r *http.Request) (*jwt.Token) {

    w.Header().Set("Content-Type", "application/json")
    //Pega token 
    tokenString := r.Header.Get("Authorization")
    fmt.Println(tokenString)
    //Verifica token da autenticação
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method")
        }
        return []byte("secret"), nil
    })

    if err != nil{
        fmt.Println("Token Inválido.")
        fmt.Println(err)
    }

    return token

}
