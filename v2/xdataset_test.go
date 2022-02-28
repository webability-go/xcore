package xcore

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func ExampleXDataset() {

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	ds := &XDataset{
		"v1":  123,
		"v2":  "abc",
		"vt":  tmp,
		"v3":  true,
		"vpi": 3.1415927,
	}
	fmt.Println(ds)

	data := &XDataset{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football", "sport": "yes"},
			&XDataset{"name": "Ping-pong", "sport": "yes"},
			&XDataset{"name": "Swimming", "sport": "yes"},
			&XDataset{"name": "Videogames", "sport": "no"},
		},
		"preferredhobby": &XDataset{
			"name":  "Baseball",
			"sport": "yes",
		},
		"metadata": &XDataset{
			"preferred-color": "blue",
			"Salary":          3568.65,
			"hiredate":        tmp,
		},
	}

	fmt.Println(data)
	// Output:
	// xcore.XDataset{v1:123 v2:abc v3:true vpi:3.1415927 vt:2020-01-01 12:00:00 +0000 UTC}
	// xcore.XDataset{clientname:Fred clientpicture:face.jpg hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:xcore.XDataset{Salary:3568.65 hiredate:2020-01-01 12:00:00 +0000 UTC preferred-color:blue} preferredhobby:xcore.XDataset{name:Baseball sport:yes}}
}

func TestXDataset_new(t *testing.T) {

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	data := map[string]interface{}{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": []map[string]interface{}{
			{"name": "Football", "sport": "yes"},
			{"name": "Ping-pong", "sport": "yes"},
			{"name": "Swimming", "sport": "yes"},
			{"name": "Videogames", "sport": "no"},
		},
		"preferredhobby": map[string]interface{}{
			"name":  "Baseball",
			"sport": "yes",
		},
		"metadata": map[string]interface{}{
			"preferred-color": "blue",
			"Salary":          3568.65,
			"hiredate":        tmp,
		},
	}

	ds := NewXDataset(data)
	str := fmt.Sprintf("%#v", ds)
	if str != "#xcore.XDataset{clientname:\"Fred\" clientpicture:\"face.jpg\" hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:#xcore.XDataset{Salary:3568.65 hiredate:time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC) preferred-color:\"blue\"} preferredhobby:#xcore.XDataset{name:\"Baseball\" sport:\"yes\"}}" {
		t.Error("Error creating and #printing new complex XDataset " + str)
		return
	}
}

func TestXDataset_simple_print(t *testing.T) {

	// 1. Create a simple XDataset
	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	ds := &XDataset{
		"v1":  123,
		"v2":  "abc",
		"vt":  tmp,
		"v3":  true,
		"vpi": 3.1415927,
	}

	// 2. print
	str := fmt.Sprintf("%v", ds)
	if str != "xcore.XDataset{v1:123 v2:abc v3:true vpi:3.1415927 vt:2020-01-01 12:00:00 +0000 UTC}" {
		t.Error("Error creating and printing simple XDataset " + str)
		return
	}

	str = fmt.Sprintf("%#v", ds)
	if str != "#xcore.XDataset{v1:123 v2:\"abc\" v3:true vpi:3.1415927 vt:time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC)}" {
		t.Error("Error creating and #printing simple XDataset " + str)
		return
	}
}

func getComplexDataset() *XDataset {

	// Create a complex XDataset with subsets
	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	data := &XDataset{
		"clientname":    "Fred",
		"clientpicture": "face.jpg",
		"hobbies": &XDatasetCollection{
			&XDataset{"name": "Football", "sport": "yes"},
			&XDataset{"name": "Ping-pong", "sport": "yes"},
			&XDataset{"name": "Swimming", "sport": "yes"},
			&XDataset{"name": "Videogames", "sport": "no"},
		},
		"preferredhobby": &XDataset{
			"name":  "Baseball",
			"sport": "yes",
		},
		"metadata": &XDataset{
			"preferred-color": "blue",
			"salary":          3568.65,
			"hiredate":        tmp,
			"hascat":          true,
			"hasdog":          false,
			"previoussalary":  0.0,
			"resume":          "",
			"numdata1":        0,
			"numdata2":        17,
		},
	}
	return data
}

