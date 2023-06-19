package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/communication/websocket"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/fetcher"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/map_manager"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/stopwatch"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/switchboard"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"github.com/google/uuid"
	netHttp "net/http"
	"regexp"
	"strconv"
)

type Instance struct {
	name string
	ws   *websocket.WebsocketHandler

	switchboard *switchboard.Switchboard

	mapManagers  map[string]*map_manager.MapManager
	stopwatches  map[string]*stopwatch.Stopwatch
	dataFetchers map[string]*fetcher.HttpFetcher

	mapName         map[string]string
	stopwatchName   map[string]string
	dataFetcherName map[string]string

	mapIndex         map[int]string
	stopwatchIndex   map[int]string
	dataFetcherIndex map[int]string
}

func NewInstance(name string) (*Instance, error) {
	if !valid(name) {
		return nil, errors.New("invalid name, must match regexp: [^\\w+$]")
	}

	return &Instance{
		name: name,
		ws:   websocket.NewWebsocketHandler(name, true),

		switchboard: switchboard.NewSwitchboard(),

		mapManagers:  make(map[string]*map_manager.MapManager),
		stopwatches:  make(map[string]*stopwatch.Stopwatch),
		dataFetchers: make(map[string]*fetcher.HttpFetcher),

		mapName:         make(map[string]string),
		stopwatchName:   make(map[string]string),
		dataFetcherName: make(map[string]string),

		mapIndex:         make(map[int]string),
		stopwatchIndex:   make(map[int]string),
		dataFetcherIndex: make(map[int]string),
	}, nil
}

func valid(name string) bool {
	re := regexp.MustCompile("^\\w+$")
	return re.MatchString(name)
}

func (i *Instance) AddResources(resources []map[string]interface{}) ([]string, error) {
	var uuids []string
	for _, resource := range resources {
		for name, r := range resource {
			switch r.(type) {
			case *map_manager.MapManager:
				if u, err := i.AddMapManager(name, r.(*map_manager.MapManager)); err != nil {
					return uuids, err
				} else {
					uuids = append(uuids, u)
				}
			case *stopwatch.Stopwatch:
				if u, err := i.AddStopwatch(name, r.(*stopwatch.Stopwatch)); err != nil {
					return uuids, err
				} else {
					uuids = append(uuids, u)
				}
			case *fetcher.HttpFetcher:
				if u, err := i.AddDataFetcher(name, r.(*fetcher.HttpFetcher)); err != nil {
					return uuids, err
				} else {
					uuids = append(uuids, u)
				}
			}
		}
	}
	return uuids, nil
}

func (i *Instance) AddMapManager(name string, mm *map_manager.MapManager) (string, error) {
	if _, notOk := i.mapName[name]; notOk {
		return "", errors.New("map name already exists")
	}
	uuid := uuid.NewString()
	i.mapManagers[uuid] = mm
	i.mapIndex[len(i.mapManagers)-1] = uuid
	i.mapName[name] = uuid
	return uuid, nil
}

func (i *Instance) RemoveMapManager(name string) error {
	uuid, ok := i.mapName[name]
	if !ok {
		return fmt.Errorf("cannot remove MapManager, none exists with name [%s]", name)
	}
	delete(i.mapManagers, uuid)

	var idx = -1
	for i, u := range i.mapIndex {
		if u == uuid {
			idx = i
		}
	}
	if idx < 0 {
		fmt.Errorf("cannot remove MapManager, no index found for name [%s]", name)
	}
	delete(i.mapIndex, idx)

	delete(i.mapName, name)
	return nil
}

func (i *Instance) AddStopwatch(name string, s *stopwatch.Stopwatch) (string, error) {
	if _, notOk := i.stopwatchName[name]; notOk {
		return "", errors.New("stopwatch name already exists")
	}
	uuid := uuid.NewString()
	i.stopwatches[uuid] = s
	i.stopwatchIndex[len(i.stopwatches)] = uuid
	i.stopwatchName[name] = uuid
	return uuid, nil
}

