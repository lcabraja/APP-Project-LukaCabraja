package map_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lcabraja/APP-Project-LukaCabraja/log"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MapManager struct {
	name string

	mapStr      map[string]string
	mapJson     map[string]string
	mapInt      map[string]int
	mapFloat64  map[string]float64
	mapBool     map[string]bool
	mapDuration map[string]time.Duration

	validStr      map[string]func(string) bool
	validJson     map[string]func(string) bool
	validInt      map[string]func(int) bool
	validFloat64  map[string]func(float64) bool
	validDuration map[string]func(time.Duration) bool

	muStr      sync.Mutex
	muJson     sync.Mutex
	muInt      sync.Mutex
	muFloat64  sync.Mutex
	muBool     sync.Mutex
	muDuration sync.Mutex

	subscriptions map[string]*MapManagerEventSubscription
}

func NewMapManager(name string) *MapManager {
	v := regexp.MustCompile("[.:'\"]")
	if r := v.Find([]byte(name)); r != nil {
		panic(fmt.Sprintf("invalid character '%s' in map manager name", string(r)))
	}
	return &MapManager{
		name: name,

		mapStr:      make(map[string]string),
		mapJson:     make(map[string]string),
		mapInt:      make(map[string]int),
		mapFloat64:  make(map[string]float64),
		mapBool:     make(map[string]bool),
		mapDuration: make(map[string]time.Duration),

		validStr:      map[string]func(string) bool{},
		validJson:     map[string]func(string) bool{},
		validInt:      map[string]func(int) bool{},
		validFloat64:  map[string]func(float64) bool{},
		validDuration: map[string]func(time.Duration) bool{},

		muStr:      sync.Mutex{},
		muJson:     sync.Mutex{},
		muInt:      sync.Mutex{},
		muFloat64:  sync.Mutex{},
		muBool:     sync.Mutex{},
		muDuration: sync.Mutex{},

		subscriptions: make(map[string]*MapManagerEventSubscription),
	}
}

func (mm *MapManager) fireEvent(e MapManagerEvent, vt ValueType, key string) {
	for _, sub := range mm.subscriptions {
		if sub.Event == e || sub.Event == EventFired {
			go sub.Handler(e, vt, key, mm)
		}
	}
}

func (mm *MapManager) GetName() string {
	return mm.name
}

func (mm *MapManager) Subscribe(sub *MapManagerEventSubscription) string {
	uuid := uuid.NewString()
	mm.subscriptions[uuid] = sub
	return uuid
}

func (mm *MapManager) Unsubscribe(uuid string) {
	delete(mm.subscriptions, uuid)
}

func (mm *MapManager) SetString(key, value string) error {
	mm.muStr.Lock()
	defer mm.muStr.Unlock()

	if validator, ok := mm.validStr[key]; ok {
		if !validator(value) {
			return fmt.Errorf("cannot set %s to [%s]", key, value)
		}
	}

	prev, ok := mm.mapStr[key]
	mm.mapStr[key] = value
	mm.fireEvent(RecordUpdated, StringType, key)
	mm.fireEvent(StringUpdated, StringType, key)
	if ok && prev != value {
		mm.fireEvent(RecordChanged, StringType, key)
		mm.fireEvent(StringChanged, StringType, key)
	}
	return nil
}

func (mm *MapManager) GetString(key string) (string, bool) {
	mm.muStr.Lock()
	defer mm.muStr.Unlock()

	value, ok := mm.mapStr[key]
	mm.fireEvent(RecordRead, StringType, key)
	mm.fireEvent(StringRead, StringType, key)
	return value, ok
}

func (mm *MapManager) SetJson(key string, value []byte) error {
	mm.muJson.Lock()
	defer mm.muJson.Unlock()

	if !json.Valid(value) {
		return errors.New("cannot set invalid json")
	}
	if validator, ok := mm.validJson[key]; ok {
		if !validator(string(value)) {
			return fmt.Errorf("cannot set %s to [%s]", key, string(value))
		}
	}

	prev, ok := mm.mapJson[key]
	mm.mapJson[key] = string(value)
	mm.fireEvent(RecordUpdated, JsonType, key)
	mm.fireEvent(JsonUpdated, JsonType, key)
	if ok && prev != string(value) {
		mm.fireEvent(RecordChanged, JsonType, key)
		mm.fireEvent(JsonChanged, JsonType, key)
	}
	return nil
}

