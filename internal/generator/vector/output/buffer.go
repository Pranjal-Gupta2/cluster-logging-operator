package output

type Buffer struct {
	SinkComponentID string
	MaxEvents       string
	Type            BufferType
	WhenFull        BufferPolicy
}

func (b Buffer) Name() string {
	return "bufferTemplate"
}

func (b Buffer) Template() string {
	return `{{define "` + b.Name() + `" -}}
[sinks.{{.SinkComponentID}}.buffer]
type = "{{ .Type }}"
max_events = {{ .MaxEvents }}
when_full = "{{ .WhenFull }}"
{{end}}`
}

type BufferType string

const (
	DiskBuffer   BufferType = "disk"
	MemoryBuffer BufferType = "memory"
)

type BufferPolicy string

const (
	BlockPolicy BufferPolicy = "block"
	DropPolicy  BufferPolicy = "drop_newest"
)
