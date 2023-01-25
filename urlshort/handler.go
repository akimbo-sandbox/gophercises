package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirect, exists := pathsToUrls[r.URL.Path]
		if exists {
			http.Redirect(w, r, redirect, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

type pathToUrl struct {
	Path string
	URL string
}

func parseYAML(y []byte) ([]pathToUrl, error) {
	var paths []pathToUrl
	err := yaml.Unmarshal(y, &paths)
	if err != nil {
		return nil, err
	}
	return paths, err
}

func buildMap(paths []pathToUrl) map[string]string {
	pMap := make(map[string]string)
	for _, p := range paths {
		pMap[p.Path] = p.URL
	}
	return pMap
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
