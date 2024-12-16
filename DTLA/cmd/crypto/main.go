package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type cliArgs struct {
	plain string
	hash  string
	cost  *int
	mode  *bool
	hex   *bool
}

func main() {
	var err error

	var args cliArgs = cliArgs{
		cost: flag.Int("c", bcrypt.DefaultCost, fmt.Sprintf("Bcrypt cost in the range [%d;%d]", bcrypt.MinCost, bcrypt.MaxCost)),
		mode: flag.Bool("m", true, "Operation mode. If true, encrypts string from first argument after the flags. If false, compares the first string (the plain text) after the flags with the second string after the flags (the hash)."),
		hex:  flag.Bool("h", false, "If in comparison mode, decodes the first input (plain text) as hexadecimal string before comparing"),
	}
	flag.Parse()

	if *args.cost < bcrypt.MinCost || *args.cost > bcrypt.MaxCost {
		fmt.Printf("Cost doesn't fit in range [%d;%d]\n", bcrypt.MinCost, bcrypt.MaxCost)
		return
	}

	if *args.mode {
		err = encrypt(&args)
	} else {
		err = cmp(&args)
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func encrypt(args *cliArgs) error {
	var err error

	args.plain = flag.Arg(0)
	if args.plain == "" {
		return errors.New("Need to pass in one string argument after the flags")
	}

	plainBytes := []byte(args.plain)
	hash, err := bcrypt.GenerateFromPassword(plainBytes, *args.cost)
	if err != nil {
		return err
	}

	fmt.Println(string(hash))

	return nil
}

func cmp(args *cliArgs) error {
	var err error

	args.plain = flag.Arg(0)
	args.hash = flag.Arg(1)
	if args.plain == "" || args.hash == "" {
		return errors.New("Need to pass in two string arguments after the flags. See --help.")
	}

	var plainBytes []byte
	if *args.hex {
		plainBytes, err = hex.DecodeString(args.plain)
		if err != nil {
			return err
		}
	} else {
		plainBytes = []byte(args.plain)
	}
	hashBytes := []byte(args.hash)
	err = bcrypt.CompareHashAndPassword(hashBytes, plainBytes)
	if err != nil {
		return err
	}

	fmt.Println("OK")

	return nil
}
