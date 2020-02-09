package strategy

import (
	"context"
	"go/types"
	"sort"

	"1pkg/gopium"
	gtypes "1pkg/gopium/types"
)

// stgmemsort defines struct fields names consistent sorting Strategy implementation
// that goes through all structure fields and uses gtypes.Extractor
// to extract gopium.Field DTO for each field
// sorts fields accordingly to their names in ascending order
// and puts it back to resulted gopium.Struct object
type stgnamesort struct {
	//nolint
	extractor gtypes.Extractor
}

// Apply stgnamesort implementation
func (stg stgnamesort) Apply(ctx context.Context, name string, st *types.Struct) (r gopium.StructError) {
	enum := stgenum(stg)
	r = enum.Apply(ctx, name, st)
	sort.SliceStable(r.Struct.Fields, func(i, j int) bool {
		return r.Struct.Fields[i].Name < r.Struct.Fields[j].Name
	})
	return
}
