package structmap

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type Foo struct {
	Bar string
	Foo string
}

type Bar struct {
	Foo  `map:",inline"`
	Anws string `map:"anws"`
}

type Bar1 struct {
	Bar
	Filtered string `map:"-"`
	Anws     string `map:"shoop"`
	Bar2     *Bar   `map:",inline"`
}

type UnmarshalStruct struct {
	L1       string
	Extended map[string]interface{} `map:",inline"`
}

func TestMarshallingSimpleStruct(t *testing.T) {
	var (
		obj interface{}
		exp string
		m   map[string]interface{}
		err error
	)

	obj = Foo{
		Foo: "foo value",
		Bar: "bar value",
	}
	exp = `{
	"Bar": "bar value",
	"Foo": "foo value"
}`
	if m, err = Marshal(obj); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}

	if m, err = Marshal(&obj); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}
}

func TestMarshallingNonStructValue(t *testing.T) {
	var (
		obj interface{}
		err error
	)

	obj = "I'm not a struct"

	if _, err = Marshal(obj); err != ErrNonStruct {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(err)

	if _, err = Marshal(&obj); err != ErrNonStruct {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(err)
}

func TestMarshallingTypicalStruct(t *testing.T) {
	var (
		obj interface{}
		exp string
		m   map[string]interface{}
		err error
	)
	foo := Foo{
		Bar: "bar value in Foo object",
		Foo: "foo value in Foo object",
	}
	obj = Bar{
		Foo:  foo,
		Anws: "42",
	}
	exp = `{
	"Bar": "bar value in Foo object",
	"Foo": "foo value in Foo object",
	"anws": "42"
}`
	if m, err = Marshal(obj); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}

	if m, err = Marshal(&obj); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}
}

func TestMarshallingComplicatedStructWithPointers(t *testing.T) {
	var (
		obj interface{}
		exp string
		m   map[string]interface{}
		err error
	)
	foo := Foo{
		Bar: "what's the anwser?",
		Foo: "foo",
	}
	obj = Bar1{
		Bar: Bar{
			Foo:  foo,
			Anws: "42",
		},
		Filtered: "I should not exist in ouput",
		Anws:     "bar",
		Bar2: &Bar{
			Foo:  foo,
			Anws: "99",
		},
	}
	exp = `{
	"Bar": "what's the anwser?",
	"Foo": "foo",
	"anws": "99",
	"shoop": "bar"
}`

	if m, err = Marshal(obj); err != nil {
		t.Fail()
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}

	if m, err = Marshal(&obj); err != nil {
		t.Fail()
	}
	if !asserMarshal(m, exp) {
		t.Fail()
	}
}

func TestUnmarshallingToSimpleStruct(t *testing.T) {

	var (
		m   map[string]interface{}
		obj interface{}
		exp string
		err error
	)

	m = map[string]interface{}{
		"Foo": "42",
		"Bar": "0",
	}
	exp = `{
	"Bar": "0",
	"Foo": "42"
}`

	obj = &Foo{}

	if err = Unmarshal(m, obj); err != nil {
		t.Fail()
	}
	if !asserMarshal(obj, exp) {
		t.Fail()
	}
}

func TestUnmarshallingToComplicatedStruct(t *testing.T) {
	var (
		m   map[string]interface{}
		obj interface{}
		exp string
		err error
	)

	m = map[string]interface{}{
		"L1":         "level 1 value",
		"additional": "some additional value",
	}
	exp = `{
	"L1": "level 1 value",
	"Extended": {
		"additional": "some additional value"
	}
}`

	obj = &UnmarshalStruct{}

	if err = Unmarshal(m, obj); err != nil {
		t.Fail()
	}
	log.Println(obj)
	if !asserMarshal(obj, exp) {
		t.Fail()
	}

}

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
