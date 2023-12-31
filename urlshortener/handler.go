package urlshortener

import (
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		newPath, found := pathsToUrls[req.URL.Path]
		if found {
			http.Redirect(rw, req, newPath, http.StatusSeeOther)
		}
		fallback.ServeHTTP(rw, req)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathsToUrls, err := parse(yml)
	if err != nil {
		return nil, err
	}
	return MapHandler(pathsToUrls, fallback), nil
}

func parse(yml []byte) (map[string]string, error) {
	var u []struct {
		Path string
		Url string	
	}
	err := yaml.Unmarshal(yml, &u)
	if err != nil {
		return nil, err
	}

	r := make(map[string]string, 0)
	for _, item := range u {
		r[item.Path] = item.Url
	}
	return r, nil
}