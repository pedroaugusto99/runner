package output

type Detail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type ResultError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []Detail `json:"details"`
}

type Result struct {
	Success   bool           `json:"success"`
	Operation string         `json:"operation"`
	Signature string         `json:"signature"`
	Valid     *bool          `json:"valid"`
	Message   string         `json:"message"`
	Output    string         `json:"output"`
	Error     *ResultError   `json:"error"`
	Metadata  map[string]any `json:"metadata"`
}

func (r *Result) ExitCode() int {
	if r.Success {
		if r.Operation == "validate" && r.Valid != nil && !*r.Valid {
			return 1
		}
		return 0
	}
	if r.Operation == "validate" {
		return 2
	}
	return 1
}
