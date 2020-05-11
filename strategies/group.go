package strategies

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"1pkg/gopium"

	"golang.org/x/sync/errgroup"
)

// list of tag presets
const (
	tskip = "-"
	tdef  = "default"
)

// list of group presets
var (
	ptgrp = group{}
)

// group defines strategy implementation
// that uses fields tags annotation
// in order to process different set of strategies
// on different groups and then combine results
// in single struct result, it should be
// able to read autogenerated tags
// generated by tag strategy.
// note: supports only next fields tags annotation formats
// `gopium:"stg,stg,stg"` processed as `default` group
// `gopium:"group:def;stg,stg,stg"` processed as named group
type group struct {
	builder Builder
}

// container carries sing group data
type container struct {
	grp string
	o   gopium.Struct
	r   gopium.Struct
	stg gopium.Strategy
}

// Curator erich group strategy with builder instance
func (stg group) Builder(builder Builder) group {
	stg.builder = builder
	return stg
}

// Apply group implementation
func (stg group) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := o
	// parse tag annotation
	// into containers groups
	containers, err := stg.parse(r)
	// in case of any error
	// just return error back
	if err != nil {
		return o, err
	}
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
	// go through all containers and apply
	// all strategies concurently on inner structs
	for i := range containers {
		container := &containers[i]
		group.Go(func() error {
			// apply strategy on struct
			rst, err := container.stg.Apply(gctx, container.o)
			// in case of any error
			// just return error back
			if err != nil {
				return err
			}
			// in case of success
			// update result on container
			container.r = rst
			// if we faced default group
			// update result comment and doc
			if container.grp == tdef {
				r = rst
			}
			return nil
		})
	}
	// wait until all strategies
	// have been applied and resolved
	if err := group.Wait(); err != nil {
		return o, err
	}
	// sort result containers lexicographicaly
	sort.SliceStable(containers, func(i, j int) bool {
		return containers[i].grp < containers[j].grp
	})
	// combine all results to single result struct
	r.Fields = nil
	for i := range containers {
		r.Fields = append(r.Fields, containers[i].r.Fields...)
	}
	return r, ctx.Err()
}

// parse helps to parse structure fields tags
// into groups container or returns parse error
// - `gopium:"stg,stg,stg"` parsed to `default` group
// - `gopium:"group:def;stg,stg,stg"` parsed to named group
// - otherwise a parse error returned
func (stg group) parse(st gopium.Struct) ([]container, error) {
	// setup temporary groups maps
	// for fields and strategies
	gfields := make(map[string][]gopium.Field)
	gstrategies := make(map[string]gopium.Strategy)
	gstrategiesnames := make(map[string]string)
	// go through all struct fields
	for _, f := range st.Fields {
		// grab the field tag
		tag, ok := reflect.StructTag(f.Tag).Lookup(tagname)
		// in case tag is empty
		// or marked as skipped
		if !ok || tag == tskip {
			gfields[tskip] = append(gfields[tskip], f)
			continue
		}
		// trim all excess separators
		tag = strings.Trim(tag, ";")
		// otherwise parse the tag
		tokens := strings.Split(tag, ";")
		switch tlen := len(tokens); tlen {
		case 1:
			stgs := tokens[0]
			// check that strategies list is consistent
			if gstg, ok := gstrategiesnames[tdef]; ok && gstg != stgs {
				return nil, fmt.Errorf(
					"inconsistent strategies list %q for field %q in group %q",
					stgs,
					f.Name,
					tdef,
				)
			}
			// collect strategies and fields
			gstrategiesnames[tdef] = stgs
			gfields[tdef] = append(gfields[tdef], f)
		case 2:
			group := tokens[0]
			stgs := tokens[1]
			// check that tag contains group anchor
			if !strings.Contains(group, "group:") {
				return nil, fmt.Errorf("tag %q can't be parsed, named group `group:` anchor wasn't found", f.Tag)
			}
			// remove group anchor
			group = strings.Replace(group, "group:", "", 1)
			// check that strategies list is consistent
			if gstg, ok := gstrategiesnames[group]; ok && gstg != stgs {
				return nil, fmt.Errorf(
					"inconsistent strategies list %q for field %q in group %q",
					stgs,
					f.Name,
					group,
				)
			}
			// collect strategies and fields
			gstrategiesnames[group] = stgs
			gfields[group] = append(gfields[group], f)
		default:
			// return parsing error msg
			return nil, fmt.Errorf("tag %q can't be parsed, neither as `default` nor named group", f.Tag)
		}
	}
	// go through all collected group strategies names
	// and build pipe strategy from them
	for grp, gstgs := range gstrategiesnames {
		// prepare strategy pipe
		names := strings.Split(gstgs, ",")
		p := make(pipe, 0, len(names))
		// go through list of strategy name
		for _, name := range names {
			// try to build new strategy by name
			stg, err := stg.builder.Build(gopium.StrategyName(name))
			// in case of any error
			// just return it back
			if err != nil {
				return nil, err
			}
			// otherwise append strategy to pipe
			p = append(p, stg)
		}
		// set group strategies val
		gstrategies[grp] = p
	}
	// setup result containers
	containers := make([]container, 0, len(gfields))
	// go through all collected group fields
	for grp, fields := range gfields {
		// prepare new empty group container
		var cnt container
		// set container group
		cnt.grp = grp
		// set container original
		// struct and its fields
		cnt.o = st
		cnt.o.Fields = fields
		// if group has strategy set it
		// otherwise set nil strategy
		if stg, ok := gstrategies[grp]; ok {
			cnt.stg = stg
		} else {
			cnt.stg = np
		}
		// append current container to result
		containers = append(containers, cnt)
	}
	// return result containers
	return containers, nil
}
