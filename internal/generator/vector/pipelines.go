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
			if input, ok := userDefined[inputRef]; ok {
				if len(input.ContainerLimitRef) > 0 {
					inputRef = fmt.Sprintf(`"throttle_%s"`, inputRef)

					if limit, ok := userDefinedLimits[input.ContainerLimitRef]; ok {
						t := Throttle{
							ComponentID: fmt.Sprintf(`throttle_%s`, input.Name),
							Inputs:      helpers.MakeInputs([]string{fmt.Sprintf(`application_routes.%s`, input.Name)}...),
							Threshold:   limit.MaxRecordsPerSecond.String(),
							KeyField:    perContainerLimitKeyField,
						}
						el = append(el, t)
					}

				} else if len(input.GroupLimitRef) > 0 {
					inputRef = fmt.Sprintf(`"throttle_%s"`, inputRef)

					if limit, ok := userDefinedLimits[input.GroupLimitRef]; ok {
						t := Throttle{
							ComponentID: fmt.Sprintf(`throttle_%s`, input.Name),
							Inputs:      helpers.MakeInputs([]string{fmt.Sprintf(`application_routes.%s`, input.Name)}...),
							Threshold:   limit.MaxRecordsPerSecond.String(),
						}
						el = append(el, t)
					}

				} else {
					inputRef = fmt.Sprintf(`"application_routes.%s"`, inputRef)
				}
			}

			modifiedInputRefs = append(modifiedInputRefs, inputRef)
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
