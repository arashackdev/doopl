// Package convert holds goverter converter interface definitions. The
// generated implementations (converter_gen.go) are checked in but should
// never be hand-edited — run `task generate` after changing a converter
// interface here.
package convert

import (
	"github.com/arashackdev/doopl/internal/apimodel"
	"github.com/arashackdev/doopl/pkg/model"
)

// APIToModel converts wire-format apimodel types into the public model
// types. One converter per resource area (translate, document, glossary,
// ...) keeps each generated file small and the diff readable when the API
// shape changes.
//
// goverter:converter
// goverter:output:file ./converter_gen.go
// goverter:output:package github.com/arashackdev/doopl/internal/convert
type APIToModel interface {
	// Translation maps a single apimodel.Translation onto model.TextResult.
	// Field names match 1:1 except DetectedSourceLang, which goverter maps
	// automatically despite the json tag difference because we map on Go
	// field names, not json tags.
	Translation(source apimodel.Translation) model.TextResult

	// Translations converts a full response's Translations slice in one call.
	Translations(source []apimodel.Translation) []model.TextResult

	// Language converts an apimodel.Language to model.Language.
	Language(source apimodel.Language) model.Language

	// Languages converts a slice of apimodel.Language to model.Language.
	Languages(source []apimodel.Language) []model.Language

	// Usage converts apimodel.UsageResponse to model.Usage.
	Usage(source apimodel.UsageResponse) model.Usage

	DocumentStatus(source apimodel.DocumentStatusResponse) model.DocumentStatusInfo

	Glossary(source apimodel.GlossaryResponse) model.Glossary

	Glossaries(source []apimodel.GlossaryResponse) []model.Glossary
}
