package walkers

import (
	"context"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"
)

// list of wout presets
var (
	jsonstd = wout{
		fmt:    fmtio.Jsonb,
		writer: fmtio.Stdout,
	}
	xmlstd = wout{
		fmt:    fmtio.Xmlb,
		writer: fmtio.Stdout,
	}
	csvstd = wout{
		fmt:    fmtio.Csvb(fmtio.Buffer()),
		writer: fmtio.Stdout,
	}
	jsonfiles = wout{
		fmt:    fmtio.Jsonb,
		writer: fmtio.File("json"),
	}
	xmlfiles = wout{
		fmt:    fmtio.Xmlb,
		writer: fmtio.File("xml"),
	}
	csvfiles = wout{
		fmt:    fmtio.Csvb(fmtio.Buffer()),
		writer: fmtio.File("csv"),
	}
)

// wout defines packages walker out implementation
type wout struct {
	// inner visiting parameters
	fmt    fmtio.Bytes
	writer fmtio.Writer
	// external visiting parameters
	parser  gopium.TypeParser
	exposer gopium.Exposer
	deep    bool
	bref    bool
}

// With erich wast walker with external visiting parameters
// parser, exposer instances and additional visiting flags
func (w wout) With(p gopium.TypeParser, exp gopium.Exposer, deep bool, bref bool) wout {
	w.parser = p
	w.exposer = exp
	w.deep = deep
	w.bref = bref
	return w
}

// Visit wout implementation uses visit function helper
// to go through all structs decls inside the package
// and applies strategy to them to get results,
// then uses struct to bytes to format strategy results
// and use writer to write results to output
func (w wout) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy) error {
	// use parser to parse types pkg data
	// we don't care about fset
	pkg, loc, err := w.parser.ParseTypes(ctx)
	if err != nil {
		return err
	}
	// create govisit func
	// using gopium.Visit helper
	// and run it on pkg scope
	ch := make(appliedCh)
	gvisit := with(w.exposer, loc, w.bref).
		visit(regex, stg, ch, w.deep)
	// prepare separate cancelation
	// context for visiting
	gctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// run visiting in separate goroutine
	go gvisit(gctx, pkg.Scope())
	// prepare struct storage
	var floc string
	f := make(collections.Flat)
	for applied := range ch {
		// in case any error happened
		// just return error back
		// it auto cancels context
		if applied.Err != nil {
			return applied.Err
		}
		// TODO hacky solution
		floc = applied.Loc
		// push struct to storage
		f[applied.ID] = applied.R
	}
	// run sync write
	// with collected strategies results
	return w.write(gctx, floc, f)
}

// write wout helps to apply struct to bytes
// to format strategy result and writer
// to write result to output
func (w wout) write(ctx context.Context, loc string, f collections.Flat) error {
	// apply formatter
	buf, err := w.fmt(f)
	// in case any error happened
	// in formatter return error back
	if err != nil {
		return err
	}
	// generate writer
	writer, err := w.writer("gopium", loc)
	if err != nil {
		return err
	}
	// write results and close writer
	// in case any error happened
	// in writer return error
	if _, err := writer.Write(buf); err != nil {
		return err
	}
	return writer.Close()
}
