// Copyright 2019 Liquidata, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"errors"
	"strings"

	"github.com/liquidata-inc/dolt/go/libraries/doltcore"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/row"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/schema"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/pipeline"
	"github.com/liquidata-inc/dolt/go/store/types"
)

type FilterFn = func(r row.Row) (matchesFilter bool)

func ParseWhere(sch schema.Schema, whereClause string) (FilterFn, error) {
	if whereClause == "" {
		return func(r row.Row) bool {
			return true
		}, nil
	} else {
		tokens := strings.Split(whereClause, "=")

		if len(tokens) != 2 {
			return nil, errors.New("'" + whereClause + "' is not in the format key=value")
		}

		key := tokens[0]
		valStr := tokens[1]

		col, ok := sch.GetAllCols().GetByName(key)

		if !ok {
			return nil, errors.New("where clause is invalid. '" + key + "' is not a known column.")
		}

		tag := col.Tag
		convFunc, err := doltcore.GetConvFunc(types.StringKind, col.Kind)
		if err != nil {
			return nil, err
		}

		val, err := convFunc(types.String(valStr))
		if err != nil {
			return nil, errors.New("unable to convert '" + valStr + "' to " + col.KindString())
		}

		return func(r row.Row) bool {
			rowVal, ok := r.GetColVal(tag)

			if !ok {
				return false
			}

			return val.Equals(rowVal)
		}, nil
	}
}

type SelectTransform struct {
	Pipeline *pipeline.Pipeline
	filter   FilterFn
	limit    int
	count    int
}

func NewSelTrans(filter FilterFn, limit int) *SelectTransform {
	return &SelectTransform{filter: filter, limit: limit}
}

func (st *SelectTransform) LimitAndFilter(inRow row.Row, props pipeline.ReadableMap) ([]*pipeline.TransformedRowResult, string) {
	if st.limit <= 0 || st.count < st.limit {
		if st.filter(inRow) {
			st.count++
			return []*pipeline.TransformedRowResult{{RowData: inRow, PropertyUpdates: nil}}, ""
		}
	} else if st.count == st.limit {
		st.Pipeline.NoMore()
	}

	return nil, ""
}
