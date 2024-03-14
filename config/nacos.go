package config

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var Client config_client.IConfigClient

func GetClient(address string) error {
	err := ViperInit(address)
	if err != nil {
		return err
	}
	ip := viper.GetString("Nacos.Ip")
	port := viper.GetInt("Nacos.port")

	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         "", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: ip,
			Port:   uint64(port),
		},
	}

	// 创建动态配置客户端的另一种方式 (推荐)
	Client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(DataId, Group string) (string, error) {

	res, err := Client.GetConfig(vo.ConfigParam{
		DataId: DataId,
		Group:  Group,
	})
	if err != nil {
		return "", err
	}
	return res, nil
}
