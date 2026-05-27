package xcore

import (
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// Push / Pop
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_PushPop(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("n", 1))
	dc.Push(newDS("n", 2))
	if dc.Count() != 2 {
		t.Fatalf("expected 2, got %d", dc.Count())
	}
	last := dc.Pop()
	v, _ := last.GetInt("n")
	if v != 2 {
		t.Fatalf("want 2, got %d", v)
	}
	if dc.Count() != 1 {
		t.Fatalf("expected 1 after Pop, got %d", dc.Count())
	}
}

func TestXDatasetCollectionTS_Pop_Empty(t *testing.T) {
	// 5.1/5.4 regression: must not panic or deadlock on empty collection
	dc := &XDatasetCollectionTS{}
	result := dc.Pop()
	if result != nil {
		t.Fatalf("expected nil from Pop on empty, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Unshift / Shift
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_UnshiftShift(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("n", 10))
	dc.Unshift(newDS("n", 1))
	first := dc.Shift()
	v, _ := first.GetInt("n")
	if v != 1 {
		t.Fatalf("want 1, got %d", v)
	}
}

func TestXDatasetCollectionTS_Shift_Empty(t *testing.T) {
	// 5.1 regression: must not panic or deadlock on empty collection
	dc := &XDatasetCollectionTS{}
	result := dc.Shift()
	if result != nil {
		t.Fatalf("expected nil from Shift on empty, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Count / Get
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_Count(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("a", 1))
	dc.Push(newDS("b", 2))
	if dc.Count() != 2 {
		t.Fatalf("expected 2, got %d", dc.Count())
	}
}

func TestXDatasetCollectionTS_Get_Valid(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("x", 99))
	el, ok := dc.Get(0)
	if !ok || el == nil {
		t.Fatal("expected element at index 0")
	}
	v, _ := el.GetInt("x")
	if v != 99 {
		t.Fatalf("want 99, got %d", v)
	}
}

func TestXDatasetCollectionTS_Get_NegativeIndex(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("x", 1))
	_, ok := dc.Get(-1)
	if ok {
		t.Fatal("expected false for negative index")
	}
}

func TestXDatasetCollectionTS_Get_OutOfBounds(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("x", 1))
	_, ok := dc.Get(99)
	if ok {
		t.Fatal("expected false for out-of-bounds index")
	}
}

// ---------------------------------------------------------------------------
// GetData — 5.2 single-snapshot race fix
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_GetData_Priority(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("k", "low"))
	dc.Push(newDS("k", "high"))
	v, ok := dc.GetData("k")
	if !ok || v != "high" {
		t.Fatalf("want high/true, got %v/%v", v, ok)
	}
}

func TestXDatasetCollectionTS_GetData_Missing(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("a", 1))
	_, ok := dc.GetData("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestXDatasetCollectionTS_GetData_Concurrent(t *testing.T) {
	// 5.2 regression: GetData must not race with Push
	dc := &XDatasetCollectionTS{}
	for i := 0; i < 10; i++ {
		dc.Push(newDS("n", i))
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dc.GetData("n")
			dc.Push(newDS("n", i+100))
		}(i)
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// GetDataString — 5.5 nil handling
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_GetDataString_Nil(t *testing.T) {
	// 5.5 regression: nil value must return ("", true), not ("nil", true)
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("k", nil))
	v, ok := dc.GetDataString("k")
	if !ok || v != "" {
		t.Fatalf("want empty/true for nil, got %q/%v", v, ok)
	}
}

func TestXDatasetCollectionTS_GetDataString_OK(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("k", "hello"))
	v, ok := dc.GetDataString("k")
	if !ok || v != "hello" {
		t.Fatalf("want hello/true, got %q/%v", v, ok)
	}
}

func TestXDatasetCollectionTS_GetDataString_Missing(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	_, ok := dc.GetDataString("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

// ---------------------------------------------------------------------------
// GetDataBool / GetDataInt / GetDataFloat / GetDataTime
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_GetDataBool(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("b", true))
	v, ok := dc.GetDataBool("b")
	if !ok || !v {
		t.Fatal("expected true")
	}
}

func TestXDatasetCollectionTS_GetDataBool_WrongType(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("b", "notabool"))
	_, ok := dc.GetDataBool("b")
	if ok {
		t.Fatal("expected false for non-bool value")
	}
}

