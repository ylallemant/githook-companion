package commit

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	nlpapi "github.com/ylallemant/githook-companion/pkg/nlp/api"
)

const (
	commitTypePrefixRegexpFmt = "^(?i)%s\\s*:{0,1}\\s*"
)

var (
	whitespaceRegexp = regexp.MustCompile(`\s+`)
)

func EnsureFormat(message, tpl string, commitTypeToken *nlpapi.Token, tokens []*nlpapi.Token) (string, error) {
	binaryTemplate := []byte(tpl)

	hasher := sha256.New()
	hasher.Write(binaryTemplate)

	templateName := hex.EncodeToString(hasher.Sum(nil))

	renderer := template.New(templateName).Funcs(sprig.FuncMap())

	var err error
	if renderer, err = renderer.Parse(string(binaryTemplate)); err != nil {
		return "", errors.Wrapf(err, "failed to parse template \"%s\"", string(binaryTemplate))
	}

	referencedTokens := templateReferencedTokens(tpl)
	cleanedMessage := cleanRawMessage(message, referencedTokens, tokens)
	templateStruct := dynamicTemplateStruct(referencedTokens)
	templateData := dynamicTemplateData(templateStruct, cleanedMessage, commitTypeToken.SourceName, referencedTokens, tokens)

	formatted := bytes.NewBuffer([]byte{})
	err = renderer.Execute(formatted, templateData)
	if err != nil {
		return "", errors.Wrap(err, "failed to render template")
	}

	result := strings.TrimSpace(formatted.String())
	result = whitespaceRegexp.ReplaceAllString(result, " ")

	return result, nil
}