func (mm *MapManager) GetJson(key string) (string, bool) {
	mm.muJson.Lock()
	defer mm.muJson.Unlock()

	value, ok := mm.mapJson[key]
	mm.fireEvent(RecordRead, JsonType, key)
	mm.fireEvent(JsonRead, JsonType, key)
	return value, ok
}

func (mm *MapManager) SetInt(key string, value int) error {
	mm.muInt.Lock()
	defer mm.muInt.Unlock()

	if validator, ok := mm.validInt[key]; ok {
		if !validator(value) {
			return fmt.Errorf("cannot set %s to [%d]", key, value)
		}
	}

	prev, ok := mm.mapInt[key]
	mm.mapInt[key] = value
	mm.fireEvent(RecordUpdated, IntType, key)
	mm.fireEvent(IntUpdated, IntType, key)
	if ok && prev != value {
		mm.fireEvent(RecordChanged, IntType, key)
		mm.fireEvent(IntChanged, IntType, key)
	}
	return nil
}

func (mm *MapManager) OffsetInt(key string, offset int) error {
	if offset == 0 {
		return nil
	}

	var (
		prev int
		ok   bool
	)

	if prev, ok = mm.GetInt(key); !ok {
		return fmt.Errorf("cannot offset %s, key not found", key)
	}

	return mm.SetInt(key, prev+offset)
}

func (mm *MapManager) GetInt(key string) (int, bool) {
	mm.muInt.Lock()
	defer mm.muInt.Unlock()
	value, ok := mm.mapInt[key]
	mm.fireEvent(RecordRead, IntType, key)
	mm.fireEvent(IntRead, IntType, key)
	return value, ok
}

func (mm *MapManager) SetFloat(key string, value float64) error {
	mm.muFloat64.Lock()
	defer mm.muFloat64.Unlock()

	if validator, ok := mm.validFloat64[key]; ok {
		if !validator(value) {
			return fmt.Errorf("cannot set %s to [%f]", key, value)
		}
	}

	prev, ok := mm.mapFloat64[key]
	mm.mapFloat64[key] = value
	mm.fireEvent(RecordUpdated, FloatType, key)
	mm.fireEvent(FloatUpdated, FloatType, key)
	if ok && prev != value {
		mm.fireEvent(RecordChanged, FloatType, key)
		mm.fireEvent(FloatChanged, FloatType, key)
	}
	return nil
}

func (mm *MapManager) GetFloat(key string) (float64, bool) {
	mm.muFloat64.Lock()
	defer mm.muFloat64.Unlock()
	value, ok := mm.mapFloat64[key]
	mm.fireEvent(RecordRead, FloatType, key)
	mm.fireEvent(FloatRead, FloatType, key)
	return value, ok
}

func (mm *MapManager) SetBool(key string, value bool) {
	mm.muBool.Lock()
	defer mm.muBool.Unlock()
	prev, ok := mm.mapBool[key]
	mm.mapBool[key] = value
	mm.fireEvent(RecordUpdated, BoolType, key)
	mm.fireEvent(BoolUpdated, BoolType, key)
	if ok && prev != value {
		mm.fireEvent(RecordChanged, BoolType, key)
		mm.fireEvent(BoolChanged, BoolType, key)
	}
}

func (mm *MapManager) GetBool(key string) (bool, bool) {
	mm.muBool.Lock()
	defer mm.muBool.Unlock()
	value, ok := mm.mapBool[key]
	mm.fireEvent(RecordRead, BoolType, key)
	mm.fireEvent(BoolRead, BoolType, key)
	return value, ok
}

func (mm *MapManager) SetDuration(key string, value time.Duration) error {
	mm.muDuration.Lock()
	defer mm.muDuration.Unlock()

	if validator, ok := mm.validDuration[key]; ok {
		if !validator(value) {
			return fmt.Errorf("cannot set %s to [%v]", key, value)
		}
	}

	prev, ok := mm.mapDuration[key]
	mm.mapDuration[key] = value
	mm.fireEvent(RecordUpdated, DurationType, key)
	mm.fireEvent(DurationUpdated, DurationType, key)
	if ok && prev != value {
		mm.fireEvent(RecordChanged, DurationType, key)
		mm.fireEvent(DurationChanged, DurationType, key)
	}
	return nil
}

