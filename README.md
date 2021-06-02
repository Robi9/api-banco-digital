# api-banco-digital
Este projeto é de uma API de transferência entre contas internas de um banco digital. Criado na linguagem Go, com a utilização do Docker e o MongoDB para armazenamento.

## Começando
O projeto utiliza as seguintes tecnologias:
  * [Linguagem Go](https://golang.org/)
  * [Docker Desktop](https://docs.docker.com/engine/install/)
  * [Docker Compose](https://docs.docker.com/compose/install/)
  * [MongoDB](https://www.mongodb.com/try/download/community) (Banco de dados)

Como o projeto usa conteinerização, não precisamos instalar a linguagem e o banco de dados.

## Desenvolvimento
Para iniciar o desenvolvimento é necessário clonar o projeto do GitHub em um diretório de sua preferência.
```
cd "diretorio-de-sua-preferencia"
git clone https://github.com/Robi9/api-banco-digital   
```
## Construção
Antes de construirmos os conteineres com o Docker Compose, precisamos criar um novo volume no Docker chamado ``dados_``.
Este projeto utiliza um volume externo para armazenamento dos dados do MongoDB. Para criarmos este volume basta dar o seguinte comando:
```
docker volume create dados_   
```
Após criarmos o volume, podemos construir nossos conteineres de acordo como estão definidos nos arquivos ``Dockerfile`` e ``docker-compose.yml``. O comando abaixo irá criar
os conteineres com as respectivas imagens do ``Go`` e ``MongoDB``.
O comando para iniciarmos a construção é o seguinte:
```
docker-compose build   
```
Ao concluir a construção com sucesso, podemos iniciar a execução de nossa api:
```
docker-compose up  
```
## Funcionalidades da Aplicação
Nossa aplicação conta com algumas rotas:
* <a name="accpost"></a>(POST) ``/accounts``: Cria uma nova conta.
* <a name="accget"></a>(GET) ``/accounts``: Retorna todas as contas cadastradas.
* <a name="accbalance"></a>(GET) ``/accounts/{ID}/balance``: Retorna o saldo da conta que pertence ao ID informado.
* <a name="transfpost"></a>(POST) ``/transfers``: Realiza uma transferência da conta autenticada para outra conta cadastrada.
* <a name="transfget"></a>(GET) ``/transfers``: Retorna todas transferências feitas pelo usuário logado.
* <a name="deppost"></a>(POST) ``/deposits``: Realiza um depósito em uma conta cadastrada.
* <a name="login"></a>(POST) ``/login``: Realiza login em uma conta cadastrada.

### [Rota de Criação de Conta](#accpost)
Esta rota recebe os dados da conta em json na solicitação, abaixo está um exemplo de solicitação com a ferramenta ``curl``. Antes, se não tiver a ferramenta instalada
poderá instalá-la com o comando ``sudo apt-get install curl`` no terminal. 
```
curl -H "Content-Type:application/json" -X POST -d '{ "ID" : 1, "name" : "Fulano", "cpf" : "111.111.111-11", "secret" : "12345", "balance" : 10.0, "created_at" : ""}' "http://localhost:5000/accounts"
```
Esta solicitação criará uma nova conta com os dados passados no json.
A informação de ``created_at`` é gerada pela aplicação.

### [Rota de Contas Cadastradas](#accget)
Esta rota retorna todas as contas cadastradas no banco. Abaixo temos o exemplo de solicitação:
```
curl http://localhost:5000/accounts
```
### [Rota de Balance](#accbalance)
Esta rota retorna o ``balance`` da conta com o ``ID`` informado, nesse exemplo a conta criada anteriormente com id 1. Abaixo temos o exemplo da solicitação:
```
curl http://localhost:5000/accounts/1/balance
```
### [Rota de Realização de Transferência](#transfpost)
Esta rota realiza transferências entre a conta logada e uma outra conta cadastrada. Abaixo temos o exemplo da solicitação:
```
curl -i POST -H "Content-Type: application/json" -H "Authorization: <token>" -d '{ "ID": "", "account_origin_id" : 0, "account_destination_id" : 2, "amount" : 10.0, "created_at" : ""}' "http://localhost:5000/transfers"
```
Para realizar a transferência você precisa estar logado, o login autentica um usuário e gera um token, com ele você pode realizar a solicitação.
O ``account_origin_id`` é extraído do token de autenticação e os outros dados que não são passados no json são gerados pela api. (O valor '0' significa vazio)
Troque ``<token>`` pelo token gerado ao realizar login.

### [Rota de Transferências Realizadas](#transfget)
Esta rota retorna todas as transferências realizadas pelo usuário logado. Abaixo temos o exemplo da solicitação:
```
curl -i POST -H "Content-Type: application/json" -H "Authorization: <token>" "http://localhost:5000/transfers"
```
### [Rota de Realização de Depósito](#deppost)
Esta rota realiza um depósito em uma conta cadastrada, basta informar o cpf do dono da conta que irá receber, o id da conta e o amount. Abaixo temos o exemplo da solicitação:
```
curl -i POST -H "Content-Type: application/json" -H "Authorization: <token>" -d `{ "ID" : "", "cpf" : "111.111.111-11", "account_destination_id" : 2, "amount" : 10.0, "created_at" : ""}` "http://localhost:5000/deposits"
```
### [Rota de Login](#login)
Esta rota realiza o login em uma conta cadastrada e autentica o usuário com a geraçãod e um token de autenticação, este é retornado após a realização de login com sucesso.
```
curl -i POST -H "Content-Type: application/json" -d `{ "cpf" : "111.111.111-11", "secret" : "12345"}` "http://localhost:5000/login"
```
## Testes
Para executar os testes execute o seguinte comando dentro da pasta do projeto:
```
go test
```
## Links
Abaixo estão algumas ferramentas utilizadas na construção da API:
* [Gorilla Mux](https://github.com/gorilla/mux)
* [MongoDB Go Driver](https://docs.mongodb.com/drivers/go/)
* [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
* [JWT Go](github.com/dgrijalva/jwt-go)
