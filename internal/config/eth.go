package config

import (
    "fmt"
    "gitlab.com/distributed_lab/figure"
    "gitlab.com/distributed_lab/kit/comfig"
    "gitlab.com/distributed_lab/kit/kv"
    "gitlab.com/distributed_lab/logan/v3/errors"
)

type EthConfiger interface {
    EthConfig() *EthConfig
}

type EthConfig struct {
    EthRPC             string `fig:"eth_rpc"`
    EthContractAddress string `fig:"eth_contract_address"`
    EthContractABI     string `fig:"eth_contract_abi"`
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
        raw := kv.MustGetStringMap(c.getter, "eth")
        fmt.Printf("Raw Config: %+v\n", raw)
        config := &EthConfig{}
		err := figure.Out(config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
        fmt.Printf("Decoded EthConfig: %+v\n", config)
        return config
    }).(*EthConfig)
}
