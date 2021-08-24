package commands

import (
	"encoding/hex"
	"log"
	"reflect"
	"testing"

	"github.com/fess932/shtrih-m-driver/pkg/driver/models"

	"github.com/fess932/shtrih-m-driver/pkg/consts"
	"github.com/stretchr/testify/assert"
)

var successCreateCloseCheckData = []byte{
	255, 69, 30, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

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
			name: "success check package",
			args: args{chk: models.CheckPackage{
				Cash:      1,
				Digital:   0,
				Rounding:  0,
				TaxSystem: consts.PSN,
			}},
			wantCmdData: successCreateCloseCheckData,
			wantErr:     false,
		},
		{
			name: "rounding > 99",
			args: args{chk: models.CheckPackage{
				Cash:      0,
				Digital:   0,
				Rounding:  100,
				TaxSystem: 0,
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

func TestCreateFNOperationV2(t *testing.T) {
	type args struct {
		o models.Operation
	}
	tests := []struct {
		name        string
		args        args
		wantCmdData []byte
		wantErr     bool
	}{
		{
			name: "first",
			args: args{
				o: models.Operation{
					Type:    0,
					Subject: 0,
					Amount:  0,
					Price:   0,
					Sum:     0,
					Name:    "",
				},
			},
			wantCmdData: successCreateCloseCheckData,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmdData, err := CreateFNOperationV2(tt.args.o)
			log.Println("got", hex.Dump(gotCmdData))

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFNOperationV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCmdData, tt.wantCmdData) {
				t.Errorf("CreateFNOperationV2() gotCmdData = %v, want %v", gotCmdData, tt.wantCmdData)
			}
		})
	}
}
