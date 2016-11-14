package rainfall

import (
	"encoding/json"
	"os"
	"sort"
)

type Location struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type LocationSorter []Location

func (s LocationSorter) Len() int           { return len(s) }
func (s LocationSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s LocationSorter) Less(i, j int) bool { return s[i].Name < s[j].Name }

type LocationStore struct {
	path   string
	locMap map[string]Location
}

func NewLocationStore(path string) *LocationStore {
	locMap := make(map[string]Location)

	return &LocationStore{
		path:   path,
		locMap: locMap,
	}
}

func (s *LocationStore) Load() error {
	if _, err := os.Stat(s.path); err != nil {
		// ファイルが存在しない
		return nil
	}

	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	d := json.NewDecoder(file)

	var locations []Location
	err = d.Decode(&locations)
	if err != nil {
		return err
	}

	s.locMap = make(map[string]Location)

	for _, loc := range locations {
		s.locMap[loc.Name] = loc
	}

	return nil
}

func (s *LocationStore) Save() error {
	file, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)

	err = e.Encode(s.Locations())

	return err
}

func (s *LocationStore) Get(name string) (Location, bool) {
	loc, ok := s.locMap[name]
	return loc, ok
}

func (s *LocationStore) Set(loc Location) {
	s.locMap[loc.Name] = loc
}

func (s *LocationStore) Del(name string) {
	delete(s.locMap, name)
}

func (s *LocationStore) Locations() []Location {
	locations := make([]Location, 0, len(s.locMap))

	for _, loc := range s.locMap {
		locations = append(locations, loc)
	}

	sort.Sort(LocationSorter(locations))

	return locations
}
