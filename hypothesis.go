package hypothesis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	)

type Client struct {
	token string
	params SearchParams
	maxSearchResults int
	httpClient http.Client
}

type SearchResult struct {
	Total int `json:"total"`
	Rows []Row `json:"rows"`
  }

type SearchParams struct {
	SearchAfter string
	User string
	Group string
}
  
type Row struct {
	ID          string        `json:"id"`
	Created     string        `json:"created"`
	Updated     string        `json:"updated"`
	User        string        `json:"user"`
	URI         string        `json:"uri"`
	Text        string        `json:"text"`
	Tags        []string      `json:"tags"`
	Group       string        `json:"group"`
	Target []struct {
		Source   string `json:"source"`
		Selector []struct {
			End    int    `json:"end,omitempty"`
			Type   string `json:"type"`
			Start  int    `json:"start,omitempty"`
			Exact  string `json:"exact,omitempty"`
			Prefix string `json:"prefix,omitempty"`
			Suffix string `json:"suffix,omitempty"`
		} `json:"selector"`
	} `json:"target"`
	Document struct {
		Title []string `json:"title"`
	} `json:"document"`
	UserInfo struct {
		DisplayName string `json:"display_name"`
	} `json:"user_info"`
}

func NewClient(token string, params SearchParams, maxSearchResults int) *Client {

	var _maxSearchResults int
	if maxSearchResults == 0 {
		_maxSearchResults = 2000
	} else {
		_maxSearchResults = maxSearchResults
	}
	client := &Client{
		token:  token,
		params: params,
		maxSearchResults: _maxSearchResults,
	}

	return client
}

func (client *Client) Search() ([]Row, error) {
	params := client.params
	url := "https://hypothes.is/api/search?limit=200&search_after=" + url.QueryEscape(params.SearchAfter) + 
		"&user=" + params.User + 
		"&group=" + params.Group
	req, err := http.NewRequest("GET", url, nil)
	if (client.token != "") {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}
	//fmt.Printf("URL: %v\n", url)
	r, err := client.httpClient.Do(req)
	if err != nil {
		return []Row{}, fmt.Errorf("error getting Hypothesis search results for %s: %v+", url, err.Error())
	}
	if r.StatusCode != 200 {
		return []Row{}, fmt.Errorf("error getting Hypothesis search results for %s: %d %s", url, r.StatusCode, r.Status)
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var searchResult SearchResult
	_ = decoder.Decode(&searchResult)
	if searchResult.Total <= client.maxSearchResults {
		client.maxSearchResults = searchResult.Total		
	}
	return searchResult.Rows, nil
}

func (client *Client) SearchAll() ([]Row, error) {
	allRows, err := client.Search()
	if len(allRows) == 0 {
		return allRows, nil
	}	
	lastRow := allRows[len(allRows)-1]
	for {
		client.params.SearchAfter = lastRow.Updated
		moreRows, err := client.Search()
		if (err != nil) {
			return allRows, err
		}
		allRows = append(allRows, moreRows...)
		if (len(allRows) >= client.maxSearchResults) {
			break
		}
		lastRow = moreRows[len(moreRows)-1]
	}  
	if (len(allRows) >= client.maxSearchResults) {
		allRows = allRows[0:client.maxSearchResults]
	}
	return allRows, err
}

