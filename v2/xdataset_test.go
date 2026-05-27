package xcore

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newDS(pairs ...interface{}) *XDataset {
	ds := &XDataset{}
	for i := 0; i < len(pairs)-1; i += 2 {
		ds.Set(pairs[i].(string), pairs[i+1])
	}
	return ds
}

// ---------------------------------------------------------------------------
// NewXDataset
// ---------------------------------------------------------------------------

func TestNewXDataset_Empty(t *testing.T) {
	ds := NewXDataset(map[string]interface{}{})
	if ds == nil {
		t.Fatal("expected non-nil XDatasetDef")
	}
}

func TestNewXDataset_Scalar(t *testing.T) {
	ds := NewXDataset(map[string]interface{}{"x": 42})
	v, ok := ds.GetInt("x")
	if !ok || v != 42 {
		t.Fatalf("want 42/true, got %d/%v", v, ok)
	}
}

func TestNewXDataset_NestedMap(t *testing.T) {
	ds := NewXDataset(map[string]interface{}{
		"inner": map[string]interface{}{"val": "hello"},
	})
	v, ok := ds.GetString("inner>val")
	if !ok || v != "hello" {
		t.Fatalf("want hello/true, got %q/%v", v, ok)
	}
}

func TestNewXDataset_SliceOfMaps(t *testing.T) {
	ds := NewXDataset(map[string]interface{}{
		"items": []map[string]interface{}{
			{"name": "a"},
			{"name": "b"},
		},
	})
	v, ok := ds.GetString("items>0>name")
	if !ok || v != "a" {
		t.Fatalf("want a/true, got %q/%v", v, ok)
	}
}

func TestNewXDataset_SliceOfInterfaces(t *testing.T) {
	ds := NewXDataset(map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"n": 1},
			map[string]interface{}{"n": 2},
		},
	})
	v, ok := ds.GetInt("items>1>n")
	if !ok || v != 2 {
		t.Fatalf("want 2/true, got %d/%v", v, ok)
	}
}

// ---------------------------------------------------------------------------
// Set / Get
// ---------------------------------------------------------------------------

func TestXDataset_SetGet_String(t *testing.T) {
	ds := &XDataset{}
	ds.Set("k", "hello")
	v, ok := ds.Get("k")
	if !ok || v != "hello" {
		t.Fatalf("want hello, got %v/%v", v, ok)
	}
}

