package hypothesis

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func Test_Uri_Has_Annos_With_References(t *testing.T) {
	client := NewClient(
		"",
		SearchParams{
			Uri: "https://web.hypothes.is/blog/introducing-search-and-profiles",
		},
		100,
	)
	maxRefs := 0
	for row := range client.SearchAll() {
		if len(row.References) > 0 {
			fmt.Printf("id: %v refs: %v\n", row.ID, row.References)
		}
		maxRefs += len(row.References)
	}
	if maxRefs < 10  {
		t.Fatalf(`expected at least %d references`, maxRefs)
	}
}
func Test_Search_Finds_Default_2000_Rows(t *testing.T) {
	expect := 2000
	client := NewClient(
		"",
		SearchParams{},
		expect,
	)
	rowCount := 0
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rowCount += 1
	}

	if rowCount != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, rowCount)
	}
}
func Test_Search_For_215_Rows_Finds_215_Rows(t *testing.T) {
	expect := 215
	client := NewClient(
		"",
		SearchParams{},
		expect,
	)
	rowCount := 0
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rowCount += 1
	}

	if rowCount != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, rowCount)
	}
}

func Test_Search_With_Token_Finds_Rows_In_Private_Group(t *testing.T) {
	expect := 1
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Group: os.Getenv("H_GROUP"),
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != 1 {
		t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}

	if rows[0].Group != client.params.Group {
		t.Fatalf(`expected group %s, got %s, `, client.params.Group, rows[0].Group)
	}
}

func Test_Search_For_Two_Tags_Finds_2_Rows(t *testing.T) {
	expect := 2
	client := NewClient(
		"",
		SearchParams{
			Tags: []string{"media", "review"},
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}

	for _, row := range rows {

		matchesMedia := stringMatchesAnyStringInSlice("media", row.Tags)
		matchesReview := stringMatchesAnyStringInSlice("review", row.Tags)

		if !(matchesMedia && matchesReview) {
			t.Fatalf(`expected "media" and "review" among tags, got %v`, row.Tags)
		}
	}

}

func Test_Search_For_Compound_Tag_Finds_3_Rows(t *testing.T) {
	expect := 3
	client := NewClient(
		"",
		SearchParams{
			Tags: []string{"social media"},
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}

	for _, row := range rows {
		if !(stringMatchesAnyStringInSlice("social media", row.Tags)) {
			t.Fatalf(`expected "social media" among tags, got %v`, row.Tags)
		}
	}

}

func Test_Search_For_User_Finds_3_Annos_For_User(t *testing.T) {
	user := "judell"
	expect := 3
	client := NewClient(
		"",
		SearchParams{
			User: user,
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}
	for _, row := range rows {
		m, _ := regexp.MatchString(user, row.User)
		if !m {
			t.Fatalf(`expected match for %s, got %s, `, user, row.User)
		}
	}
}

func Test_Search_For_Uri_Finds_3_Annotations(t *testing.T) {
	uri := "http://example.com/"
	expect := 3
	client := NewClient(
		"",
		SearchParams{
			Uri: uri,
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}
	for _, row := range rows {
		m, _ := regexp.MatchString(`example.com`, row.URI)
		if !m {
			t.Fatalf(`expected a match for %s, got %s`, uri, row.URI)
		}
	}
}

func Test_Search_For_Wildcard_Uri_Finds_3_Matching_Uris(t *testing.T) {
	wildcardUri := "https://www.nytimes/*"
	expect := 3
	client := NewClient(
		"",
		SearchParams{
			WildcardUri: wildcardUri,
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}
	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}

	for _, row := range rows {
		m, _ := regexp.MatchString(`www.nytimes.com`, row.URI)
		if !m {
			t.Fatalf(`expected a match for %s, got %s`, wildcardUri, row.URI)
		}
	}
}

func Test_GetProfile(t *testing.T) {
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{},
		0,
	)
	profile, err := client.GetProfile()

	if len(profile.Groups) == 1 {
		t.Fatalf(`%v`, err)
	}

	if err != nil {
		t.Fatalf(`%v`, err)
	}
}
func Test_Finds_A_Private_Annotation(t *testing.T) {
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Group: os.Getenv("H_GROUP"),
		},
		1,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		t.Fatalf("expected rows, got none")
	}

	if rows[0].Group == "__world__" {
		t.Fatalf(`expected %s, got __world__`, rows[0].Group)
	}
}

func Test_Search_Param_Any_Finds_One_Annotation(t *testing.T) {
	expect := 1
	client := NewClient(
		"",
		SearchParams{
			Any: "jon",
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d `, expect, len(rows))
	}
}
func Test_Limit_500_In_Search_Params_Yields_500_Rows(t *testing.T) {
	expect := 500
	client := NewClient(
		"",
		SearchParams{
			Limit: fmt.Sprintf("%d", expect),
		},
		expect,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != expect {
		t.Fatalf(`expected %d rows, got %d `, expect, len(rows))
	}

}
func Test_SearchParams_Overrides_MaxSearchResults_500_To_Yield_501_Rows(t *testing.T) {
	expect := 501
	client := NewClient(
		"",
		SearchParams{
			Limit: "501",
		},
		500,
	)

	rows := []Row{}
	for row := range client.SearchAll() {
		if row.ID == "" {
			t.Fatalf(`no ID for row %+v`, row)
		}
		rows = append(rows, row)
	}

	if len(rows) != 501 {
		t.Fatalf(`expected %d rows, got %d `, expect, len(rows))
	}

}

func Test_Selectors_To_Exact_Returns_Exact(t *testing.T) {
	expect := "the exact quote"
	text_quote_selector := Selector{"TextQuoteSelector", 0, 0, expect, "", ""}
	text_position_selector := Selector{"TextPositionSelector", 1, 2, "", "", ""}
	selectors := append([]Selector{}, text_quote_selector, text_position_selector)
	exact, err := SelectorsToExact(selectors)

	if err != nil {
		t.Fatalf(`%v`, err)
	}

	if expect != exact {
		t.Fatalf(`expected %s, got %s `, expect, exact)
	}
}

func Test_Selectors_To_Exact_Without_TextQuoteSelector_Returns_Empty(t *testing.T) {
	expect := ""
	text_position_selector := Selector{"TextPositionSelector", 1, 2, "", "", ""}
	selectors := append([]Selector{}, text_position_selector)
	exact, err := SelectorsToExact(selectors)

	if err != nil {
		t.Fatalf(`%v`, err)
	}

	if expect != exact {
		t.Fatalf(`expected %s, got %s `, expect, exact)
	}
}

func stringMatchesAnyStringInSlice(str string, strs []string) bool {
	var ret = false
	for _, _str := range strs {
		m, _ := regexp.Match(str, []byte(strings.ToLower(_str)))
		if m {
			ret = true
		}
	}
	return ret
}
