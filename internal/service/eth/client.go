package eth

import (
	"context"
	"fmt"
	"math/big"
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

type EthClient struct {
	Client          *ethclient.Client
	ContractAbi     abi.ABI
	ContractAddress common.Address
	log             *logan.Entry
}

func NewEthClient(cfg config.Config) (*EthClient, error) {
	rpcURL := fmt.Sprintf("https://mainnet.infura.io/v3/%s", cfg.EthConfig().EthRPC)
	// fmt.Println(rpcURL)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to Ethereum client")
	}

	contractAbi, err := abi.JSON(strings.NewReader(store.StoreMetaData.ABI))
	// fmt.Println(contractAbi)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse contract ABI")
	}

	return &EthClient{
		Client:          client,
		ContractAbi:     contractAbi,
		ContractAddress: common.HexToAddress(cfg.EthConfig().EthContractAddress),
		log:             cfg.Log().WithField("component", "eth_client"),
	}, nil
}

func (e *EthClient) LoadRecentBlocks(ctx context.Context, logs chan<- types.Log, startBlock uint64, endBlock uint64) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.ContractAddress},
		FromBlock: new(big.Int).SetUint64(startBlock),
		ToBlock:   new(big.Int).SetUint64(endBlock),
	}

	newlogs, err := e.Client.FilterLogs(ctx, query)
	if err != nil {
		return errors.Wrap(err, "failed to filter logs")
	}

	e.log.WithFields(logan.F{
		"fromBlock": query.FromBlock.Uint64(),
		"toBlock":   query.ToBlock.Uint64(),
		"logCount":  len(newlogs),
	}).Info("Received logs from Ethereum")

	if len(newlogs) == 0 {
		e.log.Warn("No logs found for the given block range")
	}

	for _, vLog := range newlogs {
		logs <- vLog
	}

	return nil
}

func (e *EthClient) ParseTransferEvent(vLog types.Log) (*TransferEvent, error) {
	var event TransferEvent
	err := e.ContractAbi.UnpackIntoInterface(&event, "Transfer", vLog.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack transfer event")
	}

	return &event, nil
}

type TransferEvent struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
}
