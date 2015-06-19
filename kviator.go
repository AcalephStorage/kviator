package main

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("kviator", "command line for accessing KV store (consul, etcd, zookeper)")
	kvstore = app.Flag("kvstore", "the kv storem can be consul, etcd, or zookeper").Required().String()
	client  = app.Flag("client", "the client address").Required().String()

	put    = app.Command("put", "put data to keystore")
	putKey = put.Arg("key", "the key").Required().String()
	putVal = put.Arg("val", "the value").Required().String()

	get    = app.Command("get", "get data from keystore")
	getKey = get.Arg("key", "the key").Required().String()

	del    = app.Command("del", "delete data from keystore")
	delKey = del.Arg("key", "the key").Required().String()

	cas    = app.Command("cas", "save only when no key exists")
	casKey = cas.Arg("key", "the key").Required().String()
	casVal = cas.Arg("val", "the value").Required().String()
)

func main() {
	kingpin.Version("0.1.0")

	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	var backend store.Backend
	switch *kvstore {
	case "consul":
		backend = store.CONSUL
	case "etcd":
		backend = store.ETCD
	case "zookeper":
		backend = store.ZK
	}

	kv, err := libkv.NewStore(
		backend,
		[]string{*client},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch cmd {
	case put.FullCommand():
		save(kv, *putKey, *putVal)
	case get.FullCommand():
		retrieve(kv, *getKey)
	case del.FullCommand():
		delete(kv, *delKey)
	case cas.FullCommand():
		checkAndSave(kv, *casKey, *casVal)
	}
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
