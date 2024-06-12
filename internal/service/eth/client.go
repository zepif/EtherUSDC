package eth

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zepif/EtherUSDC/internal/config"
	store "github.com/zepif/EtherUSDC/internal/store"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type EthClient struct {
	Client   *ethclient.Client
	Contract *store.Store
	log      *logan.Entry
}

func NewEthClient(cfg config.Config) (*EthClient, error) {
	rpcURL := fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", cfg.EthConfig().EthRPC)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Ethereum client")
	}

	contractAddress := common.HexToAddress(cfg.EthConfig().EthContractAddress)
	contract, err := store.NewStore(contractAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate a Store contract")
	}

	return &EthClient{
		Client:   client,
		Contract: contract,
		log:      cfg.Log().WithField("component", "eth_client"),
	}, nil
}

func (e *EthClient) LoadBlocks(ctx context.Context, logs chan<- types.Log, startBlock uint64, endBlock uint64) error {
	opts := &bind.FilterOpts{
		Start:   startBlock,
		End:     &endBlock,
		Context: ctx,
	}

	events, err := e.Contract.FilterTransfer(opts, nil, nil)
	if err != nil {
		return errors.Wrap(err, "failed to filter transfer events")
	}

	eventCount := 0
	for events.Next() {
		event := events.Event
		logs <- types.Log{
			Address:     event.Raw.Address,
			Topics:      event.Raw.Topics,
			Data:        event.Raw.Data,
			BlockNumber: event.Raw.BlockNumber,
			TxHash:      event.Raw.TxHash,
			Index:       event.Raw.Index,
		}
		eventCount++
	}

	if events.Error() != nil {
		e.log.WithError(events.Error()).Error("Error encountered while filtering transfer events")
		return errors.Wrap(events.Error(), "error encountered while filtering transfer events")
	}

	e.log.WithFields(logan.F{
		"startBlock": startBlock,
		"endBlock":   endBlock,
		"eventCount": eventCount,
	}).Info("Finished loading blocks")

	return nil
}

func (e *EthClient) ParseTransferEvent(vLog types.Log) (*store.StoreTransfer, error) {
	event, err := e.Contract.ParseTransfer(vLog)
	if err != nil {
		e.log.WithError(err).Error("Failed to parse transfer event")
		return nil, errors.Wrap(err, "failed to parse transfer event")
	}

	e.log.WithFields(logan.F{
		"from":  event.From.Hex(),
		"to":    event.To.Hex(),
		"value": event.Value.String(),
	}).Info("Parsed transfer event successfully")

	return event, nil
}
