package eth

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock EthClient
type mockEthClient struct {
	mock.Mock
	Client          *ethclient.Client
	ContractAbi     abi.ABI
	ContractAddress common.Address
}

func (m *mockEthClient) BlockNumber(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockEthClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	args := m.Called(ctx, q)
	return args.Get(0).([]types.Log), args.Error(1)
}

func (m *mockEthClient) ListenToEvents(ctx context.Context, logs chan types.Log) error {
	args := m.Called(ctx, logs)
	return args.Error(0)
}

func (m *mockEthClient) ParseTransferEvent(vLog types.Log) (*transferEvent, error) {
	args := m.Called(vLog)
	return args.Get(0).(*transferEvent), args.Error(1)
}

func TestEthClient(t *testing.T) {
	mockClient := &mockEthClient{
		ContractAddress: common.HexToAddress("0x123"),
	}

	// Test case: ListenToEvents when BlockNumber fails
	ctx := context.Background()
	logs := make(chan types.Log, 100)
	mockClient.On("BlockNumber", ctx).Return(uint64(0), assert.AnError)
	mockClient.On("ListenToEvents", ctx, logs).Return(assert.AnError)

	err := mockClient.ListenToEvents(ctx, logs)
	assert.Error(t, err)

	// Test case: ListenToEvents when FilterLogs fails
	blockNumber := uint64(1000)
	mockClient.On("BlockNumber", ctx).Return(blockNumber, nil)
	mockClient.On("FilterLogs", ctx, mock.Anything).Return(nil, assert.AnError)
	mockClient.On("ListenToEvents", ctx, logs).Return(assert.AnError)

	err = mockClient.ListenToEvents(ctx, logs)
	assert.Error(t, err)

	// Test case: ListenToEvents successful
	logs = make(chan types.Log, 100)
	mockLogs := []types.Log{
		{
			TxHash: common.HexToHash("0x456"),
			Data:   []byte{1, 2, 3},
		},
	}
	mockClient.On("BlockNumber", ctx).Return(blockNumber, nil)
	mockClient.On("FilterLogs", ctx, mock.Anything).Return(mockLogs, nil)
	mockClient.On("ListenToEvents", ctx, logs).Return(nil)

	err = mockClient.ListenToEvents(ctx, logs)
	assert.NoError(t, err)

	// Read from the logs channel
	receivedLogs := make([]types.Log, 0)
	for {
		select {
		case log := <-logs:
			receivedLogs = append(receivedLogs, log)
		default:
			goto done
		}
	}
done:
	assert.Equal(t, len(mockLogs), len(receivedLogs))

	// Test case: ParseTransferEvent successful
	mockEvent := transferEvent{
		From:   common.HexToAddress("0x789"),
		To:     common.HexToAddress("0xabc"),
		Amount: big.NewInt(1000),
	}
	mockClient.ContractAbi = mockABI(mockEvent)

	mockClient.On("ParseTransferEvent", mockLogs[0]).Return(&mockEvent, nil)

	event, err := mockClient.ParseTransferEvent(mockLogs[0])
	assert.NoError(t, err)
	assert.Equal(t, mockEvent, *event)

	// Add more test cases as needed
}

// mockABI is a helper function to create a mock ABI
func mockABI(event transferEvent) abi.ABI {
	abiString := `[{"type":"event","name":"Transfer","inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"amount","type":"uint256"}]}]`
	abiData, _ := abi.JSON(strings.NewReader(abiString))
	return abiData
}

type transferEvent struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
}
