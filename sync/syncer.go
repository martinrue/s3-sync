package sync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"strconv"
	"strings"
)

// Syncer synchronises files from a directory to a storage service.
type Syncer struct {
	Store ObjectStore
	Log   func(string, ...interface{})
}

// Run begins a sync using the specified directory.
func (s *Syncer) Run(dir string, exts string) (string, error) {
	files, err := scanDirectory(dir, exts)
	if err != nil {
		return "", err
	}

	localIndex := &Index{}

	for _, file := range files {
		filepath := path.Join(dir, file)

		key, err := CreateCASKey(filepath)
		if err != nil {
			return "", err
		}

		tags := strings.Join(parseTags(file), " ")
		localIndex.Add(key, tags, filepath, false)
	}

	if len(localIndex.Objects) == 0 {
		return "", fmt.Errorf("no files matching (%v) found in: %v", exts, dir)
	}

	remoteIndex, _ := s.fetchRemoteIndex("index.json")

	diff := remoteIndex.Diff(localIndex)

	uploaded := 0

	for _, object := range diff.Objects {
		if !object.IsNew {
			continue
		}

		s.Log("uploading: %v", object.Key)

		file, err := os.Open(object.Filepath)
		if err != nil {
			return "", err
		}

		contentType := mime.TypeByExtension(path.Ext(object.Filepath))

		if err := s.Store.Put(object.Key, contentType, file); err != nil {
			return "", err
		}

		file.Close()

		uploaded++
	}

	s.Log("processed %d objects", len(localIndex.Objects))
	s.Log("uploaded %d objects", uploaded)

	return s.saveRemoteIndex(localIndex, "index.json")
}

func (s *Syncer) fetchRemoteIndex(key string) (*Index, error) {
	index := &Index{}

	bytes, err := s.Store.Get(key)
	if err != nil {
		return index, err
	}

	if err := index.LoadJSON(bytes); err != nil {
		return index, err
	}

	return index, nil
}

func (s *Syncer) saveRemoteIndex(index *Index, key string) (string, error) {
	data, err := index.SaveJSON()
	if err != nil {
		return "", err
	}

	if err := s.Store.Put(key, "application/json", bytes.NewReader(data)); err != nil {
		return "", err
	}

	return string(data), nil
}

func scanDirectory(dir string, exts string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	matches := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !hasValidExtension(file.Name(), exts) {
			continue
		}

		matches = append(matches, file.Name())
	}

	return matches, nil
}

func hasValidExtension(filename string, exts string) bool {
	if exts == "" {
		return true
	}

	extension := strings.TrimPrefix(path.Ext(filename), ".")

	for _, ext := range strings.Split(exts, ",") {
		if extension == strings.TrimSpace(ext) {
			return true
		}
	}

	return false
}

func parseTags(filename string) []string {
	basename := strings.TrimSuffix(filename, path.Ext(filename))
	tags := make([]string, 0)

	for _, part := range strings.Split(basename, " ") {
		if _, err := strconv.Atoi(part); err == nil {
			continue
		}

		tags = append(tags, strings.TrimSpace(part))
	}

	return tags
}
