package model

// CookieMap is a map holding all the cookies and allows us to efficiently add new cookies and query on existing data.
// Instead of storing repetitive cookie names which are huge strings, it assigns an id to each cookie name and store only
// the id along with the frequency to remove duplication and only store each cookie name only once
type CookieMap struct {
	cookieId   uint64
	IDToName   map[uint64]string
	NameToID   map[string]uint64
	DateToFreq map[string]map[uint64]uint64
}

// NewCookieMap creates and returns a new CookieMap
func NewCookieMap() *CookieMap {
	return &CookieMap{
		IDToName:   make(map[uint64]string),
		NameToID:   make(map[string]uint64),
		DateToFreq: make(map[string]map[uint64]uint64),
	}
}

// Add adds the given cookie to the map
func (cm *CookieMap) Add(c Cookie) {
	if _, ok := cm.DateToFreq[c.Date]; !ok {
		cm.DateToFreq[c.Date] = make(map[uint64]uint64)
	}

	cm.DateToFreq[c.Date][cm.getId(c.Name)]++
}

// Get returns a frequency map containing all the cookie names along with their frequencies on the given date. It uses
// IDToName map to get the names of the cookies that are stored to return a mapping of [name, freq]
func (cm *CookieMap) Get(date string) map[string]uint64 {
	IDToFreq, ok := cm.DateToFreq[date]
	if !ok {
		return nil
	}

	nameToFreq := make(map[string]uint64)
	for id, freq := range IDToFreq {
		nameToFreq[cm.IDToName[id]] = freq
	}

	return nameToFreq
}

// getId returns an uint64 id for the given cookie name. If there is no id available, it allots a new id for this name
func (cm *CookieMap) getId(name string) uint64 {
	id, ok := cm.NameToID[name]
	if !ok {
		cm.cookieId++
		id = cm.cookieId
		cm.NameToID[name] = id
		cm.IDToName[id] = name
	}
	return id
}
