package vector

import (
	"encoding/json"
	"fmt"
	"strings"

	logging "github.com/openshift/cluster-logging-operator/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/internal/generator"
	. "github.com/openshift/cluster-logging-operator/internal/generator/vector/elements"
	"github.com/openshift/cluster-logging-operator/internal/generator/vector/helpers"
)

const (
	ParseJson = "json"
)

var (
	UserDefinedInput          = fmt.Sprintf("%s.%%s", RouteApplicationLogs)
	perContainerLimitKeyField = `"{{ file }}"`
)

func MakeCustomInput(input *logging.InputSpec) string {
	if input.Application != nil {
		return fmt.Sprintf(`"route_application_logs.%s"`, input.Name)
	}

	if input.Infrastructure != nil {
		return fmt.Sprintf(`"infra_routes.%s"`, input.Name)

	}

	if input.Audit != nil {
		return fmt.Sprintf(`"audit_routes.%s"`, input.Name)
	}

	return ""
}

func Pipelines(spec *logging.ClusterLogForwarderSpec, op generator.Options) []generator.Element {
	el := []generator.Element{}
	userDefined := spec.InputMap()

	for _, p := range spec.Pipelines {
		vrls := []string{}
		if p.Labels != nil && len(p.Labels) != 0 {
			s, _ := json.Marshal(p.Labels)
			vrls = append(vrls, fmt.Sprintf(".openshift.labels = %s", s))
		}
		if p.Parse == ParseJson {
			parse := `
parsed, err = parse_json(.message)
if err == null {
  .structured = parsed
}
`
			vrls = append(vrls, parse)
		}
		modifiedInputRefs := make([]string, 0)
		userDefinedLimits := spec.LimitMap()

		for _, inputRef := range p.InputRefs {
			if !logging.ReservedInputNames.Has(inputRef) {
				if input, ok := userDefined[inputRef]; ok {
					if len(input.ContainerLimitRef) > 0 {
						if limit, ok := userDefinedLimits[input.ContainerLimitRef]; ok {
							t := Throttle{
								ComponentID: fmt.Sprintf(`throttle_%s`, input.Name),
								Inputs:      helpers.MakeInputs([]string{MakeCustomInput(input)}...),
								Threshold:   limit.MaxRecordsPerSecond.AsDec().String(),
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
								Threshold:   limit.MaxRecordsPerSecond.AsDec().String(),
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
		vrl := SrcPassThrough
		if len(vrls) != 0 {
			vrl = strings.Join(helpers.TrimSpaces(vrls), "\n\n")
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
