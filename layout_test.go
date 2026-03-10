package text

import "testing"

func TestNewMetrics_MeasureAndText(t *testing.T) {
	metrics := NewMetrics(Config{})
	if metrics == nil {
		t.Fatal("NewMetrics() returned nil")
	}
	if metrics.Text() == nil {
		t.Fatal("Metrics.Text() returned nil")
	}

	advance, ascent, descent := metrics.Measure("Hello世界", TextStyle{LineHeight: 2.0})
	if advance != 9.0 {
		t.Fatalf("Measure() advance = %.1f, want 9.0", advance)
	}
	if ascent != 1.6 {
		t.Fatalf("Measure() ascent = %.1f, want 1.6", ascent)
	}
	if descent != 0.4 {
		t.Fatalf("Measure() descent = %.1f, want 0.4", descent)
	}

	if metrics.Text().Width("Hello世界") != advance {
		t.Fatalf("Metrics.Text().Width() = %.1f, want %.1f", metrics.Text().Width("Hello世界"), advance)
	}
}

func TestMetricsMeasure_DefaultLineHeight(t *testing.T) {
	metrics := NewMetrics(Config{})

	_, ascent, descent := metrics.Measure("abc", TextStyle{})
	if ascent != 0.8 {
		t.Fatalf("Measure() ascent = %.1f, want 0.8", ascent)
	}
	if descent != 0.2 {
		t.Fatalf("Measure() descent = %.1f, want 0.2", descent)
	}
}

func TestNewTerminalMetrics(t *testing.T) {
	metrics := NewTerminalMetrics()
	if metrics == nil {
		t.Fatal("NewTerminalMetrics() returned nil")
	}

	advance, _, _ := metrics.Measure("世界", TextStyle{LineHeight: 1})
	if advance != 4.0 {
		t.Fatalf("Measure() advance = %.1f, want 4.0", advance)
	}
}
