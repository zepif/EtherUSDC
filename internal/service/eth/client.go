package eth

import (
    "context"
    "log"
    "math/big"
    "strings"

    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
    "gitlab.com/distributed_lab/logan/v3/errors"
)

type EthClient struct {
    Client          *ethclient.Client
    ContractAbi     abi.ABI
    ContractAddress common.Address
}

func NewEthClient(rpcURL string, contractAddress string, contractAbiJSON string) (*EthClient, error) {
    client, err := ethclient.Dial(rpcURL)
    if err != nil {
        return nil, errors.Wrap(err, "failed to connect to Ethereum client")
    }

    contractAbi, err := abi.JSON(strings.NewReader(contractAbiJSON))
    if err != nil {
        return nil, errors.Wrap(err, "failed to parse contract ABI") 
    }

    return &EthClient{
		Client:          client,
		ContractAbi:     contractAbi,
		ContractAddress: common.HexToAddress(contractAddress),
	}, nil
}

func (e *EthClient) ListenToEvents(ctx context.Context, logs chan types.Log) error {
	query := ethereum.FilterQuery {
        Addresses: []common.Address{e.ContractAddress},
    }

    logsSub, err := e.Client.SubscribeFilterLogs(ctx, query, logs)
    if err != nil {
        return errors.Wrap(err, "failed to subscribe to contract events")
    }

    go func() {
        for {
            select {
            case err := <-logsSub.Err():
                log.Println("Error:", err)
            case vLog := <-logs:
                logs <- vLog
            }
        }
    }()

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

