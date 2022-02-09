package vector

import (
	"fmt"
	"strings"

	logging "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/internal/generator"
	. "github.com/openshift/cluster-logging-operator/internal/generator/vector/elements"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/helpers"
)

const (
	IsInfraContainer = `starts_with!(.kubernetes.pod_namespace,"kube") || starts_with!(.kubernetes.pod_namespace,"openshift") || .kubernetes.pod_namespace == "default"`

	SrcPassThrough = "."
)

var (
	AddLogTypeApp          = fmt.Sprintf(".log_type = %q", logging.InputNameApplication)
	AddLogTypeInfra        = fmt.Sprintf(".log_type = %q", logging.InputNameInfrastructure)
	AddLogTypeAudit        = fmt.Sprintf(".log_type = %q", logging.InputNameAudit)
	InfraContainerLogsExpr = fmt.Sprintf(`'%s'`, IsInfraContainer)
	AppContainerLogsExpr   = fmt.Sprintf(`'!(%s)'`, IsInfraContainer)
	InputContainerLogs     = "container_logs"
	InputJournalLogs       = "journal_logs"
)

// SourcesToInputs takes the raw log sources (container, journal, audit) and produces Inputs as defined by ClusterLogForwarder Api
func SourcesToInputs(spec *logging.ClusterLogForwarderSpec, o generator.Options) []generator.Element {
	el := []generator.Element{}

	types := generator.GatherSources(spec, o)
	// route container_logs based on type
	if types.Has(logging.InputNameApplication) || types.Has(logging.InputNameInfrastructure) {
		r := Route{
			ComponentID: "route_container_logs",
			Inputs:      helpers.MakeInputs(InputContainerLogs),
			Routes:      map[string]string{},
		}
		if types.Has(logging.InputNameApplication) {
			r.Routes["app"] = AppContainerLogsExpr
		}
		if types.Has(logging.InputNameInfrastructure) {
			r.Routes["infra"] = InfraContainerLogsExpr
		}
		//TODO Add handling of user-defined inputs
		el = append(el, r)
	}

	if types.Has(logging.InputNameApplication) {
		r := Remap{
			Desc:        `Rename log stream to "application"`,
			ComponentID: "application",
			Inputs:      helpers.MakeInputs("route_container_logs.app"),
			VRL:         AddLogTypeApp,
		}
		el = append(el, r)

		applicationRoute := Route{
			Desc:        `Add custom user inputs from CLF Spec`,
			ComponentID: "application_routes",
			Inputs:      helpers.MakeInputs("application"),
			Routes:      map[string]string{},
		}

		userDefined := spec.InputMap()
		for _, pipeline := range spec.Pipelines {
			for _, inRef := range pipeline.InputRefs {
				if input, ok := userDefined[inRef]; ok {
					if input.Application != nil {
						namespaces := make([]string, 0)
						for _, ns := range input.Application.Namespaces {
							namespaces = append(namespaces, fmt.Sprintf(`.kubernetes.pod_namespace == "%s"`, ns))
						}
						applicationRoute.Routes[input.Name] = fmt.Sprintf(`'%s'`, strings.Join(namespaces, " && "))
					}
				}
			}
		}
		el = append(el, applicationRoute)

	}
	if types.Has(logging.InputNameInfrastructure) {
		r := Remap{
			Desc:        `Rename log stream to "infrastructure"`,
			ComponentID: "infrastructure",
			Inputs:      helpers.MakeInputs("route_container_logs.infra", InputJournalLogs),
			VRL:         AddLogTypeInfra,
		}
		el = append(el, r)
	}
	if types.Has(logging.InputNameAudit) {
		r := Remap{
			Desc:        `Rename log stream to "audit"`,
			ComponentID: "audit",
			Inputs:      helpers.MakeInputs("host_audit_logs", "k8s_audit_logs", "openshift_audit_logs"),
			VRL:         AddLogTypeAudit,
		}
		el = append(el, r)
	}
	//TODO add user defined routing

	return el
}
