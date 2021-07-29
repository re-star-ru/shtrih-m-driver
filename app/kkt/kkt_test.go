package kkt

import (
	"testing"
)

func TestKKT_PrintCheck(t *testing.T) {
	type fields struct {
		state int
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{name: "tag", fields: struct{ state int }{state: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//kkt := &KKT{
			//	//state: tt.fields.state,
			//}
			//kkt.PrintCheck()
		})
	}
}
