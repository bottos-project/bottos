package abi

import (
	"encoding/json"
	//"errors"
	"sort"
	"strings"
)

//NoValueError value
//var NoValueError = errors.New("No value for this key")

//KeyIndex struct
type KeyIndex struct {
	Key   string
	Index int
}

//ByIndex struct
type ByIndex []KeyIndex

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

//Pair struct
type Pair struct {
	key   string
	value interface{}
}

//Key function
func (kv *Pair) Key() string {
	return kv.key
}

//Value function
func (kv *Pair) Value() interface{} {
	return kv.value
}

//ByPair function
type ByPair struct {
	Pairs    []*Pair
	LessFunc func(a *Pair, j *Pair) bool
}

func (a ByPair) Len() int           { return len(a.Pairs) }
func (a ByPair) Swap(i, j int)      { a.Pairs[i], a.Pairs[j] = a.Pairs[j], a.Pairs[i] }
func (a ByPair) Less(i, j int) bool { return a.LessFunc(a.Pairs[i], a.Pairs[j]) }

//StringPair struct
type StringPair struct {
	Key   string
	Value string
}

//FeildMap struct
type FeildMap struct {
	keys   []string
	values map[string]interface{}
}

//New new a FeildMap struct
func New() *FeildMap {
	o := FeildMap{}
	o.keys = []string{}
	o.values = map[string]interface{}{}
	return &o
}

//GetStringVal get string value by key
func (o *FeildMap) GetStringVal(key string) (string, bool) {
	val, exists := o.values[key]
	if !exists {
		return "", false
	}

	switch val.(type) {
	case string:
		return val.(string), true
	}
	return "", false
}

//GetStringPair get key value pairs
func (o *FeildMap) GetStringPair() []StringPair {
	pairs := make([]StringPair, len(o.keys))
	for i, k := range o.keys {
		v, _ := o.GetStringVal(k)
		pairs[i].Key = k
		pairs[i].Value = v
	}
	return pairs
}

//Get get interface value by key
func (o *FeildMap) Get(key string) (interface{}, bool) {
	val, exists := o.values[key]
	return val, exists
}

//Set set interface value by key
func (o *FeildMap) Set(key string, value interface{}) {
	_, exists := o.values[key]
	if !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

//Delete delete key
func (o *FeildMap) Delete(key string) {
	// check key is in use
	_, ok := o.values[key]
	if !ok {
		return
	}
	// remove from keys
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	// remove from values
	delete(o.values, key)
}

//Keys list all keys
func (o *FeildMap) Keys() []string {
	return o.keys
}

// SortKeys Sort the map keys using your sort func
func (o *FeildMap) SortKeys(sortFunc func(keys []string)) {
	sortFunc(o.keys)
}

// Sort Sort the map using your sort func
func (o *FeildMap) Sort(lessFunc func(a *Pair, b *Pair) bool) {
	pairs := make([]*Pair, len(o.keys))
	for i, key := range o.keys {
		pairs[i] = &Pair{key, o.values[key]}
	}

	sort.Sort(ByPair{pairs, lessFunc})

	for i, pair := range pairs {
		o.keys[i] = pair.key
	}
}

//UnmarshalJSON function
func (o *FeildMap) UnmarshalJSON(b []byte) error {
	var err error
	err = mapStringToFeildMap(string(b), o)
	if err != nil {
		return err
	}
	return nil
}

//mapStringToFeildMap function
func mapStringToFeildMap(s string, o *FeildMap) error {
	// parse string into map
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return err
	}
	// Get the order of the keys
	orderedKeys := []KeyIndex{}
	for k := range m {
		kEscaped := strings.Replace(k, `"`, `\"`, -1)
		kQuoted := `"` + kEscaped + `"`
		// Find how much content exists before this key.
		// If all content from this key and after is replaced with a close
		// brace, it should still form a valid json string.
		sTrimmed := s
		for len(sTrimmed) > 0 {
			lastIndex := strings.LastIndex(sTrimmed, kQuoted)
			if lastIndex == -1 {
				break
			}
			sTrimmed = sTrimmed[0:lastIndex]
			sTrimmed = strings.TrimSpace(sTrimmed)
			if len(sTrimmed) > 0 && sTrimmed[len(sTrimmed)-1] == ',' {
				sTrimmed = sTrimmed[0 : len(sTrimmed)-1]
			}
			maybeValidJson := sTrimmed + "}"
			testMap := map[string]interface{}{}
			err := json.Unmarshal([]byte(maybeValidJson), &testMap)
			if err == nil {
				// record the position of this key in s
				ki := KeyIndex{
					Key:   k,
					Index: len(sTrimmed),
				}
				orderedKeys = append(orderedKeys, ki)
				// shorten the string to get the next key
				startOfValueIndex := lastIndex + len(kQuoted)
				valueStr := s[startOfValueIndex : len(s)-1]
				valueStr = strings.TrimSpace(valueStr)
				if len(valueStr) > 0 && valueStr[0] == ':' {
					valueStr = valueStr[1:]
				}
				valueStr = strings.TrimSpace(valueStr)
				if valueStr[0] == '{' {
					// if the value for this key is a map
					// find end of valueStr by removing everything after last }
					// until it forms valid json
					hasValidJson := false
					i := 1
					for i < len(valueStr) && !hasValidJson {
						if valueStr[i] != '}' {
							i = i + 1
							continue
						}
						subTestMap := map[string]interface{}{}
						testValue := valueStr[0 : i+1]
						err = json.Unmarshal([]byte(testValue), &subTestMap)
						if err == nil {
							hasValidJson = true
							valueStr = testValue
							break
						}
						i = i + 1
					}
					// convert to FeildMap
					// this may be recursive it values in the map are also maps
					if hasValidJson {
						newMap := New()
						err := mapStringToFeildMap(valueStr, newMap)
						if err != nil {
							return err
						}
						m[k] = *newMap
					}
				} else if valueStr[0] == '[' {
					// if the value for this key is a slice
					// find end of valueStr by removing everything after last ]
					// until it forms valid json
					hasValidJson := false
					i := 1
					for i < len(valueStr) && !hasValidJson {
						if valueStr[i] != ']' {
							i = i + 1
							continue
						}
						subTestSlice := []interface{}{}
						testValue := valueStr[0 : i+1]
						err = json.Unmarshal([]byte(testValue), &subTestSlice)
						if err == nil {
							hasValidJson = true
							valueStr = testValue
							break
						}
						i = i + 1
					}
					// convert to slice with any map items converted to
					// FeildMaps
					// this may be recursive if values in the slice are slices
					if hasValidJson {
						var newSlice []interface{}
						err := sliceStringToSliceWithFeildMaps(valueStr, &newSlice)
						if err != nil {
							return err
						}
						m[k] = newSlice
					}
				} else {
					o.Set(k, m[k])
				}
				break
			}
		}
	}
	// Sort the keys
	sort.Sort(ByIndex(orderedKeys))
	// Convert sorted keys to string slice
	k := []string{}
	for _, ki := range orderedKeys {
		k = append(k, ki.Key)
	}
	// Set the FeildMap values
	o.values = m
	o.keys = k
	return nil
}

