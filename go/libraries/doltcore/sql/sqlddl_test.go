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

package sql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/liquidata-inc/dolt/go/libraries/doltcore/dtestutils"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/row"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/schema"
	. "github.com/liquidata-inc/dolt/go/libraries/doltcore/sql/sqltestutil"
	"github.com/liquidata-inc/dolt/go/store/types"
)

func TestExecuteCreate(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedSchema schema.Schema
		expectedErr    string
	}{
		{
			name:  "Test create single column schema",
			query: "create table testTable (id int primary key)",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{})),
		},
		{
			name:  "Test create two column schema",
			query: "create table testTable (id int primary key, age int)",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 1, types.IntKind, false)),
		},
		{
			name:        "Test syntax error",
			query:       "create table testTable id int, age int",
			expectedErr: "syntax error",
		},
		{
			name:        "Test no primary keys",
			query:       "create table testTable (id int, age int)",
			expectedErr: "at least one primary key column must be specified",
		},
		{
			name:        "Test bad table name",
			query:       "create table _testTable (id int primary key, age int)",
			expectedErr: "Invalid table name",
		},
		{
			name:        "Test bad table name begins with number",
			query:       "create table 1testTable (id int primary key, age int)",
			expectedErr: "syntax error",
		},
		{
			name:        "Test in use table name",
			query:       "create table people (id int primary key, age int)",
			expectedErr: "Table 'people' already exists",
		},
		{
			name:           "Test in use table name with if not exists",
			query:          "create table if not exists people (id int primary key, age int)",
			expectedSchema: PeopleTestSchema,
		},
		{
			name:  "Test types",
			query: "create table testTable (id int primary key, age int, first varchar, is_married boolean)",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 1, types.IntKind, false),
				schema.NewColumn("first", 2, types.StringKind, false),
				schema.NewColumn("is_married", 3, types.BoolKind, false)),
		},
		{
			name: "Test all supported types",
			query: `create table testTable (
							c0 int primary key, 
							c1 tinyint,
							c2 smallint,
							c3 mediumint,
							c4 integer,
							c5 bigint,
							c6 bool,
							c7 boolean,
							c8 bit,
							c9 text,
							c10 tinytext,
							c11 mediumtext,
							c12 longtext,
							c13 blob,
							c14 tinyblob,
							c15 mediumblob,
							c16 char,
							c17 varchar,
							c18 varchar(80),
							c19 float,
							c20 double,
							c21 decimal,
							c22 int unsigned,
							c23 tinyint unsigned,
							c24 smallint unsigned,
							c25 mediumint unsigned,
							c26 bigint unsigned,
              c27 uuid)`,
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("c0", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("c1", 1, types.IntKind, false),
				schema.NewColumn("c2", 2, types.IntKind, false),
				schema.NewColumn("c3", 3, types.IntKind, false),
				schema.NewColumn("c4", 4, types.IntKind, false),
				schema.NewColumn("c5", 5, types.IntKind, false),
				schema.NewColumn("c6", 6, types.BoolKind, false),
				schema.NewColumn("c7", 7, types.BoolKind, false),
				schema.NewColumn("c8", 8, types.BoolKind, false),
				schema.NewColumn("c9", 9, types.StringKind, false),
				schema.NewColumn("c10", 10, types.StringKind, false),
				schema.NewColumn("c11", 11, types.StringKind, false),
				schema.NewColumn("c12", 12, types.StringKind, false),
				schema.NewColumn("c13", 13, types.BlobKind, false),
				schema.NewColumn("c14", 14, types.BlobKind, false),
				schema.NewColumn("c15", 15, types.BlobKind, false),
				schema.NewColumn("c16", 16, types.StringKind, false),
				schema.NewColumn("c17", 17, types.StringKind, false),
				schema.NewColumn("c18", 18, types.StringKind, false),
				schema.NewColumn("c19", 19, types.FloatKind, false),
				schema.NewColumn("c20", 20, types.FloatKind, false),
				schema.NewColumn("c21", 21, types.FloatKind, false),
				schema.NewColumn("c22", 22, types.UintKind, false),
				schema.NewColumn("c23", 23, types.UintKind, false),
				schema.NewColumn("c24", 24, types.UintKind, false),
				schema.NewColumn("c25", 25, types.UintKind, false),
				schema.NewColumn("c26", 26, types.UintKind, false),
				schema.NewColumn("c27", 27, types.UUIDKind, false),
			),
		},
		{
			name:  "Test primary keys",
			query: "create table testTable (id int, age int, first varchar(80), is_married bool, primary key (id, age))",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 1, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("first", 2, types.StringKind, false),
				schema.NewColumn("is_married", 3, types.BoolKind, false)),
		},
		{
			name:  "Test not null constraints",
			query: "create table testTable (id int, age int, first varchar(80) not null, is_married bool, primary key (id, age))",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 1, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("first", 2, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("is_married", 3, types.BoolKind, false)),
		},
		{
			name:  "Test quoted columns",
			query: "create table testTable (`id` int, `age` int, `timestamp` varchar(80), `is married` bool, primary key (`id`, `age`))",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 1, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("timestamp", 2, types.StringKind, false),
				schema.NewColumn("is married", 3, types.BoolKind, false)),
		},
		{
			name:  "Test tag comments",
			query: "create table testTable (id int primary key comment 'tag:5', age int comment 'tag:10')",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 5, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 10, types.IntKind, false)),
		},
		{
			name:  "Test faulty tag comments",
			query: "create table testTable (id int primary key comment 'tag:a', age int comment 'this is my personal area')",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", 0, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("age", 1, types.IntKind, false)),
		},
		// Real world examples for regression testing
		// TODO: need type conversion for defaults to work here (uint to int)
		// 		{
		// 			name:  "Test ip2nation",
		// 			query: `CREATE TABLE ip2nation (
		//   ip int(11) unsigned NOT NULL default 0,
		//   country char(2) NOT NULL default '',
		//   PRIMARY KEY (ip)
		// );`,
		// 			expectedSchema: dtestutils.CreateSchema(
		// 				schema.NewColumn("ip", 0, types.UintKind, true, schema.NotNullConstraint{}),
		// 				schema.NewColumn("country", 1, types.StringKind, false, schema.NotNullConstraint{})),
		// 		},
		{
			name: "Test ip2nationCountries",
			query: `CREATE TABLE ip2nationCountries (
  code varchar(4) NOT NULL default '',
  iso_code_2 varchar(2) NOT NULL default '',
  iso_code_3 varchar(3) default '',
  iso_country varchar(255) NOT NULL default '',
  country varchar(255) NOT NULL default '',
  lat float NOT NULL default 0.0,
  lon float NOT NULL default 0.0,
  PRIMARY KEY (code)
);`,
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("code", 0, types.StringKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("iso_code_2", 1, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("iso_code_3", 2, types.StringKind, false),
				schema.NewColumn("iso_country", 3, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("country", 4, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("lat", 5, types.FloatKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("lon", 6, types.FloatKind, false, schema.NotNullConstraint{})),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			root, _ := dEnv.WorkingRoot(context.Background())

			sqlStatement, err := sqlparser.Parse(tt.query)
			if err != nil {
				if len(tt.expectedErr) > 0 {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
					return
				} else {
					// fail the test
					require.NoError(t, err)
				}
			}

			s := sqlStatement.(*sqlparser.DDL)

			updatedRoot, sch, err := ExecuteCreate(context.Background(), dEnv.DoltDB, root, s, tt.query)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NotNil(t, updatedRoot)
			assert.Equal(t, tt.expectedSchema, sch)
		})
	}
}

func TestExecuteDrop(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		tableNames  []string
		expectedErr string
	}{
		{
			name:       "drop table",
			query:      "drop table people",
			tableNames: []string{"people"},
		},
		{
			name:       "drop table if exists",
			query:      "drop table if exists people",
			tableNames: []string{"people"},
		},
		{
			name:        "drop non existent",
			query:       "drop table notfound",
			expectedErr: "Unknown table: 'notfound'",
		},
		{
			name:       "drop non existent if exists",
			query:      "drop table if exists notFound",
			tableNames: []string{"notFound"},
		},
		{
			name:       "drop many tables",
			query:      "drop table people, appearances, episodes",
			tableNames: []string{"people", "appearances", "episodes"},
		},
		{
			name:       "drop many tables, some don't exist",
			query:      "drop table if exists people, not_real, appearances, episodes",
			tableNames: []string{"people", "appearances", "not_real", "episodes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			ctx := context.Background()
			root, _ := dEnv.WorkingRoot(ctx)

			sqlStatement, err := sqlparser.Parse(tt.query)
			require.NoError(t, err)

			s := sqlStatement.(*sqlparser.DDL)

			updatedRoot, err := ExecuteDrop(ctx, dEnv.DoltDB, root, s, tt.query)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NotNil(t, updatedRoot)
			for _, tableName := range tt.tableNames {
				has, err := updatedRoot.HasTable(ctx, tableName)
				assert.NoError(t, err)
				assert.False(t, has)
			}
		})
	}
}

