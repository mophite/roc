# simple http api service

```go
func main() {
    s := roc.New(
        roc.HttpAddress("0.0.0.0:9999"),
    )

    phello.RegisterHelloWorldServer(s.Server(), &hello.Hello{})
    err := s.Run()
    if err != nil {
        rlog.Error(err)
    }
}
```

```shell
curl -H "Content-Type:application/json" -X POST -d '{"ping": "ping"}' http://127.0.0.1:9999/roc/HelloWorld/Say
```