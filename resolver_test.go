package apitools_test

import (
	"fmt"
	"reflect"
	"testing"

	apitools "github.com/sensu/sensu-api-tools"
)

type TestAPITypeA struct{}

func (*TestAPITypeA) Foo() {}

type TestAPITypeB struct{ T string }

func (*TestAPITypeB) Foo() {}

type TestAPITypeC struct{}

func init() {
	apitools.RegisterType("apis_test/v1", new(TestAPITypeA))
	apitools.RegisterType("apis_test/v1", new(TestAPITypeB))
	apitools.RegisterType("apis_test/v2", new(TestAPITypeA), apitools.WithAlias("test_api_type_a", "kazoo"))
	apitools.RegisterType("apis_test/v2", new(TestAPITypeB), apitools.WithResolveHook(func(v interface{}) {
		if b, ok := v.(*TestAPITypeB); ok {
			b.T = "flute"
		}
	}))
	apitools.RegisterType("apis_test/v1", new(TestAPITypeC))
}

func TestResolveType(t *testing.T) {
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

func TestIterTypes(t *testing.T) {
	i := 0
	apitools.IterTypes(func(apiGroup, name string, _ any) bool {
		i++
		if _, err := apitools.Resolve(apiGroup, name); err != nil {
			t.Error(err)
		}
		return true
	})
	// expect 7 types, since we alias rbac names
	if got, want := i, 7; got != want {
		t.Errorf("didn't iterate all types: got %d, want %d", got, want)
	}
}

func TestFindTypesOf(t *testing.T) {
	type Fooer interface {
		Foo()
	}

	result := apitools.FindTypesOf[Fooer]()
	if got, want := len(result), 4; got != want {
		t.Errorf("wrong number of results: got %d, want %d", got, want)
	}
}
