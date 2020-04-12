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
