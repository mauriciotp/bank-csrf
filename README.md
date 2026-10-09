# bank-csrf

API REST de um banco simples, escrita em Go, desenvolvida como desafio de estudo na Rocketseat. Permite criar contas para pessoas físicas e jurídicas, consultar saldo, depositar, sacar, transferir entre contas e encerrar contas.

## Tecnologias

- [Go](https://go.dev/) com [chi](https://github.com/go-chi/chi) para o roteamento HTTP
- [PostgreSQL](https://www.postgresql.org/) via [pgx](https://github.com/jackc/pgx)
- [sqlc](https://sqlc.dev/) para gerar código Go tipado a partir das queries SQL
- [tern](https://github.com/jackc/tern) para as migrations
- [Air](https://github.com/air-verse/air) para live reload em desenvolvimento
- [sqlfluff](https://sqlfluff.com/) para lint de SQL
- Docker Compose para subir o banco localmente

## Estrutura

```
cmd/
  api/            # ponto de entrada do servidor HTTP
  terndotenv/     # wrapper do tern que carrega o .env antes de rodar as migrations
internal/
  api/            # handlers e rotas HTTP
  services/       # regras de negócio (AccountsService)
  usecase/account # structs de request e suas validações
  validator/      # helpers de validação reutilizáveis
  jsonutils/      # encode/decode de JSON com validação
  store/pgstore/  # migrations, queries SQL e código gerado pelo sqlc
```

## Como rodar

### Pré-requisitos

- Go
- Docker e Docker Compose
- `tern`: `go install github.com/jackc/tern/v2@latest`
- `sqlc` (apenas se for alterar as queries)
- `air` (opcional, para live reload)

### 1. Variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
PORT=3080
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=postgres
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
```

### 2. Banco de dados

```sh
docker compose up -d
go run ./cmd/terndotenv
```

O segundo comando aplica as migrations, que criam as tabelas e cadastram as categorias de conta.

### 3. Servidor

```sh
go run ./cmd/api
# ou, com live reload:
air
```

### Regenerar as queries

Depois de alterar algum arquivo em `internal/store/pgstore/queries`:

```sh
cd internal/store/pgstore && sqlc generate
```

## Endpoints

Todas as rotas ficam sob `/accounts`.

| Método   | Rota                      | Descrição                     |
| -------- | ------------------------- | ----------------------------- |
| `POST`   | `/accounts`               | Cria uma conta                |
| `GET`    | `/accounts/{id}/balance`  | Consulta o saldo              |
| `POST`   | `/accounts/{id}/deposit`  | Deposita um valor             |
| `POST`   | `/accounts/{id}/withdraw` | Saca um valor                 |
| `POST`   | `/accounts/transfer`      | Transfere entre duas contas   |
| `DELETE` | `/accounts/{id}`          | Encerra a conta (saldo zero)  |

### Criar conta

```sh
curl -X POST localhost:3080/accounts \
  -H 'Content-Type: application/json' \
  -d '{
    "type": "natural_person",
    "name": "Maria Silva",
    "age": 30,
    "email": "maria@example.com",
    "income": 8000,
    "mobile_phone": "+5511999999999"
  }'
```

`type` aceita `natural_person` (pessoa física) ou `legal_entity` (pessoa jurídica). O e-mail precisa ser único.

A categoria da conta é definida pela renda informada:

| Tipo            | Renda                   | Categoria        |
| --------------- | ----------------------- | ---------------- |
| Pessoa física   | até 4.999,99 / mês      | `standard`       |
| Pessoa física   | 5.000 a 19.999,99 / mês | `premium`        |
| Pessoa física   | a partir de 20.000      | `private`        |
| Pessoa jurídica | até 359.999,99 / ano    | `micro_business` |
| Pessoa jurídica | 360.000 a 4.799.999,99  | `small_business` |
| Pessoa jurídica | a partir de 4.800.000   | `corporate`      |

### Depósito e saque

```sh
curl -X POST localhost:3080/accounts/{id}/deposit \
  -H 'Content-Type: application/json' \
  -d '{"amount": 150.00}'
```

O saque usa o mesmo corpo em `/accounts/{id}/withdraw`.

### Transferência

```sh
curl -X POST localhost:3080/accounts/transfer \
  -H 'Content-Type: application/json' \
  -d '{
    "sender_id": "<uuid-remetente>",
    "recipient_id": "<uuid-destinatario>",
    "amount": 50.00
  }'
```

## Conceitos praticados

- **Camadas separadas:** handlers HTTP, serviço com as regras de negócio e acesso a dados gerado pelo sqlc.
- **Validação de entrada:** cada request implementa a interface `Validator`. Erros de campo retornam `422` com um mapa `campo → mensagem`.
- **Transações:** a criação de conta (titular + conta) e a transferência (saque + depósito) são atômicas.
- **Concorrência:** a transferência trava as duas contas com `SELECT ... FOR UPDATE` ordenado por `id`, o que evita deadlock entre transferências simultâneas em sentidos opostos.
- **Integridade no banco:** constraints `CHECK` garantem saldo nunca negativo e que cada conta tenha exatamente um titular (pessoa física ou jurídica).

## Melhorias futuras

- Diferenciar "saldo insuficiente" de "conta não encontrada" no saque, usando a constraint `chk_balance_non_negative`.
- Ignorar contas encerradas no lock da transferência.
- Retornar `404` para conta inexistente e `422` para saldo insuficiente.
- Representar valores monetários sem `float64` (string decimal ou centavos em `int64`).
- Proteção CSRF.
- Testes automatizados.
