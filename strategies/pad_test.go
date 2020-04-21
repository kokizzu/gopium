package strategies

import (
	"context"
	"reflect"
	"testing"

	"1pkg/gopium"
	"1pkg/gopium/mocks"
)

func TestPad(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		pad     pad
		curator gopium.Curator
		ctx     context.Context
		o       gopium.Struct
		r       gopium.Struct
		err     error
	}{
		"empty struct should be applied to empty struct": {
			pad:     padsys,
			curator: mocks.Maven{SysAlignVal: 16},
			ctx:     context.Background(),
		},
		"non empty struct should be applied to explicit pad aligned struct": {
			pad:     padsys,
			curator: mocks.Maven{SysAlignVal: 6},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					gopium.PadField(4),
				},
			},
		},
		"non empty struct should be applied to explicit pad aligned struct on canceled context": {
			pad:     padtnat,
			curator: mocks.Maven{SysAlignVal: 12},
			ctx:     cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 5,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test",
						Size:  8,
						Align: 5,
					},
					gopium.PadField(2),
				},
			},
			err: cctx.Err(),
		},
		"mixed struct should be applied to explicit pad aligned struct on type natural pad": {
			pad:     padtnat,
			curator: mocks.Maven{SysAlignVal: 24},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test4",
						Size:  3,
						Align: 1,
					},
					{
						Name:  "test1",
						Size:  32,
						Align: 4,
					},
					{
						Name:  "test2",
						Size:  6,
						Align: 6,
					},
					{
						Name:  "test3",
						Size:  8,
						Align: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name:  "test4",
						Size:  3,
						Align: 1,
					},
					gopium.PadField(1),
					{
						Name:  "test1",
						Size:  32,
						Align: 4,
					},
					{
						Name:  "test2",
						Size:  6,
						Align: 6,
					},
					gopium.PadField(6),
					{
						Name:  "test3",
						Size:  8,
						Align: 8,
					},
				},
			},
		},
		"mixed struct should be applied to explicit pad aligned on same sys pad": {
			pad:     padsys,
			curator: mocks.Maven{SysAlignVal: 9},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 9,
					},
					{
						Name: "test2",
						Size: 7,
					},
					{
						Name: "test3",
						Size: 5,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 9,
					},
					{
						Name: "test2",
						Size: 7,
					},
					gopium.PadField(2),
					{
						Name: "test3",
						Size: 5,
					},
					gopium.PadField(4),
					{
						Name: "test4",
						Size: 3,
					},
					gopium.PadField(6),
				},
			},
		},
		"mixed struct should be applied to explicit pad aligned on bigger sys pad": {
			pad:     padsys,
			curator: mocks.Maven{SysAlignVal: 12},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 9,
					},
					{
						Name: "test2",
						Size: 7,
					},
					{
						Name: "test3",
						Size: 5,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 9,
					},
					gopium.PadField(3),
					{
						Name: "test2",
						Size: 7,
					},
					gopium.PadField(5),
					{
						Name: "test3",
						Size: 5,
					},
					gopium.PadField(7),
					{
						Name: "test4",
						Size: 3,
					},
					gopium.PadField(9),
				},
			},
		},
		"mixed struct should be applied to explicit pad aligned no additional aligment": {
			pad:     padsys,
			curator: mocks.Maven{SysAlignVal: 4},
			ctx:     context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 24,
					},
					{
						Name: "test2",
						Size: 12,
					},
					{
						Name: "test3",
						Size: 36,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 24,
					},
					{
						Name: "test2",
						Size: 12,
					},
					{
						Name: "test3",
						Size: 36,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			pad := tcase.pad.Curator(tcase.curator)
			r, err := pad.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}