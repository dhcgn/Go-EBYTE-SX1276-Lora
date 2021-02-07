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

func Test_hash(t *testing.T) {
	type args struct {
		op OperationParameters
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty",
			args: args{
				op: OperationParameters{},
			},
			want: "EQ6TRTPWJ4OXXLERTGM2HZ3E2IEZ36EZWKOG5BRFIXG2YPSHKQQA",
		},
		{
			name: "Changed",
			args: args{
				op: OperationParameters{
					TransmissionPowerDb: 1,
				},
			},
			want: "TTJKMPYH2R7BCJZIO7RUEZNCISYNN6GHTWLHXWSQO5R5MYSH43FQ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hash(tt.args.op); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationParameters_GetShortHash(t *testing.T) {
	type fields struct {
		Head                 Head
		Address              Address
		UartParityBit        UartParityBit
		TtlUartBaudRateBps   BaudRate
		AirDataRateKbps      AirDataRate
		Channel              Channel
		TransmissionMode     TransmissionMode
		IODriveMode          IODriveMode
		WirelessWakeUpTimeMs WirelessWakeUpTime
		FEC                  FEC
		TransmissionPowerDb  TransmissionPowerDb
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "Empty",
			fields: fields{},
			want:   "EQ6TRTPWJ4OX",
		},
		{
			name: "Set",
			fields: fields{
				FEC: FEC_On,
			},
			want: "B7GLZBUOJEPS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := OperationParameters{
				Head:                 tt.fields.Head,
				Address:              tt.fields.Address,
				UartParityBit:        tt.fields.UartParityBit,
				TtlUartBaudRateBps:   tt.fields.TtlUartBaudRateBps,
				AirDataRateKbps:      tt.fields.AirDataRateKbps,
				Channel:              tt.fields.Channel,
				TransmissionMode:     tt.fields.TransmissionMode,
				IODriveMode:          tt.fields.IODriveMode,
				WirelessWakeUpTimeMs: tt.fields.WirelessWakeUpTimeMs,
				FEC:                  tt.fields.FEC,
				TransmissionPowerDb:  tt.fields.TransmissionPowerDb,
			}
			if got := op.GetShortHash(); got != tt.want {
				t.Errorf("GetShortHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
