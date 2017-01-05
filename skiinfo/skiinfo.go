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
	resortURL         = "http://www.skiinfo.fr/%s/%s/bulletin-neige.html"
	countryResortsURL = "http://www.skiinfo.fr/%s/stations-de-ski.html"
)

type Resort struct {
	Name      string
	URL       string
	Region    string
	RegionURL string
}

type SnowDepth struct {
	Upper  string
	Middle string
	Lower  string
}

type Weather struct {
	Type          string
	Temperature   string
	WindDirection string
	WindForce     string
}

type Slopes struct {
	Beginning    bytes.Buffer
	Intermediate bytes.Buffer
	Advanced     bytes.Buffer
	Expert       bytes.Buffer
}

type SnowFall struct {
	Day  string
	Date string
	Snow string
}

type ResortDescription struct {
	Status         string
	Piste          *SnowDepth
	OffPiste       *SnowDepth
	SummmitWeather *Weather
	BaseWeather    *Weather
	SnowFallDays   []*SnowFall
	Slopes         *Slopes
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

func extractTextTag(z *html.Tokenizer) string {
	text := (string)(z.Text())
	value := strings.TrimSpace(text)
	// fmt.Printf("Val: %s\n", value)
	return value
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

func GetResort(name string, region string) (*ResortDescription, error) {
	logrus.Debugf("Retrieve resort: %s %s", name, region)

	uri := fmt.Sprintf(resortURL, region, name)
	body, err := fetch(uri, url.Values{})
	if err != nil {
		return nil, err
	}
	resortDesc := &ResortDescription{
		OffPiste:       &SnowDepth{},
		Piste:          &SnowDepth{},
		BaseWeather:    &Weather{},
		SummmitWeather: &Weather{},
		SnowFallDays:   []*SnowFall{},
		Slopes:         &Slopes{},
	}
	snowDepth := 0
	// snowFall := 0
	weather := 0
	elevationUpperState := false
	elevationMiddleState := false
	elevationLowerState := false
	slopesBeginning := false
	slopesIntermediate := false
	slopesAdvanced := false
	slopesExpert := false
	snowDay := 0
	z := html.NewTokenizer(strings.NewReader(string(body)))
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
			// fmt.Printf("Token: %s", t)
			if t.Data == "span" {
				if len(t.Attr) > 0 && strings.Contains(t.Attr[0].Val, "current_status") {
					inner := z.Next()
					if inner == html.TextToken {
						resortDesc.Status = extractTextTag(z)
					}
				}
			} else if t.Data == "ul" {
				if len(t.Attr) > 0 {
					if strings.Contains(t.Attr[0].Val, "sr_snow_depth_stations") {
						snowDepth += 1
					} else if strings.Contains(t.Attr[0].Val, "sr_weather_table") {
						weather += 1
					}
				}
			} else if t.Data == "li" {
				if len(t.Attr) > 0 {
					if t.Attr[0].Val == "elevation upper" {
						elevationUpperState = true
						elevationMiddleState = false
						elevationLowerState = false
					} else if t.Attr[0].Val == "elevation middle" {
						elevationMiddleState = true
						elevationLowerState = false
						elevationUpperState = false
					} else if t.Attr[0].Val == "elevation lower" {
						elevationLowerState = true
						elevationUpperState = false
						elevationMiddleState = false
					}
				}
			} else if t.Data == "p" {
				if len(t.Attr) > 0 {
					// fmt.Printf("=> %s %s\n", t, t.Attr)
					if strings.Contains(t.Attr[0].Val, "beginner") {
						slopesBeginning = true
						slopesIntermediate = false
						slopesAdvanced = false
						slopesExpert = false
					} else if strings.Contains(t.Attr[0].Val, "intermediate") {
						slopesIntermediate = true
						slopesBeginning = false
						slopesAdvanced = false
						slopesExpert = false
					} else if strings.Contains(t.Attr[0].Val, "advanced") {
						slopesAdvanced = true
						slopesBeginning = false
						slopesIntermediate = false
						slopesExpert = false
					} else if strings.Contains(t.Attr[0].Val, "expert") {
						slopesExpert = true
						slopesBeginning = false
						slopesAdvanced = false
						slopesIntermediate = false
					}
				}
			} else if t.Data == "div" {
				if len(t.Attr) > 0 {
					if t.Attr[0].Val == "time" {
						inner := z.Next()
						if inner == html.TextToken {
							resortDesc.SnowFallDays = append(resortDesc.SnowFallDays, &SnowFall{
								Day: extractTextTag(z),
							})
							// fmt.Printf("==> %d %s\n", snowDay, resortDesc.SnowFallDays)

						}
					} else if t.Attr[0].Val == "predicted_snowfall" {
						// snowFall += 1
					} else {
						inner := z.Next()
						if inner == html.TextToken {
							text := (string)(z.Text())
							value := strings.TrimSpace(text)
							if t.Attr[0].Val == "bluePill" {
								if elevationUpperState {
									if snowDepth == 1 {
										resortDesc.Piste.Upper = value
									} else if snowDepth == 2 {
										resortDesc.OffPiste.Upper = value
									}
									elevationUpperState = false
								} else if elevationMiddleState {
									if snowDepth == 1 {
										resortDesc.Piste.Middle = value
									} else if snowDepth == 2 {
										resortDesc.OffPiste.Middle = value
									}
									elevationMiddleState = false
								} else if elevationLowerState {
									if snowDepth == 1 {
										resortDesc.Piste.Lower = value
									} else if snowDepth == 2 {
										resortDesc.OffPiste.Lower = value
									}
									elevationLowerState = false
								} else {
									resortDesc.SnowFallDays[snowDay].Snow = value
									// fmt.Printf("Snow: %d %s\n", snowDay, resortDesc.SnowFallDays)
									snowDay += 1
								}
							}
						}
					}
				}
			}
		case html.TextToken: // text between start and end tag
			text := extractTextTag(z)
			if len(text) > 0 {
				if slopesBeginning {
					resortDesc.Slopes.Beginning.WriteString(text)
					// fmt.Printf("Begin: %s\n", text)
				} else if slopesIntermediate {
					resortDesc.Slopes.Intermediate.WriteString(text)
					// fmt.Printf("Intermediate: %s\n", text)
				} else if slopesAdvanced {
					resortDesc.Slopes.Advanced.WriteString(text)
					// fmt.Printf("Advanced: %s\n", text)
				} else if slopesExpert {
					resortDesc.Slopes.Expert.WriteString(text)
					// fmt.Printf("Expert: %s\n", text)
				}
			}
		case html.EndTagToken: // </tag>
			t := z.Token()
			if t.Data == "p" {
				slopesExpert = false
				slopesBeginning = false
				slopesAdvanced = false
				slopesIntermediate = false
			}
		case html.SelfClosingTagToken: // <tag/>
		}
	}
	return resortDesc, err
}
