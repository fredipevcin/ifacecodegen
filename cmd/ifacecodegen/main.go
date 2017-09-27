package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/Masterminds/sprig"
	"github.com/fredipevcin/ifacecodegen"
)

const defaultTemplate = `
func New{{ .Name }}Logger(s {{ .Name }}) {{ .Name }} {
	return &logger{{ .Name }}{
		s: s,
	}
}
type logger{{ .Name }} struct {
	s {{ .Name }}
}

{{ $ifaceRef := . }}
{{ range .Methods }}
func (l *logger{{ $ifaceRef.Name }}) {{ .Name }}({{ input_parameters . }}) {{ $methodRef := . }}{{ output_parameters . }} {
		defer func(begin time.Time) {
			var err error {{ if ne (output_var_error .) "" -}}
			 = {{ output_var_error . }}
			{{- end }}
			log.Print("method", "{{ snakecase .Name }}", "took", time.Since(begin), "error", err)
		}(time.Now())
		{{ return . }} l.s.{{ .Name }}({{ input_calls . }})
}

{{ end }}
`
const usageText = `ifacecodegen generates code from interface using template

Source can be passed via stdin or as an argument -source.
If data is piped than default destination is stdout.

Examples:
	cat examples/interface.go | ifacecodegen
	ifacecodegen -source=examples/interface.go -destination -

Default template:
	func New{{ .Name }}Logger(s {{ .Name }}) {{ .Name }} {
		return &logger{{ .Name }}{
			s: s,
		}
	}
	type logger{{ .Name }} struct {
		s {{ .Name }}
	}

	{{ $ifaceRef := . }}
	{{ range .Methods }}
	func (l *logger{{ $ifaceRef.Name }}) {{ .Name }}({{ input_parameters . }}) {{ $methodRef := . }}{{ output_parameters . }} {
			defer func(begin time.Time) {
				var err error {{ if ne (output_var_error .) "" -}}
				= {{ output_var_error . }}
				{{- end }}
				log.Print("method", "{{ snakecase .Name }}", "took", time.Since(begin), "error", err)
			}(time.Now())
			{{ return . }} l.s.{{ .Name }}({{ input_calls . }})
	}

	{{ end }}

Options:

`

func main() {
	var (
		source       = flag.String("source", "", "Path to Go source file.")
		templateFile = flag.String("template", "", "Path to template file which will be used to generate code.")
		destination  = flag.String("destination", "", "Output file, '-' use to output to stdout (defaults to filename.gen.go).")
		packageOut   = flag.String("package", "", "Package of the generated code; defaults to the package of the source.")
		importPkgs   = flag.String("imports", "", "Comma-separated [name=]path pairs of explicit imports to use.")
		interfaces   = flag.String("interfaces", "", "Comma-separated names of interfaces to generate (default: generates for all interfaces)")
		metas        = flag.String("meta", "", "Comma-separated key=value pairs data to use in a template.")
		debugParser  = flag.Bool("debug_parser", false, "Print out parser results only.")
	)

	flag.Usage = func() {
		io.WriteString(os.Stderr, usageText)
		flag.PrintDefaults()
	}

	flag.Parse()

	sourceName, sourceReader, isPipe := parseSource(*source)

	pkg, err := ifacecodegen.Parse(ifacecodegen.ParseOptions{
		Source: sourceReader,
	})
	if err != nil {
		log.Fatalf("failed parsing package: %v", err)
	}

	if *debugParser {
		printDebug(pkg, os.Stdout)
		return
	}

	templateReader := parseTemplate(*templateFile)

	genOpts := ifacecodegen.GenerateOptions{
		Source:          sourceName,
		Package:         pkg,
		OverridePackage: *packageOut,
		Interfaces:      parseInterfaces(*interfaces),
		Template:        templateReader,
		Meta:            parseMeta(*metas),
		Imports:         parseImports(*importPkgs),
		Functions:       sprig.TxtFuncMap(),
	}

	data, err := ifacecodegen.Generate(genOpts)
	if err != nil {
		log.Fatalf("failed generating code: %v", err)
	}

	dest := parseDestination(*destination, *source, isPipe)
	defer dest.Close()

	dest.Write(data)
}
