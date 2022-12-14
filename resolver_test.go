package apitools_test

import (
	"fmt"
	"reflect"
	"testing"

	apitools "github.com/sensu/sensu-api-tools"
)

func TestResolveType(t *testing.T) {
	type TestAPITypeA struct{}
	type TestAPITypeB struct{ T string }
	apitools.RegisterType("apis_test/v1", new(TestAPITypeA))
	apitools.RegisterType("apis_test/v1", new(TestAPITypeB))
	apitools.RegisterType("apis_test/v2", new(TestAPITypeA), apitools.WithAlias("test_api_type_a", "kazoo"))
	apitools.RegisterType("apis_test/v2", new(TestAPITypeB), apitools.WithResolveHook(func(v interface{}) {
		if b, ok := v.(*TestAPITypeB); ok {
			b.T = "flute"
		}
	}))
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
		{
			ApiVersion: "apis_test/v1.9.8",
			Type:       "TestAPITypeB",
			ExpRet:     &TestAPITypeB{},
			ExpErr:     false,
		},
		{
			ApiVersion: "apis_test/v2",
			Type:       "test_api_type_a",
			ExpRet:     &TestAPITypeA{},
			ExpErr:     false,
		},
		{
			ApiVersion: "apis_test/v2.0.0",
			Type:       "kazoo",
			ExpRet:     &TestAPITypeA{},
			ExpErr:     false,
		},
		{
			ApiVersion: "apis_test/v2.0.0",
			Type:       "TestAPITypeB",
			ExpRet:     &TestAPITypeB{T: "flute"},
			ExpErr:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s/%s", tc.ApiVersion, tc.Type), func(t *testing.T) {
			r, err := apitools.Resolve(tc.ApiVersion, tc.Type)
			if !reflect.DeepEqual(r, tc.ExpRet) {
				t.Fatalf("unexpected type: got %+v, want %+v", r, tc.ExpRet)
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
