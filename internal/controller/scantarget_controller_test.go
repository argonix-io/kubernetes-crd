package controller

import (
	"testing"

	v1alpha1 "github.com/argonix-io/kubernetes-crd/api/v1alpha1"
)

func TestScanTargetAcceptsAPIDefaultsWithoutRepeatedPatches(t *testing.T) {
	for _, tc := range []struct {
		identifier string
		ref        string
		want       string
	}{
		{"registry.example.test:8443/service", "", "latest"},
		{"registry.example.test/service:v1", "ignored", ""},
		{"registry.example.test/service@sha256:abcd", "ignored", ""},
	} {
		t.Run(tc.identifier, func(t *testing.T) {
			noDays := []int{}
			target := &v1alpha1.ScanTarget{Spec: v1alpha1.ScanTargetSpec{
				Kind: "image", Identifier: tc.identifier, DefaultRef: &tc.ref,
				EnabledScanners: []string{"trivy", "syft"}, ScheduleDays: &noDays,
			}}
			desired := buildScanTargetPayload(target)
			remote := map[string]interface{}{
				"kind": "image", "identifier": tc.identifier, "default_ref": tc.want,
				"enabled_scanners": []interface{}{"trivy", "syft"}, "schedule_days": []interface{}{float64(0)},
				"frequency": "weekly",
			}
			normalizeScanTargetSchedule(desired, remote)
			if changes := driftedFields(desired, remote); len(changes) > 0 {
				t.Fatalf("normalized target should be in sync; unexpected changes %#v", changes)
			}
		})
	}
}
