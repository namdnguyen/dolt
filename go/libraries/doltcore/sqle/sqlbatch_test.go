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
	"context"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	sqle "github.com/src-d/go-mysql-server"
	"github.com/src-d/go-mysql-server/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/liquidata-inc/dolt/go/libraries/doltcore/dtestutils"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/row"
	. "github.com/liquidata-inc/dolt/go/libraries/doltcore/sql/sqltestutil"
	"github.com/liquidata-inc/dolt/go/store/types"
)

func TestSqlBatchInserts(t *testing.T) {
	insertStatements := []string{
		`insert into people (id, first, last, is_married, age, rating, uuid, num_episodes) values
					(7, "Maggie", "Simpson", false, 1, 5.1, '00000000-0000-0000-0000-000000000007', 677)`,
		`insert into people values
					(8, "Milhouse", "VanHouten", false, 1, 5.1, '00000000-0000-0000-0000-000000000008', 677)`,
		`insert into people (id, first, last) values (9, "Clancey", "Wiggum")`,
		`insert into people (id, first, last) values
					(10, "Montgomery", "Burns"), (11, "Ned", "Flanders")`,
		`insert into episodes (id, name) values (5, "Bart the General"), (6, "Moaning Lisa")`,
		`insert into episodes (id, name) values (7, "The Call of the Simpsons"), (8, "The Telltale Head")`,
		`insert into episodes (id, name) values (9, "Life on the Fast Lane")`,
		`insert into appearances (character_id, episode_id) values (7,5), (7,6)`,
		`insert into appearances (character_id, episode_id) values (8,7)`,
		`insert into appearances (character_id, episode_id) values (9,8), (9,9)`,
		`insert into appearances (character_id, episode_id) values (10,5), (10,6)`,
		`insert into appearances (character_id, episode_id) values (11,9)`,
	}

	// Shuffle the inserts so that different tables are interleaved. We're not giving a seed here, so this is
	// deterministic.
	rand.Shuffle(len(insertStatements),
		func(i, j int) {
			insertStatements[i], insertStatements[j] = insertStatements[j], insertStatements[i]
		})

	dEnv := dtestutils.CreateTestEnv()
	ctx := context.Background()

	CreateTestDatabase(dEnv, t)
	root, _ := dEnv.WorkingRoot(ctx)

	engine := sqle.NewDefault()
	db := NewBatchedDatabase("dolt", root, dEnv)
	engine.AddDatabase(db)

	for _, stmt := range insertStatements {
		_, rowIter, err := engine.Query(sql.NewEmptyContext(), stmt)
		require.NoError(t, err)
		require.NoError(t, drainIter(rowIter))
	}

	// Before committing the batch, the database should be unchanged from its original state
	allPeopleRows, err := GetAllRows(root, PeopleTableName)
	require.NoError(t, err)
	allEpsRows, err := GetAllRows(root, EpisodesTableName)
	require.NoError(t, err)
	allAppearanceRows, err := GetAllRows(root, AppearancesTableName)
	require.NoError(t, err)

	assert.ElementsMatch(t, AllPeopleRows, allPeopleRows)
	assert.ElementsMatch(t, AllEpsRows, allEpsRows)
	assert.ElementsMatch(t, AllAppsRows, allAppearanceRows)

	// Now commit the batch and check for new rows
	err = db.Flush(ctx)
	require.NoError(t, err)

	var expectedPeople, expectedEpisodes, expectedAppearances []row.Row

	expectedPeople = append(expectedPeople, AllPeopleRows...)
	expectedPeople = append(expectedPeople,
		NewPeopleRowWithOptionalFields(7, "Maggie", "Simpson", false, 1, 5.1, uuid.MustParse("00000000-0000-0000-0000-000000000007"), 677),
		NewPeopleRowWithOptionalFields(8, "Milhouse", "VanHouten", false, 1, 5.1, uuid.MustParse("00000000-0000-0000-0000-000000000008"), 677),
		newPeopleRow(9, "Clancey", "Wiggum"),
		newPeopleRow(10, "Montgomery", "Burns"),
		newPeopleRow(11, "Ned", "Flanders"),
	)

	expectedEpisodes = append(expectedEpisodes, AllEpsRows...)
	expectedEpisodes = append(expectedEpisodes,
		newEpsRow(5, "Bart the General"),
		newEpsRow(6, "Moaning Lisa"),
		newEpsRow(7, "The Call of the Simpsons"),
		newEpsRow(8, "The Telltale Head"),
		newEpsRow(9, "Life on the Fast Lane"),
	)

	expectedAppearances = append(expectedAppearances, AllAppsRows...)
	expectedAppearances = append(expectedAppearances,
		newAppsRow(7, 5),
		newAppsRow(7, 6),
		newAppsRow(8, 7),
		newAppsRow(9, 8),
		newAppsRow(9, 9),
		newAppsRow(10, 5),
		newAppsRow(10, 6),
		newAppsRow(11, 9),
	)

	root = db.Root()
	allPeopleRows, err = GetAllRows(root, PeopleTableName)
	require.NoError(t, err)
	allEpsRows, err = GetAllRows(root, EpisodesTableName)
	require.NoError(t, err)
	allAppearanceRows, err = GetAllRows(root, AppearancesTableName)
	require.NoError(t, err)

	assertRowSetsEqual(t, expectedPeople, allPeopleRows)
	assertRowSetsEqual(t, expectedEpisodes, allEpsRows)
	assertRowSetsEqual(t, expectedAppearances, allAppearanceRows)
}

