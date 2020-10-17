## Getting Started

- Environment variables

```sh
make init
```

- Starting API

```sh
make up
```

- View logs

```sh
make logs
```

## API Request

| Endpoint        | HTTP Method           | Description       |
| --------------- | :---------------------: | :-----------------: |
| `/v1/accounts` | `POST`                | `Create account` |
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