func TestAddColumn(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedSchema schema.Schema
		expectedRows   []row.Row
		expectedErr    string
	}{
		{
			name:  "alter add column not null",
			query: "alter table people add (newColumn varchar(80) not null default 'default' comment 'tag:100')",
			expectedSchema: dtestutils.AddColumnToSchema(PeopleTestSchema,
				schema.NewColumn("newColumn", 100, types.StringKind, false, schema.NotNullConstraint{})),
			expectedRows: dtestutils.AddColToRows(t, AllPeopleRows, 100, types.String("default")),
		},
		{
			name:  "alter add column not null with expression default",
			query: "alter table people add (newColumn int not null default 2+2/2 comment 'tag:100')",
			expectedSchema: dtestutils.AddColumnToSchema(PeopleTestSchema,
				schema.NewColumn("newColumn", 100, types.IntKind, false, schema.NotNullConstraint{})),
			expectedRows: dtestutils.AddColToRows(t, AllPeopleRows, 100, types.Int(3)),
		},
		{
			name:  "alter add column not null with negative expression",
			query: "alter table people add (newColumn float not null default -1.1 comment 'tag:100')",
			expectedSchema: dtestutils.AddColumnToSchema(PeopleTestSchema,
				schema.NewColumn("newColumn", 100, types.FloatKind, false, schema.NotNullConstraint{})),
			expectedRows: dtestutils.AddColToRows(t, AllPeopleRows, 100, types.Float(-1.1)),
		},
		{
			name:        "alter add column not null with type mismatch in default",
			query:       "alter table people add (newColumn float default 'not a number' comment 'tag:100')",
			expectedErr: "Type mismatch",
		},
		{
			name:        "alter add column with tag conflict",
			query:       "alter table people add (newColumn float default 1.0 comment 'tag:1')",
			expectedErr: "A column with the tag 1 already exists",
		},
		{
			name:        "alter add column not null without default",
			query:       "alter table people add (newColumn varchar(80) not null comment 'tag:100')",
			expectedErr: "a default value must be provided",
		},
		{
			name:  "alter add column nullable",
			query: "alter table people add (newColumn bigint comment 'tag:100')",
			expectedSchema: dtestutils.AddColumnToSchema(PeopleTestSchema,
				schema.NewColumn("newColumn", 100, types.IntKind, false)),
			expectedRows: AllPeopleRows,
		},
		{
			name:  "alter add column with optional column keyword",
			query: "alter table people add column (newColumn varchar(80) comment 'tag:100')",
			expectedSchema: dtestutils.AddColumnToSchema(PeopleTestSchema,
				schema.NewColumn("newColumn", 100, types.StringKind, false)),
			expectedRows: AllPeopleRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			ctx := context.Background()
			root, _ := dEnv.WorkingRoot(ctx)

			sqlStatement, err := sqlparser.Parse(tt.query)
			require.NoError(t, err)

			s := sqlStatement.(*sqlparser.DDL)

			updatedRoot, err := ExecuteAlter(ctx, dEnv.DoltDB, root, s, tt.query)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NotNil(t, updatedRoot)
			table, _, err := updatedRoot.GetTable(ctx, PeopleTableName)
			assert.NoError(t, err)
			sch, err := table.GetSchema(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSchema, sch)

			updatedTable, ok, err := updatedRoot.GetTable(ctx, "people")
			assert.NoError(t, err)
			require.True(t, ok)

			rowData, err := updatedTable.GetRowData(ctx)
			assert.NoError(t, err)
			var foundRows []row.Row
			err = rowData.Iter(ctx, func(key, value types.Value) (stop bool, err error) {
				r, err := row.FromNoms(tt.expectedSchema, key.(types.Tuple), value.(types.Tuple))
				assert.NoError(t, err)
				foundRows = append(foundRows, r)
				return false, nil
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRows, foundRows)
		})
	}
}

