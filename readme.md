Text Generator
===

# How to use

```shell
$ go run main.go --help

-format string
    input result type ("text", "json", "bytes", "hex", "base64") (default "text")
-length int
    input password length (default 128)
-batch int
    input batch size (default 1) 
```

# TODO

- [x] Generate with length (default: 128)
- [x] Generate with result format (default: text)
- [x] Prevent nearly duplicate
- [x] Add unit test
