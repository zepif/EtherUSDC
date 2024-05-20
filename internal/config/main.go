package config

import (
    "os"
    
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	getter kv.Getter
    ethRPC             string
    ethContractAddress string
    ethContractABI     string
}

func New(getter kv.Getter) Config {
	return &config{
		getter:     getter,
		Databaser:  pgdb.NewDatabaser(getter),
		Copuser:    copus.NewCopuser(getter),
		Listenerer: comfig.NewListenerer(getter),
		Logger:     comfig.NewLogger(getter, comfig.LoggerOpts{}),
        ethRPC:              os.Getenv("ETH_RPC"),
        ethContractAddress:  os.Getenv("ETH_CONTRACT_ADDRESS"),
        ethContractABI:      os.Getenv("ETH_CONTRACT_ABI"),
	}
}
