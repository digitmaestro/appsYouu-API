package prisma

import (
	"log"
	"reflect"
	"testing"
)

func TestPrice(t *testing.T) {
	a1 := 2.55
	a2 := int32(255)
	b1 := 1.11111
	b2 := int32(111)
	c1 := float64(2)
	c2 := int32(200)
	type args struct {
		f *float64
	}
	tests := []struct {
		name string
		args args
		want *int32
	}{
		{
			name: "a",
			args: args{
				f: &a1,
			},
			want: &a2,
		},
		{
			name: "b",
			args: args{
				f: &b1,
			},
			want: &b2,
		},
		{
			name: "c",
			args: args{
				f: &c1,
			},
			want: &c2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Printf("what up")
			if got := Price(tt.args.f); !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("Price() = %v, want %v", *got, *tt.want)
			}
		})
	}
}
