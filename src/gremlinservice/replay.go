package main

import (
	"github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"io/ioutil"
	"path/filepath"
)

func replayServices() (*api_v2.GetServicesResponse, error) {
	b, err := ioutil.ReadFile(filepath.Join("data", "services"))
	if err != nil {
		return nil, err
	}
	services := &api_v2.GetServicesResponse{}
	if err := services.Unmarshal(b); err != nil {
		return nil, err
	}
	return services, nil
}

func replayChunks(chunksDir string) (before map[string]*api_v2.SpansResponseChunk, after map[string]*api_v2.SpansResponseChunk, err error) {
	if before, err = readChunks(filepath.Join(chunksDir, "before")); err != nil {
		return
	}

	if after, err = readChunks(filepath.Join(chunksDir, "after")); err != nil {
		return
	}

	return
}

func readChunks(dir string) (map[string]*api_v2.SpansResponseChunk, error) {
	chunks := make(map[string]*api_v2.SpansResponseChunk, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		chunk, err := func(filename string) (*api_v2.SpansResponseChunk, error) {
			b, err := ioutil.ReadFile(filepath.Join(dir, filename))
			if err != nil {
				return nil, err
			}

			chunk := &api_v2.SpansResponseChunk{}

			if err := chunk.Unmarshal(b); err != nil {
				return nil, err
			}

			return chunk, nil
		}(f.Name())

		if err != nil {
			return nil, err
		}

		chunks[f.Name()] = chunk
	}

	return chunks, nil
}
