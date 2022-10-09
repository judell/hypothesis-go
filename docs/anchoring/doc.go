/*
A [hypothesis.Client] returns an array of [hypothesis.Row]. Each Row includes an array of [hypothesis.Target], and each Target includes an array of [hypothesis.Selector]. These define the location of a segment within a document in several different ways. If [hypothesis.Selector.Type] is TextPositionSelector, then [hypothesis.Selector.Start] and [hypothesis.Selector.End] are significant. If [hypothesis.Selector.Type] is TextQuoteSelector, then [hypothesis.Selector.Prefix], [hypothesis.Selector.Suffix], and [hypothesis.Selector.Exact] are significant. 

When software creates annotations, it can define where they anchor in text using two data structures: TextQuoteSelector and TextPositionSelector. The latter is optional in many cases because the TextQuoteSelector alone sufficiently defines the location. TextPositionSelector matters when the combination of Prefix, Exact, and Suffix does not uniquely locate a segment. 

When software anchors annotations in text, it reads the TextQuoteSelector and/or TextPositionSelector structures from the API, and uses that info to modify the region of text they identify -- typically, by highlighting it and linking an action to it.

When software analyzes an annotated corpus it can use the selector data to, e.g., visualize overlapping annotations. [hypothesis.SelectorsToExact] is a convenience function that returns the Exact quote from an array of selectors.

[hypothesis.Client]: https://pkg.go.dev/github.com/judell/hypothesis-go#Client
[hypothesis.Row]: https://pkg.go.dev/github.com/judell/hypothesis-go#Row
[hypothesis.Selector]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector
[hypothesis.Selector.Type]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector.Type
[hypothesis.Selector.Start]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector.Start
[hypothesis.Selector.End]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector.End
[hypothesis.Selector.Prefix]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector.Prefix
[hypothesis.Selector.Suffix]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector.Suffix
[hypothesis.Selector.Exact]: https://pkg.go.dev/github.com/judell/hypothesis-go#Selector.Exact
[hypothesis.SelectorsToExact]: https://pkg.go.dev/github.com/judell/hypothesis-go#SelectorsToExact
*/
package anchoring

type ForceImport string

