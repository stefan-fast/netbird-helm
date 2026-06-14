package logic

import "testing"

func TestCheckChartVersion(t *testing.T) {
	t.Run("matching version", func(t *testing.T) {
		content := []byte("apiVersion: v2\nname: testchart\nversion: 1.2.3\n")
		if err := CheckChartVersion(content, "1.2.3", "Chart.yaml"); err != nil {
			t.Errorf("CheckChartVersion unexpected error: %v", err)
		}
	})

	t.Run("mismatched version", func(t *testing.T) {
		content := []byte("apiVersion: v2\nname: testchart\nversion: 1.2.3\n")
		if err := CheckChartVersion(content, "9.9.9", "Chart.yaml"); err == nil {
			t.Error("CheckChartVersion expected error for mismatched version, got nil")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		content := []byte("not: valid: yaml: [\n")
		if err := CheckChartVersion(content, "1.0.0", "Chart.yaml"); err == nil {
			t.Error("CheckChartVersion expected error for invalid YAML, got nil")
		}
	})
}
