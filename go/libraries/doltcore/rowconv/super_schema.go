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

package rowconv

import (
	"context"
	"errors"
	"fmt"

	"github.com/liquidata-inc/dolt/go/libraries/utils/set"

	"github.com/liquidata-inc/dolt/go/libraries/doltcore/doltdb"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/schema"
	"github.com/liquidata-inc/dolt/go/store/hash"
	"github.com/liquidata-inc/dolt/go/store/types"
)

// TagKindPair is a simple tuple that holds a tag and a NomsKind of a column
type TagKindPair struct {
	// Tag is the tag of a column
	Tag uint64

	// Kind is the NomsKind of a colum
	Kind types.NomsKind
}

// NameKindPair is a simple tuple that holds the name of a column and it's NomsKind
type NameKindPair struct {
	// Name is the name of the column
	Name string

	// Kind is the NomsKind of the column
	Kind types.NomsKind
}

// SuperSchema is an immutable schema generated by a SuperSchemaGen which defines methods for getting the schema
// and mapping another schema onto the super schema
type SuperSchema struct {
	sch       schema.Schema
	namedCols map[TagKindPair]string
}

// GetSchema gets the underlying schema.Schema object
func (ss SuperSchema) GetSchema() schema.Schema {
	if ss.sch == nil {
		panic("Bug: super schema not generated.")
	}

	return ss.sch
}

// RowConvForSchema creates a RowConverter for transforming rows with the the given schema to this super schema.
// This is done by mapping the column tag and type to the super schema column representing that tag and type.
func (ss SuperSchema) RowConvForSchema(sch schema.Schema) (*RowConverter, error) {
	inNameToOutName := make(map[string]string)
	allCols := sch.GetAllCols()
	err := allCols.Iter(func(tag uint64, col schema.Column) (stop bool, err error) {
		tkp := TagKindPair{Tag: tag, Kind: col.Kind}
		outName, ok := ss.namedCols[tkp]

		if !ok {
			return false, errors.New("failed to map columns")
		}

		inNameToOutName[col.Name] = outName
		return false, nil
	})

	if err != nil {
		return nil, err
	}

	fm, err := NewFieldMappingFromNameMap(sch, ss.sch, inNameToOutName)

	if err != nil {
		return nil, err
	}

	return NewRowConverter(fm)
}

// SuperSchemaGen is a utility class used to generate the superset of several schemas.
type SuperSchemaGen struct {
	tagKindToDestTag map[TagKindPair]uint64
	usedTags         map[uint64]struct{}
	names            map[TagKindPair]*set.StrSet
}

// NewSuperSchemaGen creates a new SuperSchemaGen
func NewSuperSchemaGen() *SuperSchemaGen {
	return &SuperSchemaGen{
		tagKindToDestTag: make(map[TagKindPair]uint64),
		usedTags:         make(map[uint64]struct{}),
		names:            make(map[TagKindPair]*set.StrSet),
	}
}

