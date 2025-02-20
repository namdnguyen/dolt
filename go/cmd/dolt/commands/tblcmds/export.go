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
	"os"

	"github.com/fatih/color"

	"github.com/liquidata-inc/dolt/go/cmd/dolt/cli"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/doltdb"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/env"
	"github.com/liquidata-inc/dolt/go/libraries/doltcore/mvdata"
	"github.com/liquidata-inc/dolt/go/libraries/utils/argparser"
	"github.com/liquidata-inc/dolt/go/libraries/utils/iohelp"
)

var exportShortDesc = `Export the contents of a table to a file.`
var exportLongDesc = `dolt table export will export the contents of <table> to <file>

See the help for <b>dolt table import</b> as the options are the same.`
var exportSynopsis = []string{
	"[-f] [-pk <field>] [-schema <file>] [-map <file>] [-continue] [-file-type <type>] <table> <file>",
}

// validateExportArgs validates the input from the arg parser, and returns the tuple:
// (table name to export, data location of table to export, data location to export to)
func validateExportArgs(apr *argparser.ArgParseResults, usage cli.UsagePrinter) (string, mvdata.TableDataLocation, mvdata.DataLocation) {
	if apr.NArg() == 0 || apr.NArg() > 2 {
		usage()
		return "", mvdata.TableDataLocation{}, nil
	}

	tableName := apr.Arg(0)
	if !doltdb.IsValidTableName(tableName) {
		cli.PrintErrln(
			color.RedString("'%s' is not a valid table name\n", tableName),
			"table names must match the regular expression:", doltdb.TableNameRegexStr)
		return "", mvdata.TableDataLocation{}, nil
	}

	path := ""
	if apr.NArg() > 1 {
		path = apr.Arg(1)
	}

	fType, _ := apr.GetValue(fileTypeParam)
	destLoc := mvdata.NewDataLocation(path, fType)

	switch val := destLoc.(type) {
	case mvdata.FileDataLocation:
		if val.Format == mvdata.InvalidDataFormat {
			cli.PrintErrln(
				color.RedString("Could not infer type file '%s'\n", path),
				"File extensions should match supported file types, or should be explicitly defined via the file-type parameter")
			return "", mvdata.TableDataLocation{}, nil
		}

	case mvdata.StreamDataLocation:
		if val.Format == mvdata.InvalidDataFormat {
			val = mvdata.StreamDataLocation{Format: mvdata.CsvFile, Reader: os.Stdin, Writer: iohelp.NopWrCloser(cli.CliOut)}
			destLoc = val
		} else if val.Format != mvdata.CsvFile && val.Format != mvdata.PsvFile {
			cli.PrintErrln(color.RedString("Cannot export this format to stdout"))
			return "", mvdata.TableDataLocation{}, nil
		}
	}

	tableLoc := mvdata.TableDataLocation{Name: tableName}

	return tableName, tableLoc, destLoc
}

func parseExportArgs(commandStr string, args []string) (bool, *mvdata.MoveOptions) {
	ap := argparser.NewArgParser()
	ap.ArgListHelp["table"] = "The table being exported."
	ap.ArgListHelp["file"] = "The file being output to."
	ap.SupportsFlag(forceParam, "f", "If data already exists in the destination, the Force flag will allow the target to be overwritten.")
	ap.SupportsFlag(contOnErrParam, "", "Continue exporting when row export errors are encountered.")
	ap.SupportsString(outSchemaParam, "s", "schema_file", "The schema for the output data.")
	ap.SupportsString(mappingFileParam, "m", "mapping_file", "A file that lays out how fields should be mapped from input data to output data.")
	ap.SupportsString(primaryKeyParam, "pk", "primary_key", "Explicitly define the name of the field in the schema which should be used as the primary key.")
	ap.SupportsString(fileTypeParam, "", "file_type", "Explicitly define the type of the file if it can't be inferred from the file extension.")

	help, usage := cli.HelpAndUsagePrinters(commandStr, exportShortDesc, exportLongDesc, exportSynopsis, ap)
	apr := cli.ParseArgs(ap, args, help)
	tableName, tableLoc, fileLoc := validateExportArgs(apr, usage)

	if fileLoc == nil || len(tableLoc.Name) == 0 {
		return false, nil
	}

	schemaFile, _ := apr.GetValue(outSchemaParam)
	mappingFile, _ := apr.GetValue(mappingFileParam)
	primaryKey, _ := apr.GetValue(primaryKeyParam)

	return apr.Contains(forceParam), &mvdata.MoveOptions{
		Operation:   mvdata.OverwriteOp,
		ContOnErr:   apr.Contains(contOnErrParam),
		TableName:   tableName,
		SchFile:     schemaFile,
		MappingFile: mappingFile,
		PrimaryKey:  primaryKey,
		Src:         tableLoc,
		Dest:        fileLoc,
	}
}

func Export(ctx context.Context, commandStr string, args []string, dEnv *env.DoltEnv) int {
	force, mvOpts := parseExportArgs(commandStr, args)

	if mvOpts == nil {
		return 1
	}

	result := executeMove(ctx, dEnv, force, mvOpts)

	if result == 0 {
		cli.PrintErrln(color.CyanString("Successfully exported data."))
	}

	return result
}
