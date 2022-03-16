package vector

import (
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/internal/generator"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/helpers"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/output"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/output/kafka"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/output/loki"

	corev1 "k8s.io/api/core/v1"
)

const (
	SinkForOtherApplicationLogs = "other"
)

func Outputs(clspec *logging.ClusterLoggingSpec, secrets map[string]*corev1.Secret, clfspec *logging.ClusterLogForwarderSpec, op generator.Options) []generator.Element {
	outputs := []generator.Element{}
	route := RouteMap(clfspec, op)

	for _, o := range clfspec.Outputs {
		secret := secrets[o.Name]
		inputs := route[o.Name].List()
		switch o.Type {
		case logging.OutputTypeKafka:
			outputs = generator.MergeElements(outputs, kafka.Conf(o, inputs, secret, op))
		case logging.OutputTypeLoki:
			outputs = generator.MergeElements(outputs, loki.Conf(o, inputs, secret, op))
		case logging.OutputTypeElasticsearch:
			// outputs = generator.MergeElements(outputs, elasticsearch.Conf(o, inputs, secret, op))
			outputs = generator.MergeElements(outputs, []generator.Element{
				output.File{
					ComponentID: o.Name,
					Desc:        "File sink for storing logs",
					Inputs:      helpers.MakeInputs(inputs...),
					Path:        `"/var/log/containers/stress.log"`,
				},
			})
		}
	}

	outputs = generator.MergeElements(outputs, []generator.Element{
		output.InternalMetricsSource{
			Desc:         "Source for generating vector's internal metrics",
			TemplateName: "inputSourceInternalMetricsTemplate",
			TemplateStr:  output.InternalMetricsSourceTemplate,
		},

		output.PromSink{
			Desc:         "Sink for exporting Prometheus metrics",
			TemplateName: "outputSinkPrometheusTemplate",
			TemplateStr:  output.PrometheusSinkTemplate,
		},
	})

	if DefaultRouteIsPresent {
		outputs = generator.MergeElements(outputs, []generator.Element{
			output.File{
				ComponentID: "other",
				Desc:        "File sink for storing un-routed application logs",
				Inputs:      helpers.MakeInputs(fmt.Sprintf("application_routes.%s", DefaultApplicationRoute)),
				Path:        `"/var/log/containers/other.log"`,
			}},
		)
	}

	return outputs
}
