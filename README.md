# gaston

[![Tests Status](https://github.com/fullpipe/gaston/workflows/Tests/badge.svg)](https://github.com/fullpipe/gaston)
[![Docker Image](https://img.shields.io/docker/image-size/fullpipe/gaston/latest)](https://cloud.docker.com/repository/docker/fullpipe/gaston)

Gaston is json-rpc 2.0 API gateway.

## Usage

First of all, you need services with json-rpc 2.0 API which will be behind gaston. For example
```yaml
# docker-compose.yaml
version: "3"
services:
  s1:
    image: fullpipe/jmock:latest
    volumes:
      - ./examples/mocks/s1:/mocks
  s2:
    image: fullpipe/jmock:latest
    volumes:
      - ./examples/mocks/s2:/mocks
```

Next we add gaston as gateway fo this services

```yaml
# docker-compose.yaml
version: "3"
services:
  ...
  gaston:
    image: fullpipe/gaston:latest
    volumes:
      - ./examples/methods:/methods
    environment:
      GASTON_JWT_HMACSECRET: qwertyuiopasdfghjklzxcvbnm123456
    ports:
      - 8080:8080
```

### methods.json example

```jsonc
[
  {
    "host": "http://user:8080/rpc", // hidden service url
    "name": "user.get", // method name for gaston clients
    "remoteName": "getUser", // method name on hidden service
    "roles": ["ROLE_USER"], // client has to have at least on of the roles, to get access to hidden service
    "paramConverters": [ // paramConverters do some work on client params
      { "type": "rename", "from": "userId", "to": "id" }
    ],
    "resultConverters": [ // resultConverters do some work on hidden service results
      { "type": "rename", "from": "firstname", "to": "name" }
      { "type": "overwrite", "name": "age", "newValue": 18 }
    ]
  },
  {
    "host": "http://email:8080/rpc",
    "name": "foo",
    "remoteName": "bar",
    "paramConverters": [
      { "type": "rename", "from": "email_input", "to": "email" }
    ]
  }
]
```

#### user.get flow

##### Client to Gaston

```jsonc
{
	"jsonrpc": "2.0",
	"method": "user.get",
	"id": 1,
	"params": {
		"userId": 123
	}
}
```

##### Gaston to Service

```jsonc
{
	"jsonrpc": "2.0",
	"method": "getUser", // method name converted to remoteName
	"id": 1,
	"params": {
		"id": 123 // userId was converted to id by paramConverter
	}
}
```

##### Service to Gaston
```jsonc
{
  "jsonrpc": "2.0",
  "result": {
    "id": 123,
    "firstname": "John",
    "lastname": "Doe",
    "age": 99
  },
  "id": 1
}
```

##### Gaston to Client

```jsonc
{
  "jsonrpc": "2.0",
  "result": {
    "id": 123,
    "name": "John", // firstname was converted to name by resultConverter
    "lastname": "Doe",
    "age": 18 // age was overwritten to 18 by resultConverter
  },
  "id": 1
}
```

## Converters

### rename

Renames param key

```json
{ "type": "rename", "from": "oldName", "to": "newName" }
```

### overwrite

Overwrite param/result by name. If param/result with name not exists, does nothing.

```json
{ "type": "overwrite", "name": "age", "newValue": 18 }
```

### snake_case

Transform param name to its snake_cased version

```jsonc
{ "type": "snakeCase", "name": "fooBar" } // will convert to `foo_bar`
// or
{ "type": "snake_case", "name": "fooBar" }
```

### remove

Removes param by name if exists

```jsonc
{ "type": "remove", "name": "fooBar" }
// or
{ "type": "delete", "name": "fooBar" }
```

### castNumber

Cast param value to number

```jsonc
{ "type": "castNumber", "name": "userId" } // { "userId: "123" } -> { "userId: 123 }
{ "type": "castNumber", "name": "price" } // { "price: "1.333" } -> { "userId: 1.333 }
```

```
true -> 1
false -> 0
null -> 0
empty string -> 0
"1" -> 1
"1.1" -> 1.1
other -> null
```

## Env var

Here available env vars, with their default values

```yaml
GASTON_SERVER_ROUTE: / # route, http://gaston/, http://gaston/v1, http://gaston/v2
GASTON_SERVER_PORT: 8080
GASTON_SERVER_METHODSPATH: "/methods/*.json" # glob path to lookup methods
GASTON_SERVER_REMOTETIMEOUT: 5 # timeout for remote requests
GASTON_JWT_HEADER: Authorization # header name to read JWT token from
GASTON_JWT_SCHEME: Bearer # Authorization scheme
GASTON_JWT_HMACSECRET: # HmacSecret to validate JWT token
GASTON_JWT_ROLESCLAIM: roles # claim key to get user roles
GASTON_JWT_USERCLAIM: sub # claim key to get user id
GASTON_JWT_REMOTEUSERHEADER: X-Verified-User # header name to pass user id to "hidden" services
GASTON_JWT_REMOTEROLESHEADER: X-Verified-Roles # header name to pass user roles to "hidden" services
```

## todo
- [ ] More converters
  - [x] rename, rename param key
  - [x] overwrite, overwrite param value if param exists
  - [x] snakeCase, convert param name to snake_case. userId -> user_id
  - [ ] cammelCase, convert param name to cammelCase. user_id -> userId
  - [x] remove, remove param by name
  - [ ] default, setup param if not exists
  - [ ] set, overwrite + default
  - [x] castNumber, cast value to number
  - [ ] castString, cast value to string
  - [ ] castBoolean, cast value to boolean
  - [ ] nullFix, convert string "null" to real null
- [ ] Examples
