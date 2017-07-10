package spider

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	iconv "github.com/djimenez/iconv-go"
	"github.com/zheng-ji/goSnowFlake"
)

type Product struct {
	ID       int64
	Name     string
	Catalog  string
	Brand    string
	Price    string
	Info     string
	Url      string
	Snapshot string
	SnapTime time.Time
}

var Worker *goSnowFlake.IdWorker

func init() {
	var err error
	Worker, err = goSnowFlake.NewIdWorker(1)
	if err != nil {
		fmt.Println(err)
	}
}

func NewID() int64 {
	if id, err := Worker.NextId(); err != nil {
		fmt.Println(err)
		return -1
	} else {
		return id
	}
}

func NewIDs(count int) []int64 {
	var ids []int64 = make([]int64, count)
	for i := 0; i < count; i++ {
		if id, err := Worker.NextId(); err != nil {
			fmt.Println(err)
		} else {
			ids[i] = id
		}
	}
	return ids
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

type Spider interface {
	Run(words []string) ([]Product, error)
}

type SpiderBase struct {
	SearchUrl string
	ItemUrl   string
	Regexp    *regexp.Regexp
}

func (s *SpiderBase) FormatUrls(words []string, url string) []string {
	words = RemoveDuplicatesAndEmpty(words)
	length := len(words)
	urls := make([]string, length)
	for i := 0; i < length; i++ {
		urls[i] = fmt.Sprintf(url, words[i])
	}
	return urls
}

func (s *SpiderBase) Body(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return string(body)
}

func (s *SpiderBase) ProductUrls(searchUrls []string, productReg *regexp.Regexp) []string {
	result := make([]string, 0, 100)
	for i := 0; i < len(searchUrls); i++ {
		body := s.Body(searchUrls[i])
		matches := productReg.FindAllStringSubmatch(body, -1)
		for i := 1; i < len(matches); i++ {
			item := fmt.Sprintf(s.ItemUrl, matches[i][1])
			result = append(result, item)
		}
	}
	return RemoveDuplicatesAndEmpty(result)
}

func (s *SpiderBase) Unicode2String(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func (s *SpiderBase) CodeConvert(str, from, to string) string {
	result, _ := iconv.ConvertString(str, from, to)
	return result
}
