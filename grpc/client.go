package grpc

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

func Client(address string) (*grpc.ClientConn, error) {
	res, err := getRpcConfig(address)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(fmt.Sprintf("consul://%v:%v/%v?wait=14s", res.App.Ip, res.Consul, res.Name), grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "round_robin"}`))
}
