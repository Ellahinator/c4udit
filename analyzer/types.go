package analyzer

import (
	"fmt"
	"regexp"
	"strings"
)

// Report is the end result of an analysis containing the files analyzed,
// the issues searched for and a map of findings per issue.
type Report struct {
	Issues        []Issue
	FilesAnalyzed []string
	// Key is Issue Identifier
	FindingsPerIssue map[string][]Finding
}

// Issue represents an Issue to search for in the codebase.
// The pattern field is a RegEx string which must compile.
type Issue struct {
	Identifier     string
	Severity       Severity
	Title          string
	Impact         string
	Pattern        string
	Recommendation string
}

// Finding represents a possible Issue found in the codebase.
type Finding struct {
	IssueIdentifier string
	File            string
	LineNumber      int
	LineContent     string
}

// Severity type defining the severity level for an Issue.
type Severity int

// The Severity Enum.
const (
	GASOP Severity = iota
	NC
	LOW
)

func slugify(s string) string {
	droppedChars := []string{
		"\"", "'", "`", ".", "/",
		"!", ",", "~", "&",
		"%", "^", "*", "#",
		"@", "|",
		"(", ")",
		"{", "}",
		"[", "]", "<", ">", "++", "==", "+", "=",
	}

	s = strings.ToLower(s)
	// output still is L01, need to add dash so its L-01
	for _, c := range droppedChars {
		s = strings.Replace(s, c, "", -1)
	}

	s = strings.Replace(s, " ", "-", -1)

	// fmt.Println("--------------------------------")
	// fmt.Println(s)

	// s = strings.Replace(s, "--", "-", -1)
	m1 := regexp.MustCompile(`[a-z]-[0-9][0-9]--`)
	m2 := regexp.MustCompile(`[a-z]-[0-9][0-9]`)
	num := m2.Find([]byte(s))
	// fmt.Println(num)
	s = m1.ReplaceAllString(s, string(num)+"-")
	// fmt.Println(s)

	// fmt.Println("--------------------------------")

	return s
}

func createLink(identifier string, title string) string {
	link := fmt.Sprintf("#%s -%s", identifier, title)
	link = slugify(link)
	return "#" + link
}

