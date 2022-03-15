package vector

import (
	"encoding/json"
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/internal/generator"
	. "github.com/openshift/cluster-logging-operator/internal/generator/vector/elements"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/helpers"
)

const (
	perContainerLimitKeyField = `"{{ file }}"`
)

func MakeCustomInput(input *logging.InputSpec) string {
	if input.Application != nil {
		return fmt.Sprintf(`"application_routes.%s"`, input.Name)
	}

	if input.Infrastructure != nil {
		return fmt.Sprintf(`"infra_routes.%s"`, input.Name)

	}

	if input.Audit != nil {
		return fmt.Sprintf(`"audit_routes.%s"`, input.Name)
	}

	return ""
}

func InputsToPipelines(spec *logging.ClusterLogForwarderSpec, op generator.Options) []generator.Element {
	el := []generator.Element{}

	for _, p := range spec.Pipelines {
		vrl := SrcPassThrough
		if p.Labels != nil && len(p.Labels) != 0 {
			s, _ := json.Marshal(p.Labels)
			vrl = fmt.Sprintf(".openshift.labels = %s", s)
		}
		modifiedInputRefs := make([]string, 0)
		userDefinedInputs := spec.InputMap()
		userDefinedLimits := spec.LimitMap()

		for _, inputRef := range p.InputRefs {
			if !logging.ReservedInputNames.Has(inputRef) {
				if input, ok := userDefinedInputs[inputRef]; ok {
					if len(input.ContainerLimitRef) > 0 {
						if limit, ok := userDefinedLimits[input.ContainerLimitRef]; ok {
							t := Throttle{
								ComponentID: fmt.Sprintf(`throttle_%s`, input.Name),
								Inputs:      helpers.MakeInputs([]string{MakeCustomInput(input)}...),
								Threshold:   limit.MaxRecordsPerSecond.String(),
								KeyField:    perContainerLimitKeyField,
							}
							el = append(el, t)
							modifiedInputRefs = append(modifiedInputRefs, fmt.Sprintf(`throttle_%s`, inputRef))
						}

					} else if len(input.GroupLimitRef) > 0 {
						if limit, ok := userDefinedLimits[input.GroupLimitRef]; ok {
							t := Throttle{
								ComponentID: fmt.Sprintf(`throttle_%s`, input.Name),
								Inputs:      helpers.MakeInputs([]string{MakeCustomInput(input)}...),
								Threshold:   limit.MaxRecordsPerSecond.String(),
							}
							el = append(el, t)
							modifiedInputRefs = append(modifiedInputRefs, fmt.Sprintf(`throttle_%s`, inputRef))
						}
					} else {
						modifiedInputRefs = append(modifiedInputRefs, MakeCustomInput(input))
					}
				}
			} else {
				modifiedInputRefs = append(modifiedInputRefs, inputRef)
			}

		}

		r := Remap{
			ComponentID: p.Name,
			Inputs:      helpers.MakeInputs(modifiedInputRefs...),
			VRL:         vrl,
		}
		el = append(el, r)

	}

	return el
}
