package main

import (
    
    "fmt"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"

)
//Retorna uma conta
func getAccount(cpf string ) Account {

    result := Account{}
    //Define a consulta do filtro para buscar um documento específico da coleção
    filter := bson.D{primitive.E{Key: "cpf", Value: cpf}}
    //Faz a conexão com o MongoDB.
    client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }
    //Cria um handle da respectiva coleção.
    collection := client.Database(DB).Collection(ACCOUNT)
    //Busca a account e faz a validação
    err = collection.FindOne(context.TODO(), filter).Decode(&result)

    if err != nil {
        fmt.Println(err)
    }

    return result
}

//Transforma o secret em hash
func SecretToHash(secret string) string {
    cost := bcrypt.DefaultCost
    hash, err := bcrypt.GenerateFromPassword([]byte(secret), cost)
    if err != nil {
        fmt.Println(err)
    }
    return string(hash)
}
//Verifica se o hash do BD é igual ao secret enviado 
func checkSecret(secretH string, secret string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(secretH), []byte(secret))
	if err != nil {
		return false
	}
	return true
}