// Markdown returns the report as string in markdown style.
func (r Report) Markdown(toc bool) string {
	// Issue output in Code4Rena format:
	// ### {{ issue.Title }}
	//
	// #### Impact
	// Issue information: [{{ issue.Identifier }}]({{ issue.Impact }})
	//
	// #### Findings
	// {{ _, finding := range findings: finding.String() }}
	//
	// #### Recommendation
	// {{issue.Recommendation}}
	//
	// #### Tools used
	//

	buf := strings.Builder{}

	buf.WriteString("# c4udit Report\n")
	buf.WriteString("\n")

	buf.WriteString("## Files analyzed\n")
	for _, f := range r.FilesAnalyzed {
		buf.WriteString("- " + f + "\n")
	}

	//low and Non-Critical
	//links

	//low
	if toc {
		buf.WriteString("# Table of Contents \n")

		canWriteLowTitle := true
		for _, issue := range r.Issues {

			findings := r.FindingsPerIssue[issue.Identifier]
			if len(findings) == 0 {
				continue
			}
			if issue.Severity == LOW {
				if canWriteLowTitle {
					buf.WriteString("Low\n")
				}
				link := createLink(issue.Identifier, issue.Title)
				buf.WriteString("- [[" + issue.Identifier + "] " + issue.Title + "](" + link + ")\n")
				canWriteLowTitle = false
			} else {
				continue
			}

		}

		//Nc
		canWriteNctitle := true
		for _, issue := range r.Issues {

			findings := r.FindingsPerIssue[issue.Identifier]
			if len(findings) == 0 {
				continue
			}
			if issue.Severity == NC {
				if canWriteNctitle {
					buf.WriteString("\nNon-Critical\n")
				}
				link := createLink(issue.Identifier, issue.Title)
				buf.WriteString("- [[" + issue.Identifier + "] " + issue.Title + "](" + link + ")\n")
				canWriteNctitle = false
			} else {
				continue
			}

		}

		buf.WriteString("\n")

	}
	buf.WriteString("## QA Issues found\n")
	buf.WriteString("\n")
	//list  low and nc
	buf.WriteString("## Low Findings\n")
	buf.WriteString("\n")
	for _, issue := range r.Issues {
		findings := r.FindingsPerIssue[issue.Identifier]
		if len(findings) == 0 {
			continue
		}
		if issue.Severity == GASOP {
			continue
		}

		if issue.Severity == NC {
			continue
		}

		buf.WriteString("### [" + issue.Identifier + "] " + issue.Title + "\n")

		// Impact
		// if issue.Severity == NC {
		if issue.Impact != "" {
			buf.WriteString("#### Impact\n")
			buf.WriteString(issue.Impact + "\n")
		}

		// Findings
		buf.WriteString("#### Findings:\n")
		buf.WriteString("```solidity\n")
		for _, finding := range findings {
			buf.WriteString(finding.String())
		}
		buf.WriteString("```\n")

		// Recommendation
		buf.WriteString("#### Recommendation\n")
		buf.WriteString(issue.Recommendation + "\n")
		buf.WriteString("\n")
	}

	buf.WriteString("## Non-Critical Findings\n")
	buf.WriteString("\n")
	for _, issue := range r.Issues {
		findings := r.FindingsPerIssue[issue.Identifier]
		if len(findings) == 0 {
			continue
		}
		if issue.Severity == GASOP {
			continue
		}

		if issue.Severity == LOW {
			continue
		}

		buf.WriteString("### [" + issue.Identifier + "] " + issue.Title + "\n")

		// Impact
		// if issue.Severity == NC {
		if issue.Impact != "" {
			buf.WriteString("#### Impact\n")
			buf.WriteString(issue.Impact + "\n")
		}

		// Findings
		buf.WriteString("#### Findings:\n")
		buf.WriteString("```solidity\n")
		for _, finding := range findings {
			buf.WriteString(finding.String())
		}
		buf.WriteString("```\n")

		// Recommendation
		buf.WriteString("#### Recommendation\n")
		buf.WriteString(issue.Recommendation + "\n")
		buf.WriteString("\n")
	}

	///gas

	if toc {

		buf.WriteString("# Table of Contents \n")

		canWriteGasTitle := true
		for _, issue := range r.Issues {

			findings := r.FindingsPerIssue[issue.Identifier]
			if len(findings) == 0 {
				continue
			}
			if issue.Severity == GASOP {
				if canWriteGasTitle {
					buf.WriteString("Gas\n")
				}
				link := createLink(issue.Identifier, issue.Title)
				buf.WriteString("- [[" + issue.Identifier + "] " + issue.Title + "](" + link + ")\n")
				canWriteGasTitle = false
			} else {
				continue
			}

		}

		buf.WriteString("\n")
	}

	buf.WriteString("## Gas Findings\n")
	buf.WriteString("\n")
	//list  gas
	for _, issue := range r.Issues {
		findings := r.FindingsPerIssue[issue.Identifier]
		if len(findings) == 0 {
			continue
		}
		if issue.Severity != GASOP {
			continue
		}

		buf.WriteString("### [" + issue.Identifier + "] " + issue.Title + "\n")

		// Impact
		if issue.Impact != "" {
			buf.WriteString("#### Impact\n")
			buf.WriteString(issue.Impact + "\n")
		}
		// Findings
		buf.WriteString("#### Findings:\n")
		buf.WriteString("```solidity\n")
		for _, finding := range findings {
			buf.WriteString(finding.String())
		}
		buf.WriteString("```\n")

		// Recommendation
		buf.WriteString("#### Recommendation\n")
		buf.WriteString(issue.Recommendation + "\n")
		buf.WriteString("\n")
	}

	// Tools used
	buf.WriteString("#### Tools used\n")
	buf.WriteString("manual, c4udit, slither" + "\n")

	buf.WriteString("\n")

	return buf.String()
}

func (r Report) String() string {
	// Build files string.
	files := "Files analyzed:\n"
	for _, f := range r.FilesAnalyzed {
		files += fmt.Sprintf("- %s\n", f)
	}
	files += "\n"

	// Build issues string.
	issues := "Issues found:\n"
	for i, issue := range r.Issues {
		// Get findings for issue
		findings := r.FindingsPerIssue[issue.Identifier]

		// Skip if no findings
		if len(findings) == 0 {
			continue
		}

		// Add findings per issue
		issues += " " + issue.Identifier + ":\n"
		for _, finding := range findings {
			issues += "  " + finding.String()
		}

		// Add newline if not last issue
		if i+1 != len(r.Issues) {
			issues += "\n"
		}
	}

	return files + issues
}

func (i Issue) String() string {
	return i.Identifier
}

func (f Finding) String() string {
	return fmt.Sprintf("%s::%d => %s\n", f.File, f.LineNumber, f.LineContent)
}

func (s Severity) String() string {
	return []string{
		"Gas Optimization",
		"Non-Critical",
		"Low Risk",
	}[s]
}