func TestXDataset_complex_print(t *testing.T) {

	data := getComplexDataset()

	str := fmt.Sprintf("%v", data)
	if str != "xcore.XDataset{clientname:Fred clientpicture:face.jpg hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:xcore.XDataset{hascat:true hasdog:false hiredate:2020-01-01 12:00:00 +0000 UTC numdata1:0 numdata2:17 preferred-color:blue previoussalary:0 resume: salary:3568.65} preferredhobby:xcore.XDataset{name:Baseball sport:yes}}" {
		t.Error("Error creating and printing complex XDataset " + str)
		return
	}

	str = fmt.Sprintf("%#v", data)
	if str != "#xcore.XDataset{clientname:\"Fred\" clientpicture:\"face.jpg\" hobbies:XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ] metadata:#xcore.XDataset{hascat:true hasdog:false hiredate:time.Date(2020, time.January, 1, 12, 0, 0, 0, time.UTC) numdata1:0 numdata2:17 preferred-color:\"blue\" previoussalary:0 resume:\"\" salary:3568.65} preferredhobby:#xcore.XDataset{name:\"Baseball\" sport:\"yes\"}}" {
		t.Error("Error creating and #printing complex XDataset " + str)
		return
	}
}

func TestXDataset_Get(t *testing.T) {

	data := getComplexDataset()

	// 2. Gets and paths
	v1, _ := data.GetString("clientname") // implicity use Get
	if v1 != "Fred" {
		t.Error("Error getting clientname from XDataset " + v1)
		return
	}

	v2, _ := data.Get("hobbies")
	strv2 := fmt.Sprintf("%v", v2)
	if strv2 != "XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ]" {
		t.Error("Error getting hobbies from XDataset " + strv2)
		return
	}

	// Existing paths
	v3, _ := data.Get("hobbies>2")
	strv3 := fmt.Sprintf("%v", v3)
	if strv3 != "xcore.XDataset{name:Swimming sport:yes}" {
		t.Error("Error getting hobbies>2 from XDataset " + strv3)
		return
	}

	v4, _ := data.GetString("hobbies>2>name") // implicity use Get
	if v4 != "Swimming" {
		t.Error("Error getting hobbies>2>name from XDataset " + v4)
		return
	}

	// Non existing paths
	v5, ok := data.Get("wrong>2>name")
	if v5 != nil || ok != false {
		t.Error("Error getting wrong>2>name from XDataset " + fmt.Sprintf("%#v", v5))
		return
	}

	v6, ok := data.Get("hobbies>xxx>name")
	if v6 != nil || ok != false {
		t.Error("Error getting hobbies>xxx>name from XDataset " + fmt.Sprintf("%#v", v6))
		return
	}

	v7, ok := data.Get("hobbies>999>name")
	if v7 != nil || ok != false {
		t.Error("Error getting hobbies>999>name from XDataset " + fmt.Sprintf("%#v", v7))
		return
	}
}

func TestXDataset_GetDataset(t *testing.T) {

	data := getComplexDataset()

	// 2. Gets a real dataset
	v1, _ := data.GetDataset("hobbies>2") // is a dataset
	strv1 := fmt.Sprintf("%v", v1)
	if strv1 != "xcore.XDataset{name:Swimming sport:yes}" {
		t.Error("Error getting XDataset " + strv1)
		return
	}

	// getdataset with errors
	v2, ok := data.GetDataset("hobbies") // is a datacollection: error
	if v2 != nil || ok != false {
		t.Error("Error getting XDataset " + fmt.Sprintf("%#v", v2))
		return
	}

	v3, ok := data.GetDataset("clientname") // is a string: error
	if v3 != nil || ok != false {
		t.Error("Error getting XDataset " + fmt.Sprintf("%#v", v3))
		return
	}

	v4, ok := data.GetDataset("xxx") // does not exists: error
	if v4 != nil || ok != false {
		t.Error("Error getting XDataset " + fmt.Sprintf("%#v", v4))
		return
	}
}

func TestXDataset_GetCollection(t *testing.T) {

	data := getComplexDataset()

	// 2. Gets a real dataset
	v1, _ := data.GetCollection("hobbies") // is a collection
	strv1 := fmt.Sprintf("%v", v1)
	if strv1 != "XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ]" {
		t.Error("Error getting XDatasetCollection " + strv1)
		return
	}

	// getdataset with errors
	v2, ok := data.GetCollection("hobbies>2") // is a dataset: error
	if v2 != nil || ok != false {
		t.Error("Error getting XDataset " + fmt.Sprintf("%#v", v2))
		return
	}

	v3, ok := data.GetCollection("clientname") // is a string: error
	if v3 != nil || ok != false {
		t.Error("Error getting XDataset " + fmt.Sprintf("%#v", v3))
		return
	}

	v4, ok := data.GetCollection("xxx") // does not exists: error
	if v4 != nil || ok != false {
		t.Error("Error getting XDataset " + fmt.Sprintf("%#v", v4))
		return
	}
}

