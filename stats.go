package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/aviate-labs/imp/internal/cmd"
)

const (
	api         = "https://ic-api.internetcomputer.org/api"
	defaultStep = 7200
)

var stats = cmd.Command{
	Name:    "stats",
	Aliases: []string{"s"},
	Summary: "some statistics",
	Commands: []cmd.Command{
		subnet,
	},
}

var subnet = cmd.Command{
	Name:    "subnet",
	Aliases: []string{"sn"},
	Summary: "gets subnet stats from the ICA",
	Commands: []cmd.Command{
		getDataCmd(
			"message-count", "mc",
			"returns the message count of a subnet",
			"metrics/pmessages-count",
		),
		getDataCmd(
			"blocks", "b",
			"returns the blocks of a subnet",
			"metrics/pblock",
		),
		getDataCmd(
			"memory-usage", "mu",
			"returns the memory usage of a subnet",
			"metrics/ic-memory-usage",
		),
		getDataCmd(
			"node-count", "nc",
			"returns the node count of a subnet",
			"metrics/ic-nodes-count",
		),
		getDataCmd(
			"finalization-rate", "fr",
			"returns the finalization rate of a subnet",
			"metrics/finalization-rate",
		),
		getDataCmd(
			"registered-canisters", "rc",
			"returns the amount of registered canisters in a subnet",
			"metrics/registered-canisters",
		),
		getDataCmd(
			"cycle-burn-rate", "br",
			"returns the cycle burn rate a subnet",
			"metrics/cycle-burn-rate",
		),
		getDataCmd(
			"message-execution-rate", "er",
			"returns the message execution rate a subnet",
			"metrics/message-execution-rate",
		),
	},
}

func getDataCmd(name string, alias string, summary string, endpoint string) cmd.Command {
	return cmd.Command{
		Name:    name,
		Aliases: []string{alias},
		Summary: summary,
		Args:    []string{"subnet"},
		Options: []cmd.Option{
			{Name: "start", HasValue: true},
			{Name: "end", HasValue: true},
			{Name: "step", HasValue: true},
		},
		Method: func(args []string, options map[string]string) error {
			subnet := args[0]
			s, e, step := getMessageCountParams(options)
			url := fmt.Sprintf(
				"%s/%s?subnet=%s&start=%d&end=%d&step=%d",
				api, endpoint, subnet, s, e, step,
			)
			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			raw, _ := io.ReadAll(resp.Body)
			fmt.Println(string(raw))
			return nil
		},
	}
}

func getMessageCountParams(options map[string]string) (int, int, int) {
	var (
		now   = int(time.Now().Unix())
		end   = now
		start = now - (2 * 24 * 60 * 60)
		step  = defaultStep
	)
	if e, ok := options["end"]; ok {
		if i, err := strconv.Atoi(e); err == nil {
			end = i
		}
	}
	if s, ok := options["start"]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			start = i
		}
	}
	if s, ok := options["step"]; ok {
		if i, err := strconv.Atoi(s); err == nil {
			step = i
		}
	}
	return start, end, step
}
