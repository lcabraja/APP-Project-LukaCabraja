package fetcher

import (
	"bytes"
	"errors"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/fetcher/parser"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type HttpFetcher struct {
	url func() (string, error)
	key string

	subscriptions map[string]*FetchEventSubscription

	data  []byte
	mime  string
	store bool
}

func NewHttpFetcher(url func() (string, error), store bool) *HttpFetcher {
	return &HttpFetcher{
		url:           url,
		subscriptions: make(map[string]*FetchEventSubscription),
		store:         store,
	}
}

func (hf *HttpFetcher) Subscribe(sub *FetchEventSubscription) string {
	uuid := uuid.NewString()
	hf.subscriptions[uuid] = sub
	return uuid
}

func (hf *HttpFetcher) Unsubscribe(uuid string) {
	delete(hf.subscriptions, uuid)
}

func (hf *HttpFetcher) fireEvent(e FetchEvent, data []byte, mime string) {
	for _, sub := range hf.subscriptions {
		if sub.Event == e {
			go hf.launchParser(sub.Name, sub.Handler, data, mime)
		}
	}
}

func (hf *HttpFetcher) launchParser(name string, parser parser.Parser, data []byte, mime string) {
	defer func() {
		if r := recover(); r != nil {
			log.Ef("Parser for [%s] panicked with error: %s\n", name, r)
		}
	}()
	if err := parser.Parse(data); err != nil {
		log.E(err.Error())
	}
}

func (hf *HttpFetcher) Fetch() error {
	url, err := hf.url()
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var mime string
	mimes := resp.Header["Content-Type"]
	if len(mimes) > 0 {
		mime = mimes[0]
	}

	var different int
	if hf.store {
		different = bytes.Compare(hf.data, data)
		hf.data = data
		hf.mime = mime
	}

	if hf.store && different == 0 {
		hf.fireEvent(DataChanged, data, mime)
	}
	hf.fireEvent(DataReceived, data, mime)
	return nil
}

func (hf *HttpFetcher) Parse() error {
	if hf.store {
		hf.fireEvent(ManualParse, hf.data, hf.mime)
		return nil
	}
	return errors.New("store is set to false, nothing to parse")
}
