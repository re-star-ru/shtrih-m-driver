package commands

import (
	"reflect"
	"testing"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"
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
			name: "first check",
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
			name: "второй check",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmdData, err := CreateFNCloseCheck(tt.args.chk)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFNCloseCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCmdData, tt.wantCmdData) {
				t.Errorf("CreateFNCloseCheck() gotCmdData = %v, want %v", gotCmdData, tt.wantCmdData)
			}
		})
	}
}
