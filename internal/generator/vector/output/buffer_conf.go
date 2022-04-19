package output

type BufferType string
type BufferOverflowAction string

const (
	DiskBuffer   BufferType           = "disk"
	MemoryBuffer BufferType           = "memory"
	BlockAction  BufferOverflowAction = "block"
	DropAction   BufferOverflowAction = "drop_newest"

	DefaultBufferType     BufferType           = DiskBuffer
	DefaultMaxEvents                           = "500"
	DefaultMaxSize                             = "8589934592" // 0x200000000, 8GB
	DefaultOverflowAction BufferOverflowAction = BlockAction
)

type Buffer struct {
	SinkComponentID string
	MaxEvents       string
	MaxSize         string
	Type            BufferType
	WhenFull        BufferOverflowAction
}

func (b Buffer) Name() string {
	return "bufferTemplate"
}

func (b Buffer) Template() string {
	return `{{define "` + b.Name() + `" -}}
[sinks.{{.SinkComponentID}}.buffer]
type = "{{ .Type }}"
{{if eq .Type "disk" -}}
max_size = {{ .MaxSize }}
{{- else -}}
max_events = {{ .MaxEvents }}
{{- end}}
when_full = "{{ .WhenFull }}"
{{end}}`
}
