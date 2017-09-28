Generating examples:

```bash
# example 1
ifacecodegen \
    -source interface.go \
    -template example1.tmpl \
    -destination example1.gen.go \
    -meta service=account \
    -imports "opentracing=github.com/opentracing/opentracing-go,tracinglog=github.com/opentracing/opentracing-go/log"

# example 2
ifacecodegen \
    -source interface.go \
    -template example2.tmpl \
    -destination example2.gen.go
```
