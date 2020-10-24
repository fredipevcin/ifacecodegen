package ifacecodegen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strconv"
)

// ParseOptions are options used to parse
type ParseOptions struct {
	Source io.Reader
}

// Parse returns parsed package
func Parse(opts ParseOptions) (*Package, error) {
	if opts.Source == nil {
		return nil, fmt.Errorf("Source should not be nil")
	}

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "", opts.Source, parser.AllErrors)
	if err != nil {
		return nil, fmt.Errorf("failed parsing source %v", err)
	}

	pkg := &Package{
		Name: file.Name.String(),
	}

	for _, decl := range file.Decls {
		var gd, ok = decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}

		for _, spec := range gd.Specs {
			var ts *ast.TypeSpec
			ts, ok = spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// Filter for only type definitions that are also interfaces
			var it *ast.InterfaceType
			it, ok = ts.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			iface := &Interface{
				Name: ts.Name.String(),
			}

			for _, attribute := range it.Methods.List {
				switch n := attribute.Type.(type) {
				case *ast.FuncType:
					m, e := parseFunc(pkg.Name, attribute.Names[0].String(), n)
					if e != nil {
						return nil, e
					}
					iface.Methods = append(iface.Methods, m)
				}
			}

			pkg.Interfaces = append(pkg.Interfaces, iface)
		}
	}
	return pkg, nil
}

func parseType(pkg string, arg ast.Expr) (Type, error) {
	switch n := arg.(type) {
	case *ast.ArrayType:
		var len = -1
		if n.Len != nil {
			len, _ = strconv.Atoi(n.Len.(*ast.BasicLit).Value)
		}
		var typ, e = parseType(pkg, n.Elt)
		if e != nil {
			return nil, e
		}
		return &TypeArray{Len: len, Type: typ}, nil
	case *ast.ChanType:
		var t, e = parseType(pkg, n.Value)
		if e != nil {
			return nil, e
		}
		var chanType = &TypeChan{Type: t}
		if n.Dir == ast.SEND {
			chanType.WriteOnly = true
		}
		if n.Dir == ast.RECV {
			chanType.ReadOnly = true
		}
		return chanType, nil
	case *ast.Ellipsis:
		var t, e = parseType(pkg, n.Elt)
		if e != nil {
			return nil, e
		}
		return &TypeVariadic{Type: t}, nil
	case *ast.FuncType:
		var method, e = parseFunc(pkg, "", n)
		if e != nil {
			return nil, e
		}
		var result = &TypeFunc{In: make([]Type, 0), Out: make([]Type, 0)}
		for _, param := range method.In {
			result.In = append(result.In, param.Type)
		}
		for _, param := range method.Out {
			result.Out = append(result.Out, param.Type)
		}
		return result, nil
	case *ast.Ident:
		if n.IsExported() {
			// assume type in this package
			return &TypeExported{Package: pkg, Type: TypeBuiltin(n.Name)}, nil
		}
		return TypeBuiltin(n.Name), nil
	case *ast.InterfaceType:
		if n.Methods != nil && len(n.Methods.List) > 0 {
			return nil, fmt.Errorf("can't handle non-empty unnamed interface types at %v", n.Pos())
		}
		return TypeBuiltin("interface{}"), nil
	case *ast.MapType:
		var key Type
		var value Type
		var e error
		key, e = parseType(pkg, n.Key)
		if e != nil {
			return nil, e
		}
		value, e = parseType(pkg, n.Value)
		if e != nil {
			return nil, e
		}
		return &TypeMap{Key: key, Value: value}, nil
	case *ast.SelectorExpr:
		var pkgName = n.X.(*ast.Ident).String()
		return &TypeExported{Package: pkgName, Type: TypeBuiltin(n.Sel.String())}, nil
	case *ast.StarExpr:
		var t, e = parseType(pkg, n.X)
		if e != nil {
			return nil, e
		}
		return &TypePointer{Type: t}, nil
	case *ast.StructType:
		if n.Fields != nil && len(n.Fields.List) > 0 {
			return nil, fmt.Errorf("can't handle non-empty unnamed struct types at %v", n.Pos())
		}
		return TypeBuiltin("struct{}"), nil
	}
	return nil, fmt.Errorf("unknown type: %T", arg)
}

func parseFunc(pkg string, name string, f *ast.FuncType) (*Method, error) {
	var method = &Method{Name: name, In: make([]*Parameter, 0), Out: make([]*Parameter, 0)}
	if f.Params != nil {
		var index int
		for _, arg := range f.Params.List {
			var t, e = parseType(pkg, arg.Type)
			if e != nil {
				return nil, e
			}

			if len(arg.Names) > 0 {
				for i := range arg.Names {
					index++
					method.In = append(method.In, &Parameter{
						Name: arg.Names[i].String(),
						Type: t,
					})
				}
			} else {
				index++
				method.In = append(method.In, &Parameter{
					Name: fmt.Sprintf("_param%d", index),
					Type: t,
				})
			}
		}
	}
	if f.Results != nil {
		var index int
		for _, arg := range f.Results.List {
			var t, e = parseType(pkg, arg.Type)
			if e != nil {
				return nil, e
			}

			if len(arg.Names) > 0 {
				for i := range arg.Names {
					index++
					method.Out = append(method.Out, &Parameter{
						Name: arg.Names[i].String(),
						Type: t,
					})
				}
			} else {
				index++
				method.Out = append(method.Out, &Parameter{
					Name: fmt.Sprintf("_result%d", index),
					Type: t,
				})
			}
		}
	}
	return method, nil
}
