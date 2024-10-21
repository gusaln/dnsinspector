package pkg

type Record struct {
	Name   string                 `json:"name"`
	Rtype  RecordType             `json:"rtype"`
	Rclass RecordClass            `json:"rclass"`
	Ttl    uint32                 `json:"ttl"`
	Data   map[string]interface{} `json:"data"`
}

func (r Record) Has(k string) bool {
	_, present := r.Data[k]

	return present
}

func (r Record) Get(k string) (interface{}, bool) {
	datum, present := r.Data[k]

	return datum, present
}
