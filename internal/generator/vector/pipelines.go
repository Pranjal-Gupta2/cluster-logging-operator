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
	UserDefinedInput = fmt.Sprintf("%s.%%s", RouteApplicationLogs)
)

func AddThrottle(spec *logging.ClusterLogForwarderSpec, op generator.Options) []generator.Element {
	el := []generator.Element{}
	userDefinedLimits := spec.LimitMap()

	for _, inputSpec := range spec.Inputs {
		if len(inputSpec.LimitRef) > 0 {
			if limit, ok := userDefinedLimits[inputSpec.LimitRef]; ok {
				t := Throttle{
					ComponentID: fmt.Sprintf(`throttle_%s`, inputSpec.Name),
					Inputs:      helpers.MakeInputs([]string{fmt.Sprintf(`application_routes.%s`, inputSpec.Name)}...),
					Threshold:   limit.MaxBytesPerSecond.String(),
				}
				el = append(el, t)
			}
		}
	}

	return el
}

func Pipelines(spec *logging.ClusterLogForwarderSpec, op generator.Options) []generator.Element {
	el := []generator.Element{}
	userDefined := spec.InputMap()
	el = append(el, AddThrottle(spec, op)...)

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
		for _, inputRef := range p.InputRefs {
			if input, ok := userDefined[inputRef]; ok {
				if len(input.LimitRef) > 0 {
					inputRef = fmt.Sprintf(`"throttle_%s"`, inputRef)
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
			// Inputs:      helpers.MakeInputs(p.InputRefs...),
			Inputs: helpers.MakeInputs(modifiedInputRefs...),
			VRL:    vrl,
		}
		el = append(el, r)

	}

	return el
}
