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

// SimpleRolloutMetaData contains all meta data concerning the SimpleRollout contract.
var SimpleRolloutMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subscriptionId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_vrfCoordinator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_trigger\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"have\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"want\",\"type\":\"address\"}],\"name\":\"OnlyCoordinatorCanFulfill\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256[]\",\"name\":\"results\",\"type\":\"uint256[]\"}],\"name\":\"DiceLanded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"}],\"name\":\"DiceRolled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"randomWords\",\"type\":\"uint256[]\"}],\"name\":\"rawFulfillRandomWords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIRolloutCallback\",\"name\":\"rolloutcb\",\"type\":\"address\"}],\"name\":\"rolloutCall\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollout_epoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rollout_results\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"s_vrfCoordinator\",\"outputs\":[{\"internalType\":\"contractIVRFCoordinatorV2Plus\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SimpleRolloutABI is the input ABI used to generate the binding from.
// Deprecated: Use SimpleRolloutMetaData.ABI instead.
var SimpleRolloutABI = SimpleRolloutMetaData.ABI

// SimpleRollout is an auto generated Go binding around an Ethereum contract.
type SimpleRollout struct {
	SimpleRolloutCaller     // Read-only binding to the contract
	SimpleRolloutTransactor // Write-only binding to the contract
	SimpleRolloutFilterer   // Log filterer for contract events
}

