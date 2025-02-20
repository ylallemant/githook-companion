package dependency

import (
	"bytes"
	"html/template"
	"runtime"

	"github.com/pkg/errors"
	"github.com/ylallemant/githook-companion/pkg/api"
)

func templateDataFromDependency(dependency *api.Dependency) *archiveContext {
	ctx := new(archiveContext)

	ctx.Version = dependency.Version
	ctx.Os = runtime.GOOS
	ctx.Arch = runtime.GOARCH
	ctx.Ext = "tar.gz"

	if ctx.Os == "windows" {
		ctx.Ext = "zip"
	}

	return ctx
}

func renderUriTempate(urlTemplate string, data *archiveContext) (string, error) {
	var err error

	tmpl := template.New("url")
	if tmpl, err = tmpl.Parse(urlTemplate); err != nil {
		return "", errors.Wrapf(err, "failed to register url template \"%s\"", urlTemplate)
	}

	var render bytes.Buffer
	err = tmpl.Execute(&render, data)
	if err != nil {
		return "", errors.Wrapf(err, "failed to render url template \"%s\"", urlTemplate)
	}

	return render.String(), nil
}
