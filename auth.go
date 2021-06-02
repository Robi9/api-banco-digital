package main

import (
    "net/http"
    "fmt"
    "github.com/dgrijalva/jwt-go"
)

//Gera o Token ao realizar um  login
func auth(account Account, result Login) (string, error) {

    //Usa o CPF e o secret para a geração
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "ID"    : account.ID,
        "cpf"   : result.CPF,
        "secret": result.Secret,
    })
    //Faz a assinatura com a key "secret" e valida
    tokenString, err := token.SignedString([]byte("secret"))
    if err != nil {
        return "",err
    }
    return tokenString, err
}

//Verifica e retorna o token da autenticação
func verifyToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {

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
        return nil, err
    }

    return token, nil

}
