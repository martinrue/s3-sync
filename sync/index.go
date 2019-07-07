package sync

import "encoding/json"

// Object holds data about a key and its tags.
type Object struct {
	Key      string `json:"key"`
	Tags     string `json:"tags"`
	Filepath string `json:"-"`
	IsNew    bool   `json:"-"`
}

// Index tracks which objects are stored in a bucket.
type Index struct {
	Objects []*Object
}

// LoadJSON unmarshals the supplied JSON data into a set of objects.
func (i *Index) LoadJSON(data []byte) error {
	return json.Unmarshal(data, &i.Objects)
}

// SaveJSON marshals the data to JSON.
func (i *Index) SaveJSON() ([]byte, error) {
	return json.Marshal(i.Objects)
}

// Add adds a new object to the index.
func (i *Index) Add(key string, tags string, filepath string, isNew bool) {
	i.Objects = append(i.Objects, &Object{key, tags, filepath, isNew})
}

// Diff builds a new index containing objects present in the supplied index but not the origin.
func (i *Index) Diff(index *Index) *Index {
	diff := &Index{}

	objects := i.buildMap()

	for _, object := range index.Objects {
		right, ok := objects[object.Key]

		if !ok || object.Tags != right.Tags {
			diff.Add(object.Key, object.Tags, object.Filepath, !ok)
		}
	}

	return diff
}

func (i *Index) buildMap() map[string]*Object {
	objects := make(map[string]*Object)

	for _, object := range i.Objects {
		objects[object.Key] = object
	}

	return objects
}
