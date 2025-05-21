package assembler

import (
	"reflect"
	"testing"
)

func TestGenerateELFHeaders(t *testing.T) {
	type args struct {
		e_entry [4]byte
		e_phnum [2]byte
	}
	tests := []struct {
		name string
		args args
		want *[0x34]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateELFHeaders(tt.args.e_entry, tt.args.e_phnum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateELFHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateSingleELFProgramHeader(t *testing.T) {
	type args struct {
		htype     byte
		offset    [4]byte
		size      [4]byte
		memoffset [4]byte
	}
	tests := []struct {
		name string
		args args
		want *[0x20]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSingleELFProgramHeader(tt.args.htype, tt.args.offset, tt.args.size, tt.args.memoffset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateSingleELFProgramHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildELFFile(t *testing.T) {
	type args struct {
		program Program
	}
	tests := []struct {
		name string
		args args
		want *[]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildELFFile(tt.args.program); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildELFFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
