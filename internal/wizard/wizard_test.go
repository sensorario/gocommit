package wizard

import "testing"

func TestCommitTypes_ContainsExpectedValues(t *testing.T) {
	expected := []string{"feat", "fix", "chore", "refactor"}
	if len(CommitTypes) != len(expected) {
		t.Fatalf("len(CommitTypes) = %d, want %d", len(CommitTypes), len(expected))
	}
	for i, v := range expected {
		if CommitTypes[i] != v {
			t.Errorf("CommitTypes[%d] = %q, want %q", i, CommitTypes[i], v)
		}
	}
}
