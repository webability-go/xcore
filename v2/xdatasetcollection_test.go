package xcore

import (
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// NewXDatasetCollection
// ---------------------------------------------------------------------------

func TestNewXDatasetCollection_Empty(t *testing.T) {
	dsc := NewXDatasetCollection([]map[string]interface{}{})
	if dsc.Count() != 0 {
		t.Fatalf("expected 0, got %d", dsc.Count())
	}
}

func TestNewXDatasetCollection_WithData(t *testing.T) {
	dsc := NewXDatasetCollection([]map[string]interface{}{
		{"a": 1},
		{"a": 2},
	})
	if dsc.Count() != 2 {
		t.Fatalf("expected 2, got %d", dsc.Count())
	}
}

// ---------------------------------------------------------------------------
// Push / Pop
// ---------------------------------------------------------------------------

func TestXDatasetCollection_PushPop(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("n", 1))
	c.Push(newDS("n", 2))
	if c.Count() != 2 {
		t.Fatalf("expected 2, got %d", c.Count())
	}
	last := c.Pop()
	v, _ := last.GetInt("n")
	if v != 2 {
		t.Fatalf("want 2, got %d", v)
	}
	if c.Count() != 1 {
		t.Fatalf("expected 1, got %d", c.Count())
	}
}

func TestXDatasetCollection_Pop_Empty(t *testing.T) {
	// 5.1 regression: must return nil, not panic
	c := &XDatasetCollection{}
	result := c.Pop()
	if result != nil {
		t.Fatalf("expected nil from Pop on empty collection, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Unshift / Shift
// ---------------------------------------------------------------------------

func TestXDatasetCollection_UnshiftShift(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("n", 10))
	c.Unshift(newDS("n", 1))
	first := c.Shift()
	v, _ := first.GetInt("n")
	if v != 1 {
		t.Fatalf("want 1, got %d", v)
	}
}

func TestXDatasetCollection_Shift_Empty(t *testing.T) {
	// 5.1 regression: must return nil, not panic
	c := &XDatasetCollection{}
	result := c.Shift()
	if result != nil {
		t.Fatalf("expected nil from Shift on empty collection, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Count / Get
// ---------------------------------------------------------------------------

func TestXDatasetCollection_Count(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("a", 1))
	c.Push(newDS("b", 2))
	c.Push(newDS("c", 3))
	if c.Count() != 3 {
		t.Fatalf("expected 3, got %d", c.Count())
	}
}

func TestXDatasetCollection_Get_Valid(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("x", 100))
	el, ok := c.Get(0)
	if !ok || el == nil {
		t.Fatal("expected element at index 0")
	}
	v, _ := el.GetInt("x")
	if v != 100 {
		t.Fatalf("want 100, got %d", v)
	}
}

func TestXDatasetCollection_Get_NegativeIndex(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("x", 1))
	_, ok := c.Get(-1)
	if ok {
		t.Fatal("expected false for negative index")
	}
}

func TestXDatasetCollection_Get_OutOfBounds(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("x", 1))
	_, ok := c.Get(5)
	if ok {
		t.Fatal("expected false for out-of-bounds index")
	}
}

// ---------------------------------------------------------------------------
// GetData / GetDataString / GetDataBool / GetDataInt / GetDataFloat / GetDataTime
// ---------------------------------------------------------------------------

func TestXDatasetCollection_GetData_Priority(t *testing.T) {
	// Last entry has highest priority
	c := &XDatasetCollection{}
	c.Push(newDS("x", "first"))
	c.Push(newDS("x", "second"))
	v, ok := c.GetData("x")
	if !ok || v != "second" {
		t.Fatalf("want second/true, got %v/%v", v, ok)
	}
}

func TestXDatasetCollection_GetData_Missing(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("a", 1))
	_, ok := c.GetData("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestXDatasetCollection_GetDataString_Nil(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("k", nil))
	v, ok := c.GetDataString("k")
	if !ok || v != "" {
		t.Fatalf("want empty/true for nil, got %q/%v", v, ok)
	}
}

func TestXDatasetCollection_GetDataString_OK(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("k", "hello"))
	v, ok := c.GetDataString("k")
	if !ok || v != "hello" {
		t.Fatalf("want hello/true, got %q/%v", v, ok)
	}
}