func TestXDataset_GetMissing(t *testing.T) {
	ds := &XDataset{}
	_, ok := ds.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestXDataset_GetOverwrite(t *testing.T) {
	ds := &XDataset{}
	ds.Set("k", "first")
	ds.Set("k", "second")
	v, _ := ds.GetString("k")
	if v != "second" {
		t.Fatalf("expected second, got %q", v)
	}
}

// ---------------------------------------------------------------------------
// Get — nested path (>-separated)
// ---------------------------------------------------------------------------

func TestXDataset_Get_NestedDataset(t *testing.T) {
	inner := newDS("val", "deep")
	outer := &XDataset{}
	outer.Set("inner", inner)
	v, ok := outer.Get("inner>val")
	if !ok || v != "deep" {
		t.Fatalf("want deep/true, got %v/%v", v, ok)
	}
}

func TestXDataset_Get_NestedCollection(t *testing.T) {
	coll := &XDatasetCollection{}
	coll.Push(newDS("x", 10))
	coll.Push(newDS("x", 20))
	outer := &XDataset{}
	outer.Set("coll", coll)
	v, ok := outer.Get("coll>1>x")
	if !ok || v.(int) != 20 {
		t.Fatalf("want 20/true, got %v/%v", v, ok)
	}
}

func TestXDataset_Get_NestedCollectionDirectElement(t *testing.T) {
	coll := &XDatasetCollection{}
	ds0 := newDS("a", 1)
	coll.Push(ds0)
	outer := &XDataset{}
	outer.Set("coll", coll)
	// path coll>0 should return the dataset itself
	v, ok := outer.Get("coll>0")
	if !ok || v == nil {
		t.Fatalf("expected element, got %v/%v", v, ok)
	}
}

func TestXDataset_Get_NestedPath_ScalarIntermediate(t *testing.T) {
	// 4.1 regression: intermediate key holds a scalar, should return not-found
	// NOT fall through to a flat lookup
	ds := &XDataset{}
	ds.Set("a", "scalar")
	_, ok := ds.Get("a>b")
	if ok {
		t.Fatal("expected false when intermediate value is not traversable (4.1 regression)")
	}
}

func TestXDataset_Get_NestedPath_IntermediateMissing(t *testing.T) {
	ds := &XDataset{}
	_, ok := ds.Get("missing>sub")
	if ok {
		t.Fatal("expected false for missing intermediate key")
	}
}

func TestXDataset_Get_NestedCollection_BadIndex(t *testing.T) {
	coll := &XDatasetCollection{}
	coll.Push(newDS("x", 1))
	outer := &XDataset{}
	outer.Set("coll", coll)
	_, ok := outer.Get("coll>notanint>x")
	if ok {
		t.Fatal("expected false for non-integer collection index")
	}
}

func TestXDataset_Get_NestedCollection_OutOfBounds(t *testing.T) {
	coll := &XDatasetCollection{}
	coll.Push(newDS("x", 1))
	outer := &XDataset{}
	outer.Set("coll", coll)
	_, ok := outer.Get("coll>99>x")
	if ok {
		t.Fatal("expected false for out-of-bounds collection index")
	}
}

// ---------------------------------------------------------------------------
// GetDataset / GetCollection
// ---------------------------------------------------------------------------

func TestXDataset_GetDataset_OK(t *testing.T) {
	inner := newDS("v", 1)
	outer := &XDataset{}
	outer.Set("inner", inner)
	got, ok := outer.GetDataset("inner")
	if !ok || got == nil {
		t.Fatal("expected dataset")
	}
}

func TestXDataset_GetDataset_WrongType(t *testing.T) {
	ds := newDS("x", 42)
	_, ok := ds.GetDataset("x")
	if ok {
		t.Fatal("expected false for non-dataset value")
	}
}

func TestXDataset_GetCollection_OK(t *testing.T) {
	coll := &XDatasetCollection{}
	coll.Push(newDS("a", 1))
	outer := &XDataset{}
	outer.Set("coll", coll)
	got, ok := outer.GetCollection("coll")
	if !ok || got == nil {
		t.Fatal("expected collection")
	}
}

func TestXDataset_GetCollection_WrongType(t *testing.T) {
	ds := newDS("x", "hello")
	_, ok := ds.GetCollection("x")
	if ok {
		t.Fatal("expected false for non-collection value")
	}
}

// ---------------------------------------------------------------------------
// GetString
// ---------------------------------------------------------------------------

func TestXDataset_GetString_String(t *testing.T) {
	ds := newDS("s", "hello")
	v, ok := ds.GetString("s")
	if !ok || v != "hello" {
		t.Fatalf("want hello/true, got %q/%v", v, ok)
	}
}

func TestXDataset_GetString_Int(t *testing.T) {
	ds := newDS("n", 42)
	v, ok := ds.GetString("n")
	if !ok || v != "42" {
		t.Fatalf("want 42/true, got %q/%v", v, ok)
	}
}

func TestXDataset_GetString_Nil(t *testing.T) {
	ds := newDS("n", nil)
	v, ok := ds.GetString("n")
	if !ok || v != "" {
		t.Fatalf("want empty/true for nil, got %q/%v", v, ok)
	}
}

func TestXDataset_GetString_Missing(t *testing.T) {
	ds := &XDataset{}
	_, ok := ds.GetString("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

// ---------------------------------------------------------------------------
// GetBool
// ---------------------------------------------------------------------------

func TestXDataset_GetBool_Bool(t *testing.T) {
	ds := newDS("b", true)
	v, ok := ds.GetBool("b")
	if !ok || !v {
		t.Fatal("expected true")
	}
}

func TestXDataset_GetBool_IntZero(t *testing.T) {
	ds := newDS("n", 0)
	v, ok := ds.GetBool("n")
	if !ok || v {
		t.Fatal("expected false for 0")
	}
}

func TestXDataset_GetBool_IntNonZero(t *testing.T) {
	ds := newDS("n", 5)
	v, ok := ds.GetBool("n")
	if !ok || !v {
		t.Fatal("expected true for 5")
	}
}

func TestXDataset_GetBool_Float(t *testing.T) {
	ds := newDS("f", float64(0.0))
	v, ok := ds.GetBool("f")
	if !ok || v {
		t.Fatal("expected false for 0.0")
	}
}

func TestXDataset_GetBool_Time_Zero(t *testing.T) {
	ds := newDS("t", time.Time{})
	v, ok := ds.GetBool("t")
	if !ok || v {
		t.Fatal("expected false for zero time")
	}
}

func TestXDataset_GetBool_Time_NonZero(t *testing.T) {
	ds := newDS("t", time.Now())
	v, ok := ds.GetBool("t")
	if !ok || !v {
		t.Fatal("expected true for non-zero time")
	}
}

func TestXDataset_GetBool_NonNilInterface(t *testing.T) {
	ds := newDS("x", "something")
	v, ok := ds.GetBool("x")
	if !ok || !v {
		t.Fatal("expected true for non-nil interface")
	}
}

func TestXDataset_GetBool_Missing(t *testing.T) {
	ds := &XDataset{}
	_, ok := ds.GetBool("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

// ---------------------------------------------------------------------------
// GetInt
// ---------------------------------------------------------------------------

func TestXDataset_GetInt_Int(t *testing.T) {
	ds := newDS("n", 7)
	v, ok := ds.GetInt("n")
	if !ok || v != 7 {
		t.Fatalf("want 7, got %d/%v", v, ok)
	}
}

func TestXDataset_GetInt_BoolTrue(t *testing.T) {
	ds := newDS("b", true)
	v, ok := ds.GetInt("b")
	if !ok || v != 1 {
		t.Fatalf("want 1, got %d/%v", v, ok)
	}
}

func TestXDataset_GetInt_BoolFalse(t *testing.T) {
	ds := newDS("b", false)
	v, ok := ds.GetInt("b")
	if !ok || v != 0 {
		t.Fatalf("want 0, got %d/%v", v, ok)
	}
}

func TestXDataset_GetInt_Float(t *testing.T) {
	ds := newDS("f", float64(3.9))
	v, ok := ds.GetInt("f")
	if !ok || v != 3 {
		t.Fatalf("want 3, got %d/%v", v, ok)
	}
}

func TestXDataset_GetInt_TimeZero(t *testing.T) {
	ds := newDS("t", time.Time{})
	v, ok := ds.GetInt("t")
	if !ok || v != 0 {
		t.Fatalf("want 0, got %d/%v", v, ok)
	}
}

func TestXDataset_GetInt_TimeNonZero(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	ds := newDS("t", ts)
	v, ok := ds.GetInt("t")
	if !ok || v != int(ts.Unix()) {
		t.Fatalf("want %d, got %d/%v", int(ts.Unix()), v, ok)
	}
}

func TestXDataset_GetInt_IntVariants(t *testing.T) {
	tests := []struct {
		val  interface{}
		want int
	}{
		{int8(10), 10},
		{int16(200), 200},
		{int32(30000), 30000},
		{int64(1 << 40), int(int64(1 << 40))},
		{uint(5), 5},
		{uint8(255), 255},
		{uint16(1000), 1000},
		{uint32(99999), 99999},
		{uint64(123), 123},
		{float32(2.5), 2},
	}
	for _, tt := range tests {
		ds := newDS("v", tt.val)
		got, ok := ds.GetInt("v")
		if !ok || got != tt.want {
			t.Errorf("GetInt(%T(%v)): want %d, got %d/%v", tt.val, tt.val, tt.want, got, ok)
		}
	}
}

func TestXDataset_GetInt_Missing(t *testing.T) {
	ds := &XDataset{}
	_, ok := ds.GetInt("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

// ---------------------------------------------------------------------------
// GetFloat
// ---------------------------------------------------------------------------

func TestXDataset_GetFloat_Float64(t *testing.T) {
	ds := newDS("f", float64(3.14))
	v, ok := ds.GetFloat("f")
	if !ok || v != 3.14 {
		t.Fatalf("want 3.14, got %f/%v", v, ok)
	}
}

func TestXDataset_GetFloat_BoolTrue(t *testing.T) {
	ds := newDS("b", true)
	v, ok := ds.GetFloat("b")
	if !ok || v != 1.0 {
		t.Fatalf("want 1.0, got %f/%v", v, ok)
	}
}

func TestXDataset_GetFloat_Missing(t *testing.T) {
	ds := &XDataset{}
	_, ok := ds.GetFloat("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

// ---------------------------------------------------------------------------
// GetTime
// ---------------------------------------------------------------------------

func TestXDataset_GetTime_OK(t *testing.T) {
	now := time.Now()
	ds := newDS("t", now)
	v, ok := ds.GetTime("t")
	if !ok || !v.Equal(now) {
		t.Fatalf("want %v, got %v/%v", now, v, ok)
	}
}

func TestXDataset_GetTime_WrongType(t *testing.T) {
	ds := newDS("t", "not-a-time")
	_, ok := ds.GetTime("t")
	if ok {
		t.Fatal("expected false for non-time value")
	}
}

// ---------------------------------------------------------------------------
// Collection getters
// ---------------------------------------------------------------------------

func TestXDataset_GetStringCollection_OK(t *testing.T) {
	ds := newDS("sc", []string{"a", "b"})
	v, ok := ds.GetStringCollection("sc")
	if !ok || len(v) != 2 || v[0] != "a" {
		t.Fatalf("unexpected: %v/%v", v, ok)
	}
}

func TestXDataset_GetBoolCollection_OK(t *testing.T) {
	ds := newDS("bc", []bool{true, false})
	v, ok := ds.GetBoolCollection("bc")
	if !ok || len(v) != 2 || !v[0] {
		t.Fatalf("unexpected: %v/%v", v, ok)
	}
}

func TestXDataset_GetIntCollection_OK(t *testing.T) {
	ds := newDS("ic", []int{1, 2, 3})
	v, ok := ds.GetIntCollection("ic")
	if !ok || v[2] != 3 {
		t.Fatalf("unexpected: %v/%v", v, ok)
	}
}

func TestXDataset_GetFloatCollection_OK(t *testing.T) {
	ds := newDS("fc", []float64{1.1, 2.2})
	v, ok := ds.GetFloatCollection("fc")
	if !ok || v[1] != 2.2 {
		t.Fatalf("unexpected: %v/%v", v, ok)
	}
}

func TestXDataset_GetTimeCollection_OK(t *testing.T) {
	now := time.Now()
	ds := newDS("tc", []time.Time{now})
	v, ok := ds.GetTimeCollection("tc")
	if !ok || !v[0].Equal(now) {
		t.Fatalf("unexpected: %v/%v", v, ok)
	}
}

// ---------------------------------------------------------------------------
// Del
// ---------------------------------------------------------------------------

func TestXDataset_Del(t *testing.T) {
	ds := newDS("k", "v")
	ds.Del("k")
	_, ok := ds.Get("k")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestXDataset_Del_Missing(t *testing.T) {
	// should not panic
	ds := &XDataset{}
	ds.Del("nonexistent")
}

// ---------------------------------------------------------------------------
// Clone
// ---------------------------------------------------------------------------

func TestXDataset_Clone_Basic(t *testing.T) {
	ds := newDS("a", 1, "b", "hello")
	cl := ds.Clone()
	v, ok := cl.GetInt("a")
	if !ok || v != 1 {
		t.Fatal("expected a=1 in clone")
	}
	// Mutate original, clone must be independent
	ds.Set("a", 99)
	v2, _ := cl.GetInt("a")
	if v2 != 1 {
		t.Fatal("clone was not independent from original")
	}
}

func TestXDataset_Clone_NestedDataset(t *testing.T) {
	inner := newDS("x", 10)
	outer := &XDataset{}
	outer.Set("inner", inner)
	cloned := outer.Clone()
	inner.Set("x", 999) // mutate original inner
	v, ok := cloned.GetInt("inner>x")
	if !ok || v != 10 {
		t.Fatalf("cloned nested dataset should be independent, got %d/%v", v, ok)
	}
}

func TestXDataset_Clone_NestedCollection(t *testing.T) {
	coll := &XDatasetCollection{}
	coll.Push(newDS("n", 5))
	outer := &XDataset{}
	outer.Set("coll", coll)
	cloned := outer.Clone()
	// Mutate original collection
	(*coll)[0].Set("n", 999)
	v, ok := cloned.GetInt("coll>0>n")
	if !ok || v != 5 {
		t.Fatalf("cloned nested collection should be independent, got %d/%v", v, ok)
	}
}

// ---------------------------------------------------------------------------
// String / GoString
// ---------------------------------------------------------------------------

func TestXDataset_String(t *testing.T) {
	ds := newDS("a", 1)
	s := ds.String()
	if s == "" {
		t.Fatal("String() returned empty")
	}
}

func TestXDataset_GoString(t *testing.T) {
	ds := newDS("a", 1)
	s := ds.GoString()
	if s == "" {
		t.Fatal("GoString() returned empty")
	}
}
