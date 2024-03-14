package grpc

import (
	"2108a-zg5/week3/shop/framework/consul"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Client(address string) (*grpc.ClientConn, error) {
	res, err := getRpcConfig(address)
	if err != nil {
		return nil, err
	}

	ip, port, err := consul.FindConsulAddress(res.Consul, res.Name)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(fmt.Sprintf("%v:%v", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
}
