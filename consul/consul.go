package consul

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"net"
)

func RegisterConsul(ConsulPort, port int64, name string) error {
	ip := GetIp()
	client, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%v:%v", ip[0], ConsulPort)})
	if err != nil {
		return err
	}

	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      uuid.NewString(),
		Name:    name,
		Tags:    []string{"GRPC"},
		Port:    int(port),
		Address: ip[0],
		Check: &api.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%v:%v", ip[0]),
			DeregisterCriticalServiceAfter: "10s",
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func FindConsulAddress(port, name string) (string, int, error) {
	ip := GetIp()
	client, err := api.NewClient(&api.Config{Address: fmt.Sprintf("%v:%v", ip[0], port)})
	if err != nil {
		return "", 0, err
	}

	byName, data, err := client.Agent().AgentHealthServiceByName(name)
	if err != nil {
		return "", 0, err
	}

	if byName == "passing" {
		return "", 0, errors.New("没有健康的服务")
	}

	return data[0].Service.Address, data[0].Service.Port, nil

}

func GetIp() (ip []string) {
	// 获取计算机上的所有网络接口地址
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	// 遍历每个网络接口地址
	for _, addr := range addrs {
		// 将接口地址转换为 IPNet 类型
		ipNet, isVailIpNet := addr.(*net.IPNet)
		// 判断是否为有效的 IPNet 类型，并且不是环回地址
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			// 判断是否为 IPv4 地址
			if ipNet.IP.To4() != nil {
				// 将 IPv4 地址转换为字符串，并添加到 ip 切片中
				ip = append(ip, ipNet.IP.String())
			}
		}
	}
	// 返回存储 IPv4 地址的切片
	return ip
}
