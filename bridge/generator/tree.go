package main

import (
	"go/ast"
)

type Param struct {
	Name  string
	Type  string
	CType string
}

// An accessible representation of a function signature used for templating.
type Signature struct {
	Name    string
	Params  []*Param
	Results []*Param

	ShouldDeref bool

	// various template conveniences

	LastResultIndex int
	LastParamIndex  int

	ParamDataName string
	ParamDataType string

	ReturnDataName string
}

// This is a full struct instead of a type alias for text-templating reasons
type ParsedSignatures struct {
	Signatures []*Signature
}

func isStringExpr(ex ast.Expr) bool {
	ident, ok := ex.(*ast.Ident)
	return ok && ident.Name == "string"
}

func isDataExpr(ex ast.Expr) (bool, string, bool) {

	// Returns if ex is a data expr, the name of the expr, and if the value should be dereferenced

	// might be an array of pointers; if so, unwrap
	array, isArray := ex.(*ast.ArrayType)
	if isArray {
		ex = array.Elt
	}

	// Might be an array of strings; if so, handle
	// Should replace with more generic "json-serializable complex primitive" logic
	// Which... is also the set of values that should be dereferenced
	if isArray && isStringExpr(ex) {
		return true, "[]string", false
	}

	// Otherwise, it must be a pointer to a struct; unwrap
	pointer, ok := ex.(*ast.StarExpr)
	if !ok {
		return false, "", true
	}

	// Grab the pointer ident
	ident, ok := pointer.X.(*ast.Ident)
	if !ok {
		return false, "", true
	}
	name := ident.Name

	// Whitelist; could replace with lexing later
	whitelist := []string{"Acquisition", "Batch", "BatchProposal", "Collection", "Client", "Config", "ContainerReference", "DeletedResponse", "Error", "FileReference", "Formula", "FormulaResult", "Gear", "GearDoc", "GearSource", "Group", "IdResponse", "Input", "Job", "JobLog", "JobLogStatement", "Key", "ModifiedResponse", "Note", "Origin", "Output", "Permission", "ProgressReader", "Project", "Result", "Session", "Subject", "Target", "UploadResponse", "UploadSource", "User", "Version"}

	if stringInSlice(name, whitelist) {
		return true, "api." + name, true
	} else {
		return false, name, true
	}
}

func isHttpRespExpr(ex ast.Expr) bool {
	pointer, ok := ex.(*ast.StarExpr)
	if !ok {
		return false
	}

	selector, ok := pointer.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "http" && selector.Sel.Name == "Response"
}

func isErrorExpr(ex ast.Expr) bool {
	ident, ok := ex.(*ast.Ident)
	return ok && ident.Name == "error"
}
