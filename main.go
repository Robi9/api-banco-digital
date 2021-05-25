// main.go
package main
import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "sync"
    "context"
    "io/ioutil"
    "encoding/json"
    //"go.mongodb.org/mongo-driver/bson"
    //"go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "time"
    "golang.org/x/crypto/bcrypt"

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

//Função de rotas
func routes() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/accounts", listaContas)
    myRouter.HandleFunc("/account", criaNovaConta).Methods("POST")

    log.Fatal(http.ListenAndServe(":5000", myRouter))
}

//Cria novo Account e armazena no BD
func criaNovaConta(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Endpoint accounts: apiCriaNovaConta")
    reqBody,_ := ioutil.ReadAll(r.Body)

    var account Account
    err := json.Unmarshal(reqBody, &account)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(account.ID)

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
    account.Balance = 1

    //Cria um handle da respectiva coleção
    collection := client.Database(DB).Collection(ACCOUNT)
    //Insere o dado e valida
    _, err = collection.InsertOne(context.TODO(), account)
    if err != nil {
        fmt.Println(err)
    }

    json.NewEncoder(w).Encode(account)
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

//Retorna a lista de contas cadastradas
func listaContas(w http.ResponseWriter, r *http.Request) {
   fmt.Println("Endpoint Acessado.")
}

func main() {
    fmt.Println("API-TEST.")
    routes()
}