func TestXDataset_GetString(t *testing.T) {

	data := getComplexDataset()

	// 2. Gets a real dataset
	v1, _ := data.GetString("hobbies") // is a collection, printed
	if v1 != "XDatasetCollection[0:xcore.XDataset{name:Football sport:yes} 1:xcore.XDataset{name:Ping-pong sport:yes} 2:xcore.XDataset{name:Swimming sport:yes} 3:xcore.XDataset{name:Videogames sport:no} ]" {
		t.Error("Error getting string " + v1)
		return
	}

	// getdataset with errors
	v2, ok := data.GetString("hobbies>2") // is a dataset: printed
	if v2 != "xcore.XDataset{name:Swimming sport:yes}" {
		t.Error("Error getting string " + v2)
		return
	}

	v3, ok := data.GetString("clientname") // is a string: error
	if v3 != "Fred" {
		t.Error("Error getting string " + v3)
		return
	}

	v4, ok := data.GetString("xxx") // does not exists: error
	if v4 != "" || ok != false {
		t.Error("Error getting XDataset " + v4)
		return
	}
}

func GetTypesDataSet() *XDataset {
	// gets ints and floats
	// Note: byte is not part of tests since it's an alias for uint8
	tm1, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	var nilptr interface{}
	data := &XDataset{
		"bool0":         false,
		"bool1":         true,
		"time0":         time.Time{},
		"time1":         tm1,
		"int0":          0,
		"int1":          1,
		"int-limit":     math.MinInt64,
		"intlimit":      math.MaxInt64,
		"int80":         int8(0),
		"int81":         int8(1),
		"int8-limit":    int8(math.MinInt8),
		"int8limit":     int8(math.MaxInt8),
		"int160":        int16(0),
		"int161":        int16(1),
		"int16-limit":   int16(math.MinInt16),
		"int16limit":    int16(math.MaxInt16),
		"int320":        int32(0),
		"int321":        int32(1),
		"int32-limit":   int32(math.MinInt32),
		"int32limit":    int32(math.MaxInt32),
		"int640":        int64(0),
		"int641":        int64(1),
		"int64-limit":   int64(math.MinInt64),
		"int64limit":    int64(math.MaxInt64),
		"uint0":         uint(0),
		"uint1":         uint(1),
		"uint80":        uint8(0),
		"uint81":        uint8(1),
		"uint8limit":    uint8(math.MaxUint8),
		"uint160":       uint16(0),
		"uint161":       uint16(1),
		"uint16limit":   uint16(math.MaxUint16),
		"uint320":       uint32(0),
		"uint321":       uint32(1),
		"uint32limit":   uint32(math.MaxUint32),
		"uint640":       uint64(0),
		"uint641":       uint64(1),
		"uint64limit":   uint64(math.MaxUint64),
		"float320":      float32(0),
		"float321":      float32(1),
		"float32-limit": math.SmallestNonzeroFloat32,
		"float32limit":  math.MaxFloat32,
		"float640":      float64(0),
		"float641":      float64(1),
		"float64-limit": math.SmallestNonzeroFloat64,
		"float64limit":  math.MaxFloat64,
		"pointer0":      nilptr,
		"pointer1":      &XDataset{},
	}
	return data
}