func drainIter(iter sql.RowIter) error {
	var returnedErr error
	for {
		_, err := iter.Next()
		if err == io.EOF {
			return returnedErr
		} else if err != nil {
			returnedErr = err
		}
	}
}

func TestSqlBatchInsertIgnoreReplace(t *testing.T) {
	t.Skip("Skipped until insert ignore statements supported in go-mysql-server")

	insertStatements := []string{
		`replace into people (id, first, last, is_married, age, rating, uuid, num_episodes) values
					(0, "Maggie", "Simpson", false, 1, 5.1, '00000000-0000-0000-0000-000000000007', 677)`,
		`insert ignore into people values
					(2, "Milhouse", "VanHouten", false, 1, 5.1, '00000000-0000-0000-0000-000000000008', 677)`,
	}

	dEnv := dtestutils.CreateTestEnv()
	ctx := context.Background()

	CreateTestDatabase(dEnv, t)
	root, _ := dEnv.WorkingRoot(ctx)

	engine := sqle.NewDefault()
	db := NewBatchedDatabase("dolt", root, dEnv)
	engine.AddDatabase(db)

	for _, stmt := range insertStatements {
		_, rowIter, err := engine.Query(sql.NewEmptyContext(), stmt)
		require.NoError(t, err)
		drainIter(rowIter)
	}

	// Before committing the batch, the database should be unchanged from its original state
	allPeopleRows, err := GetAllRows(root, PeopleTableName)
	assert.NoError(t, err)
	assert.ElementsMatch(t, AllPeopleRows, allPeopleRows)

	// Now commit the batch and check for new rows
	err = db.Flush(ctx)
	require.NoError(t, err)

	var expectedPeople []row.Row

	expectedPeople = append(expectedPeople, AllPeopleRows[1:]...) // skip homer
	expectedPeople = append(expectedPeople,
		NewPeopleRowWithOptionalFields(0, "Maggie", "Simpson", false, 1, 5.1, uuid.MustParse("00000000-0000-0000-0000-000000000007"), 677),
	)

	allPeopleRows, err = GetAllRows(root, PeopleTableName)
	assert.NoError(t, err)
	assertRowSetsEqual(t, expectedPeople, allPeopleRows)
}

