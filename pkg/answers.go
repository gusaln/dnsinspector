package pkg

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

// Utility for parsing DNS responses.
type answerReader struct {
	cursor      int
	buf         []byte
	domainIndex map[int]string
}

func (r *answerReader) next(n int) []byte {
	b := r.buf[r.cursor : r.cursor+n]

	r.cursor += n

	return b
}

func (r *answerReader) nextOne() byte {
	b := r.buf[r.cursor]

	r.cursor += 1

	return b
}

func (r *answerReader) get(from, to int) []byte {
	return r.buf[from:to]
}

func (r *answerReader) unpackLabels() string {
	labels := make([]string, 0, 8)
	offsets := make([]int, 0, 8)

	for nameLength := int(r.nextOne()); nameLength > 0; nameLength = int(r.nextOne()) {
		if (nameLength & POINTER_MARKER) > 0 {
			offsets = append(offsets, r.cursor-1)

			// We make a one byte number into a two byte one...
			offset := (nameLength << 8) | int(r.nextOne())
			// and only use the relevant information
			offset &= OFFSET_MASK
			labels = append(labels, (r.domainIndex)[offset])
			break
		} else {
			offsets = append(offsets, r.cursor-1)
			label := string(r.next(int(nameLength)))
			labels = append(labels, label)
		}
	}

	for i := range labels {
		// TODO: Prevent duplicates
		(r.domainIndex)[offsets[i]] = strings.Join(labels[i:], ".")
	}

	return strings.Join(labels, ".")
}

func (r *answerReader) parseRecord() (*Record, error) {
	name := r.unpackLabels()

	rtype := RecordType(binary.BigEndian.Uint16(r.next(2)))
	rclass := RecordClass(binary.BigEndian.Uint16(r.next(2)))
	ttl := binary.BigEndian.Uint32(r.next(4))
	dataLength := int(binary.BigEndian.Uint16(r.next(2)))

	if rclass != RC_IN {
		return nil, fmt.Errorf("Record CLASS %d, not supported", rclass)
	}

	record := &Record{
		Name:   name,
		Rtype:  rtype,
		Rclass: rclass,
		Ttl:    ttl,
		Data:   map[string]interface{}{},
	}

	switch record.Rtype {
	case RT_A:
		rawdata := r.next(dataLength)
		record.Data["ip"] = net.IPv4(rawdata[0], rawdata[1], rawdata[2], rawdata[3])
	case RT_NS:
		record.Data["ns"] = r.unpackLabels()
	case RT_MD:
		record.Data["mailbox"] = r.unpackLabels()
	case RT_MF:
		record.Data["agent"] = r.unpackLabels()
	case RT_CNAME:
		record.Data["name"] = r.unpackLabels()
	case RT_SOA:
		record.Data["main"] = r.unpackLabels()
		record.Data["mailbox"] = r.unpackLabels()
		record.Data["serial"] = binary.BigEndian.Uint32(r.next(4))
		record.Data["refresh"] = binary.BigEndian.Uint32(r.next(4))
		record.Data["retry"] = binary.BigEndian.Uint32(r.next(4))
		record.Data["expire"] = binary.BigEndian.Uint32(r.next(4))
		record.Data["minimum"] = binary.BigEndian.Uint32(r.next(4))
	case RT_MB:
		record.Data["mailbox"] = r.unpackLabels()
	case RT_MG:
		record.Data["mailbox"] = r.unpackLabels()
	case RT_MR:
		record.Data["mailbox"] = r.unpackLabels()
	case RT_NULL:
		record.Data["raw"] = r.unpackLabels()
	case RT_WKS:
		rawdata := r.next(dataLength)
		record.Data["address"] = net.IPv4(rawdata[0], rawdata[1], rawdata[2], rawdata[3])
		// see the [RFC 1010](https://datatracker.ietf.org/doc/html/rfc1010)
		record.Data["protocol"] = rawdata[4]
		record.Data["bitmap"] = rawdata[4:]
	case RT_PTR:
		record.Data["ptr"] = r.unpackLabels()
	case RT_HINFO:
		record.Data["info"] = string(r.next(dataLength))
	case RT_MINFO:
		record.Data["mailbox"] = r.unpackLabels()
		record.Data["error_mailbox"] = r.unpackLabels()
	case RT_MX:
		record.Data["preference"] = binary.BigEndian.Uint16(r.next(2))
		record.Data["exchange"] = r.unpackLabels()
	case RT_TXT:
		record.Data["data"] = string(r.next(dataLength))

	default:
		return nil, fmt.Errorf("unknown record type %d", rtype)
	}

	return record, nil
}

