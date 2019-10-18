package main

import (
  "fmt"
  "testing"
  "github.com/webability-go/xcore"
//  "unsafe"
)

/* TEST XDATASET */

func TestXDataset(t *testing.T) {
  // please log
  xcore.LOG = true
  
  // Test 1: creates 100 max, no file, expires in 1 sec
  ds := &xcore.XDataset{
    "v1": 123,
    "v2": "abc",
    "v3": true,
    "v4": &xcore.XDataset{
      "v4_p1": 456,
      "v4_p2": "def",
      "v4_p3": false,
    },
  }
  fmt.Println("ORIGINAL", ds)

  cds := ds.Clone()
  (*(*ds)["v4"].(*xcore.XDataset))["p5"] = "val5"
  
  fmt.Println("CLONED", cds)
  
  dsc := &xcore.XDatasetCollection{
    &xcore.XDataset{ "v1": 123,
      "v2": "abc",
      "v3": true,
      "v4": &xcore.XDataset{
        "v4_p1": 456,
        "v4_p2": "def",
        "v4_p3": false,
      },
    },
    &xcore.XDataset{ "v11": 123,
      "v12": "abc",
      "v13": true,
      "v14": &xcore.XDataset{
        "v14_p1": 456,
        "v14_p2": "def",
        "v14_p3": false,
      },
    },
  }
  
  fmt.Println("ORIGINAL", *dsc)

  cdsc := dsc.Clone().(*xcore.XDatasetCollection)
  
  (*((*(*cdsc)[0].(*xcore.XDataset))["v4"].(*xcore.XDataset)))["algomas"] = "el valor mas"
  
  fmt.Println("CLONED", *cdsc)
  fmt.Printf("%v %v\n", ((*(*dsc)[0].(*xcore.XDataset))["v4"].(*xcore.XDataset)), ((*(*cdsc)[0].(*xcore.XDataset))["v4"].(*xcore.XDataset)))
}


