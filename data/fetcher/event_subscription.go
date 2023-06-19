package fetcher

import "github.com/lcabraja/APP-Project-LukaCabraja/data/fetcher/parser"

type FetchEventSubscription struct {
	Name    string
	Event   FetchEvent
	Handler parser.Parser
}
