as
{{.a}}
{{range $key, $val := .apps}}
{{$key}}
{{$val.from}}
{{end}} 