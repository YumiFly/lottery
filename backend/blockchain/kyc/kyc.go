// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package kyc

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

// KYCMetaData contains all meta data concerning the KYC contract.
var KYCMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"KYCRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"customer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"KYCVerified\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"customers\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"customerAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isVerified\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"verificationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"registrationTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_customer\",\"type\":\"address\"}],\"name\":\"getKYCStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_customer\",\"type\":\"address\"}],\"name\":\"verifyKYC\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// KYCABI is the input ABI used to generate the binding from.
// Deprecated: Use KYCMetaData.ABI instead.
var KYCABI = KYCMetaData.ABI

// KYC is an auto generated Go binding around an Ethereum contract.
type KYC struct {
	KYCCaller     // Read-only binding to the contract
	KYCTransactor // Write-only binding to the contract
	KYCFilterer   // Log filterer for contract events
}

// KYCCaller is an auto generated read-only Go binding around an Ethereum contract.
type KYCCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KYCTransactor is an auto generated write-only Go binding around an Ethereum contract.
type KYCTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KYCFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type KYCFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// KYCSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type KYCSession struct {
	Contract     *KYC              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KYCCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type KYCCallerSession struct {
	Contract *KYCCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// KYCTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type KYCTransactorSession struct {
	Contract     *KYCTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// KYCRaw is an auto generated low-level Go binding around an Ethereum contract.
type KYCRaw struct {
	Contract *KYC // Generic contract binding to access the raw methods on
}

// KYCCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type KYCCallerRaw struct {
	Contract *KYCCaller // Generic read-only contract binding to access the raw methods on
}

// KYCTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type KYCTransactorRaw struct {
	Contract *KYCTransactor // Generic write-only contract binding to access the raw methods on
}

