package main

import "testing"

func TestProtectionModes(t *testing.T) {
	detected := run("detect", 1, 300, 25, 8, 10, "all")
	if detected.injected != 1 || detected.detected != 1 || detected.completed >= 300 {
		t.Fatalf("detect result: %+v", detected)
	}
	recovered := run("recover", 1, 300, 25, 8, 10, "all")
	if recovered.completed != 300 || recovered.detected == 0 || recovered.rollbacks != recovered.detected || recovered.wasted == 0 || recovered.invalid || recovered.accuracy < 0.95 {
		t.Fatalf("recover result: %+v", recovered)
	}
	clean := run("none", 1, 300, 25, 8, 0, "all")
	if clean.injected != 0 || clean.accuracy < 0.95 || clean.invalid {
		t.Fatalf("clean result: %+v", clean)
	}
}