func TestXDatasetCollectionTS_GetDataInt(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("n", 7))
	v, ok := dc.GetDataInt("n")
	if !ok || v != 7 {
		t.Fatalf("want 7/true, got %d/%v", v, ok)
	}
}

func TestXDatasetCollectionTS_GetDataFloat(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("f", float64(2.5)))
	v, ok := dc.GetDataFloat("f")
	if !ok || v != 2.5 {
		t.Fatalf("want 2.5/true, got %f/%v", v, ok)
	}
}

func TestXDatasetCollectionTS_GetDataTime(t *testing.T) {
	now := time.Now()
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("t", now))
	v, ok := dc.GetDataTime("t")
	if !ok || !v.Equal(now) {
		t.Fatalf("want %v/true, got %v/%v", now, v, ok)
	}
}

// ---------------------------------------------------------------------------
// GetCollection — 4.4 safe type assertion
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_GetCollection_OK(t *testing.T) {
	inner := &XDatasetCollectionTS{}
	inner.Push(newDS("y", 5))
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("sub", inner))
	got, ok := dc.GetCollection("sub")
	if !ok || got == nil {
		t.Fatal("expected collection")
	}
}

func TestXDatasetCollectionTS_GetCollection_WrongType(t *testing.T) {
	// 4.4 regression: must not panic when value is not XDatasetCollectionDef
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("x", "notacollection"))
	got, ok := dc.GetCollection("x")
	if ok || got != nil {
		t.Fatalf("expected (nil, false) for non-collection value, got (%v, %v)", got, ok)
	}
}

func TestXDatasetCollectionTS_GetCollection_Missing(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	_, ok := dc.GetCollection("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

// ---------------------------------------------------------------------------
// Clone — 5.3 RLock during iteration
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_Clone_Basic(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("n", 1))
	dc.Push(newDS("n", 2))
	cloned := dc.Clone()
	if cloned.Count() != 2 {
		t.Fatalf("clone count: want 2, got %d", cloned.Count())
	}
}

func TestXDatasetCollectionTS_Clone_Independent(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("n", 10))
	cloned := dc.Clone()

	// Mutate original — clone must not change
	el, _ := dc.Get(0)
	el.Set("n", 999)

	v, _ := cloned.Get(0)
	n, _ := v.GetInt("n")
	if n != 10 {
		t.Fatalf("clone was not independent from original, got %d", n)
	}
}

func TestXDatasetCollectionTS_Clone_Concurrent(t *testing.T) {
	// 5.3 regression: Clone must hold RLock so concurrent pushes don't race
	dc := &XDatasetCollectionTS{}
	for i := 0; i < 5; i++ {
		dc.Push(newDS("n", i))
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); dc.Clone() }()
		go func(i int) { defer wg.Done(); dc.Push(newDS("n", i+100)) }(i)
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// String / GoString
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_String(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("a", 1))
	s := dc.String()
	if s == "" {
		t.Fatal("String() returned empty")
	}
}

func TestXDatasetCollectionTS_GoString(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	dc.Push(newDS("a", 1))
	s := dc.GoString()
	if s == "" {
		t.Fatal("GoString() returned empty")
	}
}

// ---------------------------------------------------------------------------
// Concurrent safety — broad stress test
// ---------------------------------------------------------------------------

func TestXDatasetCollectionTS_ConcurrentStress(t *testing.T) {
	dc := &XDatasetCollectionTS{}
	var wg sync.WaitGroup
	const goroutines = 20
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			dc.Push(newDS("n", i))
			dc.GetData("n")
			dc.GetDataString("n")
			dc.Count()
			dc.Clone()
		}(i)
	}
	wg.Wait()

	// Drain without panicking (tests Shift/Pop empty guards under concurrency)
	for dc.Count() > 0 {
		dc.Shift()
	}
	// Extra calls on empty collection must be safe
	if dc.Shift() != nil {
		t.Fatal("Shift on empty returned non-nil")
	}
	if dc.Pop() != nil {
		t.Fatal("Pop on empty returned non-nil")
	}
}
