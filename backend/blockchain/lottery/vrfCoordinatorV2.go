// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package lottery

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// VRFV2PlusClientRandomWordsRequest is an auto generated low-level Go binding around an user-defined struct.
type VRFV2PlusClientRandomWordsRequest struct {
	KeyHash              [32]byte
	SubId                *big.Int
	RequestConfirmations uint16
	CallbackGasLimit     uint32
	NumWords             uint32
	ExtraArgs            []byte
}

// VRFCoordinatorV2MetaData contains all meta data concerning the VRFCoordinatorV2 contract.
var VRFCoordinatorV2MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"requester\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"name\":\"RandomWordsRequested\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"CallFullfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"keyHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"subId\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"requestConfirmations\",\"type\":\"uint16\"},{\"internalType\":\"uint32\",\"name\":\"callbackGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"numWords\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"extraArgs\",\"type\":\"bytes\"}],\"internalType\":\"structVRFV2PlusClient.RandomWordsRequest\",\"name\":\"req\",\"type\":\"tuple\"}],\"name\":\"requestRandomWords\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"s_requesters\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b50610a288061001c5f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c80637140f2541461004357806382659ac71461005f5780639b1c385e1461008f575b5f5ffd5b61005d6004803603810190610058919061043c565b6100bf565b005b61007960048036038101906100749190610499565b61026d565b6040516100869190610503565b60405180910390f35b6100a960048036038101906100a4919061053e565b61029d565b6040516100b69190610594565b60405180910390f35b5f60015f8581526020019081526020015f205f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361012c5750610268565b5f8173ffffffffffffffffffffffffffffffffffffffff1685858560405160240161015993929190610625565b6040516020818303038152906040527f1fe543e3000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516101e391906106a7565b5f604051808303815f865af19150503d805f811461021c576040519150601f19603f3d011682016040523d82523d5f602084013e610221565b606091505b5050905080610265576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161025c90610717565b60405180910390fd5b50505b505050565b6001602052805f5260405f205f915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f5f8154809291906102af90610762565b91905055503360015f5f5481526020019081526020015f205f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055507faf3c81de4b70b01b31f8facb5ee988d4217788ff467575b0360343b8b65af8ad5f54338460200135855f013586606001602081019061034491906107e2565b8760400160208101906103579190610844565b88608001602081019061036a91906107e2565b898060a0019061037a919061087b565b60405161038f9998979695949392919061096d565b60405180910390a15f549050919050565b5f5ffd5b5f5ffd5b5f819050919050565b6103ba816103a8565b81146103c4575f5ffd5b50565b5f813590506103d5816103b1565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126103fc576103fb6103db565b5b8235905067ffffffffffffffff811115610419576104186103df565b5b602083019150836020820283011115610435576104346103e3565b5b9250929050565b5f5f5f60408486031215610453576104526103a0565b5b5f610460868287016103c7565b935050602084013567ffffffffffffffff811115610481576104806103a4565b5b61048d868287016103e7565b92509250509250925092565b5f602082840312156104ae576104ad6103a0565b5b5f6104bb848285016103c7565b91505092915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6104ed826104c4565b9050919050565b6104fd816104e3565b82525050565b5f6020820190506105165f8301846104f4565b92915050565b5f5ffd5b5f60c082840312156105355761053461051c565b5b81905092915050565b5f60208284031215610553576105526103a0565b5b5f82013567ffffffffffffffff8111156105705761056f6103a4565b5b61057c84828501610520565b91505092915050565b61058e816103a8565b82525050565b5f6020820190506105a75f830184610585565b92915050565b5f82825260208201905092915050565b5f5ffd5b82818337505050565b5f6105d583856105ad565b93507f07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff831115610608576106076105bd565b5b6020830292506106198385846105c1565b82840190509392505050565b5f6040820190506106385f830186610585565b818103602083015261064b8184866105ca565b9050949350505050565b5f81519050919050565b5f81905092915050565b8281835e5f83830152505050565b5f61068182610655565b61068b818561065f565b935061069b818560208601610669565b80840191505092915050565b5f6106b28284610677565b915081905092915050565b5f82825260208201905092915050565b7f46756c66696c6c52616e646f6d576f726473206661696c65642e0000000000005f82015250565b5f610701601a836106bd565b915061070c826106cd565b602082019050919050565b5f6020820190508181035f83015261072e816106f5565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f61076c826103a8565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361079e5761079d610735565b5b600182019050919050565b5f63ffffffff82169050919050565b6107c1816107a9565b81146107cb575f5ffd5b50565b5f813590506107dc816107b8565b92915050565b5f602082840312156107f7576107f66103a0565b5b5f610804848285016107ce565b91505092915050565b5f61ffff82169050919050565b6108238161080d565b811461082d575f5ffd5b50565b5f8135905061083e8161081a565b92915050565b5f60208284031215610859576108586103a0565b5b5f61086684828501610830565b91505092915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f833560016020038436030381126108975761089661086f565b5b80840192508235915067ffffffffffffffff8211156108b9576108b8610873565b5b6020830192506001820236038313156108d5576108d4610877565b5b509250929050565b5f819050919050565b6108ef816108dd565b82525050565b6108fe816107a9565b82525050565b61090d8161080d565b82525050565b5f82825260208201905092915050565b828183375f83830152505050565b5f601f19601f8301169050919050565b5f61094c8385610913565b9350610959838584610923565b61096283610931565b840190509392505050565b5f610100820190506109815f83018c610585565b61098e602083018b6104f4565b61099b604083018a610585565b6109a860608301896108e6565b6109b560808301886108f5565b6109c260a0830187610904565b6109cf60c08301866108f5565b81810360e08301526109e2818486610941565b90509a995050505050505050505056fea26469706673582212206d23e359bf155d204f9243329f669a968ae1d6153ebd83a5748fe529670495bb64736f6c634300081c0033",
}

