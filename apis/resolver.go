package apis

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/blang/semver/v4"
)

var errAPINotFound = errors.New("api not found")

// PackageMap contains a list of packages with their Resolver func
var typeMap = map[string]map[string]reflect.Type{}

var typeMapMu = &sync.RWMutex{}

// RegisterType allows modules to register API types to be resolved.
// By default a type registered can be resolved by its type's name. Optionally
// aliases can be provided that will also resolve the type.
//
// example:
// after calling RegisterType("core/v2", new(corev2.Asset), "asset"),
// calling Resolve("core/v2", "Asset") or Resolve("core/v2", "asset") should
// return a corev2.Asset.
func RegisterType(apiGroup string, t any, typeAlias ...string) {
	rrt := reflect.ValueOf(t)
	rType := reflect.Indirect(rrt).Type()
	typeMapMu.Lock()
	defer typeMapMu.Unlock()
	if _, ok := typeMap[apiGroup]; !ok {
		typeMap[apiGroup] = map[string]reflect.Type{}
	}
	typeMap[apiGroup][rType.Name()] = rType
	for _, alias := range typeAlias {
		typeMap[apiGroup][alias] = rType
	}
}

// Resolve resolves the raw type for the requested api version and type.
func Resolve(apiVersion string, typename string) (interface{}, error) {

	// Guard read access to packageMap
	typeMapMu.RLock()
	defer typeMapMu.RUnlock()
	apiGroup, reqVer := parseAPIVersion(apiVersion)

	group, ok := typeMap[apiGroup]
	if !ok {
		return nil, fmt.Errorf("api group %s has not been registered", apiGroup)
	}
	typ, ok := group[typename]
	if !ok {
		return nil, errAPINotFound
	}

	if foundVer, err := versionOf(typ); err == nil {
		if semverGreater(reqVer, foundVer) {
			return nil, fmt.Errorf("requested version was %s, but only %s is available", reqVer, foundVer)
		}
	}

	return reflect.New(typ).Interface(), nil
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
