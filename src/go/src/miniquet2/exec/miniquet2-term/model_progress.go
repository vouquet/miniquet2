package main

import (
	"fmt"
	"sync"
)

import (
	"miniquet2/miniquet"
)

type ProgressModel struct {
	view_handler func(map[string]*miniquet.Trader)
	traders      map[string]*miniquet.Trader

	mtx  *sync.Mutex
}

func NewProgressModel() *ProgressModel {
	return &ProgressModel{
		traders:make(map[string]*miniquet.Trader, 0),
		mtx:new(sync.Mutex),
	}
}

func (self *ProgressModel) ViewHandler(f func(map[string]*miniquet.Trader)) {
	self.mtx.Lock()
	defer self.mtx.Unlock()

	self.view_handler = f
}

func (self *ProgressModel) Publish() {
	self.mtx.Lock()
	defer self.mtx.Unlock()

	self.publish()
}

func (self *ProgressModel) call_view_handler(ts map[string]*miniquet.Trader) {
	if self.view_handler == nil {
		return
	}
	self.view_handler(ts)
}

func (self *ProgressModel) publish() {
	self.call_view_handler(self.traders)
}

func (self *ProgressModel) Add(n_tr *miniquet.Trader) error {
	self.mtx.Lock()
	defer self.mtx.Unlock()

	_, ok := self.traders[n_tr.Name()]
	if ok {
		return fmt.Errorf("trader(%s) is already exsit.", n_tr.Name())
	}

	self.traders[n_tr.Name()] = n_tr
	return nil
}

func (self *ProgressModel) Remove(r_tr *miniquet.Trader) error {
	self.mtx.Lock()
	defer self.mtx.Unlock()

	_, ok := self.traders[r_tr.Name()]
	if !ok {
		return fmt.Errorf("trader(%s) does not exsit.", r_tr.Name())
	}

	delete(self.traders, r_tr.Name())
	return nil
}
