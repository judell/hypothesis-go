package hypothesis

import (
	"testing"
	"os"
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
	client := NewClient(
		"", 
		SearchParams{},
		215,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != 215 {
        t.Fatalf(`expected 215 rows, got %d, `, len(rows))
    }
}
func Test_Search_With_Token_Finds_Rows_In_Private_Group(t *testing.T) {
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Group: os.Getenv("H_GROUP"),
		},
		20,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != 20 {
        t.Fatalf(`expected 20 rows, got %d, `, len(rows))
    }
	if rows[0].Group != client.params.Group {
        t.Fatalf(`expected group %s, got %s, `, client.params.Group, rows[0].Group)
	}
}

func Test_Search_For_Two_Tags_Finds_10_Tags(t *testing.T) {
	client := NewClient(
		os.Getenv("H_TOKEN"), 
		SearchParams{
			Tags: []string{"media","review"},
		},
		10,
	)
	rows, err := client.SearchAll()

	if err != nil {
        t.Fatalf(`%v`, err)
	}
	if len(rows) != 10  {
        t.Fatalf(`expected 10 rows, got %d, `, len(rows))
    }
}


