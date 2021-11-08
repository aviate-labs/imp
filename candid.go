package main

import (
	"encoding/hex"
	"fmt"

	"github.com/aviate-labs/candid-go"
	"github.com/aviate-labs/imp/internal/cmd"
)

var candidCommand = cmd.Command{
	Name:    "candid",
	Summary: "candid utility tools",
	Commands: []cmd.Command{
		decodeCommand,
		encodeCommand,
	},
}

var decodeCommand = cmd.Command{
	Name:    "decode",
	Summary: "decode candid values",
	Args:    []string{"value"},
	Method: func(args []string, _ map[string]string) error {
		v, err := hex.DecodeString(args[0])
		if err != nil {
			return err
		}
		d, err := candid.DecodeValue(v)
		if err != nil {
			return err
		}
		fmt.Println(d)
		return nil
	},
}

var encodeCommand = cmd.Command{
	Name:    "encode",
	Summary: "encode candid values",
	Args:    []string{"value"},
	Method: func(args []string, _ map[string]string) error {
		v := args[0]
		e, err := candid.EncodeValue(v)
		if err != nil {
			return err
		}
		fmt.Printf("%x\n", e)
		return nil
	},
}
