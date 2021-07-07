package utils

import (
	"reflect"
	"testing"
)

func TestPrimeNumberDecompose(t *testing.T) {
	type args struct {
		num int32
	}
	tests := []struct {
		name    string
		args    args
		want    []int32
		wantErr bool
	}{
		{"Should Return [2,2,2,3,5]", args{120}, []int32{2, 2, 2, 3, 5}, false},
		{"Should Return [2,3,19]", args{228}, []int32{2, 2, 3, 19}, false},
		{"Should Return [2,2,3,5]", args{300}, []int32{2, 2, 3, 5, 5}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrimeNumberDecompose(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrimeNumberDecompose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrimeNumberDecompose() got = %v, want %v", got, tt.want)
			}
		})
	}
}
