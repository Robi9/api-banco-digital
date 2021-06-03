package main

import (
    "strconv"
    "fmt"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"

)
//Retorna uma conta partindo do cpf do usuário
func getAccount(cpf string) (Account, error) {

    var result Account
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
        return result, err
    }

    return result, err
}

func verifyAccountID(id int) (bool, error){
	var result Account
    //Define a consulta do filtro para buscar um documento específico da coleção
    filter := bson.D{primitive.E{Key: "_id", Value: id}}
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
        return false,err
    }

    return true,err
}

//Transforma o secret em hash
func SecretToHash(secret string) (string, error) {
    cost := bcrypt.DefaultCost
    hash, err := bcrypt.GenerateFromPassword([]byte(secret), cost)
    if err != nil {
        fmt.Println(err)
        return "",err
    }
    return string(hash), err
}
//Verifica se o hash do BD é igual ao secret enviado 
func checkSecret(secretH string, secret string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(secretH), []byte(secret))
	if err != nil {
		return false, err
	}
	return true, err
}

//Atualiza o balance de uma conta para um novo valor
func updateBalanceAccount(id int, balance float64) (error){
	v, err := verifyAccountID(id)
	if v != true {
		fmt.Println("Conta não existe.")
		return err
	}
	f := fmt.Sprintf("%.2f",balance)
	balance,_ = strconv.ParseFloat(f, 64)

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
		fmt.Println("Conta não encontrada.")
		return err
	}
	return nil
}

//Armazena a transferencia no BD
func storeTransfer(transfer Transfer) (error){

	client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }

    collection := client.Database(DB).Collection(TRANSFER)
    //Insere o dado e valida
    _, err = collection.InsertOne(context.TODO(), transfer)
    if err != nil {
        //fmt.Println(err)
        return err
    }
    return err
}

//Armazena depósitos realizados no BD
func storeDeposit(deposit Deposit) (error){

	client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }

    collection := client.Database(DB).Collection(DEPOSIT)
    //Insere o dado e valida
    _, err = collection.InsertOne(context.TODO(), deposit)
    if err != nil {
        fmt.Println(err)
        return err
    }
    return err
}