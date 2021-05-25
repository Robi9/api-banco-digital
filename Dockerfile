# Dockerfile

# start with an alpine image for small footprint
#Imagem do conteiner (Go)
FROM golang:1.12-alpine

#Diretório de trabalho no conteiner
WORKDIR /code

# Instalações de dependências
RUN apk update && apk upgrade && apk add --no-cache git 
RUN go get github.com/gorilla/mux
RUN go get -u go.mongodb.org/mongo-driver/bson
RUN go get -u go.mongodb.org/mongo-driver/bson/primitive
RUN go get -u go.mongodb.org/mongo-driver/mongo
RUN go get -u go.mongodb.org/mongo-driver/mongo/options
RUN go get golang.org/x/crypto/bcrypt

# copia o arquivo do diretório do host para o conteiner
COPY main.go main.go

#Compilação do arquivo
RUN go build main.go
RUN cp main /

#Execução do arquivo main
CMD [ "/main" ]