package config

import (
    "gitlab.com/distributed_lab/figure"
    "gitlab.com/distributed_lab/kit/comfig"
    "gitlab.com/distributed_lab/kit/kv"
    "gitlab.com/distributed_lab/logan/v3/errors"
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
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
        return &config
    }).(*EthConfig)
}
