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

```json
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

```json
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

```json
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
```json
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

```json
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
- [ ] Docs
  - [ ] config
  - [ ] converter
  - [ ] remote
  - [ ] server
- [ ] More converters
  - [x] rename, rename param key
  - [x] overwrite, overwrite param value if param exists
  - [ ] snakeCase, convert param name to snake_case. userId -> user_id
  - [ ] cammelCase, convert param name to cammelCase. user_id -> userId
  - [ ] remove, remove param by name
  - [ ] default, setup param if not exists
  - [ ] set, overwrite + default
  - [ ] castNumber, cast value to number
  - [ ] castString, cast value to string
  - [ ] castBoolean, cast value to boolean
- [ ] Examples
