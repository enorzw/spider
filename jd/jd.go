package jd

import (
	"regexp"

	"github.com/enorzw/spider"
)

type JdSpider struct {
	spider.SpiderBase
}

func NewJdSpider() {
	spider := new(JdSpider)
	spider.SearchUrl = "https://search.jd.com/Search?keyword=%s&enc=utf-8&qrst=1&rt=1&stop=1&vt=1&page=1&s=1&click=0"
	spider.ItemUrl = "https://item.jd.com/{0}.html"
	spider.Regexp = regexp.MustCompile(`data-sku="([0-9]+)"`)
}

func (s *JdSpider) Run(words []string) ([]spider.Product, error) {
	searchUrls := s.FormatUrls(words, s.SearchUrl)
	productUrls := s.ProductUrls(searchUrls, s.Regexp)
	return s.Product(productUrls), nil
}

func (s *JdSpider) Product(productUrls []string) []spider.Product {
	return nil
}
