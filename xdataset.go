package xcore

type XDataset map[string]interface{}

type XRangeDataset []XDataset

func NewXDataset() XDataset {
  return make(XDataset)
}

// makes an interface of XDataset to reuse for otrhe libraries and be sure we can call the functions


