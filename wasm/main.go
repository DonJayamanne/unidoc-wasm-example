package main

import (
	"bytes"
	"fmt"
	unipdf "github.com/unidoc/unidoc/pdf/model"
	"github.com/unidoc/unidoc/pdf/model/optimize"
	"syscall/js"
)

func makeResult(res interface{}, err error) js.Value {
	obj := map[string]interface{}{"result": res}
	if err != nil {
		obj["error"] = err.Error()
	}
	return js.ValueOf(obj)
}

func optimizePDFCall(args []js.Value) {
	result := makeResult(nil, fmt.Errorf("unknown error"))
	defer func() {
		args[1].Call("call", js.Null(), result)
	}()
	params := args[0]
	data := params.Get("data")
	buf := make([]byte, data.Length())
	for i := range buf {
		buf[i] = byte(data.Index(i).Int())
	}
	rBuf := bytes.NewReader(buf)
	reader, err := unipdf.NewPdfReader(rBuf)
	if err != nil {
		result = makeResult(nil, fmt.Errorf("reader create error %s\n", err))
		return
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		result = makeResult(nil, fmt.Errorf("failed to get number of pages"))
		return
	}

	if numPages < 1 {
		result = makeResult(nil, fmt.Errorf("empty pdf - nothing to be done"))
		return
	}

	writer := unipdf.NewPdfWriter()
	writer.SetVersion(1, 5)
	optim := optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   params.Get("CombineDuplicateDirectObjects").Bool(),
		CombineIdenticalIndirectObjects: params.Get("CombineIdenticalIndirectObjects").Bool(),
		ImageUpperPPI:                   params.Get("ImageUpperPPI").Float(),
		UseObjectStreams:                params.Get("UseObjectStreams").Bool(),
		ImageQuality:                    params.Get("ImageQuality").Int(),
		CombineDuplicateStreams:         params.Get("CombineDuplicateStreams").Bool(),
	})
	writer.SetOptimizer(optim)

	ocProps, err := reader.GetOCProperties()
	if err != nil {
		result = makeResult(nil, err)
		return
	}
	writer.SetOCProperties(ocProps)

	for j := 0; j < numPages; j++ {
		page, err := reader.GetPage(j + 1)
		if err != nil {
			result = makeResult(nil, fmt.Errorf("get page error %s", err))
			return
		}

		// Load and set outlines (table of contents).
		outlineTree := reader.GetOutlineTree()

		err = writer.AddPage(page)
		if err != nil {
			result = makeResult(nil, fmt.Errorf("add page error %s", err))
			return
		}
		writer.AddOutlineTree(outlineTree)
	}

	// Copy the forms over to the new document also.
	if reader.AcroForm != nil {
		err = writer.SetForms(reader.AcroForm)
		if err != nil {
			result = makeResult(nil, fmt.Errorf("add forms error %s", err))
			return
		}
	}

	out := bytes.NewBuffer(nil)
	err = writer.Write(out)
	if err != nil {
		result = makeResult(nil, fmt.Errorf("write error %s", err))
		return
	}
	result = makeResult(js.TypedArrayOf(out.Bytes()), nil)
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("OptimizePDFCall", js.NewCallback(optimizePDFCall))
	<-c
}
