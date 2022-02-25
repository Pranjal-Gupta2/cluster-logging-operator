package output

import (
	"github.com/openshift/cluster-logging-operator/internal/generator"
)

type File struct {
	ComponentID string
	Desc        string
	Inputs      string
	Path        string
}

func (f File) Name() string {
	return "FileTemplate"
}

func (f File) Template() string {
	return `
{{define "FileTemplate" -}}
{{- if .Desc}}
# {{.Desc}}
{{- end}}
[sinks.{{.ComponentID}}]
type = "file"
inputs = {{.Inputs}}
path = {{.Path}}
encoding.codec = "ndjson"
{{end}}
`
}

type InternalMetricsSource = generator.ConfLiteral

const InternalMetricsSourceTemplate = `
{{define "inputSourceInternalMetricsTemplate" -}}
# {{.Desc}}
[sources.internal_metrics]
type = "internal_metrics"
{{end}}
`

type PromSink = generator.ConfLiteral

const PrometheusSinkTemplate = `
{{define "outputSinkPrometheusTemplate" -}}
# {{.Desc}}
[sinks.prom_exporter]
type = "prometheus"
inputs = ["internal_metrics"]
address = "0.0.0.0:24231"
{{end}}
`
