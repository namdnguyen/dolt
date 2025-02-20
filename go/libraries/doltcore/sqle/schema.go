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

package sqle

import (
	"github.com/src-d/go-mysql-server/sql"

	"github.com/liquidata-inc/dolt/go/libraries/doltcore/schema"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/sqle/types"
)

// doltSchemaToSqlSchema returns the sql.Schema corresponding to the dolt schema given.
func doltSchemaToSqlSchema(tableName string, sch schema.Schema) (sql.Schema, error) {
	cols := make([]*sql.Column, sch.GetAllCols().Size())

	var i int
	err := sch.GetAllCols().Iter(func(tag uint64, col schema.Column) (stop bool, err error) {
		var innerErr error
		cols[i], innerErr = doltColToSqlCol(tableName, col)
		if innerErr != nil {
			return true, innerErr
		}
		i++
		return false, nil
	})

	return cols, err
}

// SqlSchemaToDoltResultSchema returns a dolt Schema from the sql schema given, suitable for use as a result set. For
// creating tables, use SqlSchemaToDoltSchema.
func SqlSchemaToDoltResultSchema(sqlSchema sql.Schema) (schema.Schema, error) {
	var cols []schema.Column
	for i, col := range sqlSchema {
		convertedCol, err := SqlColToDoltCol(uint64(i), col)
		if err != nil {
			return nil, err
		}
		cols = append(cols, convertedCol)
	}

	colColl, err := schema.NewColCollection(cols...)
	if err != nil {
		panic(err)
	}

	return schema.UnkeyedSchemaFromCols(colColl), nil
}

// SqlSchemaToDoltResultSchema returns a dolt Schema from the sql schema given, suitable for use in creating a table.
// For result set schemas, see SqlSchemaToDoltResultSchema.
func SqlSchemaToDoltSchema(sqlSchema sql.Schema) (schema.Schema, error) {
	var cols []schema.Column
	for i, col := range sqlSchema {
		convertedCol, err := SqlColToDoltCol(uint64(i), col)
		if err != nil {
			return nil, err
		}
		cols = append(cols, convertedCol)
	}

	colColl, err := schema.NewColCollection(cols...)
	if err != nil {
		return nil, err
	}

	err = schema.ValidateForInsert(colColl)
	if err != nil {
		return nil, err
	}

	return schema.SchemaFromCols(colColl), nil
}

// doltColToSqlCol returns the SQL column corresponding to the dolt column given.
func doltColToSqlCol(tableName string, col schema.Column) (*sql.Column, error) {
	colType, err := types.NomsKindToSqlType(col.Kind)
	if err != nil {
		return nil, err
	}
	return &sql.Column{
		Name:     col.Name,
		Type:     colType,
		Default:  nil,
		Nullable: col.IsNullable(),
		Source:   tableName,
	}, nil
}

// doltColToSqlCol returns the dolt column corresponding to the SQL column given
func SqlColToDoltCol(tag uint64, col *sql.Column) (schema.Column, error) {
	var constraints []schema.ColConstraint
	if !col.Nullable {
		constraints = append(constraints, schema.NotNullConstraint{})
	}
	kind, err := types.SqlTypeToNomsKind(col.Type)
	if err != nil {
		return schema.Column{}, err
	}
	return schema.NewColumn(col.Name, tag, kind, col.PrimaryKey, constraints...), nil
}
