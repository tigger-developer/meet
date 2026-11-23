// ABOUTME: Canonical CLI help embedded directly from the tracked text files.
// ABOUTME: Shared by meet and meet-helper without generated staging copies.

// Package helptext supplies the canonical CLI help prose.
package helptext

import _ "embed"

// Meet is the help text for meet.
//
//go:embed meet.txt
var Meet string

// Serve is the help text for meet-serve.
//
//go:embed meet-serve.txt
var Serve string

// Token is the help text for meet-token.
//
//go:embed meet-token.txt
var Token string

// Create is the help text for meet-create.
//
//go:embed meet-create.txt
var Create string

// Cancel is the help text for meet-cancel.
//
//go:embed meet-cancel.txt
var Cancel string

// List is the help text for meet-list.
//
//go:embed meet-list.txt
var List string

// Helper is the help text for meet-helper.
//
//go:embed meet-helper.txt
var Helper string
