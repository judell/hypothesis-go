package hypothesis

import (
	"os"
	"regexp"
	"testing"
)
func Test_Search_Finds_Default_2000_Rows(t *testing.T) {
	client := NewClient(
		"", 
		SearchParams{},
		0,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != 2000 {
        t.Fatalf(`expected 2000 rows, got %d, `, len(rows))
    }
}
func Test_Search_For_215_Rows_Finds_215_Rows(t *testing.T) {
	count := 215
	client := NewClient(
		"", 
		SearchParams{},
		count,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != count {
        t.Fatalf(`expected 215 rows, got %d, `, len(rows))
    }
}
func Test_Search_With_Token_Finds_Rows_In_Private_Group(t *testing.T) {
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

	if len(rows) != 1 {
        t.Fatalf(`expected 1 row, got %d, `, len(rows))
	}

	if rows[0].Group != client.params.Group {
        t.Fatalf(`expected group %s, got %s, `, client.params.Group, rows[0].Group)
	}
}
func Test_Search_For_Two_Tags_Finds_10_Tags(t *testing.T) {
	count := 10
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Tags: []string{"media","review"},
		},
		count,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != count  {
        t.Fatalf(`expected %d rows, got %d, `, count, len(rows))
    }
}

func Test_Search_For_User_Finds_3_Annos_For_User(t *testing.T) {
	user := "judell"
	count := 3
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			User: user,
		},
		count,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != count  {
        t.Fatalf(`expected %d rows, got %d, `, count, len(rows))
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
	count := 3
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Uri: uri,
		},
		count,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != count  {
        t.Fatalf(`expected %d rows, got %d, `, count, len(rows))
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
	count := 3
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			WildcardUri: wildcardUri,
		},
		count,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != count  {
        t.Fatalf(`expected %d rows, got %d, `, count, len(rows))
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
	count := 1
	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Any: "jon",
		},
		count,
	)

	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}

	if len(rows) != count {
        t.Fatalf(`expected %d rows, got %d `, count, len(rows))
	}
}