func (mm *MapManager) GetDuration(key string) (time.Duration, bool) {
	mm.muDuration.Lock()
	defer mm.muDuration.Unlock()
	value, ok := mm.mapDuration[key]
	mm.fireEvent(RecordRead, DurationType, key)
	mm.fireEvent(DurationRead, DurationType, key)
	return value, ok
}

func (mm *MapManager) Set(vt ValueType, key string, data interface{}) error {
	err := fmt.Errorf("provided value [%-16s]:16 not of type %s", data, vt)
	switch vt {
	case StringType:
		data, ok := data.(string)
		if !ok {
			return err
		}
		if err := mm.SetString(key, data); err != nil {
			return err
		}
	case JsonType:
		data, ok := data.([]byte)
		if !ok {
			return err
		}
		if err := mm.SetJson(key, data); err != nil {
			return err
		}
	case IntType:
		data, ok := data.(int)
		if !ok {
			return err
		}
		if err := mm.SetInt(key, data); err != nil {
			return err
		}
	case FloatType:
		data, ok := data.(float64)
		if !ok {
			return err
		}
		if err := mm.SetFloat(key, data); err != nil {
			return err
		}
	case BoolType:
		data, ok := data.(bool)
		if !ok {
			return err
		}
		mm.SetBool(key, data)
	case DurationType:
		data, ok := data.(time.Duration)
		if !ok {
			return err
		}
		if err := mm.SetDuration(key, data); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown value type %s", vt)
	}
	return nil
}

func (mm *MapManager) Get(vt ValueType, key string) (interface{}, ValueType, bool) {
	var (
		value interface{}
		ok    bool
	)

	switch vt {
	case StringType:
		if value, ok = mm.GetString(key); ok {
			return value, StringType, ok
		}

	case JsonType:
		if value, ok = mm.GetJson(key); ok {
			return value, JsonType, ok
		}

	case IntType:
		if value, ok = mm.GetInt(key); ok {
			return value, IntType, ok
		}

	case FloatType:
		if value, ok = mm.GetFloat(key); ok {
			return value, FloatType, ok
		}

	case BoolType:
		if value, ok = mm.GetBool(key); ok {
			return value, BoolType, ok
		}

	case DurationType:
		if value, ok = mm.GetDuration(key); ok {
			return value, DurationType, ok
		}
	}

	return mm.GetAny(key)
}

func (mm *MapManager) GetAny(key string) (interface{}, ValueType, bool) {
	var (
		value interface{}
		ok    bool
	)

	if value, ok = mm.GetString(key); ok {
		return value, StringType, ok
	}

	if value, ok = mm.GetJson(key); ok {
		return value, JsonType, ok
	}

	if value, ok = mm.GetInt(key); ok {
		return value, IntType, ok
	}

	if value, ok = mm.GetFloat(key); ok {
		return value, FloatType, ok
	}

	if value, ok = mm.GetBool(key); ok {
		return value, BoolType, ok
	}

	if value, ok = mm.GetDuration(key); ok {
		return value, DurationType, ok
	}

	return nil, UnknownType, false
}

