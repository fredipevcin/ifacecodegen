package ifacecodegen

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParserEmtptySource(t *testing.T) {
	opts := ParseOptions{}
	pkg, err := Parse(opts)

	if pkg != nil {
		t.Fatalf("Pkg should be nil, actual %v", pkg)
	}
	if err == nil {
		t.Fatalf("Error should not be nil")
	}
	if actual, expected := err.Error(), "Source should not be nil"; actual != expected {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}
func TestParserInvalidSource(t *testing.T) {
	opts := ParseOptions{
		Source: strings.NewReader(`foo`),
	}
	pkg, err := Parse(opts)

	if pkg != nil {
		t.Fatalf("Pkg should be nil, actual %v", pkg)
	}
	if err == nil {
		t.Fatalf("Error should not be nil")
	}

	if !strings.HasPrefix(err.Error(), "failed parsing source 1:1") {
		t.Errorf("Error should be about parsing")
	}
}

func TestParserParsePackage(t *testing.T) {
	opts := ParseOptions{
		Source: strings.NewReader(`
package foo
`),
	}
	pkg, err := Parse(opts)

	if err != nil {
		t.Fatalf("Error should be nil, actual %v", err)
	}
	if pkg == nil {
		t.Fatalf("Pkg should not be nil")
	}

	if actual, expected := pkg, (&Package{Name: "foo"}); !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestParser(t *testing.T) {
	opts := ParseOptions{
		Source: strings.NewReader(`
package foo

type Service interface {
	Foo(string) Bar
	Mult(p1, p2 string) (r1, r2 string)
}

`),
	}
	pkg, err := Parse(opts)

	if err != nil {
		t.Fatalf("Error should be nil, actual %v", err)
	}
	if pkg == nil {
		t.Fatalf("Pkg should not be nil")
	}

	expected := &Package{
		Name: "foo",
		Interfaces: []*Interface{
			&Interface{
				Name: "Service",
				Methods: []*Method{
					&Method{
						Name: "Foo",
						In: []*Parameter{
							&Parameter{
								Name: "_param1",
								Type: TypeBuiltin("string"),
							},
						},
						Out: []*Parameter{
							&Parameter{
								Name: "_result1",
								Type: &TypeExported{
									Package: "foo",
									Type:    TypeBuiltin("Bar"),
								},
							},
						},
					},
					&Method{
						Name: "Mult",
						In: []*Parameter{
							&Parameter{
								Name: "p1",
								Type: TypeBuiltin("string"),
							},
							&Parameter{
								Name: "p2",
								Type: TypeBuiltin("string"),
							},
						},
						Out: []*Parameter{
							&Parameter{
								Name: "r1",
								Type: TypeBuiltin("string"),
							},
							&Parameter{
								Name: "r2",
								Type: TypeBuiltin("string"),
							},
						},
					},
				},
			},
		},
	}
	if actual := pkg; !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func debugPkg(pkg *Package) {
	for _, iface := range pkg.Interfaces {
		fmt.Println("interface", iface.Name)
		for _, m := range iface.Methods {
			fmt.Println("  - method", m.Name)
			if len(m.In) > 0 {
				for _, mo := range m.In {
					fmt.Println("    in ", mo.Name, mo.Type)
				}
			}
			if len(m.Out) > 0 {
				for _, mo := range m.Out {
					fmt.Println("    out ", mo.Name, mo.Type)
				}
			}
		}
	}
}
