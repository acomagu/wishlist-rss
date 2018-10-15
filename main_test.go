package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var wishListerMock = `[
  {
    "num": 1,
    "name": "ｍｏｎｏ　１巻【Amazon.co.jp限定描き下ろし特典付】 (まんがタイムKRコミックス)",
    "link": "http://www.amazon.co.jp/dp/B07HNVWQTW/?coliid=I1NV433AHPN31U&colid=IQ3DUL8LXQFB&psc=0&ref_=lv_vv_lig_dp_it",
    "old-price": "N/A",
    "new-price": "<span class=\"a-offscreen\">￥864</span><span aria-hidden=\"true\"><span class=\"a-price-symbol\">￥</span><span class=\"a-price-whole\">864</span></span>",
    "date-added": "2018年10月15日に追加された商品",
    "priority": "中",
    "rating": "N/A",
    "total-ratings": "",
    "comment": "",
    "picture": "https://images-na.ssl-images-amazon.com/images/I/01MKUOLsA5L._SS135_.gif",
    "page": 1
  },
  {
    "num": 2,
    "name": "mono (1) (まんがタイムKRコミックス)",
    "link": "http://www.amazon.co.jp/dp/4832249894/?coliid=I2RM3Q4YBLTUU2&colid=IQ3DUL8LXQFB&psc=0&ref_=lv_vv_lig_dp_it",
    "old-price": "N/A",
    "new-price": "<span class=\"a-offscreen\">￥885</span><span aria-hidden=\"true\"><span class=\"a-price-symbol\">￥</span><span class=\"a-price-whole\">885</span></span>",
    "date-added": "2018年10月15日に追加された商品",
    "priority": "中",
    "rating": "N/A",
    "total-ratings": "",
    "comment": "",
    "picture": "https://images-na.ssl-images-amazon.com/images/I/91XT4laA+oL._SS135_.jpg",
    "page": 1
  },
  {
    "num": 3,
    "name": "未熟なふたりでございますが（１） (コミックＤＡＹＳコミックス)",
    "link": "http://www.amazon.co.jp/dp/B07FSM2K3F/?coliid=I12FCJFB1CTU3D&colid=IQ3DUL8LXQFB&psc=0&ref_=lv_vv_lig_dp_it",
    "old-price": "N/A",
    "new-price": "<span class=\"a-offscreen\">￥648</span><span aria-hidden=\"true\"><span class=\"a-price-symbol\">￥</span><span class=\"a-price-whole\">648</span></span>",
    "date-added": "2018年10月12日に追加された商品",
    "priority": "中",
    "rating": "N/A",
    "total-ratings": "8",
    "comment": "",
    "picture": "https://images-na.ssl-images-amazon.com/images/I/91Mlw3lvi5L._SS135_.jpg",
    "page": 1
  }
]`

var result = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
	xmlns:dc="http://purl.org/dc/elements/1.1/"
	xmlns:sy="http://purl.org/rss/1.0/modules/syndication/"
	xmlns:admin="http://webns.net/mvcb/"
	xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
	<channel>
		<title>Wishlist</title>
		<dc:language>ja</dc:language>
		<dc:date>15 Oct 18 00:00 JST</dc:date>
		<item>
			<title>ｍｏｎｏ　１巻【Amazon.co.jp限定描き下ろし特典付】 (まんがタイムKRコミックス)</title>
			<link>http://www.amazon.co.jp/dp/B07HNVWQTW/?coliid=I1NV433AHPN31U&amp;colid=IQ3DUL8LXQFB&amp;psc=0&amp;ref_=lv_vv_lig_dp_it</link>
			<dc:date>15 Oct 18 00:00 JST</dc:date>
		</item>
		<item>
			<title>mono (1) (まんがタイムKRコミックス)</title>
			<link>http://www.amazon.co.jp/dp/4832249894/?coliid=I2RM3Q4YBLTUU2&amp;colid=IQ3DUL8LXQFB&amp;psc=0&amp;ref_=lv_vv_lig_dp_it</link>
			<dc:date>15 Oct 18 00:00 JST</dc:date>
		</item>
		<item>
			<title>未熟なふたりでございますが（１） (コミックＤＡＹＳコミックス)</title>
			<link>http://www.amazon.co.jp/dp/B07FSM2K3F/?coliid=I12FCJFB1CTU3D&amp;colid=IQ3DUL8LXQFB&amp;psc=0&amp;ref_=lv_vv_lig_dp_it</link>
			<dc:date>12 Oct 18 00:00 JST</dc:date>
		</item>
	</channel>
</rss>
`

func Test(t *testing.T) {
	is := is.New(t)

	hitWishLister = func(id string) (*http.Response, error) {
		return &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Body:          ioutil.NopCloser(strings.NewReader(wishListerMock)),
			ContentLength: -1,
		}, nil
	}

	h := handler()
	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "http://localhost/code", nil)
	if err != nil {
		t.Fatal(err)
	}
	h.ServeHTTP(w, req)

	is.Equal(w.Body.String(), result)
}
