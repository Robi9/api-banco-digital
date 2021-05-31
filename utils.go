package main

import (
    
    "fmt"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"

)
//Retorna uma conta partindo do cpf do usuário
func getAccount(cpf string) (Account) {

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

//Atualiza o balance de uma conta para um novo valor
func updateBalanceAccount(id int, balance float64) (error){

	//Define a consulta do filtro para buscar um documento específico da coleção
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	//Define o atualizador para especificar a mudança a ser atualizada.
	updater := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "balance", Value: balance},
	}}}

	// //Faz a conexão com o MongoDB.
	client, err := getMongoClient()
	if err != nil {
		fmt.Println(err)
	}
	collection := client.Database(DB).Collection(ACCOUNT)

	//Executa a operação UpdateOne e valida o erro.
	_, err = collection.UpdateOne(context.TODO(), filter, updater)
	if err != nil {
		return err
	}

	return nil
}

//Armazena a transferencia no BD
func storeTransfer(transfer Transfer) {

	client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }

    collection := client.Database(DB).Collection(TRANSFER)
    //Insere o dado e valida
    _, err = collection.InsertOne(context.TODO(), transfer)
    if err != nil {
        fmt.Println(err)
        return
    }
    return
}

//Armazena depósitos realizados no BD
func storeDeposit(deposit Deposit) {

	client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }

    collection := client.Database(DB).Collection(DEPOSIT)
    //Insere o dado e valida
    _, err = collection.InsertOne(context.TODO(), deposit)
    if err != nil {
        fmt.Println(err)
        return
    }
    return
	
}