package structmap

import (
	"encoding/json"
	"fmt"
	"testing"
)

func jsonStr(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "	")
	return string(b)
}

func pJson(v interface{}) {
	fmt.Println(jsonStr(v))
}

func asserMarshal(v interface{}, s string) bool {
	pJson(v)
	return jsonStr(v) == s
}

func TestMarshal(t *testing.T) {

	var (
		obj interface{}
		exp string
		m   map[string]interface{}
		err error
	)

	type Foo struct {
		Bar string
		Foo string
	}

	obj = Foo{"bar", "foo"}
	exp = `{
	"Bar": "bar",
	"Foo": "foo"
}`

	if m, err = Marshal(obj); err != nil {
		t.Error(err)
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}

	obj = "string"
	if m, err = Marshal(obj); err != ErrNonStruct {
		t.Fail()
	}
	fmt.Println(err)

	type Bar struct {
		Foo  `json:",inline"`
		Anws string
	}

	obj = Bar{Foo: Foo{"what's the anwser?", "foo"}, Anws: "42"}
	exp = `{
	"Anws": "42",
	"Bar": "what's the anwser?",
	"Foo": "foo"
}`
	if m, err = Marshal(obj); err != nil {
		t.Error(err)
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}

	type Bar1 struct {
		Bar
		Filtered string `structmap:"-"`
		Anws     string `structmap:"shoop,inline"`
		Bar2     *Bar   `json:",inline"`
	}

	obj = Bar1{
		Bar: Bar{
			Foo:  Foo{"what's the anwser?", "foo"},
			Anws: "42",
		},
		Filtered: "blah",
		Anws:     "bar",
		Bar2: &Bar{
			Foo:  Foo{"what's the anwser?", "foo"},
			Anws: "99",
		},
	}
	exp = `{
	"Anws": "99",
	"Bar": "what's the anwser?",
	"Foo": "foo",
	"shoop": "bar"
}`
	_ = "breakpoint"
	if m, err = Marshal(&obj); err != nil {
		t.Error(err)
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}

}
