package userShoppingListEntry

import "net/http"

type JsonFormatCreateEntryRequest struct {
	Count   uint16 `json:"count,omitempty"`
	Note    string `json:"note,omitempty"`
	Checked bool   `json:"checked,omitempty"`
}

type JsonFormatUpdateEntryRequest struct {
	Count   uint16 `json:"count,omitempty"`
	Note    string `json:"note,omitempty"`
	Checked bool   `json:"checked,omitempty"`
}

type Controller interface {
	GetEntry(http.ResponseWriter, *http.Request)
	GetEntries(http.ResponseWriter, *http.Request)
	PostEntry(http.ResponseWriter, *http.Request)
	PutEntry(http.ResponseWriter, *http.Request)
	DeleteEntry(http.ResponseWriter, *http.Request)
}
