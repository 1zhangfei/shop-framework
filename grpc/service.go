package grpc

import (
	"2108a-zg5/week3/shop/framework/config"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type Conf struct {
	App struct {
		Ip   string
		Port int
	} `json:"rpc"`
	Name   string `json:"tokenName"`
	Consul string `json:"consul"`
}

func getRpcConfig(address string) (*Conf, error) {
	if err := config.ViperInit(address); err != nil {
		return nil, err
	}

	id := viper.GetString("Grpc.DataId")
	group := viper.GetString("Grpc.Group")

	res, err := config.GetConfig(id, group)
	if err != nil {
		return nil, err
	}

	var app *Conf
	if err = json.Unmarshal([]byte(res), &app); err != nil {
		return nil, err
	}

	return app, nil

}

func Service(address string, Register func(s *grpc.Server)) error {
	conf, err := getRpcConfig(address)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.App.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//支持反射机制
	reflection.Register(s)
	//健康检测
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
