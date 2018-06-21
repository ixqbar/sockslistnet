package spider

import (
	"sync"
	"net/http"
	"time"
	"github.com/Pallinder/go-randomdata"
	"github.com/PuerkitoBio/goquery"
	"net/http/cookiejar"
	"log"
	"strings"
	"fmt"
	"github.com/robertkrimen/otto"
)

type TProxy struct {
	sync.RWMutex
	httpClient *http.Client
	cookies    []*http.Cookie
}

func NewProxy() *TProxy {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	return &TProxy{
		httpClient: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				IdleConnTimeout: 10 * time.Second,
			},
		},
	}
}

func (obj *TProxy) Pull() {
	obj.Lock()
	defer obj.Unlock()

	request, err := http.NewRequest("GET", "https://sockslist.net/list/proxy-socks-5-list/", nil)
	if err != nil {
		Logger.Print(err)
		return
	}

	request.Header.Set("Referer", "https://sockslist.net/")
	request.Header.Set("User-Agent", randomdata.UserAgentString())
	response, err := obj.httpClient.Do(request)
	if err != nil {
		Logger.Print(err)
		return
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		Logger.Print(err)
		return
	}

	vm := otto.New()
	foundVar := false

	document.Find("script").Each(func(i int, selection *goquery.Selection) {
		content := selection.Eq(0).Text()
		if foundVar || strings.Contains(content, "//<![CDATA[") == false {
			return
		}

		foundVar = true
		content = strings.TrimSpace(content)
		content = strings.Replace(content, "//<![CDATA[", "", -1)
		content = strings.Replace(content, "//]]>", "", -1)
		content = strings.TrimSpace(content)

		fmt.Println(content)

		vm.Run(fmt.Sprintf(`%sfunction getPort(p){return p;}`, content))

		document.Find(".proxytbl").Eq(0).Find("tr").Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				return
			}

			address := selection.Find("td").Eq(0).Text()

			country := selection.Find("td").Eq(2).Text()
			country = strings.TrimSpace(country)

			port := selection.Find("td").Eq(1).Text();
			port = strings.TrimSpace(port)
			port = strings.Replace(port, "//<![CDATA[", "", -1)
			port = strings.Replace(port, "//]]>", "", -1)
			port = strings.Replace(port, "document.write", "getPort", -1)
			port = strings.TrimSpace(port)

			value, err := vm.Run(fmt.Sprintf("%s", port))
			if err != nil {
				Logger.Print(err)
				return
			}

			realPort, err := value.ToString()
			if err != nil {
				Logger.Print(err)
				return
			}

			GlobalVars.TQueue.Push(TQueueItem{IsSock5, &TProxyItem{"socks5", country, address, realPort, time.Now().Unix()}})
		})
	})
}
