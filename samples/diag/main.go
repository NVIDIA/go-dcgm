package main

import (
	"html/template"
	"log"
	"os"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

const diagOutput = `Software:
{{range $t := .Software}}
  {{printf "%-50s" $t.TestName}} {{$t.Status}}	{{$t.TestOutput}}
{{- end}}
{{range $g := .PerGpu}}
GPU	: {{$g.GPU}}
  {{range $t := $g.DiagResults}}
  {{printf "%-20s" $t.TestName}} {{$t.Status}}	{{$t.TestOutput}}
  {{- end}}
{{- end}}
`

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

	t := template.Must(template.New("Diag").Parse(diagOutput))
	if err = t.Execute(os.Stdout, dr); err != nil {
		log.Panicln("Template error:", err)
	}
}