func (i *Instance) RemoveStopwatch(name string) error {
	uuid, ok := i.stopwatchName[name]
	if !ok {
		return fmt.Errorf("cannot remove Stopwatch, none exists with name [%s]", name)
	}
	delete(i.stopwatches, uuid)

	var idx = -1
	for i, u := range i.stopwatchIndex {
		if u == uuid {
			idx = i
		}
	}
	if idx < 0 {
		return fmt.Errorf("cannot remove Stopwatch, no index found for name [%s]", name)
	}
	delete(i.stopwatchIndex, idx)

	delete(i.stopwatchName, name)
	return nil
}

func (i *Instance) AddDataFetcher(name string, df *fetcher.HttpFetcher) (string, error) {
	if _, notOk := i.dataFetcherName[name]; notOk {
		return "", errors.New("data fetcher name already exists")
	}
	uuid := uuid.NewString()
	i.dataFetchers[uuid] = df
	i.dataFetcherIndex[len(i.dataFetchers)] = uuid
	i.dataFetcherName[name] = uuid
	return uuid, nil
}

func (i *Instance) RemoveDataFetcher(name string) error {
	uuid, ok := i.dataFetcherName[name]
	if !ok {
		return fmt.Errorf("cannot remove DataFetcher, none exists with name [%s]", name)
	}
	delete(i.dataFetchers, uuid)

	var idx = -1
	for i, u := range i.dataFetcherIndex {
		if u == uuid {
			idx = i
		}
	}
	if idx < 0 {
		fmt.Errorf("cannot remove DataFetcher, no index found for name [%s]", name)
	}
	delete(i.dataFetcherIndex, idx)

	delete(i.dataFetcherName, name)
	return nil
}

type routeStatus int

const (
	invalidPath routeStatus = iota
	fullMatch
	shortMatch
	resourceMatch
	resourceAllMatch
)

func (i *Instance) HandleStopwatches(w netHttp.ResponseWriter, r *netHttp.Request) {
	log.Devf("HandleStopwatches: %s\n", r.URL.Path)
	stopwatchName, status := i.validateStopwatchRoute(r.URL.Path)
	if status == invalidPath {
		netHttp.Error(w, "invalid path", netHttp.StatusNotFound)
		return
	}

	sw, err := i.findStopwatch(stopwatchName)
	if err != nil {
		netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
		return
	}

	jsonData, err := sw.Json(nil)

	if err != nil {
		netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (i *Instance) validateStopwatchRoute(path string) (string, routeStatus) {
	re := regexp.MustCompile(`/stopwatch/(\w+)$`)
	match := re.FindStringSubmatch(path)
	if len(match) != 2 {
		return "", invalidPath
	}
	return match[1], fullMatch
}

func (i *Instance) findStopwatch(name string) (*stopwatch.Stopwatch, error) {
	var sw *stopwatch.Stopwatch
	if idx, err := strconv.Atoi(name); err == nil {
		var ok bool
		if sw, ok = i.stopwatches[i.mapIndex[idx]]; !ok {
			return nil, fmt.Errorf("stopwatch not found for index [%s]", name)
		}
	} else if u, ok := i.stopwatchName[name]; ok {
		var ok bool
		if sw, ok = i.stopwatches[u]; !ok {
			return nil, fmt.Errorf("stopwatch not found for name [%s]", name)
		}
	} else {
		var ok bool
		if sw, ok = i.stopwatches[name]; !ok {
			return nil, fmt.Errorf("stopwatch not found for uuid [%s]", name)
		}
	}
	return sw, nil
}

func (i *Instance) HandleSwitchboard(w netHttp.ResponseWriter, r *netHttp.Request) {
	log.Devf("HandleSwitchboard: %s\n", r.URL.Path)
	key, status := i.validateSwitchboardRoute(r.URL.Path)

	switch status {
	case fullMatch:
		jsonData, err := i.switchboard.GetAsJsonPayload(key, nil)
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	case resourceMatch:
		jsonData, err := i.switchboard.Json()
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	default:
		netHttp.Error(w, "invalid path", netHttp.StatusNotFound)
		return
	}
}

func (i *Instance) validateSwitchboardRoute(path string) (string, routeStatus) {
	reFull := regexp.MustCompile(`switchboard/(\w+)$`)
	reSwb := regexp.MustCompile(`switchboard/?$`)
	if match := reFull.FindStringSubmatch(path); match != nil {
		return match[1], fullMatch
	} else if match := reSwb.FindStringSubmatch(path); match != nil {
		return "", resourceMatch
	}
	return "", invalidPath
}

func (i *Instance) HandleMaps(w netHttp.ResponseWriter, r *netHttp.Request) {
	log.Devf("HandleMaps: %s\n", r.URL.Path)
	mapName, dataType, key, status := i.validateMapRoute(r.URL.Path)

	switch status {
	case fullMatch:
		fallthrough
	case shortMatch:
		mapManager, err := i.findMap(mapName)
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusNotFound)
			return
		}

		jsonData, err := mapManager.GetAsJsonPayload(dataType, key, nil)
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	case resourceMatch:
		mapManager, err := i.findMap(mapName)
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusNotFound)
			return
		}

		jsonData, err := mapManager.Json(false)
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	case resourceAllMatch:
		if !(len(i.mapManagers) > 0) {
			netHttp.Error(w, fmt.Sprintf("no mapManagers declared on instance [%s]", i.name), netHttp.StatusNotFound)
			return
		}

		managerData := map[string]map[string]interface{}{}
		for uuid, manager := range i.mapManagers {
			jsonData, err := manager.Json(false)
			if err != nil {
				netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
				return
			}
			var data map[string]interface{}
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
				return
			}
			managerData[uuid] = data
		}

		jsonData, err := json.Marshal(managerData)
		if err != nil {
			netHttp.Error(w, err.Error(), netHttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	default:
		netHttp.Error(w, "invalid path", netHttp.StatusBadRequest)
	}
}

