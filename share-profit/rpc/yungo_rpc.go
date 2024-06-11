package rpc

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/client"
	"log"
	"net/http"
	"share-profit/conf"
)

type YungoRpc struct {
	Rpc    api.FullNode
	Closer jsonrpc.ClientCloser
}

//云构LOTUS api
//func NewYungoRpc(lotus *conf.Lotus) func()(*YungoRpc,error) {
//
//	return func() (*YungoRpc,error){
//
//		headers := http.Header{"Authorization": []string{"Bearer " + lotus.Token}}
//		var Apis  apistruct.FullNodeStruct
//		closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+lotus.Host+"/rpc/v0", "Filecoin", []interface{}{&Apis.Internal, &Apis.CommonStruct.Internal}, headers)
//		if err != nil {
//			log.Println("connecting with lotus failed: %s", err)
//			return nil,err
//		}
//		return &YungoRpc{
//			Rpc:Apis,
//			Closer:closer,
//		},nil
//	}
//
//
//}

func NewYungoRpc(lotus *conf.Lotus) func() (*YungoRpc, error) {

	return func() (*YungoRpc, error) {

		headers := http.Header{"Authorization": []string{"Bearer " + lotus.Token}}
		//var Apis  apistruct.FullNodeStruct
		//closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+lotus.Host+"/rpc/v0", "Filecoin", []interface{}{&Apis.Internal, &Apis.CommonStruct.Internal}, headers)
		cl, stop, err := client.NewFullNodeRPCV1(context.Background(), "ws://"+lotus.Host+"/rpc/v0", headers)
		if err != nil {
			log.Println("connecting with lotus failed: %s", err)
			return nil, err
		}
		return &YungoRpc{
			Rpc:    cl,
			Closer: stop,
		}, nil
	}

}
