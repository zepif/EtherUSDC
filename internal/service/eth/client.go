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
	store "github.com/zepif/EtherUSDC/internal/store"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type EthClient struct {
	Client          *ethclient.Client
	ContractAbi     abi.ABI
	ContractAddress common.Address
}

func NewEthClient(projectID string, contractAddress string, contractAbiJSON string) (*EthClient, error) {
	rpcURL := fmt.Sprintf("https://mainnet.infura.io/v3/%s", projectID)
	print(rpcURL)
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
		ContractAddress: common.HexToAddress(contractAddress),
	}, nil
}

func (e *EthClient) ListenToEvents(ctx context.Context, logs chan<- types.Log) error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{e.ContractAddress},
	}

	latestBlock, err := e.Client.BlockNumber(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get latest block number")
	}
	startBlock := latestBlock - 2

	for startBlock <= latestBlock {
		query.FromBlock = new(big.Int).SetUint64(startBlock)
		query.ToBlock = new(big.Int).SetUint64(latestBlock)

		newlogs, err := e.Client.FilterLogs(ctx, query)
		if err != nil {
			return errors.Wrap(err, "failed to filter logs")
		}

		for _, vLog := range newlogs {
			logs <- vLog
		}

		startBlock = latestBlock + 1
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
	From  common.Address
	To    common.Address
	Value *big.Int
}
