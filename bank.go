package main

import (
	"strconv"
	"time"
	"sync"
)

type Account struct {
	sync.Mutex
	Id int `json:"id"`
	Balance float64
	Created time.Time
}

type Transfer struct {
	FromId int
	ToId int
	Sum float64
}

type Money struct {
	value float64
	currency string
}

type Bank struct {
	addAccMtx sync.Mutex
	accounts[] Account
	transfers[] Transfer
}

func (this *Bank) addAccount(balance string) int {
	x, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		return -1
	}

	this.addAccMtx.Lock()
	id := len(this.accounts)
	acc := Account{
		Id: id,
		Balance: x,
		Created: time.Now(),
	}
	this.accounts = append(this.accounts, acc)
	this.addAccMtx.Unlock()

	return id 
}

func (this *Bank) getBalance(id int) float64 {
	max := len(this.accounts)-1
	if (id > max || id < 0) {
		// wrong id
		return -1
	}
	return this.accounts[id].Balance
}

func (this *Bank) getAccount(id int) *Account {
	if (id < 0 || id > len(this.accounts)-1) {
		return nil 
	}
	return &this.accounts[id]
}

func (this *Bank) transfer(fromId int, toId int, sum float64) bool {
	max := len(this.accounts)-1
	if (max<toId || max<fromId || toId < 0 || fromId < 0 || sum <= 0) {
		// wrong ids or sum
		return false
	}

	this.accounts[fromId].Lock()
	defer this.accounts[fromId].Unlock()

	if (this.accounts[fromId].Balance >= sum) {
		this.accounts[fromId].Balance -= sum;
		this.accounts[toId].Lock()
		this.accounts[toId].Balance += sum;
		this.accounts[toId].Unlock()
		return true
	}
	// no enough money
	return false
}
