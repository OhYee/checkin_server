package main

import (
	"github.com/OhYee/blotter/http"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/checkin-server/api"
)

const (
	addr   = "0.0.0.0:50001"
	prefix = "/api/checkin/"
)

func main() {
	api.Register()
	output.Log("Server will start at http://%s", addr)
	if err := http.Server(addr, prefix); err != nil {
		output.Err(err)
	}
}
