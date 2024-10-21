package pkg

type Response struct {
	ID                 uint16     `json:"id"`
	OpCode             uint16     `json:"opcode"`
	Authoritative      bool       `json:"authoritative"`
	Truncated          bool       `json:"truncated"`
	RecursionDesired   bool       `json:"recursionDesired"`
	RecursionAvailable bool       `json:"recursionAvailable"`
	ResponseCode       uint16     `json:"responseCode"`
	Questions          []Question `json:"questions"`
	Answers            []Record   `json:"answers"`
}
