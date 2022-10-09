/*
# Create a client

Use [hypothesis.NewClient] to create a [hypothesis.Client] that searches for annotations. 

NewClient's token param is optional. If it's set to your [Hypothesis token], you'll search both public and private annotations. If it's the empty string, you'll only search public annotations. 

NewClient's [hypothesis.SearchParams] is likewise optional. If empty, your search will be unfiltered.

NewClient's maxSearchResults param determines how many annotations to fetch. If 0, the limit defaults to 400. 

To search for the most recent 10 public annotations:

	client := NewClient(
		"",
		SearchParams{},
		10,
	)

To search for the most recent 10 public or private annotations, if your token is in an env var called H_TOKEN:

	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{},
		10,
	)

To search for the most recent 10 annotations in a private group whose id in an env var called H_GROUP:

	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{
			Group: os.Getenv("H_GROUP"),
		},	
		10,
	)

To search for at most 10 public or private annotations from user 'judell', with the tag 'social media':

	client := NewClient(
		"",
		SearchParams{
			User: 'judell',
			Tags: []string{"social media"},
		},
		expect,
	)


# Search for annotations

The Hypothesis search API returns at most 200 annotations. [hypothesis.Search] encapsulates that API call, and returns an array of [hypothesis.Row]. Each Row represents one annotation. To fetch more than 200 annotations, use [hypothesis.SearchAll].

This test should find 2000 recent public or private annotations.

	func Test_Search_Finds_2000_Rows(t *testing.T) {
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

For more search examples, see the [Steampipe Hypothesis plugin].

# List private groups

If you authenticate with your token you can call the [Profile] function to list your private groups. 

	client := NewClient(
		os.Getenv("H_TOKEN"),
		SearchParams{},
		0,
	)
	profile, err := client.Profile()

Here, 'profile' is a [hypothesis.Profile] which includes [hypothesis.Profile.Groups], an array of structs that include the names and ids of your private groups.

# Anchoring

See [anchoring].

[Hypothesis token]: https://hypothes.is/account/developer
[Steampipe Hypothesis plugin]: https://hub.steampipe.io/plugins/turbot/hypothesis
[Profile]: https://
*/
package hypothesis

import (
	"github.com/judell/hypothesis-go/docs/anchoring"
)

var forceImportAnchoring anchoring.ForceImport

