# kviator
kviator is a small tool to manage data across distributed key value stores (consul, etcd)

## Build

```
$ go get github.com/AcalephStorage/kviator
$ cd $GOPATH/src/github.com/AcalephStorage/kviator
$ go build
```

## Running

```
$ ./kviator --help

for consul:

$ ./kviator --kvstore=consul --client=localhost:8500 put hello world
$ ./kviator --kvstore=consul --client=localhost:8500 get hello
#=> world
```

## TLS support

```
./kviator --kvstore=etcd --client=127.0.0.1:4001 --ca-cert=/path/to/ssl/ca-cert.pem --client-cert=/path/to/ssl/client-cert.pem --client-key=/path/to/ssl/client-key.pem
```

## TODO

- Better README
- Needs docs