package csv

import (
	"encoding/csv"
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
	mu        *sync.RWMutex
	filePath  string
}

func New(filePath string) *AddressDatabase {
	adb := new(AddressDatabase)
	adb.mu = new(sync.RWMutex)
	adb.filePath = filePath
	return adb
}

func (adb *AddressDatabase) Add(address addrdb.Address) error {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	adb.addresses = append(adb.addresses, address)
	addressesFile, err := os.OpenFile(adb.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer addressesFile.Close()
	csvWriter := gocsv.NewSafeCSVWriter(csv.NewWriter(addressesFile))
	err = gocsv.MarshalCSVWithoutHeaders([]addrdb.Address{address}, csvWriter)
	if err != nil {
		return err
	}
	return nil
}

func (adb *AddressDatabase) ReadCSV() error {
	adb.mu.Lock()
	defer adb.mu.Unlock()
	addressesFile, err := os.OpenFile(adb.filePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
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
	logs.Info("Found %d addresses in %s file.", len(adb.addresses), adb.filePath)
	return nil
}

func (adb *AddressDatabase) Len() int {
	adb.mu.RLock()
	defer adb.mu.RUnlock()
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