// SimpleRolloutCaller is an auto generated read-only Go binding around an Ethereum contract.
type SimpleRolloutCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleRolloutTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SimpleRolloutTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleRolloutFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SimpleRolloutFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleRolloutSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SimpleRolloutSession struct {
	Contract     *SimpleRollout    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SimpleRolloutCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SimpleRolloutCallerSession struct {
	Contract *SimpleRolloutCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SimpleRolloutTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SimpleRolloutTransactorSession struct {
	Contract     *SimpleRolloutTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SimpleRolloutRaw is an auto generated low-level Go binding around an Ethereum contract.
type SimpleRolloutRaw struct {
	Contract *SimpleRollout // Generic contract binding to access the raw methods on
}

// SimpleRolloutCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SimpleRolloutCallerRaw struct {
	Contract *SimpleRolloutCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleRolloutTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SimpleRolloutTransactorRaw struct {
	Contract *SimpleRolloutTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleRollout creates a new instance of SimpleRollout, bound to a specific deployed contract.
func NewSimpleRollout(address common.Address, backend bind.ContractBackend) (*SimpleRollout, error) {
	contract, err := bindSimpleRollout(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleRollout{SimpleRolloutCaller: SimpleRolloutCaller{contract: contract}, SimpleRolloutTransactor: SimpleRolloutTransactor{contract: contract}, SimpleRolloutFilterer: SimpleRolloutFilterer{contract: contract}}, nil
}

// NewSimpleRolloutCaller creates a new read-only instance of SimpleRollout, bound to a specific deployed contract.
func NewSimpleRolloutCaller(address common.Address, caller bind.ContractCaller) (*SimpleRolloutCaller, error) {
	contract, err := bindSimpleRollout(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleRolloutCaller{contract: contract}, nil
}

// NewSimpleRolloutTransactor creates a new write-only instance of SimpleRollout, bound to a specific deployed contract.
func NewSimpleRolloutTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleRolloutTransactor, error) {
	contract, err := bindSimpleRollout(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleRolloutTransactor{contract: contract}, nil
}

// NewSimpleRolloutFilterer creates a new log filterer instance of SimpleRollout, bound to a specific deployed contract.
func NewSimpleRolloutFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleRolloutFilterer, error) {
	contract, err := bindSimpleRollout(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleRolloutFilterer{contract: contract}, nil
}

// bindSimpleRollout binds a generic wrapper to an already deployed contract.
func bindSimpleRollout(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleRolloutMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleRollout *SimpleRolloutRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleRollout.Contract.SimpleRolloutCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleRollout *SimpleRolloutRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleRollout.Contract.SimpleRolloutTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleRollout *SimpleRolloutRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleRollout.Contract.SimpleRolloutTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleRollout *SimpleRolloutCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleRollout.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleRollout *SimpleRolloutTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleRollout.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleRollout *SimpleRolloutTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleRollout.Contract.contract.Transact(opts, method, params...)
}

// RequestID is a free data retrieval call binding the contract method 0x8f779201.
//
// Solidity: function requestID() view returns(uint256)
func (_SimpleRollout *SimpleRolloutCaller) RequestID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleRollout.contract.Call(opts, &out, "requestID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestID is a free data retrieval call binding the contract method 0x8f779201.
//
// Solidity: function requestID() view returns(uint256)
func (_SimpleRollout *SimpleRolloutSession) RequestID() (*big.Int, error) {
	return _SimpleRollout.Contract.RequestID(&_SimpleRollout.CallOpts)
}

// RequestID is a free data retrieval call binding the contract method 0x8f779201.
//
// Solidity: function requestID() view returns(uint256)
func (_SimpleRollout *SimpleRolloutCallerSession) RequestID() (*big.Int, error) {
	return _SimpleRollout.Contract.RequestID(&_SimpleRollout.CallOpts)
}

// RolloutEpoch is a free data retrieval call binding the contract method 0x3559f882.
//
// Solidity: function rollout_epoch() view returns(uint256)
func (_SimpleRollout *SimpleRolloutCaller) RolloutEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SimpleRollout.contract.Call(opts, &out, "rollout_epoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RolloutEpoch is a free data retrieval call binding the contract method 0x3559f882.
//
// Solidity: function rollout_epoch() view returns(uint256)
func (_SimpleRollout *SimpleRolloutSession) RolloutEpoch() (*big.Int, error) {
	return _SimpleRollout.Contract.RolloutEpoch(&_SimpleRollout.CallOpts)
}

// RolloutEpoch is a free data retrieval call binding the contract method 0x3559f882.
//
// Solidity: function rollout_epoch() view returns(uint256)
func (_SimpleRollout *SimpleRolloutCallerSession) RolloutEpoch() (*big.Int, error) {
	return _SimpleRollout.Contract.RolloutEpoch(&_SimpleRollout.CallOpts)
}

// RolloutResults is a free data retrieval call binding the contract method 0x8e741fb8.
//
// Solidity: function rollout_results(uint256 ) view returns(uint256)
func (_SimpleRollout *SimpleRolloutCaller) RolloutResults(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _SimpleRollout.contract.Call(opts, &out, "rollout_results", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RolloutResults is a free data retrieval call binding the contract method 0x8e741fb8.
//
// Solidity: function rollout_results(uint256 ) view returns(uint256)
func (_SimpleRollout *SimpleRolloutSession) RolloutResults(arg0 *big.Int) (*big.Int, error) {
	return _SimpleRollout.Contract.RolloutResults(&_SimpleRollout.CallOpts, arg0)
}

// RolloutResults is a free data retrieval call binding the contract method 0x8e741fb8.
//
// Solidity: function rollout_results(uint256 ) view returns(uint256)
func (_SimpleRollout *SimpleRolloutCallerSession) RolloutResults(arg0 *big.Int) (*big.Int, error) {
	return _SimpleRollout.Contract.RolloutResults(&_SimpleRollout.CallOpts, arg0)
}

// SVrfCoordinator is a free data retrieval call binding the contract method 0x9eccacf6.
//
// Solidity: function s_vrfCoordinator() view returns(address)
func (_SimpleRollout *SimpleRolloutCaller) SVrfCoordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleRollout.contract.Call(opts, &out, "s_vrfCoordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SVrfCoordinator is a free data retrieval call binding the contract method 0x9eccacf6.
//
// Solidity: function s_vrfCoordinator() view returns(address)
func (_SimpleRollout *SimpleRolloutSession) SVrfCoordinator() (common.Address, error) {
	return _SimpleRollout.Contract.SVrfCoordinator(&_SimpleRollout.CallOpts)
}

// SVrfCoordinator is a free data retrieval call binding the contract method 0x9eccacf6.
//
// Solidity: function s_vrfCoordinator() view returns(address)
func (_SimpleRollout *SimpleRolloutCallerSession) SVrfCoordinator() (common.Address, error) {
	return _SimpleRollout.Contract.SVrfCoordinator(&_SimpleRollout.CallOpts)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_SimpleRollout *SimpleRolloutTransactor) RawFulfillRandomWords(opts *bind.TransactOpts, requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _SimpleRollout.contract.Transact(opts, "rawFulfillRandomWords", requestId, randomWords)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_SimpleRollout *SimpleRolloutSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _SimpleRollout.Contract.RawFulfillRandomWords(&_SimpleRollout.TransactOpts, requestId, randomWords)
}

// RawFulfillRandomWords is a paid mutator transaction binding the contract method 0x1fe543e3.
//
// Solidity: function rawFulfillRandomWords(uint256 requestId, uint256[] randomWords) returns()
func (_SimpleRollout *SimpleRolloutTransactorSession) RawFulfillRandomWords(requestId *big.Int, randomWords []*big.Int) (*types.Transaction, error) {
	return _SimpleRollout.Contract.RawFulfillRandomWords(&_SimpleRollout.TransactOpts, requestId, randomWords)
}

// RolloutCall is a paid mutator transaction binding the contract method 0xe0a742a8.
//
// Solidity: function rolloutCall(address rolloutcb) returns()
func (_SimpleRollout *SimpleRolloutTransactor) RolloutCall(opts *bind.TransactOpts, rolloutcb common.Address) (*types.Transaction, error) {
	return _SimpleRollout.contract.Transact(opts, "rolloutCall", rolloutcb)
}

// RolloutCall is a paid mutator transaction binding the contract method 0xe0a742a8.
//
// Solidity: function rolloutCall(address rolloutcb) returns()
func (_SimpleRollout *SimpleRolloutSession) RolloutCall(rolloutcb common.Address) (*types.Transaction, error) {
	return _SimpleRollout.Contract.RolloutCall(&_SimpleRollout.TransactOpts, rolloutcb)
}

// RolloutCall is a paid mutator transaction binding the contract method 0xe0a742a8.
//
// Solidity: function rolloutCall(address rolloutcb) returns()
func (_SimpleRollout *SimpleRolloutTransactorSession) RolloutCall(rolloutcb common.Address) (*types.Transaction, error) {
	return _SimpleRollout.Contract.RolloutCall(&_SimpleRollout.TransactOpts, rolloutcb)
}

// SimpleRolloutDiceLandedIterator is returned from FilterDiceLanded and is used to iterate over the raw logs and unpacked data for DiceLanded events raised by the SimpleRollout contract.
type SimpleRolloutDiceLandedIterator struct {
	Event *SimpleRolloutDiceLanded // Event containing the contract specifics and raw log

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
func (it *SimpleRolloutDiceLandedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleRolloutDiceLanded)
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
		it.Event = new(SimpleRolloutDiceLanded)
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
func (it *SimpleRolloutDiceLandedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleRolloutDiceLandedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleRolloutDiceLanded represents a DiceLanded event raised by the SimpleRollout contract.
type SimpleRolloutDiceLanded struct {
	RequestId *big.Int
	Results   []*big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDiceLanded is a free log retrieval operation binding the contract event 0xbb49a05d9ec394e79a221aa845d28284598dad4c765b725d0cd3054f801c7972.
//
// Solidity: event DiceLanded(uint256 indexed requestId, uint256[] indexed results)
func (_SimpleRollout *SimpleRolloutFilterer) FilterDiceLanded(opts *bind.FilterOpts, requestId []*big.Int, results [][]*big.Int) (*SimpleRolloutDiceLandedIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var resultsRule []interface{}
	for _, resultsItem := range results {
		resultsRule = append(resultsRule, resultsItem)
	}

	logs, sub, err := _SimpleRollout.contract.FilterLogs(opts, "DiceLanded", requestIdRule, resultsRule)
	if err != nil {
		return nil, err
	}
	return &SimpleRolloutDiceLandedIterator{contract: _SimpleRollout.contract, event: "DiceLanded", logs: logs, sub: sub}, nil
}

// WatchDiceLanded is a free log subscription operation binding the contract event 0xbb49a05d9ec394e79a221aa845d28284598dad4c765b725d0cd3054f801c7972.
//
// Solidity: event DiceLanded(uint256 indexed requestId, uint256[] indexed results)
func (_SimpleRollout *SimpleRolloutFilterer) WatchDiceLanded(opts *bind.WatchOpts, sink chan<- *SimpleRolloutDiceLanded, requestId []*big.Int, results [][]*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var resultsRule []interface{}
	for _, resultsItem := range results {
		resultsRule = append(resultsRule, resultsItem)
	}

	logs, sub, err := _SimpleRollout.contract.WatchLogs(opts, "DiceLanded", requestIdRule, resultsRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleRolloutDiceLanded)
				if err := _SimpleRollout.contract.UnpackLog(event, "DiceLanded", log); err != nil {
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

// ParseDiceLanded is a log parse operation binding the contract event 0xbb49a05d9ec394e79a221aa845d28284598dad4c765b725d0cd3054f801c7972.
//
// Solidity: event DiceLanded(uint256 indexed requestId, uint256[] indexed results)
func (_SimpleRollout *SimpleRolloutFilterer) ParseDiceLanded(log types.Log) (*SimpleRolloutDiceLanded, error) {
	event := new(SimpleRolloutDiceLanded)
	if err := _SimpleRollout.contract.UnpackLog(event, "DiceLanded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SimpleRolloutDiceRolledIterator is returned from FilterDiceRolled and is used to iterate over the raw logs and unpacked data for DiceRolled events raised by the SimpleRollout contract.
type SimpleRolloutDiceRolledIterator struct {
	Event *SimpleRolloutDiceRolled // Event containing the contract specifics and raw log

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
func (it *SimpleRolloutDiceRolledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleRolloutDiceRolled)
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
		it.Event = new(SimpleRolloutDiceRolled)
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
func (it *SimpleRolloutDiceRolledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleRolloutDiceRolledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleRolloutDiceRolled represents a DiceRolled event raised by the SimpleRollout contract.
type SimpleRolloutDiceRolled struct {
	RequestId *big.Int
	Epoch     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDiceRolled is a free log retrieval operation binding the contract event 0xc8aa159ba2f8c0865ab5ea1cdb5b70ce807eebfc517029c8668458d69a4059da.
//
// Solidity: event DiceRolled(uint256 indexed requestId, uint256 indexed epoch)
func (_SimpleRollout *SimpleRolloutFilterer) FilterDiceRolled(opts *bind.FilterOpts, requestId []*big.Int, epoch []*big.Int) (*SimpleRolloutDiceRolledIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _SimpleRollout.contract.FilterLogs(opts, "DiceRolled", requestIdRule, epochRule)
	if err != nil {
		return nil, err
	}
	return &SimpleRolloutDiceRolledIterator{contract: _SimpleRollout.contract, event: "DiceRolled", logs: logs, sub: sub}, nil
}

// WatchDiceRolled is a free log subscription operation binding the contract event 0xc8aa159ba2f8c0865ab5ea1cdb5b70ce807eebfc517029c8668458d69a4059da.
//
// Solidity: event DiceRolled(uint256 indexed requestId, uint256 indexed epoch)
func (_SimpleRollout *SimpleRolloutFilterer) WatchDiceRolled(opts *bind.WatchOpts, sink chan<- *SimpleRolloutDiceRolled, requestId []*big.Int, epoch []*big.Int) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _SimpleRollout.contract.WatchLogs(opts, "DiceRolled", requestIdRule, epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleRolloutDiceRolled)
				if err := _SimpleRollout.contract.UnpackLog(event, "DiceRolled", log); err != nil {
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

// ParseDiceRolled is a log parse operation binding the contract event 0xc8aa159ba2f8c0865ab5ea1cdb5b70ce807eebfc517029c8668458d69a4059da.
//
// Solidity: event DiceRolled(uint256 indexed requestId, uint256 indexed epoch)
func (_SimpleRollout *SimpleRolloutFilterer) ParseDiceRolled(log types.Log) (*SimpleRolloutDiceRolled, error) {
	event := new(SimpleRolloutDiceRolled)
	if err := _SimpleRollout.contract.UnpackLog(event, "DiceRolled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
