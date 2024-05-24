package service

import (
	"net"
	"net/http"

	"github.com/zepif/EtherUSDC/internal/config"
    "github.com/zepif/EtherUSDC/internal/data/pg"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
    "github.com/zepif/EtherUSDC/internal/service/eth"
    "github.com/zepif/EtherUSDC/internal/service/workers"
    //"github.com/zepif/EtherUSDC/internal/service/handlers"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
}

func (s *service) run(cfg config.Config) error {
	s.log.Info("Service started")
	r := s.router(cfg)

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) {
    log := cfg.Log()
    db := pg.NewMasterQ(cfg.DB())
    
    ethConfig := cfg.EthConfig()
    ethClient, err := eth.NewEthClient(ethConfig.EthRPC, ethConfig.EthContractAddress, ethConfig.EthContractABI)
    if err != nil {
        log.WithError(err).Fatal("failed to create Ethereum client")
    }

    transactionWorker := workers.NewTransactionWorker(log, db, ethClient)
    err = transactionWorker.Start()
    if err != nil {
        log.WithError(err).Fatal("failed to start transaction worker")
    }
    defer transactionWorker.Stop()


	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}
}
