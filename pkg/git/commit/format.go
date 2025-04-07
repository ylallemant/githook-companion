package commit

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
	"github.com/ylallemant/githook-companion/pkg/api"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

var tokenReferenceRegexp = regexp.MustCompile(fmt.Sprintf(
	`\.(%s)`,
	nlpapi.TokenNameRegexp.String(),
))

func TokenNameTemplateFormat(name string) string {
	return strcase.ToCamel(name)
}

func TokenNameStructFormat(name string) string {
	return strcase.ToSnake(name)
}

func templateReferencedTokens(template string) []string {
	matches := tokenReferenceRegexp.FindAllStringSubmatch(template, -1)

	references := make([]string, 0)
	for _, match := range matches {
		reference := TokenNameStructFormat(match[1])

		if !slices.Contains(references, reference) && reference != api.CommitMessageKey {
			references = append(references, reference)
		}
	}

	return references
}

func cleanRawMessage(message string, tokenNameReferences []string, tokens []*nlpapi.Token) string {
	for _, tokenNameReference := range tokenNameReferences {
		for _, token := range tokens {
			if token.Name == tokenNameReference && token.Source == nlpapi.TokenSourceLexeme {
				message = strings.ReplaceAll(message, token.Word.Raw, "")
			} else if token.Name == tokenNameReference && token.Source == nlpapi.TokenSourceLexemeComposite {
				message = strings.ReplaceAll(message, token.Word.FromComposite, "")
			}
		}
	}

	return message
}

func dynamicTemplateStruct(tokenNameReferences []string) reflect.Type {
	fields := make([]reflect.StructField, 0)

	for _, tokenName := range tokenNameReferences {
		fields = append(fields, reflect.StructField{
			Name: TokenNameTemplateFormat(tokenName),
			Type: reflect.TypeOf(""),
		})
	}

	fields = append(fields, reflect.StructField{
		Name: TokenNameTemplateFormat(api.CommitMessageKey),
		Type: reflect.TypeOf(""),
	})

	dynamicStruct := reflect.StructOf(fields)
	log.Debug().Msgf("dynamic template struct %+v)", dynamicStruct)
	return dynamicStruct
}

func dynamicTemplateData(dynamicStruct reflect.Type, message, commitTypeTokenSourceName string, tokenNameReferences []string, tokens []*nlpapi.Token) reflect.Value {
	instance := reflect.New(dynamicStruct)
	log.Debug().Msgf("generate dynamic template struct instance")

	for _, tokenNameReference := range tokenNameReferences {
		log.Debug().Msgf("- search token reference \"%+v\"", tokenNameReference)
		for _, token := range tokens {
			if token.Name == tokenNameReference {

				if token.Name != api.CommitTypeTokenName || (token.SourceName == commitTypeTokenSourceName && token.Name == api.CommitTypeTokenName) {
					log.Debug().Msgf("  match token \"%+v\": %+v", token.Name, token)
					log.Debug().Msgf("  add value \"%+v\"", token.Value)
					instance.Elem().FieldByName(TokenNameTemplateFormat(token.Name)).Set(reflect.ValueOf(token.Value))
				}
			}
		}
	}

	instance.Elem().FieldByName(TokenNameTemplateFormat(api.CommitMessageKey)).Set(reflect.ValueOf(message))

	log.Debug().Msgf("dynamic template struct instance %+v", instance)
	return instance
}
