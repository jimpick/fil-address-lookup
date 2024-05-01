package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

const url = "https://api.node.glif.io/rpc/v1"

func getAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx, url, "Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal}, nil)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	tipset, err := api.ChainHead(ctx)
	if err != nil {
		log.Fatalf("calling chain head: %s", err)
	}
	fmt.Printf("Current chain head is: %d %s\n", tipset.Height(), tipset.String())
	io.WriteString(w, "This is my website!\n")
}

func main() {
	http.HandleFunc("/", getAddress)

	err := http.ListenAndServe(":3000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
