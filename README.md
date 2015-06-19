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

## TODO

- Better README
- Needs docs