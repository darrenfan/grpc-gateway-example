1. with curl
```
curl http://localhost:50001/example/echo
```

2. with grpc_cli

```
grpc_cli call localhost:50001 example.Example.Echo "value:'test'"
```
