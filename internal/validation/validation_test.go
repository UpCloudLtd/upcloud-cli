package validation

import "testing"

func TestUuid4(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"61cc19f4-5702-4120-b0e4-82dce36e7b88"}, false},
		{"case insensitive", args{"61CC19F4-5702-4120-B0E4-82DCE36E7B88"}, false},
		{"wrong version", args{"61cc19f4-5702-5120-b0e4-82dce36e7b88"}, true},
		{"wrong length", args{"61cc19f4-5702-4120-b0e4-82dce36e7b888"}, true},
		{"unsupported delimiter", args{"61cc19f4-5702-4120fb0e4-82dce36e7b88"}, true},
		{"invalid chars", args{"61cX19f4-5702-4120-b0e4-82dce36e7b88"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UUID4(tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("UUID4() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
