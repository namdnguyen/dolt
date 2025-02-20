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

package tblcmds

import (
	"context"

	"github.com/fatih/color"

	"github.com/liquidata-inc/dolt/go/cmd/dolt/cli"
	"github.com/liquidata-inc/dolt/go/cmd/dolt/commands"
	"github.com/liquidata-inc/dolt/go/cmd/dolt/errhand"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/doltdb"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/env"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/row"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/rowconv"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/schema"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/pipeline"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/typed"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/typed/noms"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/untyped"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/untyped/fwt"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/untyped/nullprinter"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/table/untyped/tabular"
	"github.com/liquidata-inc/dolt/go/libraries/utils/argparser"
	"github.com/liquidata-inc/dolt/go/libraries/utils/iohelp"
	"github.com/liquidata-inc/dolt/go/store/types"
)

const (
	whereParam        = "where"
	limitParam        = "limit"
	hideConflictsFlag = "hide-conflicts"
	defaultLimit      = -1
	cnfColName        = "Cnf"
)

var fwtStageName = "fwt"

var cnfTag = schema.ReservedTagMin

var selShortDesc = "print a selection of a table"
var selLongDesc = `The dolt table select command selects rows from a table and prints out some or all of the table's columns`
var selSynopsis = []string{
	"[--limit <record_count>] [--where <col1=val1>] [--hide-conflicts] [<commit>] <table> [<column>...]",
}

type SelectArgs struct {
	tblName       string
	colNames      []string
	whereClause   string
	limit         int
	hideConflicts bool
}

func Select(ctx context.Context, commandStr string, args []string, dEnv *env.DoltEnv) int {
	ap := newArgParser()
	help, usage := cli.HelpAndUsagePrinters(commandStr, selShortDesc, selLongDesc, selSynopsis, ap)
	apr := cli.ParseArgs(ap, args, help)
	args = apr.Args()

	if len(args) == 0 {
		usage()
		return 1
	}

	root, verr := commands.GetWorkingWithVErr(dEnv)

	if verr == nil {
		var cm *doltdb.Commit
		cm, verr = commands.MaybeGetCommitWithVErr(dEnv, args[0])

		if verr == nil {
			if cm != nil {
				args = args[1:]

				var err error
				root, err = cm.GetRootValue()

				if err != nil {
					cli.PrintErrln(color.RedString("error: failed to get root value: " + err.Error()))
					return 1
				}
			}

			if len(args) == 0 {
				cli.Println("No tables specified")
				usage()
				return 1
			}

			tblName := args[0]

			var colNames []string
			if len(args) > 1 {
				colNames = args[1:]
			}

			selArgs := &SelectArgs{
				tblName,
				colNames,
				apr.GetValueOrDefault(whereParam, ""),
				apr.GetIntOrDefault(limitParam, defaultLimit),
				apr.Contains(hideConflictsFlag)}

			verr = printTable(ctx, root, selArgs)
		}
	}

	if verr != nil {
		cli.PrintErrln(verr.Verbose())
		return 1
	}

	return 0
}

func newArgParser() *argparser.ArgParser {
	ap := argparser.NewArgParser()
	ap.ArgListHelp["table"] = "List of tables to be printed."
	ap.ArgListHelp["column"] = "List of columns to be printed"
	ap.SupportsString(whereParam, "", "column", "")
	ap.SupportsInt(limitParam, "", "record_count", "")
	ap.SupportsFlag(hideConflictsFlag, "", "")
	return ap
}

// Runs the selection pipeline and prints the table of resultant values, returning any error encountered.
func printTable(ctx context.Context, root *doltdb.RootValue, selArgs *SelectArgs) errhand.VerboseError {
	var verr errhand.VerboseError
	if has, err := root.HasTable(ctx, selArgs.tblName); err != nil {
		return errhand.BuildDError("error: failed to read tables").AddCause(err).Build()
	} else if !has {
		return errhand.BuildDError("error: unknown table '%s'", selArgs.tblName).Build()
	}

	tbl, _, err := root.GetTable(ctx, selArgs.tblName)

	if err != nil {
		return errhand.BuildDError("error: failed to read to get table '%s'", selArgs.tblName).AddCause(err).Build()
	}

	tblSch, err := tbl.GetSchema(ctx)

	if err != nil {
		return errhand.BuildDError("error: failed to get schema").AddCause(err).Build()
	}

	whereFn, err := commands.ParseWhere(tblSch, selArgs.whereClause)

	if err != nil {
		return errhand.BuildDError("error: failed to parse where cause").AddCause(err).SetPrintUsage().Build()
	}

	selTrans := commands.NewSelTrans(whereFn, selArgs.limit)
	transforms := pipeline.NewTransformCollection(pipeline.NewNamedTransform("select", selTrans.LimitAndFilter))
	sch, err := maybeAddCnfColTransform(ctx, transforms, tbl, tblSch)

	if err != nil {
		return errhand.BuildDError("error: failed to setup pipeline").AddCause(err).Build()
	}

	outSch, verr := addMapTransform(selArgs, sch, transforms)

	if verr != nil {
		return verr
	}

	p, err := createPipeline(ctx, tbl, tblSch, outSch, transforms)

	if err != nil {
		return errhand.BuildDError("error: failed to setup pipeline").AddCause(err).Build()
	}

	selTrans.Pipeline = p

	p.Start()
	err = p.Wait()

	if err != nil {
		return errhand.BuildDError("error: error processing results").AddCause(err).Build()
	}

	return nil
}