const POINTER_MARKER = 0b11000000
const OFFSET_MASK = 0b0011111111111111

func ParsedResponse(buf []byte) (*Response, error) {
	domainIndex := make(map[int]string)
	r := answerReader{
		cursor:      0,
		buf:         buf,
		domainIndex: domainIndex,
	}

	id := binary.BigEndian.Uint16(r.next(2))
	header := binary.BigEndian.Uint16(r.next(2))

	qr := ((header >> 15) & 1)
	opcode := ((header >> 11) & 0b1111)
	authoritativeAnswer := ((header >> 10) & 1) > 0
	truncated := ((header >> 9) & 1) > 0
	recursionDesired := ((header >> 8) & 1) > 0
	recursionAvailable := ((header >> 7) & 1) > 0
	z := ((header >> 4) & 0b111)
	responseCode := ((header) & 0b1111)

	if qr != 1 {
		return nil, fmt.Errorf("response has the header of a query (first bit should be 1, and was %d)", qr)
	}

	if z != 0 {
		return nil, fmt.Errorf("the Z part of the header was not zero")
	}

	questionCount := binary.BigEndian.Uint16(r.next(2))
	answerCount := binary.BigEndian.Uint16(r.next(2))
	nsCount := binary.BigEndian.Uint16(r.next(2))
	additionalRecordsCount := binary.BigEndian.Uint16(r.next(2))

	questions := make([]Question, questionCount)
	answersRR := make([]Record, answerCount)
	authoritativeRR := make([]Record, nsCount)
	additionalRR := make([]Record, additionalRecordsCount)

	// Questions
	for i := 0; i < int(questionCount); i++ {
		name := r.unpackLabels()

		rtype := RecordType(binary.BigEndian.Uint16(r.next(2)))
		rclass := RecordClass(binary.BigEndian.Uint16(r.next(2)))

		if rclass != RC_IN {
			return nil, fmt.Errorf("Record CLASS %d, not supported", rclass)
		}

		questions[i] = Question{
			Name:  name,
			Rtype: rtype,
		}
	}

	// Answers
	for i := 0; i < int(answerCount); i++ {
		record, err := r.parseRecord()
		if err != nil {
			return nil, fmt.Errorf("could not parse record at %d: %w", r.cursor, err)
		}
		answersRR[i] = *record
	}

	// Authoritative
	for i := 0; i < int(nsCount); i++ {
		record, err := r.parseRecord()
		if err != nil {
			return nil, fmt.Errorf("could not parse record at %d: %w", r.cursor, err)
		}
		authoritativeRR[i] = *record
	}

	// Additional
	for i := 0; i < int(additionalRecordsCount); i++ {
		record, err := r.parseRecord()
		if err != nil {
			return nil, fmt.Errorf("could not parse record at %d: %w", r.cursor, err)
		}
		additionalRR[i] = *record
	}

	return &Response{
		ID:                   id,
		OpCode:               opcode,
		IsAuthoritative:      authoritativeAnswer,
		Truncated:            truncated,
		RecursionDesired:     recursionDesired,
		RecursionAvailable:   recursionAvailable,
		ResponseCode:         responseCode,
		Questions:            questions,
		Answers:              answersRR,
		AuthoritativeRecords: authoritativeRR,
		AdditionalRecords:    additionalRR,
	}, nil
}
