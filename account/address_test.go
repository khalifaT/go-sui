package account

import "testing"

func Test_IsValidAddress(t *testing.T) {
	tests := []struct {
		address   string
		wantValid bool
	}{
		{"MDpQc 1IIzkie1dJdj nfm85XmRCJmk KHVUU05Abg==", false},
		{"addr1qyjc2543ctgs08k5ha872y4z458jwjsta7eldqmz5eudvhc7yflv0e0g7urz7n88ls0kx7d0m4eqp5tkyf7cfdat0yxqugvpua", false},
		{"X09wJFxwQDdTU1tzMy5NJXdSTnknPCh9J0tNUCdmIw  ", false},
		{"5f713bef531629b47dd1bdbb382a", false},
		{"f1e2a6d12cd5e62a3ce9b2c12e9e2d37d81c", false},
		{"0X5f713bef531629b47dd1bdbb382acec5224fc9abc16133e3", false},
		{"0x503ff67d9291215ffccafddbd08d86e86b3425c6356c9679", false},
		{"9edd26f2ef1c1796f9feaa703c8628e5a70618c8", true},
		{"5f713bef531629b47dd1bdbb382acec5224fc9ab", true},
		{"0Xdce47e3e523b5e52a36d74295c0d83d91f80b47c", true},
		{"0x4288ba9932cc115784794fcfb709213f30d40a54", true},
	}

	for _, tt := range tests {

		if gotValid, _ := IsValidAddress(tt.address); gotValid != tt.wantValid {
			t.Errorf("address %s IsValidAddress() = %v, want %v", tt.address, gotValid, tt.wantValid)
		}
	}
}
