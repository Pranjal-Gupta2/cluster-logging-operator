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
