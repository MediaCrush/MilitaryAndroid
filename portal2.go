package main

import (
	"encoding/json"
	"fmt"
	"github.com/jdiez17/go-irc"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

const (
	QUERY       string = "http://www.portal2sounds.com/list.php?%s"
	QUERY_P1           = "http://p1.portal2sounds.com/list.php?%s"
	PORTALSOUND        = "http://%s.portal2sounds.com/%d"
)

type portalSoundEntry struct {
	Id     string
	Text   string
	Who    string
	Domain string
}

func (p portalSoundEntry) Url() string {
	id, _ := strconv.Atoi(p.Id)
	return fmt.Sprintf(PORTALSOUND, p.Domain, id)
}

type portal2SoundsResponse struct {
	Numbers []int
	Content []portalSoundEntry
	Domain  string
}

// I hate you, portal2sounds.com, for your awful JSON.   
func (p *portal2SoundsResponse) UnmarshalJSON(data []byte) error {
	var content []interface{}
	err := json.Unmarshal(data, &content)
	if err != nil {
		return err
	}

	for _, v := range content {
		switch t := v.(type) {
		case float64:
			p.Numbers = append(p.Numbers, int(t))
		case map[string]interface{}:
			var entry portalSoundEntry

			entry.Domain = p.Domain
			for k, v := range t {
				value := v.(string)
				entry := reflect.ValueOf(&entry)
				k = strings.Title(k)
				f := entry.Elem().FieldByName(k)
				if f.IsValid() {
					f.SetString(value)
				}
			}
			p.Content = append(p.Content, entry)
		}
	}

	return nil
}

func getSoundEntries(data []byte, domain string) ([]portalSoundEntry, error) {
	var response portal2SoundsResponse
	response.Domain = domain
	err := json.Unmarshal([]byte(data), &response)

	if err != nil {
		return nil, err
	}

	return response.Content, nil
}

func fetchEntries(rqurl string) ([]portalSoundEntry, error) {
	response, err := http.Get(rqurl)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	purl, _ := url.Parse(rqurl)
	domain := strings.Split(purl.Host, ".")[0]
	content, err := getSoundEntries(bytes, domain)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func findQuote(text string, who string) ([]portalSoundEntry, error) {
	content := make([]portalSoundEntry, 0)

	v := url.Values{}
	v.Add("quote", text)
	v.Add("who", who)

	url := fmt.Sprintf(QUERY, v.Encode())
	fmt.Println(url)
	quotes, err := fetchEntries(url)
	if err != nil {
		return nil, err
	}

	url = fmt.Sprintf(QUERY_P1, v.Encode())
	quotes_p1, err := fetchEntries(url)
	if err != nil {
		return nil, err
	}

	content = append(content, quotes...)
	content = append(content, quotes_p1...)
	return content, nil
}

func portalCommandHandler(c *irc.Connection, e *irc.Event) {
	var quotes []portalSoundEntry
	text := ""
	who := ""

	if !strings.Contains(e.Payload["message"], "\"") {
		text = strings.Join(e.Params, " ")
	} else {
		if len(e.Params) == 1 {
			text = e.Params[0]
		} else {
			text = e.Params[0]
			who = e.Params[1]
		}
	}

	quotes, err := findQuote(text, who)
	if err != nil {
		panic(err)
	}

	if len(quotes) == 0 {
		e.React(c, "Nothing found.")
		return
	}

	e.React(c, quotes[0].Url())
}
