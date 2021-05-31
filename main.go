// main.go
package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "sync"
    "context"
    "strconv"
    "io/ioutil"
    "encoding/json"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
    "github.com/google/uuid"
    "github.com/dgrijalva/jwt-go"

)

/* Usado para criar um objeto único do cliente MongoDB.
Inicializado e exposto por meio de GetMongoClient(). */
var clientInstance *mongo.Client

//Usado durante a criação do objeto cliente único em GetMongoClient().
var clientInstanceError error

//Usado para executar o procedimento de criação do cliente apenas uma vez.
var mongoOnce sync.Once

//Dados de configuração do BD
const (
    CONNECTIONSTRING = "mongodb://mongodb:27017" //localhost
    DB               = "api-banco-digital"
    ACCOUNT          = "accounts"
    TRANSFER         = "transfers"
    DEPOSIT          = "deposits"
)

//GetMongoClient - Retorne a conexão com mongodb
func getMongoClient() (*mongo.Client, error) {
    //Executa a operação de criação de conexão apenas uma vez.
    mongoOnce.Do(func() {
        // Define as opções do cliente
        clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
        // Conecta ao Mongodb
        client, err := mongo.Connect(context.TODO(), clientOptions)
        if err != nil {
            clientInstanceError = err
        }
        // Verifica a conexão
        err = client.Ping(context.TODO(), nil)
        if err != nil {
            clientInstanceError = err
        }
        clientInstance = client
    })
    return clientInstance, clientInstanceError
   
}

// Estrutura de Account
type Account struct {
    ID         int     `json:"ID"  bson:"_id,omitempty"`
    Name       string  `json:"name" bson:"name"`
    CPF        string  `json:"cpf" bson:"cpf"`
    Secret     string  `json:"secret" bson:"secret"`
    Balance    float64 `json:"balance" bson:"balance"`
    Created_At string  `json:"created_at" bson:"created_at"`
}

//Estrutura de Transfer
type Transfer struct{
    ID                     string     `json:"ID"  bson:"_id,omitempty"`
    Account_Origin_Id      int        `json:"account_origin_id"  bson:"account_origin_id"`
    Account_Destination_Id int        `json:"account_destination_id"  bson:"account_destination_id"`
    Amount                 float64      `json:"amount"  bson:"amount"` 
    Created_At             string     `json:"created_at"  bson:"created_at"`

}

//Estrutura de Deposit
type Deposit struct{
    ID                     string     `json:"ID"  bson:"_id,omitempty"`
    CPF                    string     `json:"cpf"  bson:"cpf"`
    Account_Destination_Id int        `json:"account_destination_id"  bson:"account_destination_id"`
    Amount                 float64    `json:"amount"  bson:"amount"` 
    Created_At             string     `json:"created_at"  bson:"created_at"`

}

//Função de rotas
func routes() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/login", authLogin).Methods("POST") //Faz login
    myRouter.HandleFunc("/transfers", newTransfer).Methods("POST") //Realiza transferência
    //myRouter.HandleFunc("/transfers", getAllTransfers) //Retorna todas transferências feitas pelo usuário logado
    myRouter.HandleFunc("/accounts", newAccount).Methods("POST") //Cria nova conta
    myRouter.HandleFunc("/accounts", getAllAccounts) //Retorna todas as contas cadastradas
    myRouter.HandleFunc("/accounts/{ID}/balance", getBalance) //Retorna o saldo da conta que pertence ao ID informado
    //myRouter.HandleFunc("/deposits", newDeposit).Methods("POST") //Realiza um depósito em uma conta cadastrada

    log.Fatal(http.ListenAndServe(":5000", myRouter))
}

//Cria novo Account e armazena no BD
func newAccount(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Endpoint: apiNewAccount")
    reqBody,_ := ioutil.ReadAll(r.Body)

    var account Account
    err := json.Unmarshal(reqBody, &account)
    if err != nil {
        fmt.Println(err)
    }

    client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }

    //Pega hora/data de agora
    start := time.Now()
    //Formata hora/data e adiciona em Created_At
    account.Created_At = start.Format(("02/01/2006 15:04:05"))
    //Transforma secret em hash
    account.Secret = SecretToHash(account.Secret)

    //Zera valor da conta
    account.Balance =  0.0

    //Cria um handle da respectiva coleção
    collection := client.Database(DB).Collection(ACCOUNT)
    //Insere o dado e valida
    _, err = collection.InsertOne(context.TODO(), account)
    if err != nil {
        fmt.Println(err)
    }

    json.NewEncoder(w).Encode(account)
}

//Retorna a lista de contas cadastradas
func getAllAccounts(w http.ResponseWriter, r *http.Request) {
   fmt.Println("Endpoint: apiGetAllAccounts")

   //Defina a consulta do filtro para buscar um documento específico da coleção
    filter := bson.D{{}} //bson.D{{}} especifica 'todos os documentos'
    accounts := []Account{}
    //Faz a conexão com o MongoDB
    client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }
    //Cria um handle da respectiva coleção
    collection := client.Database(DB).Collection(ACCOUNT)

    //Executa a operação Localizar e valide o erro.
    cur, findError := collection.Find(context.TODO(), filter)
    if findError != nil {
        fmt.Println(findError)
    }
    //Map de resultados para a slice
    for cur.Next(context.TODO()) {
        t := Account{}
        err := cur.Decode(&t)
        if err != nil {
            fmt.Println(err)
        }
        accounts = append(accounts, t)
    }
    // quando terminado fecha o cursor
    cur.Close(context.TODO())
    json.NewEncoder(w).Encode(accounts)
}

