package elements

type Route struct {
	ComponentID string
	Desc        string
	Inputs      string
	Routes      map[string]string
}

func (r Route) Name() string {
	return "routeTemplate"
}

func (r Route) Template() string {
	return `
{{define "routeTemplate" -}}
{{- if .Desc}}
# {{.Desc}}
{{- end}}
[transforms.{{.ComponentID}}]
type = "route"
inputs = {{.Inputs}}
{{- range $route_name, $route_expr := .Routes}}
route.{{$route_name}} = {{$route_expr}}
{{- end}}
{{end}}
`
}

type Remap struct {
	ComponentID string
	Desc        string
	Inputs      string
	VRL         string
}

func (r Remap) Name() string {
	return "remapTemplate"
}

func (r Remap) Template() string {
	return `
{{define "remapTemplate" -}}
{{- if .Desc}}
# {{.Desc}}
{{- end}}
[transforms.{{.ComponentID}}]
type = "remap"
inputs = {{.Inputs}}
source = """
{{.VRL}}
"""
{{end}}
`
}

type Throttle struct {
	ComponentID string
	Desc        string
	Inputs      string
	Threshold   string
}

func (t Throttle) Name() string {
	return "throttleTemplate"
}

func (t Throttle) Template() string {
	return `
{{define "throttleTemplate" -}}
{{- if .Desc}}
# {{.Desc}}
{{- end}}
[transforms.{{.ComponentID}}]
type = "throttle"
inputs = {{.Inputs}}
window_secs = 1
threshold = {{.Threshold}}
key_field = '.kubernetes.container_name'
{{end}}
`
}