// VRFCoordinatorV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use VRFCoordinatorV2MetaData.ABI instead.
var VRFCoordinatorV2ABI = VRFCoordinatorV2MetaData.ABI

// VRFCoordinatorV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VRFCoordinatorV2MetaData.Bin instead.
var VRFCoordinatorV2Bin = VRFCoordinatorV2MetaData.Bin

// DeployVRFCoordinatorV2 deploys a new Ethereum contract, binding an instance of VRFCoordinatorV2 to it.
func DeployVRFCoordinatorV2(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VRFCoordinatorV2, error) {
	parsed, err := VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VRFCoordinatorV2Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VRFCoordinatorV2{VRFCoordinatorV2Caller: VRFCoordinatorV2Caller{contract: contract}, VRFCoordinatorV2Transactor: VRFCoordinatorV2Transactor{contract: contract}, VRFCoordinatorV2Filterer: VRFCoordinatorV2Filterer{contract: contract}}, nil
}

// VRFCoordinatorV2 is an auto generated Go binding around an Ethereum contract.
type VRFCoordinatorV2 struct {
	VRFCoordinatorV2Caller     // Read-only binding to the contract
	VRFCoordinatorV2Transactor // Write-only binding to the contract
	VRFCoordinatorV2Filterer   // Log filterer for contract events
}

// VRFCoordinatorV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type VRFCoordinatorV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type VRFCoordinatorV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VRFCoordinatorV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VRFCoordinatorV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VRFCoordinatorV2Session struct {
	Contract     *VRFCoordinatorV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VRFCoordinatorV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VRFCoordinatorV2CallerSession struct {
	Contract *VRFCoordinatorV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// VRFCoordinatorV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VRFCoordinatorV2TransactorSession struct {
	Contract     *VRFCoordinatorV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// VRFCoordinatorV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type VRFCoordinatorV2Raw struct {
	Contract *VRFCoordinatorV2 // Generic contract binding to access the raw methods on
}

