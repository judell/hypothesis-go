package hypothesis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	)

type Client struct {
	token string
	params SearchParams
	maxSearchResults int
	httpClient http.Client
}

type Profile struct {
	Userid    string `json:"userid"`
	Authority string `json:"authority"`
	Groups    []struct {
		Name   string `json:"name"`
		ID     string `json:"id"`
		Public bool   `json:"public"`
		URL    string `json:"url,omitempty"`
	} `json:"groups"`
	Features struct {
		NotebookLaunch     bool `json:"notebook_launch"`
		EmbedCachebuster   bool `json:"embed_cachebuster"`
		ClientDisplayNames bool `json:"client_display_names"`
	} `json:"features"`
	Preferences struct {
	} `json:"preferences"`
	UserInfo struct {
		DisplayName string `json:"display_name"`
	} `json:"user_info"`
}

type SearchResult struct {
	Total int `json:"total"`
	Rows []Row `json:"rows"`
  }

type SearchParams struct {
	SearchAfter string
	Any string
	User string
	Group string
	Uri string
	WildcardUri string
	Tags[] string
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
	tagArray := apply(params.Tags, tagParamWrap)
	tags := strings.Join(tagArray, "")
	url := "https://hypothes.is/api/search?limit=200&search_after=" + url.QueryEscape(params.SearchAfter) + 
		"&user=" + params.User + 
		"&group=" + params.Group +
		tags
	if client.params.Any != "" {
		url += "&any=" + params.Any
	}
	if client.params.Uri != "" {
		url += "&uri=" + params.Uri
	} else {
		if client.params.WildcardUri != "" {
			url += "&wildcard_uri=" + params.WildcardUri
		}
	}
	req, _ := http.NewRequest("GET", url, nil)
	if (client.token != "") {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}

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

func (client *Client) Profile() (Profile, error) {
	url := "https://hypothes.is/api/profile"
	req, _ := http.NewRequest("GET", url, nil)
	if (client.token != "") {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}
	r, err := client.httpClient.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("error getting Hypothesis profile %d %v+", r.StatusCode, err.Error())
	}
	if r.StatusCode != 200 {
		return Profile{}, fmt.Errorf("error getting Hypothesis profile %d %v+",  r.StatusCode, r.Status)
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var profile Profile
	_ = decoder.Decode(&profile)
	return profile, err
}

func apply(strings []string, fn func(string) string) []string {
	var result []string
	for _, item := range strings {
		result = append(result, fn(item))
	}
	return result
}

func tagParamWrap(str string) string {
  return fmt.Sprintf(`&tag=%s`, strings.Replace(str, " ", "%20", -1))
}
