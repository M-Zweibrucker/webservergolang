# Desafio Client-Server-API

Sistema de cotação de dólar em Go com requisitos de timeout e persistência.

## Requisitos

- Go 1.16 ou superior
- SQLite3

## Instalação

```bash
# Instalar dependências
go get github.com/mattn/go-sqlite3
go mod tidy
```

## Como usar

### 1. Iniciar o servidor

```bash
go run server.go
```

O servidor irá:
- Rodar na porta 8080
- Expor o endpoint `/cotacao`
- Consumir a API de cotação em https://economia.awesomeapi.com.br/json/last/USD-BRL
- Timeout de 200ms para chamar a API
- Timeout de 10ms para persistir no banco de dados SQLite
- Salvar cada cotação no arquivo `cotacoes.db`

### 2. Executar o cliente

Em outro terminal:

```bash
go run client.go
```

O cliente irá:
- Fazer requisição HTTP para o servidor em http://localhost:8080/cotacao
- Timeout de 300ms para receber resposta do servidor
- Salvar a cotação no arquivo `cotacao.txt` no formato: `Dólar: {valor}`

## Arquivos gerados

- `cotacoes.db` - Banco de dados SQLite com histórico de cotações
- `cotacao.txt` - Arquivo com a última cotação no formato: `Dólar: {valor}`

## Timeouts implementados

Conforme especificado no desafio:

1. **Server -> API Externa**: 200ms
2. **Server -> Banco de dados**: 10ms
3. **Client -> Server**: 300ms

Todos os contextos retornam erro nos logs caso o tempo seja insuficiente.

## Estrutura do projeto

```
.
├── client.go       # Cliente HTTP
├── server.go       # Servidor HTTP + API + Banco
├── go.mod          # Gerenciamento de dependências
├── go.sum          # Checksums das dependências
└── README.md       # Este arquivo
```

