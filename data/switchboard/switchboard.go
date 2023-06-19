package switchboard

import (
	"encoding/json"
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/formatter"
	"github.com/lcabraja/APP-Project-LukaCabraja/data/map_manager"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"github.com/google/uuid"
	"sync"
)

type resolverType int

const (
	customResolver resolverType = iota
	formatResolver
	proxyResolver
)

type Switchboard struct {
	resolvers map[string]resolverType
	custom    map[string]ValueCalculator
	format    map[string]formatter.DataFormatter
	proxies   map[string]*Dependency
	deps      map[string][]*Dependency

	muKeys sync.RWMutex

	subscriptions map[string]*SwitchboardEventSubscription
	muSub         sync.RWMutex
}

func NewSwitchboard() *Switchboard {
	return &Switchboard{
		resolvers: make(map[string]resolverType),
		custom:    make(map[string]ValueCalculator),
		format:    make(map[string]formatter.DataFormatter),
		proxies:   make(map[string]*Dependency),
		deps:      make(map[string][]*Dependency),

		subscriptions: make(map[string]*SwitchboardEventSubscription),
	}
}

func (sb *Switchboard) Subscribe(sev *SwitchboardEventSubscription) string {
	sb.muSub.Lock()
	defer sb.muSub.Unlock()

	uuid := uuid.NewString()
	sb.subscriptions[uuid] = sev
	return uuid
}

func (sb *Switchboard) Unsubscribe(uuid string) {
	sb.muSub.Lock()
	defer sb.muSub.Unlock()

	delete(sb.subscriptions, uuid)
}

func (sb *Switchboard) fireEvent(e SwitchboardEvent, key string) {
	sb.muSub.RLock()
	defer sb.muSub.RUnlock()

	for _, sev := range sb.subscriptions {
		if sev.Keys == nil {
			sev.Handler(e, key, sb)
			continue
		}
		for _, k := range sev.Keys {
			if k == key {
				sev.Handler(e, key, sb)
			}
		}
	}
}

func (sb *Switchboard) HandleMapManagerEvent(e map_manager.MapManagerEvent, _ map_manager.ValueType, key string, mm *map_manager.MapManager) {
	log.Df("received event [%s] for key %s\n", e, key)
	sb.muSub.RLock()
	defer sb.muSub.RUnlock()

	for k, dep := range sb.deps {
		for _, d := range dep {
			if d.key == key && d.mm == mm {
				sb.fireEvent(KeyChanged, k)
			}
		}
	}
}

func (sb *Switchboard) canRegisterFor(key string) bool {
	sb.muKeys.RLock()
	defer sb.muKeys.RUnlock()

	var contains bool
	_, contains = sb.resolvers[key]
	_, contains = sb.custom[key]
	_, contains = sb.format[key]
	_, contains = sb.proxies[key]
	return !contains
}

func (sb *Switchboard) RegisterCustom(key string, c ValueCalculator, dep ...*Dependency) bool {
	if !sb.canRegisterFor(key) {
		return false
	}

	for _, d := range dep {
		if !d.isValid() {
			return false
		}
	}

	sb.muKeys.Lock()
	defer sb.muKeys.Unlock()

	sb.resolvers[key] = customResolver
	sb.custom[key] = c
	sb.deps[key] = append(sb.deps[key], dep...)

	return true
}

func (sb *Switchboard) RegisterFormatter(key string, f formatter.DataFormatter, dep *Dependency) bool {
	if !sb.canRegisterFor(key) {
		return false
	}

	if !dep.isValid() {
		return false
	}

	sb.muKeys.Lock()
	defer sb.muKeys.Unlock()

	sb.resolvers[key] = formatResolver
	sb.format[key] = f
	sb.deps[key] = append(sb.deps[key], dep)

	return true
}

func (sb *Switchboard) RegisterProxy(key string, d *Dependency) bool {
	if !sb.canRegisterFor(key) {
		return false
	}

	sb.muKeys.Lock()
	defer sb.muKeys.Unlock()

	sb.resolvers[key] = proxyResolver
	sb.proxies[key] = d
	sb.deps[key] = append(sb.deps[key], d)

	return true
}

