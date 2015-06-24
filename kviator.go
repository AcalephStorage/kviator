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

	Note:

	    kviator can also read the value from Stdin. The syntax would look like this:

	        cmd | kviator ... put -
	        kviator ... put - < val.file

	    The - character is necessary to force kviator to read from Stdin. Without the -, Stdin
	    is ignored.
	`
)

var (
	kvstore string
	client  string
)

func init() {
	flag.StringVar(&kvstore, "kvstore", "", "the kvstore to connect to. Can be consul, etcd, or zookeper.")
	flag.StringVar(&client, "client", "", "the client IP address")
	flag.Usage = help
	flag.Parse()

	switch flag.Arg(0) {
	case "put":
		key := flag.Arg(1)
		val := strings.Join(flag.Args()[2:], " ")
		val = parseVal(val)
		save(key, val)
	case "get":
		key := flag.Arg(1)
		retrieve(key)
	case "del":
		key := flag.Arg(1)
		delete(key)
	case "cas":
		key := flag.Arg(1)
		val := strings.Join(flag.Args()[2:], " ")
		val = parseVal(val)
		checkAndSave(key, val)
	case "exists":
		key := flag.Arg(1)
		keyExists(key)
	default:
		help()
		os.Exit(8)
	}

}

func main() {

}

func parseVal(arg string) string {
	arg = strings.TrimSpace(arg)
	if arg == "-" {
		stat, err := os.Stdin.Stat()
		if err == nil {
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				in, _ := ioutil.ReadAll(os.Stdin)
				inStr := strings.TrimSuffix(string(in), "\n")
				inStr = strings.TrimSpace(inStr)
				return inStr
			}
		}
	}
	return arg
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

func save(key, val string) {
	kv := kvstoreConn(kvstore, client)
	err := kv.Put(key, []byte(val), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func retrieve(key string) {
	kv := kvstoreConn(kvstore, client)
	kvPair, err := kv.Get(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	} else {
		fmt.Fprintln(os.Stdout, string(kvPair.Value))
	}
}

func delete(key string) {
	kv := kvstoreConn(kvstore, client)
	err := kv.Delete(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4)
	}
}

func checkAndSave(key, val string) {
	kv := kvstoreConn(kvstore, client)
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

func keyExists(key string) {
	kv := kvstoreConn(kvstore, client)
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