// Creates a pipeline to select and print rows from the table given. Adds a fixed-width printing transform to the
// collection of transformations given.
func createPipeline(ctx context.Context, tbl *doltdb.Table, tblSch schema.Schema, outSch schema.Schema, transforms *pipeline.TransformCollection) (*pipeline.Pipeline, error) {
	colNames, err := schema.ExtractAllColNames(outSch)

	if err != nil {
		return nil, err
	}

	addSizingTransform(outSch, transforms)

	rowData, err := tbl.GetRowData(ctx)

	if err != nil {
		return nil, err
	}

	rd, err := noms.NewNomsMapReader(ctx, rowData, tblSch)

	if err != nil {
		return nil, err
	}

	wr, err := tabular.NewTextTableWriter(iohelp.NopWrCloser(cli.CliOut), outSch)

	if err != nil {
		return nil, err
	}

	badRowCallback := func(tff *pipeline.TransformRowFailure) (quit bool) {
		cli.PrintErrln(color.RedString("error: failed to transform row %s.", row.Fmt(ctx, tff.Row, outSch)))
		return true
	}

	rdProcFunc := pipeline.ProcFuncForReader(ctx, rd)
	wrProcFunc := pipeline.ProcFuncForWriter(ctx, wr)

	p := pipeline.NewAsyncPipeline(rdProcFunc, wrProcFunc, transforms, badRowCallback)
	p.RunAfter(func() { rd.Close(ctx) })
	p.RunAfter(func() { wr.Close(ctx) })

	// Insert the table header row at the appropriate stage
	r, err := untyped.NewRowFromTaggedStrings(tbl.Format(), outSch, colNames)

	if err != nil {
		return nil, err
	}

	p.InjectRow(fwtStageName, r)

	return p, nil
}

func addSizingTransform(outSch schema.Schema, transforms *pipeline.TransformCollection) {
	nullPrinter := nullprinter.NewNullPrinter(outSch)
	transforms.AppendTransforms(pipeline.NewNamedTransform(nullprinter.NULL_PRINTING_STAGE, nullPrinter.ProcessRow))

	autoSizeTransform := fwt.NewAutoSizingFWTTransformer(outSch, fwt.PrintAllWhenTooLong, 10000)
	transforms.AppendTransforms(pipeline.NamedTransform{Name: fwtStageName, Func: autoSizeTransform.TransformToFWT})
}

func addMapTransform(selArgs *SelectArgs, sch schema.Schema, transforms *pipeline.TransformCollection) (schema.Schema, errhand.VerboseError) {
	colColl := sch.GetAllCols()

	if len(selArgs.colNames) > 0 {
		cols := make([]schema.Column, 0, len(selArgs.colNames)+1)

		if !selArgs.hideConflicts {
			if col, ok := sch.GetAllCols().GetByName(cnfColName); ok {
				cols = append(cols, col)
			}
		}

		for _, name := range selArgs.colNames {
			if col, ok := colColl.GetByName(name); !ok {
				return nil, errhand.BuildDError("error: unknown column '%s'", name).Build()
			} else {
				cols = append(cols, col)
			}
		}

		colColl, _ = schema.NewColCollection(cols...)
	}

	outSch := schema.UnkeyedSchemaFromCols(colColl)
	untypedSch, err := untyped.UntypeUnkeySchema(outSch)

	if err != nil {
		return nil, errhand.BuildDError("error: failed to create untyped schema").AddCause(err).Build()
	}

	mapping, err := rowconv.TagMapping(sch, untypedSch)

	if err != nil {
		return nil, errhand.BuildDError("error: failed to create mapping").AddCause(err).Build()
	}

	rConv, _ := rowconv.NewRowConverter(mapping)
	transform := pipeline.NewNamedTransform("map", rowconv.GetRowConvTransformFunc(rConv))
	transforms.AppendTransforms(transform)

	return mapping.DestSch, nil
}

func maybeAddCnfColTransform(ctx context.Context, transColl *pipeline.TransformCollection, tbl *doltdb.Table, tblSch schema.Schema) (schema.Schema, error) {
	if has, err := tbl.HasConflicts(); err != nil {
		return nil, err
	} else if has {
		// this is so much code to add a column
		const transCnfSetName = "set cnf col"

		_, confSchema := untyped.NewUntypedSchemaWithFirstTag(cnfTag, cnfColName)
		schWithConf, _ := typed.TypedSchemaUnion(confSchema, tblSch)

		_, confData, _ := tbl.GetConflicts(ctx)

		cnfTransform := pipeline.NewNamedTransform(transCnfSetName, CnfTransformer(ctx, tblSch, schWithConf, confData))
		transColl.AppendTransforms(cnfTransform)

		return schWithConf, nil
	}

	return tblSch, nil
}

var confLabel = types.String(" ! ")
var noConfLabel = types.String("   ")

func CnfTransformer(ctx context.Context, inSch, outSch schema.Schema, conflicts types.Map) func(inRow row.Row, props pipeline.ReadableMap) (rowData []*pipeline.TransformedRowResult, badRowDetails string) {
	return func(inRow row.Row, props pipeline.ReadableMap) ([]*pipeline.TransformedRowResult, string) {
		key := inRow.NomsMapKey(inSch)

		v, err := key.Value(ctx)

		if err != nil {
			panic(err)
		}

		if has, err := conflicts.Has(ctx, v); err != nil {
			panic(err)
		} else if has {
			inRow, err = inRow.SetColVal(cnfTag, confLabel, outSch)

			if err != nil {
				panic(err)
			}
		} else {
			inRow, err = inRow.SetColVal(cnfTag, noConfLabel, outSch)

			if err != nil {
				panic(err)
			}
		}

		return []*pipeline.TransformedRowResult{{RowData: inRow, PropertyUpdates: nil}}, ""
	}
}
