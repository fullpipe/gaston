# gaston

[![Tests Status](https://github.com/fullpipe/gaston/workflows/Tests/badge.svg)](https://github.com/fullpipe/gaston)

```yaml
host: ok.ru/rpc
version: v2
name: user.get
rename: get
convertParams:
    - rename:
        from: userId
        to: id
    - rename: [userId, id]
    - snakeCase: user_id
    - remove: user_id
    - castNumber: user_id
convertResult:
    - cammelCase: firstName
    - remove: first_name
    - castString: first_name
```

```json
{
    "host": "host.ru/rpc",
    "name": "user.get",
    "rename": "get",
    "convertParams": [
        { "type": "rename", "from": "userId", "to": "id" },
        { "type": "snakeCase", "name": "user_id" },
        { "type": "remove", "name": "user_id" },
        { "type": "castNumber", "name": "user_id" },
        { "type": "setValue", "name": "user_id", "value": 2 },
    ]
```

### todo

- [ ] More converters
  - [x] rename, rename param key
  - [x] overwrite, overwrite param value if param exists
  - [ ] snakeCase, convert param key to snake_case. userId -> user_id
  - [ ] cammelCase, convert param key to cammelCase. user_id -> userId
  - [ ] remove, remove param by key
  - [ ] default, setup param if not exists
  - [ ] set, overwrite + default
  - [ ] castNumber, cast value to number
  - [ ] castString, cast value to string
  - [ ] castBoolean, cast value to boolean