func TestXDataset_GetBool(t *testing.T) {

	data := GetTypesDataSet()

	tests := map[string]bool{
		"bool0":    false,
		"bool1":    true,
		"time0":    false,
		"time1":    true,
		"int0":     false,
		"int1":     true,
		"int80":    false,
		"int81":    true,
		"int160":   false,
		"int161":   true,
		"int320":   false,
		"int321":   true,
		"int640":   false,
		"int641":   true,
		"uint0":    false,
		"uint1":    true,
		"uint80":   false,
		"uint81":   true,
		"uint160":  false,
		"uint161":  true,
		"uint320":  false,
		"uint321":  true,
		"uint640":  false,
		"uint641":  true,
		"float320": false,
		"float321": true,
		"float640": false,
		"float641": true,
		"pointer0": false,
		"pointer1": true,
	}

	for id, res := range tests {
		r, ok := data.GetBool(id)
		if r != res || !ok {
			t.Error("Error getting bool " + fmt.Sprintf("%s %#+v %#+v", id, r, ok))
			return
		}
	}
	r, ok := data.GetBool("xxx")
	if r != false || ok {
		t.Error("Error getting bool " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetInt(t *testing.T) {

	data := GetTypesDataSet()

	tests := map[string]int{
		"bool0":         0,
		"bool1":         1,
		"time0":         0,
		"time1":         1577880000,
		"int0":          0,
		"int1":          1,
		"int-limit":     math.MinInt64,
		"intlimit":      math.MaxInt64,
		"int80":         0,
		"int81":         1,
		"int8-limit":    math.MinInt8,
		"int8limit":     math.MaxInt8,
		"int160":        0,
		"int161":        1,
		"int16-limit":   math.MinInt16,
		"int16limit":    math.MaxInt16,
		"int320":        0,
		"int321":        1,
		"int32-limit":   math.MinInt32,
		"int32limit":    math.MaxInt32,
		"int640":        0,
		"int641":        1,
		"int64-limit":   math.MinInt64,
		"int64limit":    math.MaxInt64,
		"uint0":         0,
		"uint1":         1,
		"uint80":        0,
		"uint81":        1,
		"uint160":       0,
		"uint161":       1,
		"uint320":       0,
		"uint321":       1,
		"uint640":       0,
		"uint641":       1,
		"float320":      0,
		"float321":      1,
		"float32-limit": 0,
		"float32limit":  math.MinInt64, // deberia de ser max ???
		"float640":      0,
		"float641":      1,
		"float64-limit": 0,
		"float64limit":  math.MinInt64,
	}

	for id, res := range tests {
		r, ok := data.GetInt(id)
		if r != res || !ok {
			t.Error("Error getting int " + fmt.Sprintf("%s %#+v %#+v", id, r, ok))
			return
		}
	}
	r, ok := data.GetInt("xxx")
	if r != 0 || ok {
		t.Error("Error getting int " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetFloat(t *testing.T) {

	data := GetTypesDataSet()

	tests := map[string]float64{
		"bool0":         0,
		"bool1":         1,
		"time0":         0,
		"time1":         1577880000,
		"int0":          0,
		"int1":          1,
		"int-limit":     -9223372036854775808,
		"intlimit":      9223372036854775807,
		"int80":         0,
		"int81":         1,
		"int8-limit":    -128,
		"int8limit":     127,
		"int160":        0,
		"int161":        1,
		"int16-limit":   -32768,
		"int16limit":    32767,
		"int320":        0,
		"int321":        1,
		"int32-limit":   -2147483648,
		"int32limit":    2147483647,
		"int640":        0,
		"int641":        1,
		"int64-limit":   -9223372036854775808,
		"int64limit":    9223372036854775807,
		"uint0":         0,
		"uint1":         1,
		"uint80":        0,
		"uint81":        1,
		"uint160":       0,
		"uint161":       1,
		"uint320":       0,
		"uint321":       1,
		"uint640":       0,
		"uint641":       1,
		"float320":      0,
		"float321":      1,
		"float32-limit": math.SmallestNonzeroFloat32,
		"float32limit":  math.MaxFloat32,
		"float640":      0,
		"float641":      1,
		"float64-limit": math.SmallestNonzeroFloat64,
		"float64limit":  math.MaxFloat64,
	}

	for id, res := range tests {
		r, ok := data.GetFloat(id)
		if r != res || !ok {
			t.Error("Error getting float " + fmt.Sprintf("%s %#+v %#+v", id, r, ok))
			return
		}
	}
	r, ok := data.GetFloat("xxx")
	if r != 0 || ok {
		t.Error("Error getting float " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetTime(t *testing.T) {

	data := GetTypesDataSet()

	tmp, _ := time.Parse(time.RFC3339, "2020-01-01T12:00:00.0Z")
	tests := map[string]time.Time{
		"time0": {},
		"time1": tmp,
	}

	for id, res := range tests {
		r, ok := data.GetTime(id)
		if r != res || !ok {
			t.Error("Error getting time " + fmt.Sprintf("%s %#+v %#+v", id, r, ok))
			return
		}
	}
	r, ok := data.GetTime("float641")
	if !r.Equal(time.Time{}) || ok {
		t.Error("Error getting time " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
	r, ok = data.GetTime("xxx")
	if !r.Equal(time.Time{}) || ok {
		t.Error("Error getting time " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetStringCollection(t *testing.T) {

	data := &XDataset{
		"sc": []string{
			"s1", "s2", "s3", "s4",
		},
	}

	r, ok := data.GetStringCollection("sc")
	strv1 := fmt.Sprintf("%v", r)
	if strv1 != "[s1 s2 s3 s4]" || !ok {
		t.Error("Error getting stringcollection " + strv1)
		return
	}
	r, ok = data.GetStringCollection("xxx")
	if r != nil || ok {
		t.Error("Error getting stringcollection " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetBoolCollection(t *testing.T) {

	data := &XDataset{
		"sc": []bool{
			false, false, true, false,
		},
	}

	r, ok := data.GetBoolCollection("sc")
	strv1 := fmt.Sprintf("%v", r)
	if strv1 != "[false false true false]" || !ok {
		t.Error("Error getting stringcollection " + strv1)
		return
	}
	r, ok = data.GetBoolCollection("xxx")
	if r != nil || ok {
		t.Error("Error getting stringcollection " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetIntCollection(t *testing.T) {

	data := &XDataset{
		"sc": []int{
			0, 1, math.MaxInt32, math.MinInt64,
		},
	}

	r, ok := data.GetIntCollection("sc")
	strv1 := fmt.Sprintf("%v", r)
	if strv1 != "[0 1 2147483647 -9223372036854775808]" || !ok {
		t.Error("Error getting stringcollection " + strv1)
		return
	}
	r, ok = data.GetIntCollection("xxx")
	if r != nil || ok {
		t.Error("Error getting stringcollection " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetFloatCollection(t *testing.T) {

	data := &XDataset{
		"sc": []float64{
			0, 1, math.Pi, math.MaxFloat64,
		},
	}

	r, ok := data.GetFloatCollection("sc")
	strv1 := fmt.Sprintf("%v", r)
	if strv1 != "[0 1 3.141592653589793 1.7976931348623157e+308]" || !ok {
		t.Error("Error getting stringcollection " + strv1)
		return
	}
	r, ok = data.GetFloatCollection("xxx")
	if r != nil || ok {
		t.Error("Error getting stringcollection " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_GetTimeCollection(t *testing.T) {

	data := &XDataset{
		"sc": []time.Time{
			{}, {},
		},
	}

	r, ok := data.GetTimeCollection("sc")
	strv1 := fmt.Sprintf("%v", r)
	if strv1 != "[0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC]" || !ok {
		t.Error("Error getting stringcollection " + strv1)
		return
	}
	r, ok = data.GetTimeCollection("xxx")
	if r != nil || ok {
		t.Error("Error getting stringcollection " + fmt.Sprintf("xxx %#v %#v", r, ok))
		return
	}
}

func TestXDataset_Del(t *testing.T) {

	data := GetTypesDataSet()

	r, ok := data.GetInt("int1")
	if r != 1 || !ok {
		t.Error("Error getting int " + fmt.Sprintf("int1 %#v %#v", r, ok))
		return
	}
	data.Del("int1")
	r, ok = data.GetInt("int1")
	if r == 1 || ok {
		t.Error("Error getting int " + fmt.Sprintf("int1 %#v %#v", r, ok))
		return
	}
}

func TestXDataset_Clone(t *testing.T) {

	ds := &XDataset{
		"v1": 123,
		"v2": "abc",
		"v3": true,
		"v4": &XDataset{
			"v4_p1": 456,
			"v4_p2": "def",
			"v4_p3": false,
		},
		"v5": &XDatasetCollection{
			&XDataset{"v1": "p1"},
			&XDataset{"v1": "p2"},
			&XDataset{"v1": "p3"},
			&XDataset{"v1": "p4"},
		},
	}
	original := fmt.Sprintf("%v", ds)

	cds := ds.Clone()
	cloned := fmt.Sprintf("%v", cds)

	if original != cloned {
		t.Error("Error during clonation of XDataset " + cloned)
		return
	}

	// Verify it has been really cloned
	(*(*ds)["v4"].(*XDataset))["p5"] = "val5"
	val5o, _ := ds.GetString("v4>p5")  // should be "val5"
	val5c, _ := cds.GetString("v4>p5") // should be ""
	if val5o == val5c {
		t.Error("Error during clonation of XDataset " + val5c)
		return
	}
}
