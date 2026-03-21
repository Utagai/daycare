package picker

import (
	"testing"
)

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		text    string
		pattern string
		want    bool
	}{
		{"hello", "hello", true},
		{"hello", "helo", true},
		{"hello", "hll", true},
		{"hello", "xyz", false},
		{"hello", "el", true},
		{"hello", "le", false},
		{"", "anything", false},
		{"anything", "", true},
		{"", "", true},
		{"project", "proj", true},
		{"project", "prj", true},
		{"my-project-name", "mpn", true},
		{"alpha", "al", true},
		{"alpha", "ga", false},
		{"beta", "be", true},
		{"gamma", "ga", true},
	}

	for _, tt := range tests {
		t.Run(tt.text+"_"+tt.pattern, func(t *testing.T) {
			got := fuzzyMatch(tt.text, tt.pattern)
			if got != tt.want {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tt.text, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestStringItem(t *testing.T) {
	item := StringItem("test-project")
	if item.FilterValue() != "test-project" {
		t.Errorf("StringItem.FilterValue() = %v, want test-project", item.FilterValue())
	}
	if item.DisplayValue() != "test-project" {
		t.Errorf("StringItem.DisplayValue() = %v, want test-project", item.DisplayValue())
	}
}

func TestProjectItem(t *testing.T) {
	item := ProjectItem{Name: "myproject", Path: "/path/to/myproject"}
	if item.FilterValue() != "myproject" {
		t.Errorf("ProjectItem.FilterValue() = %v, want myproject", item.FilterValue())
	}
	if item.DisplayValue() != "myproject" {
		t.Errorf("ProjectItem.DisplayValue() = %v, want myproject", item.DisplayValue())
	}
}

func TestWorktreeItem(t *testing.T) {
	item := WorktreeItem{Name: "feature1", Branch: "feature1", Path: "/path/feature1", IsMain: false}
	if item.FilterValue() != "feature1" {
		t.Errorf("WorktreeItem.FilterValue() = %v, want feature1", item.FilterValue())
	}
	if item.DisplayValue() != "feature1" {
		t.Errorf("WorktreeItem.DisplayValue() = %v, want feature1", item.DisplayValue())
	}

	mainItem := WorktreeItem{Name: "main", Branch: "main", Path: "/path/main", IsMain: true}
	if mainItem.DisplayValue() != "main (main)" {
		t.Errorf("WorktreeItem.DisplayValue() for main = %v, want main (main)", mainItem.DisplayValue())
	}
}

func TestModelFilter(t *testing.T) {
	items := []Item{
		StringItem("alpha"),
		StringItem("beta"),
		StringItem("gamma"),
	}

	model := New(items, "Test", -1)

	// Initially all items should be available
	if len(model.filtered) != 3 {
		t.Errorf("Initial filtered length = %d, want 3", len(model.filtered))
	}

	// Filter by "alp"
	model.textInput.SetValue("alp")
	model.filter()
	if len(model.filtered) != 1 {
		t.Errorf("Filtered length with 'alp' = %d, want 1", len(model.filtered))
	}
	if model.filtered[0].FilterValue() != "alpha" {
		t.Errorf("Filtered item = %v, want alpha", model.filtered[0].FilterValue())
	}

	// Filter by "bet" - should match beta
	model.textInput.SetValue("bet")
	model.filter()
	if len(model.filtered) != 1 {
		t.Errorf("Filtered length with 'bet' = %d, want 1", len(model.filtered))
	}
	if model.filtered[0].FilterValue() != "beta" {
		t.Errorf("Filtered item = %v, want beta", model.filtered[0].FilterValue())
	}

	// Filter by non-matching
	model.textInput.SetValue("xyz")
	model.filter()
	if len(model.filtered) != 0 {
		t.Errorf("Filtered length with 'xyz' = %d, want 0", len(model.filtered))
	}
}

func TestModelCursorBounds(t *testing.T) {
	items := []Item{
		StringItem("one"),
		StringItem("two"),
		StringItem("three"),
	}

	// Start with invalid cursor (beyond bounds) - should clamp to 0
	model := New(items, "Test", 10)

	// Cursor should be clamped to 0 (not last index in current implementation)
	if model.cursor != 0 {
		t.Errorf("Initial cursor = %d, want 0", model.cursor)
	}

	// Start with valid cursor
	model = New(items, "Test", 1)
	if model.cursor != 1 {
		t.Errorf("Initial cursor = %d, want 1", model.cursor)
	}
}

func TestNew(t *testing.T) {
	items := []Item{StringItem("a"), StringItem("b"), StringItem("c")}

	model := New(items, "Select Item", -1)

	if model.title != "Select Item" {
		t.Errorf("title = %v, want Select Item", model.title)
	}

	if len(model.items) != 3 {
		t.Errorf("items length = %d, want 3", len(model.items))
	}

	if len(model.filtered) != 3 {
		t.Errorf("filtered length = %d, want 3", len(model.filtered))
	}

	if model.cursor != 0 {
		t.Errorf("cursor = %d, want 0", model.cursor)
	}
}

func TestModelSelected(t *testing.T) {
	items := []Item{StringItem("one"), StringItem("two")}

	model := New(items, "Test", 0)

	// Initially no selection
	if model.Selected() != nil {
		t.Error("Selected() should be nil initially")
	}

	// Set selected
	model.selected = items[0]
	if model.Selected() == nil {
		t.Error("Selected() should return the selected item")
	}
}

func TestModelCancelled(t *testing.T) {
	items := []Item{StringItem("one")}

	model := New(items, "Test", 0)

	// Initially not cancelled
	if model.Cancelled() {
		t.Error("Cancelled() should be false initially")
	}

	// Set cancelled
	model.cancelled = true
	if !model.Cancelled() {
		t.Error("Cancelled() should return true when set")
	}
}
