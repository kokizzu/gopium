package strategy

import (
	"context"
	"go/types"
	"sort"

	"1pkg/gopium"
	gtypes "1pkg/gopium/types"
)

// stgmemsort defines struct optimal memory fields sorting Strategy implementation
// that goes through all structure fields and uses gtypes.Extractor
// to extract gopium.Field DTO for each field
// sorts fields accordingly to their sizes in descending order
// and puts it back to resulted gopium.Struct object
type stgmemsort struct {
	//nolint
	extractor gtypes.Extractor
}

// Apply stgmemsort implementation
func (stg stgmemsort) Apply(ctx context.Context, name string, st *types.Struct) (r gopium.StructError) {
	enum := stgenum(stg)
	r = enum.Apply(ctx, name, st)
	sort.SliceStable(r.Struct.Fields, func(i, j int) bool {
		return r.Struct.Fields[j].Size < r.Struct.Fields[i].Size
	})
	return
}
