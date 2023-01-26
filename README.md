## Genealogy Tree Challenge
Este desafio consiste em criar uma aplicação web para armazenar e visualizar informações de uma árvore genealógica.

## Compose Stack
- Golang para o desenvolvimento da API Rest
- PostgreSQL como banco de dados
- Swagger para documentação e testes da API

## Executando o projeto com Docker
Para executar o projeto, é necessário ter o Docker instalado em sua máquina.

- Clone este repositório para sua máquina local.
- Na raiz do projeto, execute o comando docker-compose build para construir as imagens dos containers:
```
docker-compose build
```
- Em seguida, execute o comando docker-compose up para iniciar os containers.
```
docker-compose up
```
- A aplicação ficará disponível em http://localhost:9000

## Documentação da API
A documentação da API está disponível através do Swagger.

**Swagger**
Link para acesso: http://localhost:9000/genealogy/swagger/index.html
![Swagger Representation](/assets/swagger.png)

## Arquitetura
A arquitetura deste projeto segue o padrão Clean Architecture, que separa as responsabilidades em camadas distintas.
Projeto usado como referência: https://github.com/bxcodec/go-clean-arch