func getBalance(w http.ResponseWriter, r *http.Request) {
   
    fmt.Println("Endpoint: apiGetBalance")

    vars := mux.Vars(r)
    id := vars["ID"]
    _id, _ := strconv.Atoi(id)

    result := Account{}
    //Define a consulta do filtro para buscar um documento específico da coleção
    filter := bson.D{primitive.E{Key: "_id", Value: _id}}
    //Faz a conexão com o MongoDB.
    client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }
    //Cria um handle da respectiva coleção.
    collection := client.Database(DB).Collection(ACCOUNT)

    err = collection.FindOne(context.TODO(), filter).Decode(&result)

    if err == nil {
        json.NewEncoder(w).Encode(result.Balance)
    } else {
        fmt.Println(err)
    }
}

func newTransfer(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint: apiNewTransfer")

    //Pega hora/data de agora
    start := time.Now()

    //Verifica e pega token
    token := verifyToken(w,r)
    if token == nil{
        json.NewEncoder(w).Encode("Token inválido, entre em sua conta novamente!")
        return
    }
   
    //Pega dados de transferência
    transfer := Transfer{}
    reqBody,_ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(reqBody, &transfer)
    if err != nil {
        fmt.Println(err)
    }

    var result, accountOrigin, accountDestination Account

    //Extrai dados do token
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        result.ID = int(claims["ID"].(float64))
        result.CPF = claims["cpf"].(string)
    } else {
        fmt.Println("Erro ao recuperar dados do Token")
        return
    }

    //Gera um ID para a transferencia
    idTransfer := uuid.New()
    transfer.ID = idTransfer.String()

    //Pega o ID de origem que estava no Token
    transfer.Account_Origin_Id = result.ID

    //Formata hora/data e adiciona em Created_At
    transfer.Created_At = start.Format(("02/01/2006 15:04:05"))

    //Busca accountOrigin partindo do CPF
    accountOrigin = getAccount(result.CPF)

    //Define a consulta do filtro para buscar um documento específico da coleção
    filter := bson.D{primitive.E{Key: "_id", Value: transfer.Account_Destination_Id}}

    //Faz a conexão com o MongoDB.
    client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }

    //Cria um handle da respectiva coleção.
    collection := client.Database(DB).Collection(ACCOUNT)

    err = nil
    //Busca a accountDestination e faz a validação
    err = collection.FindOne(context.TODO(), filter).Decode(&accountDestination)
    if err != nil {
        fmt.Println(err)
    }
    //Verificação de balance disponível na accountOrigem
    if accountOrigin.Balance >= transfer.Amount {
        //Calculo do Balance novo de conta de origem e destino
        newBalanceAccountOrigin  := accountOrigin.Balance - transfer.Amount
        newBalanceAccountDestination := accountDestination.Balance + transfer.Amount

        //Atualiza Balance de ambas as contas
        updateBalanceAccount(accountOrigin.ID, newBalanceAccountOrigin)
        updateBalanceAccount(accountDestination.ID, newBalanceAccountDestination)

        json.NewEncoder(w).Encode("Transferência realizada com sucesso!")
    }else{
        json.NewEncoder(w).Encode("Saldo insuficiente.")
    }

    //Armazena a transferência no BD
    storeTransfer(transfer)
    return
}

//Retorna todas as Transferências do usuário autenticado
func getAllTransfers(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint: apiGetAllTransfers")

    token := verifyToken(w,r)
    if token == nil{
        json.NewEncoder(w).Encode("Token inválido, entre em sua conta novamente!")
        return
    }

    var ID_Account int
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        ID_Account = int(claims["ID"].(float64))
    } else {
         fmt.Println("Erro ao recuperar dados do Token")
        return
    }

   //Defina a consulta do filtro para buscar um documento específico da coleção
    filter := bson.D{primitive.E{Key: "account_origin_id", Value: ID_Account}}
    transfers := []Transfer{}
    //Faz a conexão com o MongoDB
    client, err := getMongoClient()
    if err != nil {
        fmt.Println(err)
    }
    //Cria um handle da respectiva coleção
    collection := client.Database(DB).Collection(TRANSFER)

    //Executa a operação Localizar e valide o erro.
    cur, findError := collection.Find(context.TODO(), filter)
    if findError != nil {
        fmt.Println(findError)
    }
    //Map de resultados para a slice
    for cur.Next(context.TODO()) {
        t := Transfer{}
        err := cur.Decode(&t)
        if err != nil {
            fmt.Println(err)
        }
        transfers = append(transfers, t)
    }
    // quando terminado fecha o cursor
    cur.Close(context.TODO())
    json.NewEncoder(w).Encode(transfers)
}

func newDeposit(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint: apiMakeDeposit")

    //Pega hora/data de agora
    start := time.Now()

    reqBody,_ := ioutil.ReadAll(r.Body)
    var deposit Deposit
    err := json.Unmarshal(reqBody, &deposit)
    if err != nil {
        fmt.Println(err)
    }

    //Gera um ID para a transferencia
    idDeposit := uuid.New()
    deposit.ID = idDeposit.String()

    //Formata hora/data e adiciona em Created_At
    deposit.Created_At = start.Format(("02/01/2006 15:04:05"))

    //Pega conta para depósito
    accountDestination := getAccount(deposit.CPF)

    deposit.Account_Destination_Id = accountDestination.ID

    //Calcula novo balance da conta de depósito
    balance := accountDestination.Balance+deposit.Amount

    //Atualiza balance da conta de depósito e valida
    err = nil
    err = updateBalanceAccount(accountDestination.ID, balance)

    if err != nil{
        json.NewEncoder(w).Encode("Erro no depósito, tente novamente!")
        return
    }else{
        storeDeposit(deposit)
        json.NewEncoder(w).Encode("Depósito realizado com sucesso!")
    }    
}

func main() {
    fmt.Println("API-BANCO-DIGITAL.")
    routes()
}

