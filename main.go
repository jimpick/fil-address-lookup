package main

import (
	"context"
	"fmt"
	"log"

	jsonrpc "github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
)

func main() {
	url := "https://api.node.glif.io/rpc/v1"

	var api lotusapi.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), url, "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, nil)
	if err != nil {
		log.Fatalf("connecting with lotus failed: %s", err)
	}
	defer closer()

	tipset, err := api.ChainHead(context.Background())
	if err != nil {
		log.Fatalf("calling chain head: %s", err)
	}
	fmt.Printf("Current chain head is: %d %s\n", tipset.Height(), tipset.String())
}
