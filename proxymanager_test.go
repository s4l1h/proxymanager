package proxymanager_test

import (
	"testing"

	"fmt"

	"github.com/akmyazilim/proxymanager"
)

func ExampleProxy_String() {
	p := proxymanager.Proxy{
		Type:     "http",
		Host:     "10.0.0.1",
		Port:     "1080",
		Username: "username",
		Password: "password",
	}
	fmt.Println(p)
	// Output: http://username:password@10.0.0.1:1080
}
func ExampleManager_Has() {
	u := "http://1.1.1.1:1010"
	p := proxymanager.Proxy{
		Host: "10.0.0.1",
		Port: "1080",
	}
	plist := proxymanager.New(3)
	plist.AddFromURL(u)
	if plist.Has(u) == true {
		fmt.Printf("%s Exists", u)
	}
	if plist.Has("http://proxyhost:333") == true {
		fmt.Println("http://proxyhost:333 Exists")
	}
	if plist.Has(p) == true {
		fmt.Printf("Exists %s:%s", p.Host, p.Port)
	}
	// Output: http://1.1.1.1:1010 Exists

}
func TestAddFromURL(t *testing.T) {

	u := "socks5://username:password@host:1020"
	plist := proxymanager.New(3)
	plist.AddFromURL(u)
	url := plist.GiveMeProxyURL().String()

	if url != u {
		t.Errorf("AddFromURL error %s!=%s", url, u)
	}
}
func ExampleManager_AddFromURL() {

	plist := proxymanager.New(3)
	plist.AddFromURL("http://username:password@192.168.1.11:1020")
	proxy := plist.GiveMeProxy()
	fmt.Printf("Host:%s Port:%s Username:%s Password:%s Type:%s", proxy.Host, proxy.Port, proxy.Username, proxy.Password, proxy.Type)
	// Output: Host:192.168.1.11 Port:1020 Username:username Password:password Type:http
}
func ExampleManager_GiveMeProxyURL() {

	plist := proxymanager.New(3)
	plist.Add(proxymanager.Proxy{
		Host:     "192.168.1.11",
		Port:     "1000",
		Username: "user",
		Password: "pass",
	})
	url := plist.GiveMeProxyURL()
	fmt.Println(url)
	// Output: http://user:pass@192.168.1.11:1000

}

func ExampleManager_Add() {

	plist := proxymanager.New(3)
	plist.Add(proxymanager.Proxy{
		Host:     "192.168.1.11",
		Port:     "1000",
		Username: "none",
		Password: "none",
	})
	proxy := plist.GiveMeProxy()
	fmt.Println(proxy.Host)
	// Output: 192.168.1.11
}
func ExampleManager_Remove() {

	plist := proxymanager.New(2)
	plist.Add(proxymanager.Proxy{
		Host:     "192.168.1.10",
		Port:     "1000",
		Username: "none",
		Password: "none",
	})
	plist.Add(proxymanager.Proxy{
		Host:     "192.168.1.11",
		Port:     "1000",
		Username: "none",
		Password: "none",
	})
	plist.Add(proxymanager.Proxy{
		Host:     "192.168.1.12",
		Port:     "1000",
		Username: "none",
		Password: "none",
	})
	plist.Remove("192.168.1.12")
	plist.Remove(proxymanager.Proxy{Host: "192.168.11"})

}
func TestProxy(t *testing.T) {
	example := []string{
		"192.168.1.100",
		"192.168.1.101",
		"192.168.1.102",
	}

	total := len(example)
	limit := 3

	plist := proxymanager.New(limit)

	for _, p := range example {
		plist.Add(
			proxymanager.Proxy{
				Host:     p,
				Port:     "1000",
				Username: "none",
				Password: "none",
			})
	}
	if plist.WriteIndex != total {
		t.Errorf("WriteIndex Error: %d!=%d", plist.WriteIndex, total)
	}
	if len(plist.List) != total {
		t.Errorf("List size error")
	}

	if plist.Limit != limit {
		t.Errorf("Limit Error")
	}

	if plist.StepIndex != 0 {
		t.Errorf("StepIndex error")
	}
	if plist.ReadIndex != 0 {
		t.Errorf("ReadIndex error")
	}
	proxy := plist.GiveMeProxy()
	if proxy.Host != example[0] {
		t.Errorf("GiveMeProxy return wrong proxy")
	}
	if plist.ReadIndex != 0 {
		t.Errorf("ReadIndex error %d", plist.ReadIndex)
	}
	if plist.StepIndex != 1 {
		t.Errorf("StepIndex error %d", plist.StepIndex)
	}
	proxy = plist.GiveMeProxy()
	proxy = plist.GiveMeProxy()
	proxy = plist.GiveMeProxy() // index 1 step 1
	if proxy.Host != example[1] {
		t.Errorf("GiveMeProxy return wrong proxt")
	}
	if plist.ReadIndex != 1 {
		t.Errorf("ReadIndex error %d", plist.ReadIndex)
	}
	if plist.StepIndex != 1 {
		t.Errorf("StepIndex error %d", plist.StepIndex)
	}

	proxy = plist.GiveMeProxy() // index 1 step 2
	proxy = plist.GiveMeProxy() // index 1 step 3
	proxy = plist.GiveMeProxy() // index 2 step 1

	if proxy.Host != example[2] {
		t.Errorf("GiveMeProxy return wrong proxt")
	}
	if plist.ReadIndex != 2 {
		t.Errorf("ReadIndex error %d", plist.ReadIndex)
	}
	if plist.StepIndex != 1 {
		t.Errorf("StepIndex error %d", plist.StepIndex)
	}

	proxy = plist.GiveMeProxy() // index 2 step 2
	proxy = plist.GiveMeProxy() // index 2 step 3
	if plist.ReadIndex != plist.WriteIndex {
		t.Errorf("ReadIndex error %d", plist.ReadIndex)
		t.Errorf("WriteIndex error %d", plist.WriteIndex)
	}
	proxy = plist.GiveMeProxy() // index 0 step 1
	if plist.ReadIndex != 0 {
		t.Errorf("ReadIndex error %d", plist.ReadIndex)
	}

	plist.Remove("192.168.1.100")
	if plist.WriteIndex != 2 {
		t.Errorf("WriteIndex error %d", plist.WriteIndex)
	}
	plist.Remove(proxymanager.Proxy{Host: "192.168.1.101"})
	if plist.WriteIndex != 1 {
		t.Errorf("WriteIndex error %d", plist.WriteIndex)
	}

}
