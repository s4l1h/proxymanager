s4l1hpackage checker_test

import (
	"testing"

	"errors"

	"fmt"

	"github.com/akmyazilim/proxymanager"
	"github.com/akmyazilim/proxymanager/checker"
)

func ExampleChecker_Check() {
	c := checker.New()
	c.Add(checker.Function{
		Name: "checkPort1024",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1024" {
				return errors.New("Port is Not 1024"), false
			}
			return nil, true
		},
	})
	if c.Check(proxymanager.Proxy{
		Host: "1.1.1.1",
		Port: "1024",
	}) == true {
		fmt.Println("1.1.1.1 Matched")
	}
	if c.Check(proxymanager.Proxy{
		Host: "1.1.1.2",
		Port: "1026",
	}) == true {
		fmt.Println("1.1.1.2 Matched")
	}

	// Output: 1.1.1.1 Matched
}
func ExampleChecker_Run() {

	manager := proxymanager.New(3)
	manager.AddFromURL("http://u:p@host1.com:1024")
	manager.AddFromURL("http://u:p@host2.com:1024")
	manager.AddFromURL("http://u:p@host3.com:1025")
	manager.AddFromURL("http://u:p@host4.com:1020")
	manager.AddFromURL("socks5://u:p@host1.com:1024")

	c := checker.New()
	c.Add(checker.Function{
		Name: "checkPort1024",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1024" {
				return errors.New("Port is Not 1024"), false
			}
			return nil, true
		},
	})
	// Run all checkers and delete unmatched
	manager = c.Run(manager)

	if manager.Has("http://u:p@host3.com:1025") == true {
		fmt.Println("http://u:p@host3.com:1025 Exists")
	}
	if manager.Has("socks5://u:p@host1.com:1024") == true {
		fmt.Println("socks5://u:p@host1.com:1024 Exists")
	}

	// Output: socks5://u:p@host1.com:1024 Exists

}

func ExampleChecker_Add() {

	fn := checker.Function{
		Name: "checkPort1030",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1030" {
				return errors.New("Port is Not 1030"), false
			}
			return nil, true
		},
	}

	c := checker.New()
	c.Add(checker.Function{
		Name: "checkPort1024",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1024" {
				return errors.New("Port is Not 1024"), false
			}
			return nil, true
		},
	})
	c.Add(fn)
}

func ExampleChecker_Remove() {

	fn := checker.Function{
		Name: "checkPort1030",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1030" {
				return errors.New("Port is Not 1030"), false
			}
			return nil, true
		},
	}

	c := checker.New()
	c.Add(checker.Function{
		Name: "checkPort1024",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1024" {
				return errors.New("Port is Not 1024"), false
			}
			return nil, true
		},
	})
	c.Add(fn)
	c.Remove(fn)
	c.Remove("checkPort1024")
}
func ExampleChecker_Has() {

	fn := checker.Function{
		Name: "checkPort1030",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1030" {
				return errors.New("Port is Not 1030"), false
			}
			return nil, true
		},
	}

	c := checker.New()

	c.Add(fn)
	if c.Has("checkPort1030") {
		fmt.Println("we have checkerPort1030")
	} else {
		fmt.Println("we don't have checkPort1030")
	}
	// Output: we have checkerPort1030
}
func TestChecker(t *testing.T) {

	FN1 := checker.Function{
		Name: "checkWork",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			return nil, true
		},
	}

	c := checker.New()
	c.Add(FN1)

	c.Add(checker.Function{
		Name: "checkHTTPS",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			return nil, true
		},
	})
	c.Add(checker.Function{
		Name: "checkPort1024",
		Fn: func(proxy proxymanager.Proxy) (error, bool) {
			if proxy.Port != "1024" {
				return errors.New("Port is Not 1024"), false
			}
			return nil, true
		},
	})

	if c.Has(FN1) != true {
		t.Error("Has Function Error")
	}
	if c.Has("checkWork") != true {
		t.Error("Has Function Error")
	}
	if c.Has("checkPort1024") != true {
		t.Error("Has Function Error")
	}
	if c.Has("checkHTTPSvvvvvvvvWW") != false {
		t.Error("Has Function Error")
	}

	c.Remove("checkHTTPS")

	if c.Has("checkHTTPS") == true {
		t.Error("Has Function Error")
	}

	c.Remove(FN1)

	if c.Has("checkWork") == true {
		t.Error("Has Function Error")
	}

}
