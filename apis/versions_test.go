package apis_test

import (
	"testing"

	"runtime/debug"

	"github.com/sensu/sensu-api-tools/apis"
)

// TODO(eric): This test doesn't work yet because of https://github.com/golang/go/issues/33976
func TestAPIModuleVersions(t *testing.T) {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok && buildInfo.Deps != nil {
		t.Fatal("remove this if block, the test works now")
	} else {
		t.Skip()
	}
	modVersions := apis.APIModuleVersions()
	if _, ok := modVersions["core/v2"]; !ok {
		t.Errorf("missing core/v2 module version")
	}
	if _, ok := modVersions["core/v3"]; !ok {
		t.Errorf("missing core/v3 module version")
	}
}
