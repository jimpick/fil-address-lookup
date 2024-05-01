package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/filecoin-project/go-address"
	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
)

const url = "https://api.node.glif.io/rpc/v1"

func getAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("Request: %+v\n", r)

	addr := r.URL.Path[1:]
	fmt.Printf("Address: %+v\n", addr)

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx, url, "Filecoin",
		[]interface{}{&api.Internal, &api.CommonStruct.Internal}, nil)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	var queryAddr address.Address

	ethAddr, err := ethtypes.ParseEthAddress(addr)
	if err == nil {
		// Eth address, Get Agent delegated address
		queryAddr, err = ethAddr.ToFilecoinAddress()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// It might be a Fil address
		queryAddr, err = address.NewFromString(addr)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Query addr: %v\n", queryAddr)

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
