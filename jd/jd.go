package jd

import (
	"encoding/json"
	"fmt"
	"regexp"

	"strings"

	"time"

	"github.com/enorzw/spider"
)

type JdSpider struct {
	spider.SpiderBase
}

var (
	nameReg  *regexp.Regexp = regexp.MustCompile(`[[:space:]]*name:[[:space:]]+\'(?P<name>.*)\',`)
	brandReg *regexp.Regexp = regexp.MustCompile("<li title='(?P<brand>[[:space:]\\(\\)（）\\w\u4e00-\u9fa5]*)'>品牌")
)

func NewJdSpider() *JdSpider {
	spider := new(JdSpider)
	spider.SearchUrl = "https://search.jd.com/Search?keyword=%s&enc=utf-8&qrst=1&rt=1&stop=1&vt=1&page=1&s=1&click=0"
	spider.ItemUrl = "https://item.jd.com/%s.html"
	spider.Regexp = regexp.MustCompile(`data-sku="([0-9]+)"`)
	return spider
}

func (s *JdSpider) Run(words []string) ([]*spider.Product, error) {
	fmt.Println(words)
	searchUrls := s.FormatUrls(words, s.SearchUrl)

	fmt.Println(searchUrls)

	productUrls := s.ProductUrls(searchUrls, s.Regexp)
	return s.Product(productUrls), nil
}

func (s *JdSpider) Product(productUrls []string) []*spider.Product {
	products := make([]*spider.Product, 0, len(productUrls))
	ids := spider.NewIDs(len(productUrls))
	for i, url := range productUrls {
		fmt.Println(url)
		body := s.Body(url)
		body = s.CodeConvert(body, "gbk", "utf-8")
		product := new(spider.Product)

		if nameReg.MatchString(body) {
			name := nameReg.FindStringSubmatch(string(body))[1]
			product.Name, _ = s.Unicode2String(strings.TrimSpace(name))
		}
		if brandReg.MatchString(body) {
			brand := brandReg.FindStringSubmatch(string(body))[1]
			product.Brand = brand
		}

		product.ID = ids[i]
		product.Price = s.Price(body)
		product.Url = url
		product.SnapTime = time.Now()
		products = append(products, product)
	}
	return products
}

func (s *JdSpider) Price(body string) string {
	url := "http://p.3.cn/prices/mgets?skuIds=J_%s&type=%s"
	var SkuidReg = regexp.MustCompile(`[[:space:]]*skuid:[[:space:]]+(?P<skuid>[0-9]+),`)
	var SkuidkeyReg = regexp.MustCompile(`[[:space:]]*skuidkey:[[:space:]]*\'(?P<skuidkey>.*)\',`)

	skuid := SkuidReg.FindStringSubmatch(string(body))[1]
	skuidkey := SkuidkeyReg.FindStringSubmatch(string(body))[1]
	pbody := s.Body(fmt.Sprintf(url, skuid, skuidkey))

	m := make([]map[string]interface{}, 10)
	e := json.Unmarshal([]byte(pbody), &m)
	if e != nil {
		panic(e.Error())
	}

	if val, ok := m[0]["p"].(string); ok {
		return val
	}
	return "0"
}
