package consul

import (
	"reflect"
	"testing"
)

func TestMarshalCommandMeta(t *testing.T) {
	type args struct {
		data []CommandMeta
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Normal usage", args{data: []CommandMeta{{Command: "cmd", Queue: "queue"}}}, `[{"command":"cmd","queue":"queue"}]`, false},
		{"No command", args{data: []CommandMeta{{Queue: "queue"}}}, "", true},
		{"No queue", args{data: []CommandMeta{{Command: "cmd"}}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalCommandMeta(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalCommandMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MarshalCommandMeta() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCommandMeta(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *CommandMetaInfo
		wantErr bool
	}{
		{"Noraml usage", args{data: []byte(`[{"command":"cmd","queue":"queue"}]`)}, &CommandMetaInfo{{Command: "cmd", Queue: "queue"}}, false},
		{"No command", args{data: []byte(`[{"queue":"queue"}]`)}, nil, true},
		{"No queue", args{data: []byte(`[{"command":"cmd"}]`)}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommandMeta(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommandMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommandMeta() got = %v, want %v", got, tt.want)
			}
		})
	}
}
