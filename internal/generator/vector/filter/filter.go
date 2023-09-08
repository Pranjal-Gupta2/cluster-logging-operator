package filter

import (
	"fmt"

	loggingv1 "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/internal/generator"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/filter/apiaudit"
)

// RemapVRL returns a VRL expression to add to the remap program of a pipeline containing this filter.
func RemapVRL(filterSpec *loggingv1.FilterSpec, spec *loggingv1.ClusterLogForwarderSpec, op generator.Options) (string, error) {
	types := generator.GatherSources(spec, op)
	switch filterSpec.Type {

	case loggingv1.FilterAPIAudit:
		if types.Has(loggingv1.InputNameAudit) {
			return apiaudit.PolicyToVRL(filterSpec.APIAudit)
		}

	default:
		return "", fmt.Errorf("unknown filter type: %v", filterSpec.Type)
	}
	return "", nil // Nothing to do
}
