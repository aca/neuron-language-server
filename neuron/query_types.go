package neuron

type QueryResult struct {
	Skipped Skipped  `json:"skipped"`
	Result  []Result `json:"result"`
	// Query   []NeuronQuery `json:"query"`
}

type QueryClass struct {
	ZettelsViewLinkView   ZettelsViewLinkView `json:"zettelsViewLinkView"`
	ZettelsViewGroupByTag bool                `json:"zettelsViewGroupByTag"`
}

type Result struct {
	ZettelTags    []string      `json:"zettelTags"`
	ZettelDay     string        `json:"zettelDay"`
	ZettelID      string        `json:"zettelID"`
	ZettelError   ZettelError   `json:"zettelError"`
	ZettelContent []interface{} `json:"zettelContent"`
	// ZettelQueries     [][]ResultZettelQuery `json:"zettelQueries"`
	ZettelFormat      ZettelFormat `json:"zettelFormat"`
	ZettelPath        string       `json:"zettelPath"`
	ZettelTitle       string       `json:"zettelTitle"`
	ZettelTitleInBody bool         `json:"zettelTitleInBody"`
}

type ZettelError struct {
	Right []interface{} `json:"Right"`
}

type Skipped struct {
}

type ZettelsViewLinkView string

const (
	LinkViewDefault  ZettelsViewLinkView = "LinkView_Default"
	LinkViewShowDate ZettelsViewLinkView = "LinkView_ShowDate"
)

type QueryEnum string

const (
	ZettelQueryZettelByID   QueryEnum = "ZettelQuery_ZettelByID"
	ZettelQueryZettelsByTag QueryEnum = "ZettelQuery_ZettelsByTag"
)

type ZettelFormat string

const (
	Markdown ZettelFormat = "markdown"
)

type NeuronQuery struct {
	Enum       *QueryEnum
	UnionArray []QueryQueryUnion
}

type QueryQueryUnion struct {
	AnythingArray []interface{}
	QueryClass    *QueryClass
}

type ResultZettelQuery struct {
	Enum       *QueryEnum
	UnionArray []ZettelQueryZettelQuery
}

type ZettelQueryZettelQuery struct {
	QueryClass  *QueryClass
	String      *string
	StringArray []string
}