func TestUnsupportedAlterStatements(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		expectedErr string
	}{
		{
			name:        "alter add index",
			query:       "alter table people add index myidx on (id, first)",
			expectedErr: "Unsupported",
		},
		{
			name:        "create index",
			query:       "create index myidx on people (id, first)",
			expectedErr: "Unsupported",
		},
		{
			name:        "alter drop index",
			query:       "alter table people drop index myidx",
			expectedErr: "Unsupported",
		},
		{
			name:        "drop index",
			query:       "drop index myidx on people",
			expectedErr: "Unsupported",
		},
		{
			name:        "alter change column",
			query:       "alter table people change id newId (varchar(80) not null)",
			expectedErr: "Unsupported",
		},
		{
			name:        "alter add foreign key",
			query:       "alter table appearances add constraint people_id_ref foreign key (id) references people (id)",
			expectedErr: "Unsupported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			ctx := context.Background()
			root, _ := dEnv.WorkingRoot(ctx)

			sqlStatement, err := sqlparser.Parse(tt.query)
			require.NoError(t, err)

			s := sqlStatement.(*sqlparser.DDL)

			_, err = ExecuteAlter(ctx, dEnv.DoltDB, root, s, tt.query)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestDropColumn(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedSchema schema.Schema
		expectedRows   []row.Row
		expectedErr    string
	}{
		{
			name:           "alter drop column",
			query:          "alter table people drop rating",
			expectedSchema: dtestutils.RemoveColumnFromSchema(PeopleTestSchema, RatingTag),
			expectedRows:   dtestutils.ConvertToSchema(dtestutils.RemoveColumnFromSchema(PeopleTestSchema, RatingTag), AllPeopleRows...),
		},
		{
			name:           "alter drop column with optional column keyword",
			query:          "alter table people drop column rating",
			expectedSchema: dtestutils.RemoveColumnFromSchema(PeopleTestSchema, RatingTag),
			expectedRows:   dtestutils.ConvertToSchema(dtestutils.RemoveColumnFromSchema(PeopleTestSchema, RatingTag), AllPeopleRows...),
		},
		{
			name:        "drop primary key",
			query:       "alter table people drop column id",
			expectedErr: "Cannot drop column in primary key",
		},
		{
			name:        "table not found",
			query:       "alter table notFound drop column id",
			expectedErr: "Unknown table: 'notFound'",
		},
		{
			name:        "column not found",
			query:       "alter table people drop column notFound",
			expectedErr: "Unknown column: 'notFound'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			ctx := context.Background()
			root, _ := dEnv.WorkingRoot(ctx)

			sqlStatement, err := sqlparser.Parse(tt.query)
			require.NoError(t, err)

			s := sqlStatement.(*sqlparser.DDL)

			updatedRoot, err := ExecuteAlter(ctx, dEnv.DoltDB, root, s, tt.query)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NotNil(t, updatedRoot)
			table, _, err := updatedRoot.GetTable(ctx, PeopleTableName)
			assert.NoError(t, err)
			sch, err := table.GetSchema(ctx)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSchema, sch)

			updatedTable, ok, err := updatedRoot.GetTable(ctx, "people")
			assert.NoError(t, err)
			require.True(t, ok)

			rowData, err := updatedTable.GetRowData(ctx)
			assert.NoError(t, err)
			var foundRows []row.Row
			err = rowData.Iter(ctx, func(key, value types.Value) (stop bool, err error) {
				updatedSch, err := updatedTable.GetSchema(ctx)
				assert.NoError(t, err)
				r, err := row.FromNoms(updatedSch, key.(types.Tuple), value.(types.Tuple))
				assert.NoError(t, err)
				foundRows = append(foundRows, r)
				return false, nil
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRows, foundRows)
		})
	}
}

