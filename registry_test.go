package huma

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Output[T any] struct{}

type Embedded[P any] struct{}

type EmbeddedTwo[P, V any] struct{}

type S struct{}

type ü struct{}

type MP4 struct{}

func TestDefaultSchemaNamer(t *testing.T) {
	type Renamed Output[*[]Embedded[time.Time]]

	for _, example := range []struct {
		typ  any
		name string
		hint string
	}{
		{int(0), "Int", ""},
		{int64(0), "Int64", ""},
		{S{}, "S", ""},
		{time.Time{}, "Time", ""},
		{Output[int]{}, "OutputInt", ""},
		{Output[*int]{}, "OutputInt", ""},
		{Output[[]int]{}, "OutputListInt", ""},
		{Output[[]*int]{}, "OutputListInt", ""},
		{Output[[][]int]{}, "OutputListListInt", ""},
		{Output[map[string]int]{}, "OutputMapStringInt", ""},
		{Output[map[string][]*int]{}, "OutputMapStringListInt", ""},
		{Output[S]{}, "OutputS", ""},
		{Output[ü]{}, "OutputÜ", ""},
		{Output[MP4]{}, "OutputMP4", ""},
		{Output[Embedded[*time.Time]]{}, "OutputEmbeddedTime", ""},
		{Output[*[]Embedded[time.Time]]{}, "OutputListEmbeddedTime", ""},
		{Output[EmbeddedTwo[[]time.Time, **url.URL]]{}, "OutputEmbeddedTwoListTimeURL", ""},
		{Renamed{}, "Renamed", ""},
		{struct{}{}, "SomeGenericThing", "Some[pkg.Generic]Thing"},
		{struct{}{}, "Type1Type2Type3", "pkg1.Type1[path/to/pkg2.Type2]pkg3.Type3"},
	} {
		t.Run(example.name, func(t *testing.T) {
			name := DefaultSchemaNamer(reflect.TypeOf(example.typ), example.hint)
			assert.Equal(t, example.name, name)
		})
	}
}

func TestSchemaAlias(t *testing.T) {
	type StringContainer struct {
		Value string
	}
	type StructWithStringContainer struct {
		Name StringContainer `json:"name,omitempty"`
	}
	type StructWithString struct {
		Name string `json:"name,omitempty"`
	}
	registry := NewMapRegistry("#/components/schemas", DefaultSchemaNamer)
	registry.RegisterTypeAlias(reflect.TypeOf(StringContainer{}), reflect.TypeOf(""))
	schemaWithContainer := registry.Schema(reflect.TypeOf(StructWithStringContainer{}), false, "")
	schemaWithString := registry.Schema(reflect.TypeOf(StructWithString{}), false, "")
	assert.Equal(t, schemaWithString, schemaWithContainer)
}

func TestRefForNamedPrimitives(t *testing.T) {
	type CustomHeader string

	type FirstRequest struct {
		Header CustomHeader `json:"header" description:"A custom header"`
	}

	type SecondRequest struct {
		AnotherHeader CustomHeader `json:"another_header" description:"Another custom header"`
	}

	// Create registry with reuse for named primitives disabled (default)
	registry := NewMapRegistry("#/components/schemas", DefaultSchemaNamer)

	// Get schema for both requests
	first := SchemaFromType(registry, reflect.TypeOf(FirstRequest{}))
	second := SchemaFromType(registry, reflect.TypeOf(SecondRequest{}))

	// Verify that the schemas contain the CustomHeader type inline
	if first.Properties["header"].Ref != "" {
		t.Errorf("Expected header to be defined inline, but got a ref: %s", first.Properties["header"].Ref)
	}
	if second.Properties["another_header"].Ref != "" {
		t.Errorf("Expected another_header to be defined inline, but got a ref: %s", second.Properties["another_header"].Ref)
	}

	// Create registry with reuse for named primitives enabled
	registry = NewMapRegistry("#/components/schemas", DefaultSchemaNamer)
	registry.(interface {
		EnableRefForNamedPrimitives(bool)
	}).EnableRefForNamedPrimitives(true)

	// Get schema for both requests
	first = SchemaFromType(registry, reflect.TypeOf(FirstRequest{}))
	second = SchemaFromType(registry, reflect.TypeOf(SecondRequest{}))

	// Verify that the schemas reference the CustomHeader type
	if first.Properties["header"].Ref == "" {
		t.Errorf("Expected header to use a ref, but it's defined inline")
	}
	if second.Properties["another_header"].Ref == "" {
		t.Errorf("Expected another_header to use a ref, but it's defined inline")
	}

	// Verify that both references point to the same schema
	if first.Properties["header"].Ref != second.Properties["another_header"].Ref {
		t.Errorf("Expected both properties to use the same ref, but got %s and %s",
			first.Properties["header"].Ref, second.Properties["another_header"].Ref)
	}
}
