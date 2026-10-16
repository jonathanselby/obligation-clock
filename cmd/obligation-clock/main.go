package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonathanselby/obligation-clock/internal/clock"
	"os"
)

const version = "0.1.0"

func main() {
	var in string
	var v bool
	flag.StringVar(&in, "input", "", "request JSON")
	flag.BoolVar(&v, "version", false, "print version")
	flag.Parse()
	if v {
		fmt.Println("obligation-clock " + version)
		return
	}
	if in == "" {
		flag.Usage()
		os.Exit(2)
	}
	b, err := os.ReadFile(in)
	if err != nil {
		bad(err)
	}
	var r clock.Request
	if err := json.Unmarshal(b, &r); err != nil {
		bad(err)
	}
	out, err := clock.Calculate(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		bad(err)
	}
}

func bad(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(3)
}