func TestRenameColumn(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		expectedSchema schema.Schema
		expectedRows   []row.Row
		expectedErr    string
	}{
		{
			name:  "alter rename column with column and as keywords",
			query: "alter table people rename column rating as newRating",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", IdTag, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("first", FirstTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("last", LastTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("is_married", IsMarriedTag, types.BoolKind, false),
				schema.NewColumn("age", AgeTag, types.IntKind, false),
				schema.NewColumn("newRating", RatingTag, types.FloatKind, false),
				schema.NewColumn("uuid", UuidTag, types.UUIDKind, false),
				schema.NewColumn("num_episodes", NumEpisodesTag, types.UintKind, false),
			),
			expectedRows: AllPeopleRows,
		},
		{
			name:  "alter rename column with column and to keyword",
			query: "alter table people rename column rating to newRating",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", IdTag, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("first", FirstTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("last", LastTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("is_married", IsMarriedTag, types.BoolKind, false),
				schema.NewColumn("age", AgeTag, types.IntKind, false),
				schema.NewColumn("newRating", RatingTag, types.FloatKind, false),
				schema.NewColumn("uuid", UuidTag, types.UUIDKind, false),
				schema.NewColumn("num_episodes", NumEpisodesTag, types.UintKind, false),
			),
			expectedRows: AllPeopleRows,
		},
		{
			name:  "alter rename primary key column",
			query: "alter table people rename column id to newId",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("newId", IdTag, types.IntKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("first", FirstTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("last", LastTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("is_married", IsMarriedTag, types.BoolKind, false),
				schema.NewColumn("age", AgeTag, types.IntKind, false),
				schema.NewColumn("rating", RatingTag, types.FloatKind, false),
				schema.NewColumn("uuid", UuidTag, types.UUIDKind, false),
				schema.NewColumn("num_episodes", NumEpisodesTag, types.UintKind, false),
			),
			expectedRows: AllPeopleRows,
		},
		{
			name:        "table not found",
			query:       "alter table notFound rename column id to newId",
			expectedErr: "Unknown table: 'notFound'",
		},
		{
			name:        "column not found",
			query:       "alter table people rename column notFound to newNotFound",
			expectedErr: "Unknown column: 'notFound'",
		},
		{
			name:        "column name collision",
			query:       "alter table people rename column id to age",
			expectedErr: "A column with the name 'age' already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			ctx := context.Background()
			root, _ := dEnv.WorkingRoot(ctx)

			sqlStatement, err := sqlparser.Parse(tt.query)
			require.NoError(t, err)

			s := sqlStatement.(*sqlparser.DDL)

			updatedRoot, err := ExecuteAlter(ctx, dEnv.DoltDB, root, s, tt.query)

			if tt.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			require.NotNil(t, updatedRoot)
			table, _, err := updatedRoot.GetTable(ctx, PeopleTableName)
			assert.NoError(t, err)
			sch, err := table.GetSchema(ctx)
			assert.Equal(t, tt.expectedSchema, sch)

			updatedTable, ok, err := updatedRoot.GetTable(ctx, "people")
			assert.NoError(t, err)
			require.True(t, ok)

			rowData, err := updatedTable.GetRowData(ctx)
			assert.NoError(t, err)
			var foundRows []row.Row
			err = rowData.Iter(ctx, func(key, value types.Value) (stop bool, err error) {
				updatedSch, err := updatedTable.GetSchema(ctx)
				assert.NoError(t, err)
				r, err := row.FromNoms(updatedSch, key.(types.Tuple), value.(types.Tuple))
				assert.NoError(t, err)
				foundRows = append(foundRows, r)
				return false, nil
			})

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedRows, foundRows)
		})
	}
}

