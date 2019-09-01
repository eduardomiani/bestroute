# Best Route

[![Build Status](https://travis-ci.org/eduardomiani/bestroute.svg?branch=master)](https://travis-ci.org/eduardomiani/bestroute)

Programa desenvolvido com o desafio de encontrar a melhor rota para um viajante, dado um arquivo de rotas cadastradas.
O programa foi desenvolvido na linguagem [Go](https://golang.org/), devido a produtividade e os recursos disponibilizados pelas bibliotecas nativas da linguagem.
Para a interface rest, não foi utilizado nenhum framework ou biblioteca externa, apenas as bibliotecas nativas.

Os arquivos estão na seguinte estrutura:

* /rest - *pacote com os arquivos relacionados a interface rest*
* /route - *pacote com os arquivos core da aplicação, tratamento das rotas e manipulação do arquivo*
* main.go - *arquivo principal, executa o programa, disponibilizando uma das duas interfaces*


 **Tipos de arquivos**:

 * *.go - *Código fonte da aplicação*
 * *_test.go - *Arquivos contendo os testes unitários (por convenção da comunidade, eles ficam no mesmo pacote do código fonte que ele irá testar)*
 * *.csv - *Arquivo criado especialmente para alguns testes unitários*

## Como começar

Caso queira baixar o código, compilar e rodar por conta própria, será necessário fazer a instalação da [SDK do Go](https://golang.org/dl/).
Após instalada a SDK seguindo as instruções da página, executar o seguinte comando no terminal:

```sh
$ go get -u github.com/eduardomiani/bestroute
```

Se a instalação da SDK foi realizada com sucesso seguindo as instruções da página oficial, após rodar o comando acima, o binário do programa deverá estar disponível no Path do sistema operacional.
Caso contrário, basta acessar a raiz do projeto e rodar o comando para compilar:

```sh
$ go build
```

Com isso, o binário estará disponível na raiz do projeto.


### Alternativa

Caso não queira instalar e configurar o Go, os binários deste desafio *(bestroute)* foram disponibilizados para download nos links abaixo:

   - [Distribuição Linux](https://drive.google.com/open?id=1LrinL8rgqwXEBVC_ZBR72veHEHhXoqY3)
   - [Distribuição Mac OS](https://drive.google.com/open?id=1B7YbaVl5d1YKvBviKyA-AgZgl5SeNNyk)
   - [Distribuição Windows](https://drive.google.com/open?id=1Vc0CL1vSvF0n6uW-AHe8PvF1nt3S5zkj)

## Execução

### Interface de Console

Para executar a interface de console, basta executar o binário passando um arquivo de rotas como argumento:

```sh
$ bestroute input-file.txt
```

### Interface Rest

Para executar a interface rest, basta executar o binário passando a opção '-it rest' e o arquivo de rotas:

```sh
$ bestroute -it rest input-file.txt
```

Opcionalmente, a porta do servidor pode ser informada:

```sh
$ bestroute -it rest -p 8000 input-file.txt
```
*A porta padrão da interface rest é 8080*

#### Rest endpoints:

Endpoint: http://localhost:8080/api/v1/routes

**GET: Procura e retorna as melhores rotas, dado origem e destino**

Parâmetros de URL:
* from - *origem da viagem*
* to - *destino da viagem*
* limit - *opcional: quantidade de resultados a retornar*

Exemplo:

```
GET /api/v1/routes
      ?from=GRU
      &to=CDG
      &limit
```

**POST: Adiciona uma nova rota ao arquivo ou atualiza o preço caso  a rota já exista**

O Payload deve ser enviado em uma estrutura no formato json com os seguintes campos:

* from *string*
* to *string*
* price *int*

Exemplo:

```
POST /api/v1/routes

{
    "from": "ABC",
    "to": "DEF",
    "price": 70
}
```

Em caso de sucesso da requisição POST, dois status de resposta diferentes podem ser retornados:

* 201 - *Caso a nova rota tenha sido criada com sucesso*
* 200 - *Caso a rota já exista no arquivo, ela será atualizada*

## Testes Unitários

Para rodar os testes unitários basta rodar o seguinte comando da raiz do projeto:

```
$ go test ./... -v
```

## Autor

* **Eduardo Miani Ferreira**