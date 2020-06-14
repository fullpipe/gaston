# gaston

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

- [x] restructure to pkg
- [x] config.MethodsFromJson
- [ ] server.Builder
    - [ ] NewServer
    - [ ]
- [ ] server.NewServer
- [x] remote.Remote
- [x] remote.Method
- [x] remote.Collection
- [x] remote.Middleware
- [x] converter...
- [ ] middleware.NewAuthenticationMiddleware
- [ ] converters: snakeCase, cammelCase, remove, default, set, castNumber,
  castString, castBoolean
