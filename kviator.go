package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/mgutz/str"
)

const (
	appName = "kviator"
	version = "0.0.3"

	helpText = `
	kviator is a cli client for accessing consul, etcd, or zookeper KV.

	Syntax:

	    kviator --kvstore [consul|etcd|zzookeper] --client <kv_addr> <command> <key> [<val>]

	Options:
	    --kvstore     The kvstore to connect to. Can be consul, etcd, or zookeper.
	    --client      The url of the kvstore. (eg. localhost:8500)

	Commands:
	    put           put a key value pair in the kvstore
	    get           retrieve a key value pair from the kvstore
	    del           removes a key value pair from the kvstore
	    cas           put a key value pair in the keystore only when it's empty
	    exists        returns true when key value pair exists

	Arguments:
	    key           The key. Required for all commands.
	    val           The value. required for put and cas.
	`
)

var (
	kvstore string
	client  string
)

func init() {
	// append STDIN data to ARGV
	stat, err := os.Stdin.Stat()
	if err == nil {
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			extraArgs, _ := ioutil.ReadAll(os.Stdin)
			args := strings.Replace(string(extraArgs), "\n", "", -1)
			argv := str.ToArgv(args)
			os.Args = append(os.Args, argv...)
		}
	}

	flag.StringVar(&kvstore, "kvstore", "", "the kvstore to connect to. Can be consul, etcd, or zookeper.")
	flag.StringVar(&client, "client", "", "the client IP address")
	flag.Usage = help
	flag.Parse()

	kv := kvstoreConn(kvstore, client)

	switch flag.Arg(0) {
	case "put":
		key := flag.Arg(1)
		val := flag.Arg(2)
		save(kv, key, val)
	case "get":
		key := flag.Arg(1)
		retrieve(kv, key)
	case "del":
		key := flag.Arg(1)
		delete(kv, key)
	case "cas":
		key := flag.Arg(1)
		val := flag.Arg(2)
		checkAndSave(kv, key, val)
	case "exists":
		key := flag.Arg(1)
		keyExists(kv, key)
	default:
		help()
		os.Exit(8)
	}

}

func main() {

}

func kvstoreConn(kvstore, client string) store.Store {
	var backend store.Backend
	switch kvstore {
	case "consul":
		backend = store.CONSUL
	case "etcd":
		backend = store.ETCD
	case "zookeper":
		backend = store.ZK
	}
	kv, err := libkv.NewStore(
		backend,
		[]string{client},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return kv
}

func save(kv store.Store, key, val string) {
	err := kv.Put(key, []byte(val), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func retrieve(kv store.Store, key string) {
	kvPair, err := kv.Get(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	} else {
		fmt.Fprintln(os.Stdout, string(kvPair.Value))
	}
}

func delete(kv store.Store, key string) {
	err := kv.Delete(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
}

func checkAndSave(kv store.Store, key, val string) {
	_, err := kv.Get(key)
	if err != nil {
		err := kv.Put(key, []byte(val), nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(6)
		}
	} else {
		fmt.Fprintln(os.Stderr, "key is already set")
		os.Exit(5)
	}
}

func keyExists(kv store.Store, key string) {
	_, err := kv.Get(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "false")
		os.Exit(7)
	} else {
		fmt.Fprintln(os.Stdout, "true")
	}
}

func help() {
	fmt.Fprintf(os.Stdout, "\n\t%s %s\n\n%s\n", appName, version, helpText)
}