// VRFCoordinatorV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VRFCoordinatorV2CallerRaw struct {
	Contract *VRFCoordinatorV2Caller // Generic read-only contract binding to access the raw methods on
}

// VRFCoordinatorV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VRFCoordinatorV2TransactorRaw struct {
	Contract *VRFCoordinatorV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewVRFCoordinatorV2 creates a new instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2(address common.Address, backend bind.ContractBackend) (*VRFCoordinatorV2, error) {
	contract, err := bindVRFCoordinatorV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2{VRFCoordinatorV2Caller: VRFCoordinatorV2Caller{contract: contract}, VRFCoordinatorV2Transactor: VRFCoordinatorV2Transactor{contract: contract}, VRFCoordinatorV2Filterer: VRFCoordinatorV2Filterer{contract: contract}}, nil
}

// NewVRFCoordinatorV2Caller creates a new read-only instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2Caller(address common.Address, caller bind.ContractCaller) (*VRFCoordinatorV2Caller, error) {
	contract, err := bindVRFCoordinatorV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Caller{contract: contract}, nil
}

// NewVRFCoordinatorV2Transactor creates a new write-only instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2Transactor(address common.Address, transactor bind.ContractTransactor) (*VRFCoordinatorV2Transactor, error) {
	contract, err := bindVRFCoordinatorV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Transactor{contract: contract}, nil
}

// NewVRFCoordinatorV2Filterer creates a new log filterer instance of VRFCoordinatorV2, bound to a specific deployed contract.
func NewVRFCoordinatorV2Filterer(address common.Address, filterer bind.ContractFilterer) (*VRFCoordinatorV2Filterer, error) {
	contract, err := bindVRFCoordinatorV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2Filterer{contract: contract}, nil
}

// bindVRFCoordinatorV2 binds a generic wrapper to an already deployed contract.
func bindVRFCoordinatorV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VRFCoordinatorV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinatorV2 *VRFCoordinatorV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.VRFCoordinatorV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VRFCoordinatorV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.contract.Transact(opts, method, params...)
}

// SRequesters is a free data retrieval call binding the contract method 0x82659ac7.
//
// Solidity: function s_requesters(uint256 ) view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Caller) SRequesters(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _VRFCoordinatorV2.contract.Call(opts, &out, "s_requesters", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SRequesters is a free data retrieval call binding the contract method 0x82659ac7.
//
// Solidity: function s_requesters(uint256 ) view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) SRequesters(arg0 *big.Int) (common.Address, error) {
	return _VRFCoordinatorV2.Contract.SRequesters(&_VRFCoordinatorV2.CallOpts, arg0)
}

// SRequesters is a free data retrieval call binding the contract method 0x82659ac7.
//
// Solidity: function s_requesters(uint256 ) view returns(address)
func (_VRFCoordinatorV2 *VRFCoordinatorV2CallerSession) SRequesters(arg0 *big.Int) (common.Address, error) {
	return _VRFCoordinatorV2.Contract.SRequesters(&_VRFCoordinatorV2.CallOpts, arg0)
}

// CallFullfillRandomWords is a paid mutator transaction binding the contract method 0x7140f254.
//
// Solidity: function CallFullfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) CallFullfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "CallFullfillRandomWords", requestId, randomWords)
}

// CallFullfillRandomWords is a paid mutator transaction binding the contract method 0x7140f254.
//
// Solidity: function CallFullfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) CallFullfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CallFullfillRandomWords(&_VRFCoordinatorV2.TransactOpts, requestId, randomWords)
}

// CallFullfillRandomWords is a paid mutator transaction binding the contract method 0x7140f254.
//
// Solidity: function CallFullfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) CallFullfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.CallFullfillRandomWords(&_VRFCoordinatorV2.TransactOpts, requestId, randomWords)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x9b1c385e.
//
// Solidity: function requestRandomWords((bytes32,uint256,uint16,uint32,uint32,bytes) req) returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Transactor) RequestRandomWords(opts *bind.TransactOpts, req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2.contract.Transact(opts, "requestRandomWords", req)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x9b1c385e.
//
// Solidity: function requestRandomWords((bytes32,uint256,uint16,uint32,uint32,bytes) req) returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Session) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestRandomWords(&_VRFCoordinatorV2.TransactOpts, req)
}

