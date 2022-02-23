package output

import (
	"github.com/openshift/cluster-logging-operator/internal/generator"
)

type FileSink = generator.ConfLiteral

const FileSinkTemplate = `
{{define "outputSinkFileTemplate" -}}
# {{.Desc}}
[sinks.{{.ComponentID}}]
type = "file"
inputs = {{.InLabel}}
path = "var/log/containers/stress.log"
encoding.codec = "ndjson"
{{end}}
`
