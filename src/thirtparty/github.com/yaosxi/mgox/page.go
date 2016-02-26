package mgox

type Page struct {
	//	Request int `json:"request"`
	Cursor int `json:"cursor"`
	Count  int `json:"count"`
	Total  int `json:"total"`
	Next   int `json:"next"`
}

var PAGE_RECORD_COUNT = 20