// AddSchema will add a schema which will be incorporated into the superset of schemas
func (ssg *SuperSchemaGen) AddSchema(sch schema.Schema) error {
	err := sch.GetAllCols().Iter(func(tag uint64, col schema.Column) (stop bool, err error) {
		tagKind := TagKindPair{Tag: tag, Kind: col.Kind}
		_, exists := ssg.tagKindToDestTag[tagKind]

		if !exists {
			destTag := tag

			for {
				_, collides := ssg.usedTags[destTag]
				if !collides {
					ssg.tagKindToDestTag[tagKind] = destTag
					ssg.usedTags[destTag] = struct{}{}
					ssg.names[tagKind] = set.NewStrSet([]string{col.Name})
					return false, nil
				}

				if destTag == tag {
					destTag = schema.ReservedTagMin
				} else {
					destTag++
				}
			}
		} else {
			ssg.names[tagKind].Add(col.Name)
		}

		return false, nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (ssg *SuperSchemaGen) nameCols() map[TagKindPair]string {
	colNames := make(map[string][]TagKindPair)
	for tagKind, names := range ssg.names {
		name := fmt.Sprintf("%d_%s", tagKind.Tag, tagKind.Kind.String())
		if names.Size() == 1 {
			name = names.AsSlice()[0]
		}
		colNames[name] = append(colNames[name], tagKind)
	}

	results := make(map[TagKindPair]string)
	for name, tagKinds := range colNames {
		if len(tagKinds) == 1 {
			results[tagKinds[0]] = name
			continue
		}

		for _, tagKind := range tagKinds {
			name := fmt.Sprintf("%s_%s_%d", name, tagKind.Kind.String(), tagKind.Tag)
			results[tagKind] = name
		}
	}

	return results
}

// GenerateSuperSchema takes all the accumulated schemas and generates a schema which is the superset of all of
// those schemas.
func (ssg *SuperSchemaGen) GenerateSuperSchema(additionalCols ...NameKindPair) (SuperSchema, error) {
	namedCols := ssg.nameCols()

	colColl, _ := schema.NewColCollection()
	for tagKind, colName := range namedCols {
		destTag, ok := ssg.tagKindToDestTag[tagKind]

		if !ok {
			panic("mismatch between namedCols and tagKindToDestTag")
		}

		col := schema.NewColumn(colName, destTag, tagKind.Kind, false)

		var err error
		colColl, err = colColl.Append(col)

		if err != nil {
			return SuperSchema{}, err
		}
	}

	if len(additionalCols) > 0 {
		nextReserved := schema.ReservedTagMin

		for _, nameKindPair := range additionalCols {
			if _, ok := colColl.GetByName(nameKindPair.Name); ok {
				return SuperSchema{}, errors.New("Additional column name collision: " + nameKindPair.Name)
			}

			for {
				if _, ok := ssg.usedTags[nextReserved]; !ok {
					break
				}
				nextReserved++
			}

			var err error
			ssg.usedTags[nextReserved] = struct{}{}
			colColl, err = colColl.Append(schema.NewColumn(nameKindPair.Name, nextReserved, nameKindPair.Kind, false))

			if err != nil {
				return SuperSchema{}, err
			}
		}
	}

	sch := schema.UnkeyedSchemaFromCols(colColl)
	return SuperSchema{sch: sch, namedCols: namedCols}, nil
}

// AddHistoryOfTableAtCommit will traverse a commit graph adding all versions of a tables schema to the schemas being
// supersetted.
func (ssg *SuperSchemaGen) AddHistoryOfTableAtCommit(ctx context.Context, tblName string, ddb *doltdb.DoltDB, cm *doltdb.Commit) error {
	addedSchemas := make(map[hash.Hash]struct{})
	processedCommits := make(map[hash.Hash]struct{})
	return ssg.addHistoryOfTableAtCommit(ctx, tblName, addedSchemas, processedCommits, ddb, cm)
}

func (ssg *SuperSchemaGen) addHistoryOfTableAtCommit(ctx context.Context, tblName string, addedSchemas, processedCommits map[hash.Hash]struct{}, ddb *doltdb.DoltDB, cm *doltdb.Commit) error {
	cmHash, err := cm.HashOf()

	if err != nil {
		return err
	}

	if _, ok := processedCommits[cmHash]; ok {
		return nil
	}

	processedCommits[cmHash] = struct{}{}

	root, err := cm.GetRootValue()

	if err != nil {
		return err
	}

	tbl, ok, err := root.GetTable(ctx, tblName)

	if err != nil {
		return err
	}

	if ok {
		schRef, err := tbl.GetSchemaRef()

		if err != nil {
			return err
		}

		h := schRef.TargetHash()

		if _, ok = addedSchemas[h]; !ok {
			sch, err := tbl.GetSchema(ctx)

			if err != nil {
				return err
			}

			err = ssg.AddSchema(sch)

			if err != nil {
				return err
			}
		}
	}

	numParents, err := cm.NumParents()

	if err != nil {
		return err
	}

	for i := 0; i < numParents; i++ {
		cm, err := ddb.ResolveParent(ctx, cm, i)

		if err != nil {
			return err
		}

		err = ssg.addHistoryOfTableAtCommit(ctx, tblName, addedSchemas, processedCommits, ddb, cm)

		if err != nil {
			return err
		}
	}

	return nil
}

// AddHistoryOfTable will traverse all commit graphs which have local branches associated with them and add all
// passed versions of a table's schema to the schemas being supersetted
func (ssg *SuperSchemaGen) AddHistoryOfTable(ctx context.Context, tblName string, ddb *doltdb.DoltDB) error {
	refs, err := ddb.GetRefs(ctx)

	if err != nil {
		return err
	}

	addedSchemas := make(map[hash.Hash]struct{})
	processedCommits := make(map[hash.Hash]struct{})

	for _, ref := range refs {
		cs, err := doltdb.NewCommitSpec("HEAD", ref.String())

		if err != nil {
			return err
		}

		cm, err := ddb.Resolve(ctx, cs)

		if err != nil {
			return err
		}

		err = ssg.addHistoryOfTableAtCommit(ctx, tblName, addedSchemas, processedCommits, ddb, cm)

		if err != nil {
			return err
		}
	}

	return nil
}