// RequestRandomWords is a paid mutator transaction binding the contract method 0x9b1c385e.
//
// Solidity: function requestRandomWords((bytes32,uint256,uint16,uint32,uint32,bytes) req) returns(uint256)
func (_VRFCoordinatorV2 *VRFCoordinatorV2TransactorSession) RequestRandomWords(req VRFV2PlusClientRandomWordsRequest) (*types.Transaction, error) {
	return _VRFCoordinatorV2.Contract.RequestRandomWords(&_VRFCoordinatorV2.TransactOpts, req)
}

// VRFCoordinatorV2RandomWordsRequestedIterator is returned from FilterRandomWordsRequested and is used to iterate over the raw logs and unpacked data for RandomWordsRequested events raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2RandomWordsRequestedIterator struct {
	Event *VRFCoordinatorV2RandomWordsRequested // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VRFCoordinatorV2RandomWordsRequested)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(VRFCoordinatorV2RandomWordsRequested)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VRFCoordinatorV2RandomWordsRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VRFCoordinatorV2RandomWordsRequested represents a RandomWordsRequested event raised by the VRFCoordinatorV2 contract.
type VRFCoordinatorV2RandomWordsRequested struct {
	RequestId            *big.Int
	Requester            common.Address
	SubId                *big.Int
	KeyHash              [32]byte
	CallbackGasLimit     uint32
	RequestConfirmations uint16
	NumWords             uint32
	ExtraArgs            []byte
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterRandomWordsRequested is a free log retrieval operation binding the contract event 0xaf3c81de4b70b01b31f8facb5ee988d4217788ff467575b0360343b8b65af8ad.
//
// Solidity: event RandomWordsRequested(uint256 requestId, address requester, uint256 subId, bytes32 keyHash, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes extraArgs)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) FilterRandomWordsRequested(opts *bind.FilterOpts) (*VRFCoordinatorV2RandomWordsRequestedIterator, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.FilterLogs(opts, "RandomWordsRequested")
	if err != nil {
		return nil, err
	}
	return &VRFCoordinatorV2RandomWordsRequestedIterator{contract: _VRFCoordinatorV2.contract, event: "RandomWordsRequested", logs: logs, sub: sub}, nil
}

// WatchRandomWordsRequested is a free log subscription operation binding the contract event 0xaf3c81de4b70b01b31f8facb5ee988d4217788ff467575b0360343b8b65af8ad.
//
// Solidity: event RandomWordsRequested(uint256 requestId, address requester, uint256 subId, bytes32 keyHash, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes extraArgs)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) WatchRandomWordsRequested(opts *bind.WatchOpts, sink chan<- *VRFCoordinatorV2RandomWordsRequested) (event.Subscription, error) {

	logs, sub, err := _VRFCoordinatorV2.contract.WatchLogs(opts, "RandomWordsRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VRFCoordinatorV2RandomWordsRequested)
				if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRandomWordsRequested is a log parse operation binding the contract event 0xaf3c81de4b70b01b31f8facb5ee988d4217788ff467575b0360343b8b65af8ad.
//
// Solidity: event RandomWordsRequested(uint256 requestId, address requester, uint256 subId, bytes32 keyHash, uint32 callbackGasLimit, uint16 requestConfirmations, uint32 numWords, bytes extraArgs)
func (_VRFCoordinatorV2 *VRFCoordinatorV2Filterer) ParseRandomWordsRequested(log types.Log) (*VRFCoordinatorV2RandomWordsRequested, error) {
	event := new(VRFCoordinatorV2RandomWordsRequested)
	if err := _VRFCoordinatorV2.contract.UnpackLog(event, "RandomWordsRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
