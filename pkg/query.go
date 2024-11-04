package pkg

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
)

type Query struct {
	id               uint16
	recursionDesired bool
	questions        []Question
}

func (q Query) Id() string {
	return fmt.Sprintf("%x", q.id)
}

func (q Query) IsRecursionDesired() bool {
	return q.recursionDesired
}

func (q Query) QuestionCount() int {
	return len(q.questions)
}

func (q Query) Questions() []Question {
	cp := make([]Question, len(q.questions))
	copy(cp, q.questions)
	return q.questions
}

func appendDomainAsLabels(buf []byte, domain string) []byte {
	parts := strings.Split(domain, ".")
	for _, p := range parts {
		buf = append(buf, byte(len(p)))
		buf = append(buf, []byte(p)...)
	}

	buf = append(buf, 0)

	return buf
}

func (q Query) AsBytes() []byte {
	domainMap := make(map[string]int)

	msg := make([]byte, 0, 64)
	msg = binary.BigEndian.AppendUint16(msg, q.id)

	// Header
	var header uint16
	if q.recursionDesired {
		header |= 1 << 8
	}
	msg = binary.BigEndian.AppendUint16(msg, header)

	// QDCOUNT
	msg = binary.BigEndian.AppendUint16(msg, uint16(len(q.questions)))
	// ANCOUNT
	msg = append(msg, 0, 0)
	// NSCOUNT
	msg = append(msg, 0, 0)
	// ARCOUNT
	msg = append(msg, 0, 0)

	for _, qstn := range q.questions {
		if offset, exists := domainMap[qstn.Name]; exists {
			msg = binary.BigEndian.AppendUint16(msg, uint16(offset))
			// mark it as a label
			msg[len(msg)-2] |= 0xc0
		} else {
			domainMap[qstn.Name] = len(msg)
			msg = appendDomainAsLabels(msg, qstn.Name)
		}
		// QTYPE
		msg = binary.BigEndian.AppendUint16(msg, uint16(qstn.Rtype))
		// QCLASS is always IN = 1
		msg = binary.BigEndian.AppendUint16(msg, uint16(RC_IN))
	}

	return msg
}

func (q Query) AsTcpBytes() []byte {
	rawmsg := q.AsBytes()

	// if len(rawmsg)-2 > 1<<16-1 {
	// 	return nil, fmt.Errorf("the message is too long to send through TCP")
	// }

	msg := append([]byte{0, 0}, rawmsg...)
	binary.BigEndian.PutUint16(msg, uint16(len(msg)-2))

	return msg
}

/// Query building

// Facilitates the task of building a query
type QueryBuilder struct {
	recursionDesired bool
	questions        []Question
}

// Creates a new QueryBuilder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		questions: make([]Question, 0),
	}
}

func (b *QueryBuilder) RecursionDesired() *QueryBuilder {
	b.recursionDesired = true

	return b
}

func (b *QueryBuilder) RecursionNotDesired() *QueryBuilder {
	b.recursionDesired = false

	return b
}

func (b *QueryBuilder) AddQuestion(domain string, qtype RecordType) *QueryBuilder {
	b.questions = append(b.questions, Question{
		Name:  strings.ToLower(strings.TrimSpace(domain)),
		Rtype: qtype,
	})

	return b
}

// Builds the query with a unique ID
func (b QueryBuilder) Build() *Query {
	id := [2]byte{}
	rand.Read(id[:])

	questions := make([]Question, len(b.questions))
	copy(questions, b.questions)

	return &Query{
		id:               binary.NativeEndian.Uint16(id[:]),
		recursionDesired: b.recursionDesired,
		questions:        b.questions,
	}
}

// Creates a new query without recursion
func NewQuery(domain string, qtype RecordType) *Query {
	return NewQueryBuilder().AddQuestion(domain, qtype).Build()
}

// Creates a new query with recursion
func NewQueryWithRecursion(domain string, qtype RecordType) *Query {
	return NewQueryBuilder().RecursionDesired().AddQuestion(domain, qtype).Build()
}
