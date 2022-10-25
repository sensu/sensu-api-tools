package apitools_test

import (
	"testing"

	"runtime/debug"

	apitools "github.com/sensu/sensu-api-tools"
)

// TODO(ck): Decide what to do about this test.
// This is difficult behavior to spec, not only due to https://github.com/golang/go/issues/33976
// but because this module intentionally excludes dependencies on the core modules.
func TestAPIModuleVersions(t *testing.T) {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok && buildInfo.Deps != nil {
		t.Fatal("remove this if block, the test works now")
	} else {
		t.Skip()
	}
	modVersions := apitools.APIModuleVersions()
	if _, ok := modVersions["core/v2"]; !ok {
		t.Errorf("missing core/v2 module version")
	}
	if _, ok := modVersions["core/v3"]; !ok {
		t.Errorf("missing core/v3 module version")
	}
}
