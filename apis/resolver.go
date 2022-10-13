package apis

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/blang/semver/v4"
)

var ErrAPINotFound = errors.New("api not found")

// Resolver function for resolving type name
type Resolver func(string) (interface{}, error)

// PackageMap contains a list of packages with their Resolver func
var packageMap = map[string]Resolver{}

var packageMapMu = &sync.RWMutex{}

// RegisterResolver adds a package to packageMap with its resolver.
func RegisterResolver(key string, resolver Resolver) {
	packageMapMu.Lock()
	defer packageMapMu.Unlock()
	packageMap[key] = resolver
}

// Resolveaw resolves the raw type for the requested type.
func Resolve(apiVersion string, typename string) (interface{}, error) {
	availableModules := APIModuleVersions()

	// Guard read access to packageMap
	packageMapMu.RLock()
	defer packageMapMu.RUnlock()
	apiGroup, reqVer := ParseAPIVersion(apiVersion)
	foundVer, ok := availableModules[apiGroup]
	if ok {
		if semverGreater(reqVer, foundVer) {
			return nil, fmt.Errorf("requested version was %s, but only %s is available", reqVer, foundVer)
		}
	}
	resolver, ok := packageMap[apiGroup]
	if !ok {
		return nil, ErrAPINotFound
	}
	return resolver(typename)
}

func ApiVersion(version string) string {
	parts := strings.Split(version, "/")
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return path.Join(parts[len(parts)-2], parts[len(parts)-1])
}

func semverGreater(s1, s2 string) bool {
	s1Ver, err := semver.ParseTolerant(s1)
	if err != nil {
		// semver should be validated before being passed here
		return false
	}
	s2Ver, err := semver.ParseTolerant(s2)
	if err != nil {
		// semver should be validated before being passed here
		return false
	}
	return s1Ver.GT(s2Ver)
}