func (i *Instance) validateMapRoute(path string) (string, map_manager.ValueType, string, routeStatus) {
	reFull := regexp.MustCompile(`map\/(\w+)\/(\w+)\/(\w+)`)
	reShort := regexp.MustCompile(`map\/(\w+)\/(\w+)`)
	reMap := regexp.MustCompile(`map\/(\w+)`)
	reAll := regexp.MustCompile(`map\/?`)

	var (
		mapName  string
		dataType map_manager.ValueType
		key      string
		r        routeStatus
	)

	if matches := reFull.FindStringSubmatch(path); matches != nil {
		mapName = matches[1]
		dataType = map_manager.GetType(matches[2])
		key = matches[3]
		r = fullMatch
	} else if matches := reShort.FindStringSubmatch(path); matches != nil {
		mapName = matches[1]
		key = matches[2]
		r = shortMatch
	} else if matches := reMap.FindStringSubmatch(path); matches != nil {
		mapName = matches[1]
		r = resourceMatch
	} else if matches := reAll.FindStringSubmatch(path); matches != nil {
		r = resourceAllMatch
	} else {
		return "", map_manager.UnknownType, "", invalidPath
	}

	return mapName, dataType, key, r
}

func (i *Instance) findMap(mapName string) (*map_manager.MapManager, error) {
	var mapManager *map_manager.MapManager
	if idx, err := strconv.Atoi(mapName); err == nil {
		var ok bool
		if mapManager, ok = i.mapManagers[i.mapIndex[idx]]; !ok {
			return nil, fmt.Errorf("map not found for index [%s]", mapName)
		}
	} else if u, ok := i.mapName[mapName]; ok {
		var ok bool
		if mapManager, ok = i.mapManagers[u]; !ok {
			return nil, fmt.Errorf("map not found for name [%s]", mapName)
		}
	} else {
		var ok bool
		if mapManager, ok = i.mapManagers[mapName]; !ok {
			return nil, fmt.Errorf("map not found for uuid [%s]", mapName)
		}
	}
	return mapManager, nil
}
