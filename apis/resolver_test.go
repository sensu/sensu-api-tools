package apis_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/sensu/sensu-api-tools/apis"
)

func TestResolveType(t *testing.T) {
	type TestAPITypeA struct{}
	type TestAPITypeB struct{}
	resolver := func(s string) (interface{}, error) {
		switch s {
		case "TestAPITypeA":
			return new(TestAPITypeA), nil
		case "TestAPITypeB":
			return new(TestAPITypeB), nil
		}
		return nil, apis.ErrAPINotFound
	}
	apis.RegisterResolver("apis_test/v1", resolver)
	testCases := []struct {
		ApiVersion string
		Type       string
		ExpRet     interface{}
		ExpErr     bool
	}{
		{
			ApiVersion: "apis_test/v1",
			Type:       "TestAPITypeA",
			ExpRet:     &TestAPITypeA{},
			ExpErr:     false,
		},
		{
			ApiVersion: "apis_test/v1",
			Type:       "TestAPITypeB",
			ExpRet:     &TestAPITypeB{},
			ExpErr:     false,
		},
		{
			ApiVersion: "apis_test/v1.0.0",
			Type:       "TestAPITypeB",
			ExpRet:     &TestAPITypeB{},
			ExpErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.ApiVersion, tc.Type), func(t *testing.T) {
			r, err := apis.Resolve(tc.ApiVersion, tc.Type)
			if !reflect.DeepEqual(r, tc.ExpRet) {
				t.Fatalf("unexpected type: got %T, want %T", r, tc.ExpRet)
			}
			if err != nil && !tc.ExpErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil && tc.ExpErr {
				t.Fatal("expected an error")
			}
		})
	}
}
