package text

import (
	"testing"

	"github.com/SCKelemen/unicode/uax50"
)

func TestDefaultVerticalTextStyle(t *testing.T) {
	style := DefaultVerticalTextStyle()

	if style.WritingMode != WritingModeHorizontalTB {
		t.Fatalf("WritingMode = %v, want %v", style.WritingMode, WritingModeHorizontalTB)
	}
	if style.TextOrientation != TextOrientationMixed {
		t.Fatalf("TextOrientation = %v, want %v", style.TextOrientation, TextOrientationMixed)
	}
	if style.TextCombineUpright != TextCombineUprightNone {
		t.Fatalf("TextCombineUpright = %v, want %v", style.TextCombineUpright, TextCombineUprightNone)
	}
	if style.GlyphOrientationVertical != 0 {
		t.Fatalf("GlyphOrientationVertical = %.1f, want 0", style.GlyphOrientationVertical)
	}
}

func TestCharOrientation(t *testing.T) {
	txt := NewTerminal()

	if got := txt.CharOrientation('A', TextOrientationUpright); got != uax50.Upright {
		t.Fatalf("CharOrientation(upright) = %v, want %v", got, uax50.Upright)
	}
	if got := txt.CharOrientation('A', TextOrientationSideways); got != uax50.Rotated {
		t.Fatalf("CharOrientation(sideways) = %v, want %v", got, uax50.Rotated)
	}
	if got := txt.CharOrientation('A', TextOrientationMixed); got != uax50.LookupOrientation('A') {
		t.Fatalf("CharOrientation(mixed) = %v, want %v", got, uax50.LookupOrientation('A'))
	}
}

func TestIsUprightAndIsRotated(t *testing.T) {
	txt := NewTerminal()

	uprightStyle := VerticalTextStyle{TextOrientation: TextOrientationUpright}
	if !txt.IsUpright('A', uprightStyle) {
		t.Fatal("IsUpright() = false, want true for upright style")
	}
	if txt.IsRotated('A', uprightStyle) {
		t.Fatal("IsRotated() = true, want false for upright style")
	}

	sidewaysStyle := VerticalTextStyle{TextOrientation: TextOrientationSideways}
	if !txt.IsRotated('A', sidewaysStyle) {
		t.Fatal("IsRotated() = false, want true for sideways style")
	}
	if txt.IsUpright('A', sidewaysStyle) {
		t.Fatal("IsUpright() = true, want false for sideways style")
	}
}

func TestWritingModeHelpers(t *testing.T) {
	if !IsVerticalWritingMode(WritingModeVerticalRL) {
		t.Fatal("IsVerticalWritingMode(VerticalRL) = false, want true")
	}
	if !IsVerticalWritingMode(WritingModeSidewaysLR) {
		t.Fatal("IsVerticalWritingMode(SidewaysLR) = false, want true")
	}
	if IsVerticalWritingMode(WritingModeHorizontalTB) {
		t.Fatal("IsVerticalWritingMode(HorizontalTB) = true, want false")
	}

	if !IsHorizontalWritingMode(WritingModeHorizontalTB) {
		t.Fatal("IsHorizontalWritingMode(HorizontalTB) = false, want true")
	}
	if IsHorizontalWritingMode(WritingModeVerticalLR) {
		t.Fatal("IsHorizontalWritingMode(VerticalLR) = true, want false")
	}
}
