package hypothesis

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)
func Test_Search_Finds_Default_2000_Rows(t *testing.T) {
	expect := 2000
	client := NewClient(
		"", 
		SearchParams{},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != expect {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
    }
}
func Test_Search_For_215_Rows_Finds_215_Rows(t *testing.T) {
	expect := 215
	client := NewClient(
		"", 
		SearchParams{},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != expect {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
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
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != 1 {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}

	if rows[0].Group != client.params.Group {
        t.Fatalf(`expected group %s, got %s, `, client.params.Group, rows[0].Group)
	}
}
func Test_Search_For_Two_Tags_Finds_10_Tags(t *testing.T) {
	expect := 10
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Tags: []string{"media","review"},
		},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != expect  {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
    }

	for _, row := range rows {

		matchesMedia := stringMatchesAnyStringInSlice("media", row.Tags)
		matchesReview := stringMatchesAnyStringInSlice("review", row.Tags)

		if ! ( matchesMedia && matchesReview ) {
			t.Fatalf(`expected "media" and "review" among tags, got %v`, row.Tags)
		}
	}

}

func Test_Search_For_Compound_Tag_Finds_3_Tags(t *testing.T) {
	expect := 3 
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Tags: []string{"social media"},
		},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != expect  {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
    }

	for _, row := range rows {
		if ! ( stringMatchesAnyStringInSlice("social media", row.Tags) ) {
			t.Fatalf(`expected "social media" among tags, got %v`, row.Tags)
		}
	}

}

func Test_Search_For_User_Finds_3_Annos_For_User(t *testing.T) {
	user := "judell"
	expect := 3
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			User: user,
		},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != expect  {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
    }
	for _, row := range rows {
		m, _ := regexp.MatchString(user, row.User)
		if ! m {
			t.Fatalf(`expected match for %s, got %s, `, user, row.User)
		}
	}
}

func Test_Search_For_Uri_Finds_3_Annotations(t *testing.T) {
	uri := "http://example.com/"
	expect := 3
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Uri: uri,
		},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != expect  {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
    }
	for _, row := range rows {
		m, _ := regexp.MatchString(`example.com`, row.URI)
		if ! m {
			t.Fatalf(`expected a match for %s, got %s`, uri, row.URI)
		}
    }
}

func Test_Search_For_Wildcard_Uri_Finds_3_Matching_Uris(t *testing.T) {
	wildcardUri := "https://www.nytimes/*"
	expect := 3
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			WildcardUri: wildcardUri,
		},
		expect,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != expect  {
        t.Fatalf(`expected %d rows, got %d, `, expect, len(rows))
	}
   
	for _, row := range rows {
		m, _ := regexp.MatchString(`www.nytimes.com`, row.URI)
		if ! m {
			t.Fatalf(`expected a match for %s, got %s`, wildcardUri, row.URI)
		}
    }
}

func Test_Profile(t *testing.T) {
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{},
		0,
	)
	profile, err := client.Profile()

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
	
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) == 0 {
        t.Fatalf(`%v`, err)
	}

	if rows[0].Group == "__world__" {
        t.Fatalf(`%v`, err)
	} 
}

func Test_Search_Param_Any_Finds_One_Annotation(t *testing.T) {
	expect := 1
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Any: "jon",
		},
		expect,
	)

	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != expect {
        t.Fatalf(`expected %d rows, got %d `, expect, len(rows))
	}
}
func Test_Limit_500_In_Search_Params_Yields_500_Rows(t *testing.T) {
	expect := 500
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Limit: fmt.Sprintf("%d", expect),
		},
		expect,
	)

	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != expect {
        t.Fatalf(`expected %d rows, got %d `, expect, len(rows))
	}
	
}
func Test_SearchParams_Overrides_MaxSearchResults_500_To_Yield_501_Rows(t *testing.T) {
	expect := 501
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Limit: "501",
		},
		500,
	)

	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != 501 {
        t.Fatalf(`expected %d rows, got %d `, expect, len(rows))
	}
	
}



func stringMatchesAnyStringInSlice(str string, strings []string) bool {
	var ret = false
    for _, _str := range strings {
		m, _ := regexp.Match(str, []byte(_str))
		if m {
			ret = true
		}
     }
	return ret
}