func TestSqlBatchInsertErrors(t *testing.T) {
	dEnv := dtestutils.CreateTestEnv()
	ctx := context.Background()

	CreateTestDatabase(dEnv, t)
	root, _ := dEnv.WorkingRoot(ctx)

	engine := sqle.NewDefault()
	db := NewBatchedDatabase("dolt", root, dEnv)
	engine.AddDatabase(db)

	_, rowIter, err := engine.Query(sql.NewEmptyContext(), `insert into people (id, first, last, is_married, age, rating, uuid, num_episodes) values
					(0, "Maggie", "Simpson", false, 1, 5.1, '00000000-0000-0000-0000-000000000007', 677)`)
	// This won't generate an error until we commit the batch (duplicate key)
	assert.NoError(t, err)
	assert.NoError(t, drainIter(rowIter))

	// This generates an error at insert time because of the bad type for the uuid column
	_, _, err = engine.Query(sql.NewEmptyContext(), `insert into people values
					(2, "Milhouse", "VanHouten", false, 1, 5.1, true, 677)`)
	assert.Error(t, err)

	// Error from the first statement appears here
	assert.Error(t, db.Flush(ctx))
}

func assertRowSetsEqual(t *testing.T, expected, actual []row.Row) {
	equal, diff := rowSetsEqual(expected, actual)
	assert.True(t, equal, diff)
}

// Returns whether the two slices of rows contain the same elements using set semantics (no duplicates), and an error
// string if they aren't.
func rowSetsEqual(expected, actual []row.Row) (bool, string) {
	if len(expected) != len(actual) {
		return false, fmt.Sprintf("Sets have different sizes: expected %d, was %d", len(expected), len(actual))
	}

	for _, ex := range expected {
		if !containsRow(actual, ex) {
			return false, fmt.Sprintf("Missing row: %v", ex)
		}
	}

	return true, ""
}

func containsRow(rs []row.Row, r row.Row) bool {
	for _, r2 := range rs {
		equal, _ := rowsEqual(r, r2)
		if equal {
			return true
		}
	}
	return false
}

func rowsEqual(expected, actual row.Row) (bool, string) {
	er, ar := make(map[uint64]types.Value), make(map[uint64]types.Value)
	_, err := expected.IterCols(func(t uint64, v types.Value) (bool, error) {
		er[t] = v
		return false, nil
	})

	if err != nil {
		panic(err)
	}

	_, err = actual.IterCols(func(t uint64, v types.Value) (bool, error) {
		ar[t] = v
		return false, nil
	})

	if err != nil {
		panic(err)
	}

	opts := cmp.Options{cmp.AllowUnexported(), dtestutils.FloatComparer}
	eq := cmp.Equal(er, ar, opts)
	var diff string
	if !eq {
		diff = cmp.Diff(er, ar, opts)
	}
	return eq, diff
}

func newPeopleRow(id int, first, last string) row.Row {
	vals := row.TaggedValues{
		IdTag:    types.Int(id),
		FirstTag: types.String(first),
		LastTag:  types.String(last),
	}

	r, err := row.New(types.Format_7_18, PeopleTestSchema, vals)

	if err != nil {
		panic(err)
	}

	return r
}

func newEpsRow(id int, name string) row.Row {
	vals := row.TaggedValues{
		EpisodeIdTag: types.Int(id),
		EpNameTag:    types.String(name),
	}

	r, err := row.New(types.Format_7_18, EpisodesTestSchema, vals)

	if err != nil {
		panic(err)
	}

	return r
}

func newAppsRow(charId int, epId int) row.Row {
	vals := row.TaggedValues{
		AppCharacterTag: types.Int(charId),
		AppEpTag:        types.Int(epId),
	}

	r, err := row.New(types.Format_7_18, AppearancesTestSchema, vals)

	if err != nil {
		panic(err)
	}

	return r
}
