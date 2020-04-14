# Examples

Here are example how to use the h2conn package. This project here is just for test.

## get certificate

```
// openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt
make certs
```

## [push example](./push)

Run server:

```
make push-server.go
```

Run client:

```
make push-client.go
```

Simple connection test

```
$ go run server.go 
2020/04/10 21:09:37 Serving on https://localhost:8000
```

```
$ go run client.go 
Got response 200: HTTP/2.0 Hello
```

```
$ go run client.go -version 1
Got response 200: HTTP/1.1 Hello
```

Server push

```
$ go run client.go -version 1
Got response 200: HTTP/2.0 Hello
```

```
Got connection: HTTP/1.1
Handling 1st
Can't push to client
```

```
$ go run client.go
Got response 200: HTTP/2.0 Hello
```

```
Got connection: HTTP/2.0
Handling 1st
Got connection: HTTP/2.0
Handling 2nd
```

## [echo example](./chat)

Run server:

    make echo-server

Run client:

    make echo-client

## [chat example](./chat)

Run server:

    make chat-server

Run client:

    make chat-client
