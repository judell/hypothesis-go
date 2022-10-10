package hypothesis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-Create-a-client
type Client struct {
	token            string
	params           SearchParams
	maxSearchResults int
	httpClient       http.Client
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-List_private_groups
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

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-Search_for_annotations
type SearchResult struct {
	Total int   `json:"total"`
	Rows  []Row `json:"rows"`
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-Search_for_annotations
type SearchParams struct {
	SearchAfter string
	Limit       string
	Any         string
	User        string
	Group       string
	Uri         string
	WildcardUri string
	Tags        []string
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go/docs/anchoring
type Selector = struct {
	Type   string `json:"type"`
	Start  int    `json:"start,omitempty"`
	End    int    `json:"end,omitempty"`
	Exact  string `json:"exact,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go/docs/anchoring
type Target = struct {
	Source   string `json:"source"`
	Selector []Selector
}
type Row struct {
	ID       string   `json:"id"`
	Created  string   `json:"created"`
	Updated  string   `json:"updated"`
	User     string   `json:"user"`
	URI      string   `json:"uri"`
	Text     string   `json:"text"`
	Tags     []string `json:"tags"`
	References  []string `json:"references"`
	Group    string   `json:"group"`
	Target   []Target `json:"target"`
	Document struct {
		Title []string `json:"title"`
	} `json:"document"`
	UserInfo struct {
		DisplayName string `json:"display_name"`
	} `json:"user_info"`
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-Create_a_client
func NewClient(token string, params SearchParams, maxSearchResults int) *Client {

	var _maxSearchResults int

	if maxSearchResults == 0 {
		_maxSearchResults = 400
	}
	if maxSearchResults > 0 {
		_maxSearchResults = maxSearchResults
	}
	if params.Limit != "" {
		_maxSearchResults, _ = strconv.Atoi(params.Limit)
	}

	client := &Client{
		token:            token,
		params:           params,
		maxSearchResults: _maxSearchResults,
	}

	return client
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-Search_for_annotations
func (client *Client) Search() ([]Row, error) {
	params := client.params
	if params.Group == "" {
		params.Group = "__world__" // fix breaking change in api, group now required
	}
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
	} else if client.params.WildcardUri != "" {
		url += "&wildcard_uri=" + params.WildcardUri
	}
	req, _ := http.NewRequest("GET", url, nil)
	if client.token != "" {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}

	r, err := client.httpClient.Do(req)
	if err != nil {
		return []Row{}, fmt.Errorf("error getting Hypothesis search results for %s: %s", url, err.Error())
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

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-Search_for_annotations
func (client *Client) SearchAll() <-chan Row {
	channel := make(chan Row)
	initialRows, err := client.Search()
	if err != nil {
		fmt.Printf("%+v\n", err)
		close(channel)
		return channel
	}
	if len(initialRows) == 0 {
		close(channel)
		return channel
	}
	lastRow := initialRows[len(initialRows)-1]
	allRows := 0
	go func() {
		for _, row := range initialRows {
			if allRows >= client.maxSearchResults {
				break
			}
			channel <- row
			allRows += 1
		}
		for allRows < client.maxSearchResults {
			client.params.SearchAfter = lastRow.Updated
			moreRows, err := client.Search()
			if err != nil {
				fmt.Printf("%+v\n", err)
			}
			lastRow = moreRows[len(moreRows)-1]
			for _, row := range moreRows {
				channel <- row
				allRows += 1
				if allRows >= client.maxSearchResults {
					break
				}
			}
		}
		close(channel)
	}()
	return channel
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go#hdr-List_private_groups
func (client *Client) GetProfile() (Profile, error) {
	url := "https://hypothes.is/api/profile"
	req, _ := http.NewRequest("GET", url, nil)
	if client.token != "" {
		req.Header.Add("Authorization", "Bearer "+client.token)
	}
	r, err := client.httpClient.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("error getting Hypothesis profile %d %v+", r.StatusCode, err.Error())
	}
	if r.StatusCode != 200 {
		return Profile{}, fmt.Errorf("error getting Hypothesis profile %d %v+", r.StatusCode, r.Status)
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var profile Profile
	_ = decoder.Decode(&profile)
	return profile, err
}

// See https://pkg.go.dev/github.com/judell/hypothesis-go/docs/anchoring
func SelectorsToExact(selectors []Selector) (string, error) {
	empty := ""
	if len(selectors) == 0 {
		return empty, nil
	}
	for _, sel := range selectors {
		if sel.Type == "TextQuoteSelector" {
			return sel.Exact, nil
		}
	}
	return empty, nil
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