// NewKYC creates a new instance of KYC, bound to a specific deployed contract.
func NewKYC(address common.Address, backend bind.ContractBackend) (*KYC, error) {
	contract, err := bindKYC(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &KYC{KYCCaller: KYCCaller{contract: contract}, KYCTransactor: KYCTransactor{contract: contract}, KYCFilterer: KYCFilterer{contract: contract}}, nil
}

// NewKYCCaller creates a new read-only instance of KYC, bound to a specific deployed contract.
func NewKYCCaller(address common.Address, caller bind.ContractCaller) (*KYCCaller, error) {
	contract, err := bindKYC(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &KYCCaller{contract: contract}, nil
}

// NewKYCTransactor creates a new write-only instance of KYC, bound to a specific deployed contract.
func NewKYCTransactor(address common.Address, transactor bind.ContractTransactor) (*KYCTransactor, error) {
	contract, err := bindKYC(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &KYCTransactor{contract: contract}, nil
}

// NewKYCFilterer creates a new log filterer instance of KYC, bound to a specific deployed contract.
func NewKYCFilterer(address common.Address, filterer bind.ContractFilterer) (*KYCFilterer, error) {
	contract, err := bindKYC(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &KYCFilterer{contract: contract}, nil
}

// bindKYC binds a generic wrapper to an already deployed contract.
func bindKYC(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := KYCMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KYC *KYCRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KYC.Contract.KYCCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KYC *KYCRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KYC.Contract.KYCTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KYC *KYCRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KYC.Contract.KYCTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_KYC *KYCCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _KYC.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_KYC *KYCTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KYC.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_KYC *KYCTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _KYC.Contract.contract.Transact(opts, method, params...)
}

// Customers is a free data retrieval call binding the contract method 0x336989ae.
//
// Solidity: function customers(address ) view returns(address customerAddress, bool isVerified, address verifier, uint256 verificationTime, uint256 registrationTime)
func (_KYC *KYCCaller) Customers(opts *bind.CallOpts, arg0 common.Address) (struct {
	CustomerAddress  common.Address
	IsVerified       bool
	Verifier         common.Address
	VerificationTime *big.Int
	RegistrationTime *big.Int
}, error) {
	var out []interface{}
	err := _KYC.contract.Call(opts, &out, "customers", arg0)

	outstruct := new(struct {
		CustomerAddress  common.Address
		IsVerified       bool
		Verifier         common.Address
		VerificationTime *big.Int
		RegistrationTime *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CustomerAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.IsVerified = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Verifier = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.VerificationTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.RegistrationTime = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Customers is a free data retrieval call binding the contract method 0x336989ae.
//
// Solidity: function customers(address ) view returns(address customerAddress, bool isVerified, address verifier, uint256 verificationTime, uint256 registrationTime)
func (_KYC *KYCSession) Customers(arg0 common.Address) (struct {
	CustomerAddress  common.Address
	IsVerified       bool
	Verifier         common.Address
	VerificationTime *big.Int
	RegistrationTime *big.Int
}, error) {
	return _KYC.Contract.Customers(&_KYC.CallOpts, arg0)
}

// Customers is a free data retrieval call binding the contract method 0x336989ae.
//
// Solidity: function customers(address ) view returns(address customerAddress, bool isVerified, address verifier, uint256 verificationTime, uint256 registrationTime)
func (_KYC *KYCCallerSession) Customers(arg0 common.Address) (struct {
	CustomerAddress  common.Address
	IsVerified       bool
	Verifier         common.Address
	VerificationTime *big.Int
	RegistrationTime *big.Int
}, error) {
	return _KYC.Contract.Customers(&_KYC.CallOpts, arg0)
}

// GetKYCStatus is a free data retrieval call binding the contract method 0x000c8df0.
//
// Solidity: function getKYCStatus(address _customer) view returns(bool, uint256, address)
func (_KYC *KYCCaller) GetKYCStatus(opts *bind.CallOpts, _customer common.Address) (bool, *big.Int, common.Address, error) {
	var out []interface{}
	err := _KYC.contract.Call(opts, &out, "getKYCStatus", _customer)

	if err != nil {
		return *new(bool), *new(*big.Int), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return out0, out1, out2, err

}

// GetKYCStatus is a free data retrieval call binding the contract method 0x000c8df0.
//
// Solidity: function getKYCStatus(address _customer) view returns(bool, uint256, address)
func (_KYC *KYCSession) GetKYCStatus(_customer common.Address) (bool, *big.Int, common.Address, error) {
	return _KYC.Contract.GetKYCStatus(&_KYC.CallOpts, _customer)
}

// GetKYCStatus is a free data retrieval call binding the contract method 0x000c8df0.
//
// Solidity: function getKYCStatus(address _customer) view returns(bool, uint256, address)
func (_KYC *KYCCallerSession) GetKYCStatus(_customer common.Address) (bool, *big.Int, common.Address, error) {
	return _KYC.Contract.GetKYCStatus(&_KYC.CallOpts, _customer)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_KYC *KYCTransactor) Register(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _KYC.contract.Transact(opts, "register")
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_KYC *KYCSession) Register() (*types.Transaction, error) {
	return _KYC.Contract.Register(&_KYC.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x1aa3a008.
//
// Solidity: function register() returns()
func (_KYC *KYCTransactorSession) Register() (*types.Transaction, error) {
	return _KYC.Contract.Register(&_KYC.TransactOpts)
}

// VerifyKYC is a paid mutator transaction binding the contract method 0x38d16011.
//
// Solidity: function verifyKYC(address _customer) returns()
func (_KYC *KYCTransactor) VerifyKYC(opts *bind.TransactOpts, _customer common.Address) (*types.Transaction, error) {
	return _KYC.contract.Transact(opts, "verifyKYC", _customer)
}

// VerifyKYC is a paid mutator transaction binding the contract method 0x38d16011.
//
// Solidity: function verifyKYC(address _customer) returns()
func (_KYC *KYCSession) VerifyKYC(_customer common.Address) (*types.Transaction, error) {
	return _KYC.Contract.VerifyKYC(&_KYC.TransactOpts, _customer)
}

// VerifyKYC is a paid mutator transaction binding the contract method 0x38d16011.
//
// Solidity: function verifyKYC(address _customer) returns()
func (_KYC *KYCTransactorSession) VerifyKYC(_customer common.Address) (*types.Transaction, error) {
	return _KYC.Contract.VerifyKYC(&_KYC.TransactOpts, _customer)
}

// KYCKYCRegisteredIterator is returned from FilterKYCRegistered and is used to iterate over the raw logs and unpacked data for KYCRegistered events raised by the KYC contract.
type KYCKYCRegisteredIterator struct {
	Event *KYCKYCRegistered // Event containing the contract specifics and raw log

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
func (it *KYCKYCRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KYCKYCRegistered)
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
		it.Event = new(KYCKYCRegistered)
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
func (it *KYCKYCRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KYCKYCRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KYCKYCRegistered represents a KYCRegistered event raised by the KYC contract.
type KYCKYCRegistered struct {
	Customer  common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterKYCRegistered is a free log retrieval operation binding the contract event 0xf9939794799d745065fb0b05ca2172b35e1d42e6869932affe454edd17e50842.
//
// Solidity: event KYCRegistered(address indexed customer, uint256 timestamp)
func (_KYC *KYCFilterer) FilterKYCRegistered(opts *bind.FilterOpts, customer []common.Address) (*KYCKYCRegisteredIterator, error) {

	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}

	logs, sub, err := _KYC.contract.FilterLogs(opts, "KYCRegistered", customerRule)
	if err != nil {
		return nil, err
	}
	return &KYCKYCRegisteredIterator{contract: _KYC.contract, event: "KYCRegistered", logs: logs, sub: sub}, nil
}

// WatchKYCRegistered is a free log subscription operation binding the contract event 0xf9939794799d745065fb0b05ca2172b35e1d42e6869932affe454edd17e50842.
//
// Solidity: event KYCRegistered(address indexed customer, uint256 timestamp)
func (_KYC *KYCFilterer) WatchKYCRegistered(opts *bind.WatchOpts, sink chan<- *KYCKYCRegistered, customer []common.Address) (event.Subscription, error) {

	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}

	logs, sub, err := _KYC.contract.WatchLogs(opts, "KYCRegistered", customerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KYCKYCRegistered)
				if err := _KYC.contract.UnpackLog(event, "KYCRegistered", log); err != nil {
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

// ParseKYCRegistered is a log parse operation binding the contract event 0xf9939794799d745065fb0b05ca2172b35e1d42e6869932affe454edd17e50842.
//
// Solidity: event KYCRegistered(address indexed customer, uint256 timestamp)
func (_KYC *KYCFilterer) ParseKYCRegistered(log types.Log) (*KYCKYCRegistered, error) {
	event := new(KYCKYCRegistered)
	if err := _KYC.contract.UnpackLog(event, "KYCRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// KYCKYCVerifiedIterator is returned from FilterKYCVerified and is used to iterate over the raw logs and unpacked data for KYCVerified events raised by the KYC contract.
type KYCKYCVerifiedIterator struct {
	Event *KYCKYCVerified // Event containing the contract specifics and raw log

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
func (it *KYCKYCVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(KYCKYCVerified)
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
		it.Event = new(KYCKYCVerified)
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
func (it *KYCKYCVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *KYCKYCVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// KYCKYCVerified represents a KYCVerified event raised by the KYC contract.
type KYCKYCVerified struct {
	Customer  common.Address
	Verifier  common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterKYCVerified is a free log retrieval operation binding the contract event 0xd8b879062b804f0a2bd0b75ba02d6dbd9999fc9971ece871c8d1aedb447c5f93.
//
// Solidity: event KYCVerified(address indexed customer, address indexed verifier, uint256 timestamp)
func (_KYC *KYCFilterer) FilterKYCVerified(opts *bind.FilterOpts, customer []common.Address, verifier []common.Address) (*KYCKYCVerifiedIterator, error) {

	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}
	var verifierRule []interface{}
	for _, verifierItem := range verifier {
		verifierRule = append(verifierRule, verifierItem)
	}

	logs, sub, err := _KYC.contract.FilterLogs(opts, "KYCVerified", customerRule, verifierRule)
	if err != nil {
		return nil, err
	}
	return &KYCKYCVerifiedIterator{contract: _KYC.contract, event: "KYCVerified", logs: logs, sub: sub}, nil
}

// WatchKYCVerified is a free log subscription operation binding the contract event 0xd8b879062b804f0a2bd0b75ba02d6dbd9999fc9971ece871c8d1aedb447c5f93.
//
// Solidity: event KYCVerified(address indexed customer, address indexed verifier, uint256 timestamp)
func (_KYC *KYCFilterer) WatchKYCVerified(opts *bind.WatchOpts, sink chan<- *KYCKYCVerified, customer []common.Address, verifier []common.Address) (event.Subscription, error) {

	var customerRule []interface{}
	for _, customerItem := range customer {
		customerRule = append(customerRule, customerItem)
	}
	var verifierRule []interface{}
	for _, verifierItem := range verifier {
		verifierRule = append(verifierRule, verifierItem)
	}

	logs, sub, err := _KYC.contract.WatchLogs(opts, "KYCVerified", customerRule, verifierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(KYCKYCVerified)
				if err := _KYC.contract.UnpackLog(event, "KYCVerified", log); err != nil {
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

// ParseKYCVerified is a log parse operation binding the contract event 0xd8b879062b804f0a2bd0b75ba02d6dbd9999fc9971ece871c8d1aedb447c5f93.
//
// Solidity: event KYCVerified(address indexed customer, address indexed verifier, uint256 timestamp)
func (_KYC *KYCFilterer) ParseKYCVerified(log types.Log) (*KYCKYCVerified, error) {
	event := new(KYCKYCVerified)
	if err := _KYC.contract.UnpackLog(event, "KYCVerified", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
