package eth

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zepif/EtherUSDC/internal/config"
	store "github.com/zepif/EtherUSDC/internal/store"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var TransferEventID common.Hash

func init() {
	const StoreABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`
	parsedABI, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		panic(err)
	}
	TransferEventID = parsedABI.Events["Transfer"].ID
}

type EthClient struct {
	Client          *ethclient.Client
	Contract        *store.Store
	Address         common.Address
	TransferEventID common.Hash
	log             *logan.Entry
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

	parsedABI, err := abi.JSON(strings.NewReader(store.StoreMetaData.ABI))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse contract ABI")
	}
	transferEventID := parsedABI.Events["Transfer"].ID

	return &EthClient{
		Client:          client,
		Contract:        contract,
		Address:         contractAddress,
		TransferEventID: transferEventID,
		log:             cfg.Log().WithField("component", "eth_client"),
	}, nil
}

func (e *EthClient) SubscribeLogs(ctx context.Context, logs chan<- types.Log) (ethereum.Subscription, error) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.Address},
		Topics:    [][]common.Hash{{e.TransferEventID}},
	}

	sub, err := e.Client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe to logs")
	}

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	return sub, nil
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
