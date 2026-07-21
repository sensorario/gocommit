package commit

import "testing"

func TestExtractJiraTicket(t *testing.T) {
	cases := []struct {
		branch string
		want   string
	}{
		// exactly one dash → matched
		{"feature/PROJ-123", "PROJ-123"},
		{"feature/AB-99", "AB-99"},
		{"hotfix/JIRA-1", "JIRA-1"},
		// multiple dashes in the ticket segment → not matched
		{"feature/PROJ-123-add-login", ""},
		{"hotfix/JIRA-1-fix-something", ""},
		// no slash, no ticket
		{"main", ""},
		{"develop", ""},
		// empty
		{"", ""},
	}

	for _, c := range cases {
		got := ExtractJiraTicket(c.branch)
		if got != c.want {
			t.Errorf("ExtractJiraTicket(%q) = %q, want %q", c.branch, got, c.want)
		}
	}
}
