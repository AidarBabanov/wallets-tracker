package csv

import (
	"fmt"
	"github.com/AidarBabanov/wallets-tracker/internal/addrdb"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gocarina/gocsv"
	"os"
	"sync"
)

type AddressDatabase struct {
	addresses []addrdb.Address
	index     int
	mu        sync.RWMutex
}

func New() *AddressDatabase {
	adb := new(AddressDatabase)
	adb.mu = sync.RWMutex{}
	return adb
}

func (adb *AddressDatabase) ReadCSV(path string) error {
	adb.mu.Lock()
	adb.mu.Unlock()
	addressesFile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	err = gocsv.UnmarshalFile(addressesFile, &adb.addresses)
	if err != nil {
		return err
	}
	err = addressesFile.Close()
	if err != nil {
		return err
	}
	logs.Info("Found %d addresses in %s file.", len(adb.addresses), path)
	return nil
}

func (adb *AddressDatabase) Len() int {
	adb.mu.RLock()
	adb.mu.RUnlock()
	return len(adb.addresses)
}

func (adb *AddressDatabase) IsEmpty() bool {
	adb.mu.RLock()
	defer adb.mu.RUnlock()
	return adb.isEmpty()
}

func (adb *AddressDatabase) isEmpty() bool {
	return len(adb.addresses) == 0
}

func (adb *AddressDatabase) Next() (addrdb.Address, error) {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	if adb.isEmpty() {
		return addrdb.Address{}, fmt.Errorf("address database is empty")
	}
	if adb.index >= len(adb.addresses) {
		adb.index = 0
	}
	addr := adb.addresses[adb.index]
	adb.index++
	if adb.index >= len(adb.addresses) {
		adb.index = 0
	}
	return addr, nil
}

func (adb *AddressDatabase) All() []addrdb.Address {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	return adb.addresses
}

func (adb *AddressDatabase) Index() int {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	return adb.index
}
