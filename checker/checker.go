package checker

import (
	"sync"

	"github.com/akmyazilim/proxymanager"
)

// FunctionType checker fn type
type FunctionType func(proxy proxymanager.Proxy) (error, bool)

// Function checker type
type Function struct {
	Name string
	Fn   FunctionType
}

// New return checker
func New() *Checker {
	return &Checker{}
}

// Checker object
type Checker struct {
	sync.Mutex // Embeding olarak ekleyelim.
	Functions  []Function
}

//Add checker function
func (checker *Checker) Add(c Function) {
	checker.Lock()         // Diğer goroutines'lerin erişmesini engelleyelim.
	defer checker.Unlock() // İşlem bittikten sonra erişim engelini kaldıralım
	checker.Functions = append(checker.Functions, c)
}

// Remove checker function
func (checker *Checker) Remove(c interface{}) {
	checker.Lock()         // Diğer goroutines'lerin erişmesini engelleyelim.
	defer checker.Unlock() // İşlem bittikten sonra erişim engelini kaldıralım

	name := ""
	switch c.(type) {
	case Function:
		name = c.(Function).Name
	case string:
		name = c.(string)
	}
	for i, v := range checker.Functions {
		if v.Name == name {
			checker.Functions = append(checker.Functions[:i], checker.Functions[i+1:]...)
			break
		}
	}
}

// Has checker function
func (checker *Checker) Has(c interface{}) bool {
	checker.Lock()         // Diğer goroutines'lerin erişmesini engelleyelim.
	defer checker.Unlock() // İşlem bittikten sonra erişim engelini kaldıralım
	name := ""
	switch c.(type) {
	case Function:
		name = c.(Function).Name
	case string:
		name = c.(string)
	}
	for _, v := range checker.Functions {
		if v.Name == name {
			return true
		}
	}
	return false

}

// Check proxy is matched
func (checker *Checker) Check(p proxymanager.Proxy) bool {
	for _, c := range checker.Functions {
		if _, result := c.Fn(p); result == false {
			return false
		}
	}
	return true
}

// Run run all checkers and return proxy manager
func (checker *Checker) Run(manager *proxymanager.Manager) *proxymanager.Manager {
	wg := sync.WaitGroup{}
	for _, p := range manager.List {
		wg.Add(1)
		pp := p // range :S
		go func(pp proxymanager.Proxy) {
			defer wg.Done()
			if checker.Check(pp) == false {
				manager.Remove(pp)
			}
		}(pp)
	}
	wg.Wait()
	return manager
}
