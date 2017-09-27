package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/fredipevcin/ifacecodegen"
)

func parseSource(source string) (string, io.Reader, bool) {
	fi, err := os.Stdin.Stat()
	if err == nil && fi.Mode()&os.ModeNamedPipe > 0 {
		return "stdin", os.Stdin, true
	}

	file, err := os.Open(source)
	if err != nil {
		log.Fatalf("failed opening source file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed reading from source file: %v", err)
	}

	return source, bytes.NewReader(data), false
}

func parseDestination(destination string, source string, isPipe bool) io.WriteCloser {
	if destination == "-" || (isPipe && destination == "") {
		return os.Stdout
	}

	if destination == "" {
		destination = strings.TrimSuffix(source, ".go") + ".gen.go"
	}

	f, err := os.Create(destination)
	if err != nil {
		log.Fatalf("failed opening destination file: %v", err)
	}

	return f
}

func parseMeta(meta string) map[string]string {
	metas := make(map[string]string)
	if meta != "" {
		for _, kv := range strings.Split(meta, ",") {
			eq := strings.Index(kv, "=")
			if eq == -1 {
				log.Fatalf("invalid meta value: %v", kv)
			}
			k, v := kv[:eq], kv[eq+1:]
			metas[k] = v
		}
	}
	return metas
}

func parseImports(imports string) []*ifacecodegen.Import {
	var result []*ifacecodegen.Import

	if imports != "" {
		for _, kv := range strings.Split(imports, ",") {
			eq := strings.Index(kv, "=")
			k, v := "", kv[eq+1:]
			if eq != -1 {
				k = kv[:eq]
			}

			result = append(result, &ifacecodegen.Import{
				Path:    k,
				Package: v,
			})
		}
	}

	return result
}

func parseInterfaces(interfaces string) []string {
	if interfaces != "" {
		return strings.Split(interfaces, ",")
	}

	return []string{}
}

func printDebug(pkg *ifacecodegen.Package, w io.Writer) {
	fmt.Fprintln(w, "package", pkg.Name)

	for _, iface := range pkg.Interfaces {
		fmt.Fprintln(w, "interface", iface.Name)
		for _, m := range iface.Methods {
			fmt.Fprintln(w, "  - method", m.Name)
			if len(m.In) > 0 {
				for _, mo := range m.In {
					fmt.Fprintln(w, "    in ", mo.Name, mo.Type)
				}
			}
			if len(m.Out) > 0 {
				for _, mo := range m.Out {
					fmt.Fprintln(w, "    out ", mo.Name, mo.Type)
				}
			}
		}
	}
}

func parseTemplate(templateFile string) io.Reader {
	if templateFile != "" {
		file, err := os.Open(templateFile)
		if err != nil {
			log.Fatalf("failed opening template file: %v", err)
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("failed reading from template file: %v", err)
		}

		return bytes.NewReader(data)
	}

	return strings.NewReader(defaultTemplate)
}