func (sb *Switchboard) Get(key string) (string, map_manager.ValueType, error) {
	sb.muKeys.RLock()
	defer sb.muKeys.RUnlock()

	r, ok := sb.resolvers[key]
	if !ok {
		return "", map_manager.UnknownType, fmt.Errorf("key %s not found", key)
	}

	switch r {
	case customResolver:
		return sb.getCustom(key)
	case formatResolver:
		return sb.getFormat(key)
	case proxyResolver:
		return sb.getProxy(key)
	}

	return "", map_manager.UnknownType, fmt.Errorf("key %s not found", key)
}

func (sb *Switchboard) getProxy(key string) (string, map_manager.ValueType, error) {
	proxy, ok := sb.proxies[key]
	if !ok {
		return "", map_manager.UnknownType, fmt.Errorf("missing resolver for key: [%s]", key)
	}

	value, vt, ok := proxy.mm.Get(proxy.vt, proxy.key)
	if !ok {
		return "", map_manager.UnknownType, fmt.Errorf("key %s not found", key)
	}

	return fmt.Sprintf("%v", value), vt, nil
}

func (sb *Switchboard) getFormat(key string) (string, map_manager.ValueType, error) {
	format, ok := sb.format[key]
	if !ok {
		return "", map_manager.UnknownType, fmt.Errorf("missing resolver for key: [%s]", key)
	}

	dep, ok := sb.deps[key]
	if !ok || len(dep) == 0 {
		return "", map_manager.UnknownType, fmt.Errorf("missing dependencies for key: [%s]", key)
	}

	data, vt, ok := dep[0].mm.Get(dep[0].vt, dep[0].key)
	if !ok {
		return "", map_manager.UnknownType, fmt.Errorf("key %s not found in source mapmanager", key)
	}

	res, err := format.Format(data)

	return fmt.Sprintf("%s", res), vt, err
}

func (sb *Switchboard) getCustom(key string) (string, map_manager.ValueType, error) {
	custom, ok := sb.custom[key]
	if !ok {
		return "", map_manager.UnknownType, fmt.Errorf("missing resolver for key: [%s]", key)
	}

	deps, err := sb.getCustomDependencies(key)
	if err != nil {
		return "", map_manager.UnknownType, err
	}

	val, err := sb.getCustomResult(key, custom, deps)
	if err != nil {
		return "", map_manager.UnknownType, err
	}
	return val, map_manager.StringType, nil
}

func (sb *Switchboard) getCustomDependencies(key string) (map[string]interface{}, error) {
	dep, ok := sb.deps[key]
	if !ok || len(dep) == 0 {
		return nil, fmt.Errorf("missing dependencies for key: [%s]", key)
	}
	data := make(map[string]interface{})
	for _, d := range dep {
		val, _, ok := d.mm.Get(d.vt, d.key)
		if !ok {
			return nil, fmt.Errorf("key %s not found in source mapmanager", key)
		}
		data[fmt.Sprintf("%s:%s", d.mm.GetName(), d.key)] = val
		data[d.key] = val
	}
	return data, nil
}

func (sb *Switchboard) getCustomResult(k string, c ValueCalculator, d map[string]interface{}) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Ef("ValueCalculator for key [%s] panicked with error: %s\n", k, r)
		}
	}()
	return c.Result(d)
}

func (sb *Switchboard) GetAsJsonPayload(key string, additional map[string]interface{}) ([]byte, error) {
	val, vt, err := sb.Get(key)
	if err != nil {
		return nil, fmt.Errorf("key %s not found", key)
	}

	jsonData := map[string]interface{}{
		"key":       key,
		"value":     val,
		"valueType": vt.String(),
		"ok":        true,
	}

	for k, v := range additional {
		if _, present := jsonData[k]; !present {
			jsonData[k] = v
		}
	}

	return json.Marshal(jsonData)
}

func (sb *Switchboard) Json() ([]byte, error) {
	data := map[string]interface{}{}
	for k := range sb.resolvers {
		value, vt, err := sb.Get(k)
		if err != nil {
			data[k] = map[string]interface{}{
				"key":       k,
				"value":     "",
				"valueType": map_manager.UnknownType.String(),
				"ok":        false,
				"err":       err.Error(),
			}
		}
		data[k] = map[string]interface{}{
			"key":       k,
			"value":     value,
			"valueType": vt.String(),
			"ok":        true,
		}
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
