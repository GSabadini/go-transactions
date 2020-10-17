## Requirements/dependencies
- Docker
- Docker-compose

## Getting Started

- Start application

```sh
make start
```

- Run tests in container

```sh
make test
```

- Run tests local (it is necessary to have golang installed)

```sh
make test-local
```

- View logs

```sh
make logs
```

- Stop application

```sh
make down
```

## API Endpoint

| Endpoint        | HTTP Method           | Description       |
| --------------- | :---------------------: | :-----------------: |
| `/v1/accounts` | `POST`                | `Create account` |
| `/v1/accounts` | `GET`                | `Find account by ID` |
| `/v1/health`| `GET`                 | `Health check`  |


#### Test endpoints API using curl

- Creating new account

`Request`
```bash
curl -i --request POST 'http://localhost:3001/v1/accounts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "document": {
        "number": "12345678900"
    }
}'
```

`Response`
```json
{
    "id": "1a4028ea-3c18-4714-b650-d1058ae7a053",
    "document": {
        "number": "12345678900"
    },
    "created_at": "2020-10-17T02:28:05Z"
}
```

- Find account by ID

`Request`
```bash
curl -i --request GET 'http://localhost:3001/v1/accounts/1a4028ea-3c18-4714-b650-d1058ae7a053'
```

`Response`
```json
{
    "id": "1a4028ea-3c18-4714-b650-d1058ae7a053",
    "document": {
        "number": "12345678900"
    },
    "created_at": "2020-10-17T02:28:05Z"
}
```