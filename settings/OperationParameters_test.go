package settings

import (
	"reflect"
	"testing"
)

func TestNewOperationParametersFromData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want *OperationParameters
	}{
		{
			name: "from docs",
			args: args{
				data: []byte{0xC0, 0x00, 0x00, 0x1A, 0x17, 0x44},
			},
			want: &OperationParameters{
				Head: HeadSavePowerDown,
				Address: [2]byte{
					0x00,
					0x00,
				},
				UartParityBit:        UartParityBit_8N1,
				TtlUartBaudRateBps:   BaudRateBps_9600,
				AirDataRateKbps:      AirDataRate_2_4k,
				Channel:              0x17, // TODO Double Check!
				TransmissionMode:     TransmissionModeTransparent,
				IODriveMode:          IODriveModePushPull,
				WirelessWakeUpTimeMs: WirelessWakeUpTime_250ms,
				FEC:                  FEC_On,
				TransmissionPowerDb:  TransmissionPowerDb_30db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOperationParametersFromData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Log("got: ", got)
				t.Log("want:", tt.want)
				t.Errorf("NewOperationParametersFromData() = %v, want %v", got, tt.want)
			}
		})
	}
}