func TestXDatasetCollection_GetDataBool(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("b", true))
	v, ok := c.GetDataBool("b")
	if !ok || !v {
		t.Fatal("expected true")
	}
}

func TestXDatasetCollection_GetDataBool_WrongType(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("b", "notabool"))
	_, ok := c.GetDataBool("b")
	if ok {
		t.Fatal("expected false for non-bool value")
	}
}

func TestXDatasetCollection_GetDataInt(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("n", 42))
	v, ok := c.GetDataInt("n")
	if !ok || v != 42 {
		t.Fatalf("want 42/true, got %d/%v", v, ok)
	}
}

func TestXDatasetCollection_GetDataFloat(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("f", float64(1.5)))
	v, ok := c.GetDataFloat("f")
	if !ok || v != 1.5 {
		t.Fatalf("want 1.5/true, got %f/%v", v, ok)
	}
}

func TestXDatasetCollection_GetDataTime(t *testing.T) {
	now := time.Now()
	c := &XDatasetCollection{}
	c.Push(newDS("t", now))
	v, ok := c.GetDataTime("t")
	if !ok || !v.Equal(now) {
		t.Fatalf("want %v/true, got %v/%v", now, v, ok)
	}
}

// ---------------------------------------------------------------------------
// GetCollection
// ---------------------------------------------------------------------------

func TestXDatasetCollection_GetCollection_OK(t *testing.T) {
	inner := &XDatasetCollection{}
	inner.Push(newDS("y", 5))
	c := &XDatasetCollection{}
	c.Push(newDS("sub", inner))
	got, ok := c.GetCollection("sub")
	if !ok || got == nil {
		t.Fatal("expected collection")
	}
}

func TestXDatasetCollection_GetCollection_WrongType(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("x", "notacollection"))
	_, ok := c.GetCollection("x")
	if ok {
		t.Fatal("expected false for non-collection value")
	}
}

// ---------------------------------------------------------------------------
// Clone
// ---------------------------------------------------------------------------

func TestXDatasetCollection_Clone(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("n", 1))
	c.Push(newDS("n", 2))
	cloned := c.Clone()
	if cloned.Count() != 2 {
		t.Fatalf("clone count: want 2, got %d", cloned.Count())
	}
	// Mutate original
	(*c)[0].Set("n", 99)
	v, _ := cloned.Get(0)
	n, _ := v.GetInt("n")
	if n != 1 {
		t.Fatalf("clone was not independent from original, got %d", n)
	}
}

// ---------------------------------------------------------------------------
// String / GoString
// ---------------------------------------------------------------------------

func TestXDatasetCollection_String(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("a", 1))
	s := c.String()
	if s == "" {
		t.Fatal("String() returned empty")
	}
}

func TestXDatasetCollection_GoString(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("a", 1))
	s := c.GoString()
	if s == "" {
		t.Fatal("GoString() returned empty")
	}
}

// ---------------------------------------------------------------------------
// Edge cases: interleaved Push/Shift/Pop
// ---------------------------------------------------------------------------

func TestXDatasetCollection_Interleaved(t *testing.T) {
	c := &XDatasetCollection{}
	c.Push(newDS("n", 1))
	c.Push(newDS("n", 2))
	c.Push(newDS("n", 3))

	first := c.Shift()
	n1, _ := first.GetInt("n")
	if n1 != 1 {
		t.Fatalf("expected 1, got %d", n1)
	}

	c.Unshift(newDS("n", 0))
	front := c.Shift()
	n0, _ := front.GetInt("n")
	if n0 != 0 {
		t.Fatalf("expected 0, got %d", n0)
	}

	last := c.Pop()
	n3, _ := last.GetInt("n")
	if n3 != 3 {
		t.Fatalf("expected 3, got %d", n3)
	}

	if c.Count() != 1 {
		t.Fatalf("expected 1 remaining, got %d", c.Count())
	}
}
