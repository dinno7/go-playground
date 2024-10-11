package main

import "net/url"

type GoogleTranslator struct {
	from string
	to   string
	text string // The input text
}

func (gt *GoogleTranslator) GetAPIUrl() string {
	URL := url.URL{
		Scheme:   "https",
		Host:     "translate.googleapis.com",
		Path:     "/translate_a/single",
		RawQuery: "client=gtx&dt=t",
	}
	q := URL.Query()
	q.Add("sl", gt.from)
	q.Add("tl", gt.to)
	q.Add("q", gt.text)
	URL.RawQuery = q.Encode()

	return URL.String()
}
