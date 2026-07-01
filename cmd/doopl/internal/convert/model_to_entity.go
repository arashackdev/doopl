// Package convert holds the goverter converter interface for turning public
// model types into cmd/doopl's CLI-only entity types. Mirrors
// internal/convert at the library root — same pattern, one layer further
// out. Run `task generate` after editing the interface below.
//
//go:generate go run github.com/jmattheis/goverter/cmd/goverter gen .
package convert

import (
	"github.com/arashackdev/doopl/cmd/doopl/internal/entity"
	"github.com/arashackdev/doopl/pkg/model"
)

// ModelToEntity converts public model types into cmd/doopl's CLI-only
// entity types.
//
// goverter:converter
// goverter:output:file ./converter_gen.go
// goverter:output:package github.com/arashackdev/doopl/cmd/doopl/internal/convert
type ModelToEntity interface {
	// goverter:map DetectedSourceLang DetectedSourceLang
	TranslationRow(source model.TextResult) entity.TranslationRow

	TranslationRows(source []model.TextResult) []entity.TranslationRow

	LanguageRow(source model.Language) entity.LanguageRow

	LanguageRows(source []model.Language) []entity.LanguageRow

	UsageRow(source model.Usage) entity.UsageRow
}
