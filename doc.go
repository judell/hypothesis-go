/*
# A Go client for the Hypothesis API

[Hypothesis] is software for annotating the web.

This package implements [hypothesis.Client] which searches the [Hypothesis API].

See [anchoring].

[Hypothesis]: https://web.hypothes.is
[Hypothesis API]: https://h.readthedocs.io/en/latest/api/
*/
package hypothesis

import (
	"github.com/judell/hypothesis-go/docs/anchoring"
)

var forceImportAnchoring anchoring.ForceImport

