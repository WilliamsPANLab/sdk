package main

import (
	// "encoding/json"
	"errors"
	. "fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// GenerateSignatures gets a bunch of ast.FuncDecls and scans them into a struct for templating.
func GenerateSignatures(path string) (*token.FileSet, *ParsedSignatures) {

	fset, funcs := GetRelevantFunctionSignatures(path)

	sigs := []*Signature{}

	for _, fn := range funcs {
		// DebugPrintFuncSig(fn)
		// ast.Print(fset, fn)

		name := fn.Name.Name

		signature := &Signature{
			Name:           name,
			Params:         []*Param{},
			Results:        []*Param{},
			ParamDataName:  "",
			ReturnDataName: "nil",
			ShouldDeref:    true,
		}

		parameters := fn.Type.Params.List
		for _, paramSet := range parameters {

			pt := "unknown"
			cgpt := "unknown"
			ct := "unknown"

			dExpr, dname, shouldDeref := isDataExpr(paramSet.Type)

			// Check param type
			if isStringExpr(paramSet.Type) {
				pt = "string"
				cgpt = "*C.char"
				ct = "char*"
			} else if dExpr {
				pt = "data"
				cgpt = "*C.char"
				ct = "char*"

				if !shouldDeref {
					signature.ShouldDeref = false
				}
			} else {
				Println("Function", name, "has an unknown parameter:")
				ast.Print(fset, paramSet.Type)
				Println("---")
				Println()
			}

			for _, param := range paramSet.Names {
				signature.Params = append(signature.Params, &Param{
					Name:    param.Name,
					Type:    pt,
					CgoType: cgpt,
					CType:   ct,
				})

				if pt == "data" {
					signature.ParamDataName = param.Name
					signature.ParamDataType = dname
				}
			}
		}
		signature.LastParamIndex = len(signature.Params) - 1

		results := fn.Type.Results.List
		for _, field := range results {
			pt := "unknown"
			name := "unknown"

			dExpr, _, _ := isDataExpr(field.Type)

			if isStringExpr(field.Type) {
				pt = "string"
				name = "data"
				signature.ReturnDataName = "data"
			} else if dExpr {
				pt = "data"
				name = "data"
				signature.ReturnDataName = "data"
			} else if isHttpRespExpr(field.Type) {
				pt = "http"
				name = "_"
			} else if isErrorExpr(field.Type) {
				pt = "error"
				name = "err"
			} else {

				ast.Print(fset, field.Type)
			}

			// Println(pt, name)

			signature.Results = append(signature.Results, &Param{
				Name: name,
				Type: pt,
			})
		}
		signature.LastResultIndex = len(signature.Results) - 1

		// x, _ := json.MarshalIndent(signature, "", "\t")
		// Println(string(x))

		sigs = append(sigs, signature)
	}

	psigs := &ParsedSignatures{
		Signatures: sigs,
	}

	// sort.Sort(psigs) // alphasort by function name

	return fset, psigs
}

// Scan & parse a package of source code and return the set of function signatures we want to compile.
func GetRelevantFunctionSignatures(path string) (*token.FileSet, []*ast.FuncDecl) {
	fset := token.NewFileSet()

	// Wall-clock warning: discover and parse source code from disk
	packages, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	check(err)

	apiPackage, ok := packages["api"]
	if !ok {
		check(errors.New("Package API not found. Wrong directory path?"))
	}

	funcs := []*ast.FuncDecl{}

	// Inspect the AST, finding all function declarations and filtering
	ast.Inspect(apiPackage, func(n ast.Node) bool {

		// Filter on function declarations
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// Ignore closures
		if fn.Name == nil || fn.Type.Params == nil || fn.Type.Results == nil {
			return true
		}

		// Ignore functions that are not attached to a receiver
		if fn.Recv == nil || len(fn.Recv.List) != 1 {
			return true
		}

		// Check that the function receiver is *api.Client
		pointer, ok := fn.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return true
		}
		identifier, ok := pointer.X.(*ast.Ident)
		if !ok {
			return true
		}
		if identifier.Name != "Client" {
			return true
		}

		name := fn.Name.Name

		// Ignore unexported functions
		if strings.ToLower(string(name[0])) == string(name[0]) {
			return true
		}

		// Troublesome functions that need to either be fixed or ignored
		blacklist := []string{
			// boolean parameter
			"ModifyJob",

			// variadic string array
			"StartNextPendingJob",

			// JobState enum
			"ChangeJobState",

			// map string -> interface
			"ProposeBatch",

			// api.JobLogStatement instead of []*JobLogStatement
			"AddJobLogs",

			// Doesn't detect int return
			"CancelBatch",

			// Doesn't detect map[string]interface{} return
			"GetGearInvocation",

			// Progress reporting
			"Upload",
			"UploadSimple",
			"UploadToProject",
			"UploadToSession",
			"UploadToAcquisition",
			"UploadToCollection",
			"UploadToAnalysis",
			"Download",
			"DownloadSimple",
			"DownloadFromProject",
			"DownloadFromSession",
			"DownloadFromAcquisition",
			"DownloadFromCollection",
			"DownloadFromAnalysis",
			"DownloadInputFromAnalysis",
		}
		if stringInSlice(name, blacklist) {
			return false
		}

		// Println(name)

		funcs = append(funcs, fn)
		return false // don't need to go deeper; we found what we're looking for.
	})

	return fset, funcs
}

func DetectSDKVersion(path string) string {

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	check(err)

	for _, decl := range file.Decls {

		// Filter on general declarations
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range gen.Specs {
			// Filter on value specifications
			value, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			// Version decl only has one name and one value
			if len(value.Names) != 1 || len(value.Values) != 1 {
				continue
			}

			if value.Names[0].Name == "Version" {
				escapedVersion := value.Values[0].(*ast.BasicLit).Value

				// AST includes quotations; remove those so we can have a 'raw' string value
				return strings.Replace(escapedVersion, "\"", "", -1)
			}
		}
	}

	Println("Warning: Could not detect SDK version.")
	return "unknown"
}
