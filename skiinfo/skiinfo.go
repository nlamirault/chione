// Copyright (C) 2017 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package skiinfo

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	baseURL           = "http://www.skiinfo.fr"
	countryResortsURL = "http://www.skiinfo.fr/%s/stations-de-ski.html"
)

type Resort struct {
	Name      string
	URL       string
	Region    string
	RegionURL string
}

func fetch(uri string, data url.Values) ([]byte, error) {
	u, _ := url.ParseRequestURI(uri)
	urlStr := fmt.Sprintf("%v", u)

	client := &http.Client{}
	logrus.Debugf("URI: %s %s", urlStr, data)

	r, _ := http.NewRequest("GET", urlStr, bytes.NewBufferString(data.Encode()))
	resp, err := client.Do(r)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Http request to %s failed: %s", r.URL, err.Error())
	}
	logrus.Debugf("HTTP Status: %s", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, fmt.Errorf("errorination happened reading the body: %s", err.Error())
	}
	return body, nil
}

func ListResorts(country string) (map[string]*Resort, error) {
	logrus.Debugf("Retrieve resorts for country: %s", country)
	resorts := map[string]*Resort{}
	uri := fmt.Sprintf(countryResortsURL, country)
	body, err := fetch(uri, url.Values{})
	if err != nil {
		return nil, err
	}
	z := html.NewTokenizer(strings.NewReader(string(body)))
	var name string
	var url string
	for {
		// token type
		tokenType := z.Next()
		if tokenType == html.ErrorToken {
			break
		}
		// token := z.Token()
		switch tokenType {
		case html.StartTagToken: // <tag>
			t := z.Token()

			if t.Data == "div" {
				if len(t.Attr) > 0 && t.Attr[0].Val == "name" {
					z.Next()
					link := z.Token()
					// fmt.Printf("Resort link: %s %s\n", link, link.Attr)
					if len(link.Attr) > 1 && link.Attr[1].Key == "title" {
						name = link.Attr[1].Val
						url = link.Attr[0].Val
						// fmt.Printf("Resort name: %s\n", name)
					}
				} else if len(t.Attr) > 0 && t.Attr[0].Val == "rRegion" {
					z.Next()
					link := z.Token()
					if len(link.Attr) > 1 && link.Attr[1].Key == "title" {
						regionName := link.Attr[1].Val
						regionUrl := link.Attr[0].Val
						// fmt.Printf("Resort region: %s %s\n", name, regionName)
						resorts[name] = &Resort{
							Name:      name,
							URL:       url,
							Region:    regionName,
							RegionURL: regionUrl,
						}
					}

				}

			}
		case html.TextToken: // text between start and end tag
		case html.EndTagToken: // </tag>
		case html.SelfClosingTagToken: // <tag/>
		}
	}

	return resorts, nil
}
