package commit

import "strings"

func ExtractJiraTicket(branchName string) string {
	branchParts := strings.Split(branchName, "/")
	for _, part := range branchParts {
		if len(part) > 0 {
			if len(part) > 4 && strings.Contains(part, "-") {
				pieces := strings.Split(part, "-")
				if len(pieces) == 2 && len(pieces[0]) >= 2 && len(pieces[1]) > 0 {
					return part
				}
			}
		}
	}
	return ""
}
