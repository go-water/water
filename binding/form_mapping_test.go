package binding

import (
	"net/http"
	"net/url"
	"testing"
)

func TestMapFormByTag_EmptyValuesOverridePreset(t *testing.T) {
	type payload struct {
		Name  string `query:"name"`
		Count int    `query:"count"`
		Flag  bool   `query:"flag"`
	}

	p := payload{Name: "alice", Count: 5, Flag: true}
	err := mapFormByTag(&p, map[string][]string{
		"name":  {""},
		"count": {""},
		"flag":  {""},
	}, "query")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "" || p.Count != 0 || p.Flag {
		t.Fatalf("got %+v, want empty values to override preset", p)
	}
}

func TestMapFormByTag_NonEmptyValuesOverridePreset(t *testing.T) {
	type payload struct {
		Name  string `query:"name"`
		Count int    `query:"count"`
	}

	p := payload{Name: "alice", Count: 5}
	err := mapFormByTag(&p, map[string][]string{
		"name":  {"bob"},
		"count": {"10"},
	}, "query")
	if err != nil {
		t.Fatal(err)
	}

	if p.Name != "bob" || p.Count != 10 {
		t.Fatalf("got %+v, want non-empty override", p)
	}
}

func TestQueryBinding_EmptyValuesOverridePreset(t *testing.T) {
	type payload struct {
		N int `query:"n"`
	}

	p := payload{N: 7}
	req := &http.Request{URL: &url.URL{RawQuery: "n="}}
	if err := Query.Bind(req, &p); err != nil {
		t.Fatal(err)
	}

	if p.N != 0 {
		t.Fatalf("N=%d, want 0", p.N)
	}
}

func TestSetFormMap_EmptyLastValueOverrides(t *testing.T) {
	m := map[string]string{"a": "keep", "b": "old"}
	form := map[string][]string{
		"a": {"new"},
		"b": {"x", ""},
	}

	if err := setFormMap(m, form); err != nil {
		t.Fatal(err)
	}

	if m["a"] != "new" || m["b"] != "" {
		t.Fatalf("got %#v", m)
	}
}

func TestSetFormMap_ZeroLengthValuesOverrideEmptyString(t *testing.T) {
	m := map[string]string{"k": "v"}
	form := map[string][]string{"k": {}}

	if err := setFormMap(m, form); err != nil {
		t.Fatal(err)
	}

	if m["k"] != "" {
		t.Fatalf("k=%q, want empty string", m["k"])
	}
}
