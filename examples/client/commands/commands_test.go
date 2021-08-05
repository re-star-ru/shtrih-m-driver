package commands

import (
	"reflect"
	"testing"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateFNCloseCheck(t *testing.T) {
	type args struct {
		chk models.CheckPackage
	}

	tests := []struct {
		name        string
		args        args
		wantCmdData []byte
		wantErr     bool
	}{
		{
			name: "rounding > 99",
			args: args{chk: models.CheckPackage{
				CashierINN: "",
				Operations: nil,
				Cash:       0,
				Casheless:  0,
				BottomLine: "",
				Rounding:   100,
				TaxSystem:  0,
				Electronic: false,
			}},
			wantCmdData: nil,
			wantErr:     true,
		},
		{
			name: "wrong inn",
			args: args{chk: models.CheckPackage{
				CashierINN: "",
				Operations: nil,
				Cash:       0,
				Casheless:  0,
				BottomLine: "",
				Rounding:   0,
				TaxSystem:  0,
				Electronic: false,
			}},
			wantCmdData: nil,
			wantErr:     true,
		},
		{
			name: "wrong inn",
			args: args{chk: models.CheckPackage{
				CashierINN: "",
				Operations: nil,
				Cash:       0,
				Casheless:  0,
				BottomLine: "",
				Rounding:   0,
				TaxSystem:  0,
				Electronic: false,
			}},
			wantCmdData: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmdData, err := CreateFNCloseCheck(tt.args.chk)

			if tt.wantErr {
				if assert.Error(t, err) {
					assert.Nil(t, gotCmdData)
				}
				return
			}

			assert.Len(t, gotCmdData, 182)

			if !reflect.DeepEqual(gotCmdData, tt.wantCmdData) {
				t.Errorf("CreateFNCloseCheck() gotCmdData = \n%v\nwant\n%v", gotCmdData, tt.wantCmdData)
			}
		})
	}
}
