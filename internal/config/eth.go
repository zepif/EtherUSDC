package config

import (
    "gitlab.com/distributed_lab/kit/comfig"
    "gitlab.com/distributed_lab/kit/kv"
)

type EthConfiger interface {
    EthConfig() *EthConfig
}

type EthConfig struct {
    EthRPC             string
    EthContractAddress string
    EthContractABI     string
}

func NewEthConfiger(getter kv.Getter) EthConfiger {
    return &ethConfig{
        getter: getter,
    }
}

type ethConfig struct {
    getter kv.Getter
    once   comfig.Once
}

func (c *ethConfig) EthConfig() *EthConfig {
    return c.once.Do(func() interface{} {
        config := &EthConfig{}
        raw := kv.MustGetStringMap(c.getter, "eth")
        config.EthRPC = raw["ETH_RPC"]
        config.EthContractAddress = raw["ETH_CONTRACT_ADDRESS"]
        config.EthContractABI = raw["ETH_CONTRACT_ABI"]
        return config
    }).(*EthConfig)
}
