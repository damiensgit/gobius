package ipfs

import (
	"gobius/bindings/arbiusrouterv1"
	"gobius/client"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ArbiusRouterContract defines the interface for interacting with IPFS incentive functions
type ArbiusRouterContract interface {
	// Check if incentive is available for a task
	Incentives(opts *bind.CallOpts, taskId [32]byte) (*big.Int, error)

	// Claim IPFS incentive
	ClaimIncentive(opts *bind.TransactOpts, taskid_ [32]byte, sigs_ []arbiusrouterv1.Signature) (*types.Transaction, error)
}

// RouterContractWrapper implements the interface using actual contract binding
type RouterContractWrapper struct {
	arbiusRouter *arbiusrouterv1.ArbiusRouterV1
	client       *client.Client
}

func NewRouterContractWrapper(arbiusRouter *arbiusrouterv1.ArbiusRouterV1, client *client.Client) *RouterContractWrapper {
	return &RouterContractWrapper{
		arbiusRouter: arbiusRouter,
		client:       client,
	}
}

func (r *RouterContractWrapper) Incentives(opts *bind.CallOpts, taskId [32]byte) (*big.Int, error) {
	// Call the actual contract method
	return r.arbiusRouter.Incentives(opts, taskId)
}

func (r *RouterContractWrapper) ClaimIncentive(opts *bind.TransactOpts, taskid_ [32]byte, sigs_ []arbiusrouterv1.Signature) (*types.Transaction, error) {
	// Call the actual contract method (assuming it exists)
	return r.arbiusRouter.ClaimIncentive(opts, taskid_, sigs_)
}

// MockIPFSIncentiveContract implements the interface with mock responses
type MockRouterContract struct {
	incentiveAmounts map[[32]byte]*big.Int
	claimResults     map[[32]byte]error
	defaultAmount    *big.Int
}

func NewMockRouterContract() *MockRouterContract {
	return &MockRouterContract{
		incentiveAmounts: make(map[[32]byte]*big.Int),
		claimResults:     make(map[[32]byte]error),
		defaultAmount:    big.NewInt(1000000000000000000), // Default 1 ETH incentive
	}
}

// SetIncentiveAmount lets you configure mock responses for specific task IDs
func (m *MockRouterContract) SetIncentiveAmount(taskId [32]byte, amount *big.Int) {
	m.incentiveAmounts[taskId] = amount
}

// SetClaimResult lets you configure mock claim results for specific task IDs
func (m *MockRouterContract) SetClaimResult(taskId [32]byte, err error) {
	m.claimResults[taskId] = err
}

func (m *MockRouterContract) Incentives(opts *bind.CallOpts, taskId [32]byte) (*big.Int, error) {
	// Return configured mock amount or default
	if amount, exists := m.incentiveAmounts[taskId]; exists {
		return amount, nil
	}
	return m.defaultAmount, nil
}

func (m *MockRouterContract) ClaimIncentive(opts *bind.TransactOpts, taskid_ [32]byte, sigs_ []arbiusrouterv1.Signature) (*types.Transaction, error) {
	// Return configured error or success
	if err, exists := m.claimResults[taskid_]; exists {
		if err != nil {
			return nil, err
		}
	}

	// Mock transaction for successful claims
	return types.NewTransaction(
		0,                // nonce
		common.Address{}, // to
		big.NewInt(0),    // amount
		0,                // gas limit
		big.NewInt(0),    // gas price
		nil,              // data
	), nil
}
