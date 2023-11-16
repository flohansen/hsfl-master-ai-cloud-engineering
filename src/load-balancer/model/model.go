package model

import "net/url"

type Target struct {
	ContainerId string
	Url         *url.URL
}