func sliceStringToSliceWithFeildMaps(valueStr string, newSlice *[]interface{}) error {
	// if the value for this key is a []interface, convert any map items to an FeildMap.
	// find end of valueStr by removing everything after last ]
	// until it forms valid json
	itemsStr := strings.TrimSpace(valueStr)
	itemsStr = itemsStr[1 : len(itemsStr)-1]
	// get next item in the slice
	itemIndex := 0
	startItem := 0
	endItem := 0
	for endItem <= len(itemsStr) {
		couldBeItemEnd := false
		couldBeItemEnd = couldBeItemEnd || endItem == len(itemsStr)
		couldBeItemEnd = couldBeItemEnd || (endItem < len(itemsStr) && itemsStr[endItem] == ',')
		if !couldBeItemEnd {
			endItem = endItem + 1
			continue
		}
		// if this substring compiles to json, it's the next item
		possibleItemStr := strings.TrimSpace(itemsStr[startItem:endItem])
		var possibleItem interface{}
		err := json.Unmarshal([]byte(possibleItemStr), &possibleItem)
		if err != nil {
			endItem = endItem + 1
			continue
		}
		if possibleItemStr[0] == '{' {
			// if item is map, convert to FeildMap
			oo := New()
			err := mapStringToFeildMap(possibleItemStr, oo)
			if err != nil {
				return err
			}
			// add new FeildMap item to new slice
			slice := *newSlice
			slice = append(slice, *oo)
			*newSlice = slice
		} else if possibleItemStr[0] == '[' {
			// if item is slice, convert to slice with FeildMaps
			var newItem []interface{}
			err := sliceStringToSliceWithFeildMaps(possibleItemStr, &newItem)
			if err != nil {
				return err
			}
			// replace original slice item with new slice
			slice := *newSlice
			slice = append(slice, newItem)
			*newSlice = slice
		} else {
			// any non-slice and non-map item, just add json parsed item
			slice := *newSlice
			slice = append(slice, possibleItem)
			*newSlice = slice
		}
		// remove this item from itemsStr
		startItem = endItem + 1
		endItem = endItem + 1
		itemIndex = itemIndex + 1
	}
	return nil
}

//MarshalJSON function
func (o FeildMap) MarshalJSON() ([]byte, error) {
	s := "{"
	for _, k := range o.keys {
		// add key
		kEscaped := strings.Replace(k, `"`, `\"`, -1)
		s = s + `"` + kEscaped + `":`
		// add value
		v := o.values[k]
		vBytes, err := json.Marshal(v)
		if err != nil {
			return []byte{}, err
		}
		s = s + string(vBytes) + ","
	}
	if len(o.keys) > 0 {
		s = s[0 : len(s)-1]
	}
	s = s + "}"
	return []byte(s), nil
}
