package pool

import (
	"fmt"
	"strconv"
	"sync"
)

type Cookie struct {
	Name      string
	Value     string
	Raw       string
	Domain    string
	Channel   string
	UsedCount int
}

func (c *Cookie) GetAttr(attr string) string {
	switch attr {
	case "name":
		return c.Name
	case "raw":
		return c.Raw
	case "domain":
		return c.Domain
	case "channel":
		return c.Channel
	case "usedcount":
		return strconv.Itoa(c.UsedCount)
	default:
		return ""
	}
}

type CookiePool struct {
	cookies []*Cookie
	domain  string
	size    int
	mapping map[string]int
	mu      sync.Mutex
}

func (cp *CookiePool) Cap() int {
	return len(cp.cookies)
}

func (cp *CookiePool) Add(cookie *Cookie) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.cookies = append(cp.cookies, cookie)
	cp.size = len(cp.cookies)
}

func (cp *CookiePool) Get() *Cookie {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if cp.size == 0 {
		return nil
	}
	cookie := cp.cookies[0]
	cp.cookies = cp.cookies[1:]
	cookie.UsedCount++
	return cookie
}

func testCase1() {
	cookiePool := &CookiePool{}

	c1 := &Cookie{Domain: "QQ", Raw: "ab=1;bc=3"}
	c2 := &Cookie{Domain: "QQ", Raw: "ab=1;bc=3"}
	c3 := &Cookie{Domain: "QQ", Raw: "ab=1;bc=3"}

	cookiePool.Add(c1)
	cookiePool.Add(c3)
	cookiePool.Add(c2)

	fmt.Println(cookiePool)
	fmt.Println(cookiePool.Get())
	fmt.Println(cookiePool.Get())
	fmt.Println(cookiePool.Cap())
	lastCookie := cookiePool.Get()
	fmt.Println("Attr:", lastCookie.GetAttr("raw"))
	fmt.Println(lastCookie.UsedCount)
	fmt.Println(lastCookie.GetAttr("usedcount"))
	cookiePool.Add(lastCookie)

	fmt.Println(cookiePool.Cap())
	fmt.Println(lastCookie.GetAttr("usedcount"))
	nextCookie := cookiePool.Get()
	fmt.Println("NextCookie:", nextCookie)
}