func (mm *MapManager) GetAsJsonPayload(vt ValueType, key string, additional map[string]interface{}) ([]byte, error) {
	value, vt, ok := mm.Get(vt, key)
	if !ok {
		return nil, errors.New("key not found")
	}
	jsonData := map[string]interface{}{
		"key":       key,
		"value":     value,
		"valueType": vt.String(),
		"ok":        true,
	}

	for k, v := range additional {
		if _, present := jsonData[k]; !present {
			jsonData[k] = v
		}
	}

	bytes, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (mm *MapManager) EnsureExists(vt ValueType, key string, defaultValue interface{}, validator interface{}) error {
	switch vt {
	case StringType:
		if _, exists := mm.mapStr[key]; exists {
			return fmt.Errorf("string key %s already exists", key)
		}

		mm.mapStr[key] = ""

		if defaultValue != nil && vt.Validate(defaultValue) == nil {
			log.Devf("invalid default value %v for type %s\n\n\n", defaultValue, vt)
			d := defaultValue.(string)
			mm.mapStr[key] = d
		}

		switch validator.(type) {
		case func(string) bool:
			mm.validStr[key] = validator.(func(string) bool)
		default:
			return fmt.Errorf("invalid validator type %T", validator)
		}
	case JsonType:
		if _, exists := mm.mapJson[key]; exists {
			return fmt.Errorf("json key %s already exists", key)
		}

		mm.mapJson[key] = ""

		if defaultValue != nil && vt.Validate(defaultValue) == nil {
			d := defaultValue.(string)
			if !json.Valid([]byte(d)) {
				return errors.New("cannot set invalid json")
			}
			mm.mapJson[key] = d
		}

		switch validator.(type) {
		case func(string) bool:
			mm.validJson[key] = validator.(func(string) bool)
		default:
			return fmt.Errorf("invalid validator type %T", validator)
		}

	case IntType:
		if _, exists := mm.mapInt[key]; exists {
			return fmt.Errorf("integer key %s already exists", key)
		}

		mm.mapInt[key] = 0

		if defaultValue != nil && vt.Validate(defaultValue) == nil {
			d := defaultValue.(int)
			mm.mapInt[key] = d
		}

		switch validator.(type) {
		case func(int) bool:
			mm.validInt[key] = validator.(func(int) bool)
		default:
			return fmt.Errorf("invalid validator type %T", validator)
		}
	case FloatType:
		if _, exists := mm.mapFloat64[key]; exists {
			return fmt.Errorf("float key %s already exists", key)
		}

		mm.mapFloat64[key] = 0

		if defaultValue != nil && vt.Validate(defaultValue) == nil {
			d := defaultValue.(float64)
			mm.mapFloat64[key] = d
		}

		switch validator.(type) {
		case func(float64) bool:
			mm.validFloat64[key] = validator.(func(float64) bool)
		default:
			return fmt.Errorf("invalid validator type %T", validator)
		}
	case BoolType:
		if _, exists := mm.mapBool[key]; exists {
			return fmt.Errorf("bool key %s already exists", key)
		}

		mm.mapBool[key] = false

		if defaultValue != nil && vt.Validate(defaultValue) == nil {
			d := defaultValue.(bool)
			mm.mapBool[key] = d
		}
	case DurationType:
		if _, exists := mm.mapDuration[key]; exists {
			return fmt.Errorf("duration key %s already exists", key)
		}

		mm.mapDuration[key] = 0

		if defaultValue != nil && vt.Validate(defaultValue) == nil {
			d := defaultValue.(time.Duration)
			mm.mapDuration[key] = d
		}

		switch validator.(type) {
		case func(time.Duration) bool:
			mm.validDuration[key] = validator.(func(time.Duration) bool)
		default:
			return fmt.Errorf("invalid validator type %T", validator)
		}
	}
	return fmt.Errorf("invalid value type %s", vt.String())
}

func (mm *MapManager) Reset() {
	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		mm.muStr.Lock()
		mm.mapStr = make(map[string]string)
		mm.muStr.Unlock()
		wg.Done()
	}()

	go func() {
		mm.muJson.Lock()
		mm.mapJson = make(map[string]string)
		mm.muJson.Unlock()
		wg.Done()
	}()

	go func() {
		mm.muInt.Lock()
		mm.mapInt = make(map[string]int)
		mm.muInt.Unlock()
		wg.Done()
	}()

	go func() {
		mm.muFloat64.Lock()
		mm.mapFloat64 = make(map[string]float64)
		mm.muFloat64.Unlock()
		wg.Done()
	}()

	go func() {
		mm.muBool.Lock()
		mm.mapBool = make(map[string]bool)
		mm.muBool.Unlock()
		wg.Done()
	}()

	go func() {
		mm.muDuration.Lock()
		mm.mapDuration = make(map[string]time.Duration)
		mm.muDuration.Unlock()
		wg.Done()
	}()

	wg.Wait()
}

func (mm *MapManager) Json(indent bool) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	var d = map[string]interface{}{
		"string":   mm.mapStr,
		"json":     mm.mapJson,
		"int":      mm.mapInt,
		"float":    mm.mapFloat64,
		"bool":     mm.mapBool,
		"duration": mm.mapDuration,
	}

	if indent {
		data, err = json.MarshalIndent(d, "", "  ")
	} else {
		data, err = json.Marshal(d)
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}
