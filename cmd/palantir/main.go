package main

import (
	"fmt"
	"os"
	"github.com/alecthomas/kingpin"
	"github.com/ravishi/palantir/server"
)

var (
	debug    = kingpin.Flag("debug", "Run in debug mode").Short('d').Default("false").Bool()
	address  = kingpin.Flag("address", "TCP address to listen on").Short('a').Default(":8080").String()
)

func main() {
	kingpin.Parse()

	s := server.New(&server.ServerConfig{
		Debug: *debug,
	})

	if err := s.Start(*address); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
