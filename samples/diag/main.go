package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

const diagOutput = `DCGM {{.DCGMVersion}}	Driver {{.DriverVersion}}
{{range .Entities -}}
Entity {{entity .Entity}}	Serial {{.SerialNumber}}	SKU {{.SKUDeviceID}}
{{end -}}
{{range .SystemErrors -}}
System error	code={{.Code}}	category={{.Category}}	severity={{.Severity}}	{{.Message}}
{{end -}}
{{range .Tests -}}
[{{.Category}}] {{.Name}}	{{.Status}}	Plugin {{.PluginName}}
{{range .Results}}  Result {{entity .Entity}}	{{.Status}}
{{end -}}
{{range .Errors}}  Error {{entity .Entity}}	code={{.Code}}	category={{.Category}}	severity={{.Severity}}	{{.Message}}
{{end -}}
{{range .Info}}  Info {{entity .Entity}}	{{.Message}}
{{end -}}
{{if .AuxData}}  Auxiliary data	version={{.AuxDataVersion}}	{{.AuxData}}
{{end -}}
{{end -}}
`

func entityLabel(entity dcgm.GroupEntityPair) string {
	if entity.EntityGroupId == dcgm.FE_NONE {
		return "global"
	}
	return fmt.Sprintf("%s %d", entity.EntityGroupId, entity.EntityId)
}

func newDiagOutputTemplate() *template.Template {
	return template.Must(template.New("Diag").Funcs(template.FuncMap{"entity": entityLabel}).Parse(diagOutput))
}

func main() {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		log.Panicln(err)
	}
	defer cleanup()

	dr, err := dcgm.RunDiag(dcgm.DiagQuick, dcgm.GroupAllGPUs())
	if err != nil {
		log.Panicln(err)
	}

	t := newDiagOutputTemplate()
	if err = t.Execute(os.Stdout, dr); err != nil {
		log.Panicln("Template error:", err)
	}
}
