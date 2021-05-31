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
RUN go get -u golang.org/x/crypto/bcrypt
RUN go get -u github.com/dgrijalva/jwt-go
#RUN go get -u github.com/satori/go.uuid
RUN go get -u github.com/google/uuid

# copia o arquivo do diretório do host para o conteiner
COPY main.go main.go
COPY utils.go utils.go
COPY auth.go auth.go

#Compilação do arquivo
RUN go build main.go auth.go utils.go
RUN cp main /

#Execução do arquivo main
CMD [ "/main" ]