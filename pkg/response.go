package pkg

type Response struct {
	ID                   uint16     `json:"id"`
	OpCode               uint16     `json:"opcode"`
	IsAuthoritative      bool       `json:"is_authoritative"`
	Truncated            bool       `json:"truncated"`
	RecursionDesired     bool       `json:"recursionDesired"`
	RecursionAvailable   bool       `json:"recursionAvailable"`
	ResponseCode         uint16     `json:"responseCode"`
	Questions            []Question `json:"questions"`
	Answers              []Record   `json:"answers"`
	AuthoritativeRecords []Record   `json:"authoritativeRecords"`
	AdditionalRecords    []Record   `json:"additionalRecords"`
}

func (res Response) IsSuccess() bool {
	return res.OpCode == 0
}

type RecordMap map[RecordType][]Record

// NewRecordMap creates an empty RecordMap
func NewRecordMap() RecordMap {
	return map[RecordType][]Record{}
}

// RecordsByType returns a map that groups records by type
func (res Response) RecordsByType() RecordMap {
	dict := map[RecordType][]Record{}

	for _, r := range res.Answers {
		if s, exists := dict[r.Rtype]; exists {
			dict[r.Rtype] = append(s, r)
		} else {
			dict[r.Rtype] = []Record{r}
		}
	}

	return dict
}

// Merge two RecordMaps into a new RecordMap
func (rm RecordMap) Merge(other RecordMap) RecordMap {
	return mergeRecordMaps(NewRecordMap(), rm, other)
}

// Merge a RecordMaps into this one
func (rm *RecordMap) MergeInto(other RecordMap) RecordMap {
	return mergeRecordMaps(*rm, other)
}

// mergeRecordMaps merges RecordMaps into a RecordMap
func mergeRecordMaps(into RecordMap, others ...RecordMap) RecordMap {
	for _, other := range others {
		for rT, otherSlice := range other {
			if s, exists := into[rT]; exists {
				into[rT] = append(s, otherSlice...)
			} else {
				s := make([]Record, len(otherSlice))
				copy(s, otherSlice)
				into[rT] = s
			}
		}
	}

	return into
}
