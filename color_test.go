package main

import "testing"

func TestNewColor(t *testing.T) {
	tt := []struct {
		in  string
		err error
	}{
		{in: "#010101"},
		{in: "#AABBCC"},
		{in: "#ab0194"},
		{in: "invalid", err: ErrInvalidColorString},
		{in: "9123", err: ErrInvalidColorString},
	}
	for _, tc := range tt {
		t.Run(tc.in, func(t *testing.T) {
			_, err := NewColor(tc.in)
			if err == tc.err {
				return
			}
			if err == nil {
				t.Fatalf("expected error %q for input %q", tc.err, tc.in)
			}
			t.Fatalf("did not expect error %q for input %q", err, tc.in)
		})
	}
}

func TestRGB(t *testing.T) {
	tt := []struct {
		c       string
		r, g, b uint8
	}{
		{
			c: "#FF0000",
			r: 0xFF,
		},
		{
			c: "#00FE00",
			g: 0xFE,
		},
		{
			c: "#0000BB",
			b: 0xBB,
		},
	}
	for _, tc := range tt {
		t.Run(tc.c, func(t *testing.T) {
			c, err := NewColor(tc.c)
			if err != nil {
				t.Fatalf("unexpected test error: %v", err)
			}
			if r := c.R(); r != tc.r {
				t.Errorf("expected red to be %#02x, but got %#02x", tc.r, r)
			}
			if g := c.G(); g != tc.g {
				t.Errorf("expected green to be %#02x, but got %#02x", tc.g, g)
			}
			if b := c.B(); b != tc.b {
				t.Errorf("expected blue to be %#02x, but got %#02x", tc.b, b)
			}
		})
	}
}