func TestRenameTable(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		oldTableName   string
		newTableName   string
		expectedSchema schema.Schema
		expectedRows   []row.Row
		expectedErr    string
	}{
		{
			name:           "alter rename table",
			query:          "rename table people to newPeople",
			oldTableName:   "people",
			newTableName:   "newPeople",
			expectedSchema: PeopleTestSchema,
			expectedRows:   AllPeopleRows,
		},
		{
			name:           "alter rename table with alter syntax",
			query:          "alter table people rename to newPeople",
			oldTableName:   "people",
			newTableName:   "newPeople",
			expectedSchema: PeopleTestSchema,
			expectedRows:   AllPeopleRows,
		},
		{
			name:           "rename multiple tables",
			query:          "rename table people to newPeople, appearances to newAppearances",
			oldTableName:   "appearances",
			newTableName:   "newAppearances",
			expectedSchema: AppearancesTestSchema,
			expectedRows:   AllAppsRows,
		},
		{
			name:        "table not found",
			query:       "rename table notFound to newNowFound",
			expectedErr: "Unknown table: 'notFound'",
		},
		{
			name:        "table name in use",
			query:       "rename table people to appearances",
			expectedErr: "A table with the name 'appearances' already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := dtestutils.CreateTestEnv()
			CreateTestDatabase(dEnv, t)
			ctx := context.Background()
			root, _ := dEnv.WorkingRoot(ctx)

			sqlStatement, err := sqlparser.Parse(tt.query)
			if err != nil {
				if len(tt.expectedErr) > 0 {
					require.Error(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)
					return
				} else {
					// fail the test
					require.NoError(t, err)
				}
			}

			s := sqlStatement.(*sqlparser.DDL)

			updatedRoot, err := ExecuteAlter(ctx, dEnv.DoltDB, root, s, tt.query)
			if len(tt.expectedErr) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			} else {
				require.NoError(t, err)
			}
			require.NotNil(t, updatedRoot)

			has, err := updatedRoot.HasTable(ctx, tt.oldTableName)
			require.NoError(t, err)
			assert.False(t, has)
			newTable, ok, err := updatedRoot.GetTable(ctx, tt.newTableName)
			require.NoError(t, err)
			require.True(t, ok)

			sch, err := newTable.GetSchema(ctx)
			require.NoError(t, err)
			require.Equal(t, tt.expectedSchema, sch)

			rowData, err := newTable.GetRowData(ctx)
			require.NoError(t, err)
			var foundRows []row.Row
			err = rowData.Iter(ctx, func(key, value types.Value) (stop bool, err error) {
				r, err := row.FromNoms(tt.expectedSchema, key.(types.Tuple), value.(types.Tuple))
				require.NoError(t, err)
				foundRows = append(foundRows, r)
				return false, nil
			})

			require.NoError(t, err)

			// Some test cases deal with rows declared in a different order than noms returns them, so use an order-
			// insensitive comparison here.
			assert.ElementsMatch(t, tt.expectedRows, foundRows)
		})
	}
}
