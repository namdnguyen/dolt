// Copyright 2016 Attic Labs, Inc. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/attic-labs/noms/go/d"
	"github.com/attic-labs/noms/go/types"
)

func getElemDesc(s types.Collection, index int) types.StructDesc {
	t := s.Type().Desc.(types.CompoundDesc).ElemTypes[index]
	d.PanicIfTrue(types.StructKind != t.Kind(), "Expected StructKind, found %s", types.KindToString[t.Type().Kind()])
	return t.Desc.(types.StructDesc)
}

// GetListElemDesc ensures that l is a types.List of structs, pulls the types.StructDesc that describes the elements of l out of vr, and returns the StructDesc.
func GetListElemDesc(l types.List, vr types.ValueReader) types.StructDesc {
	return getElemDesc(l, 0)
}

// GetMapElemDesc ensures that m is a types.Map of structs, pulls the types.StructDesc that describes the elements of m out of vr, and returns the StructDesc.
func GetMapElemDesc(m types.Map, vr types.ValueReader) types.StructDesc {
	return getElemDesc(m, 1)
}

func writeValuesFromChan(structChan chan types.Struct, sd types.StructDesc, comma rune, output io.Writer) {
	fieldNames := getFieldNamesFromStruct(sd)
	csvWriter := csv.NewWriter(output)
	csvWriter.Comma = comma
	d.PanicIfTrue(csvWriter.Write(fieldNames) != nil, "Failed to write header %v", fieldNames)
	record := make([]string, len(fieldNames))
	for s := range structChan {
		for i, f := range fieldNames {
			record[i] = fmt.Sprintf("%v", s.Get(f))
		}
		d.PanicIfTrue(csvWriter.Write(record) != nil, "Failed to write record %v", record)
	}

	csvWriter.Flush()
	d.PanicIfTrue(csvWriter.Error() != nil, "error flushing csv")
}

// Write takes a types.List l of structs (described by sd) and writes it to output as comma-delineated values.
func WriteList(l types.List, sd types.StructDesc, comma rune, output io.Writer) {
	structChan := make(chan types.Struct, 1024)
	go func() {
		l.IterAll(func(v types.Value, index uint64) {
			structChan <- v.(types.Struct)
		})
		close(structChan)
	}()
	writeValuesFromChan(structChan, sd, comma, output)
}

// Write takes a types.Map m of structs (described by sd) and writes it to output as comma-delineated values.
func WriteMap(m types.Map, sd types.StructDesc, comma rune, output io.Writer) {
	structChan := make(chan types.Struct, 1024)
	go func() {
		m.IterAll(func(k, v types.Value) {
			structChan <- v.(types.Struct)
		})
		close(structChan)
	}()
	writeValuesFromChan(structChan, sd, comma, output)
}

func getFieldNamesFromStruct(structDesc types.StructDesc) (fieldNames []string) {
	structDesc.IterFields(func(name string, t *types.Type) {
		d.PanicIfTrue(!types.IsPrimitiveKind(t.Kind()), "Expected primitive kind, found %s", types.KindToString[t.Kind()])
		fieldNames = append(fieldNames, name)
	})
	return
}
