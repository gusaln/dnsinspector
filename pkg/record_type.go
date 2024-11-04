package pkg

import (
	"encoding/json"
	"fmt"
)

type RecordType uint8

const (
	// A host address
	RT_A RecordType = 1

	// An authoritative name server
	RT_NS = 2

	// A mail destination (Obsolete - use MX)
	RT_MD = 3
	// A mail forwarder (Obsolete - use MX)
	RT_MF = 4

	// The canonical name for an alias
	RT_CNAME = 5

	// Marks the start of a zone of authority
	RT_SOA = 6

	// A mailbox domain name (EXPERIMENTAL)
	RT_MB = 7

	// A mail group member (EXPERIMENTAL)
	RT_MG = 8

	// A mail rename domain name (EXPERIMENTAL)
	RT_MR = 9

	// A null RR (EXPERIMENTAL)
	RT_NULL = 10

	// A well known service description
	RT_WKS = 11

	// A domain name pointer
	RT_PTR = 12

	// Host information
	RT_HINFO = 13

	// Mailbox or mail list information
	RT_MINFO = 14

	// Mail exchange
	RT_MX = 15

	// Text strings
	RT_TXT = 16

	// A request for a transfer of an entire zone - Only for Questions
	RT_AXFR = 252

	// A request for mailbox-related records (MB, MG or MR) - Only for Questions
	RT_MAILB = 253

	// A request for mail agent RRs (Obsolete - see MX) - Only for Questions
	RT_MAILA = 254

	// A request for all records - Only for Questions
	RT_ALL = 255
)

func (rt RecordType) String() string {
	switch rt {
	case RT_A:
		return "A"
	case RT_NS:
		return "NS"
	case RT_MD:
		return "MD"
	case RT_MF:
		return "MF"
	case RT_CNAME:
		return "CNAME"
	case RT_SOA:
		return "SOA"
	case RT_MB:
		return "MB"
	case RT_MG:
		return "MG"
	case RT_MR:
		return "MR"
	case RT_NULL:
		return "NULL"
	case RT_WKS:
		return "WKS"
	case RT_PTR:
		return "PTR"
	case RT_HINFO:
		return "HINFO"
	case RT_MINFO:
		return "MINFO"
	case RT_MX:
		return "MX"
	case RT_TXT:
		return "TXT"
	case RT_AXFR:
		return "AXFR"
	case RT_MAILA:
		return "MAILA"
	case RT_MAILB:
		return "MAILB"
	case RT_ALL:
		return "ALL"
	default:
		return "UNKNOWN"
	}
}

func (rt RecordType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rt.String())
}

func (rt *RecordType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "A":
		*rt = RT_A
	case "NS":
		*rt = RT_NS
	case "MD":
		*rt = RT_MD
	case "MF":
		*rt = RT_MF
	case "CNAME":
		*rt = RT_CNAME
	case "SOA":
		*rt = RT_SOA
	case "MB":
		*rt = RT_MB
	case "MG":
		*rt = RT_MG
	case "MR":
		*rt = RT_MR
	case "NULL":
		*rt = RT_NULL
	case "WKS":
		*rt = RT_WKS
	case "PTR":
		*rt = RT_PTR
	case "HINFO":
		*rt = RT_HINFO
	case "MINFO":
		*rt = RT_MINFO
	case "MX":
		*rt = RT_MX
	case "TXT":
		*rt = RT_TXT
	case "AXFR":
		*rt = RT_AXFR
	case "MAILA":
		*rt = RT_MAILA
	case "MAILB":
		*rt = RT_MAILB
	case "ALL":
		*rt = RT_ALL
	default:
		return fmt.Errorf("unknown record type '%s'", s)
	}

	return nil
}

type RecordClass uint16

const (
	// Internet
	RC_IN RecordClass = 1

	// CSNET - (Obsolete - used only for examples in some obsolete RFCs)
	RC_CS = 2

	// CHAOS
	RC_CH = 3

	// Hesiod [Dyer 87]
	RC_HS = 4

	// RC_ALL - Only for Questions
	RC_ALL = 255
)

func (rc RecordClass) String() string {
	switch rc {
	case RC_IN:
		return "IN"
	case RC_CS:
		return "CS"
	case RC_CH:
		return "CH"
	case RC_HS:
		return "HS"
	case RC_ALL:
		return "ALL"
	default:
		return "UNKNOWN"
	}
}

func (rc RecordClass) MarshalJSON() ([]byte, error) {
	return json.Marshal(rc.String())
}

func (rc *RecordClass) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "IN":
		*rc = RC_IN
	case "CS":
		*rc = RC_CS
	case "CH":
		*rc = RC_CH
	case "HS":
		*rc = RC_HS
	case "ALL":
		*rc = RC_ALL
	default:
		return fmt.Errorf("unknown record class '%s'", s)
	}

	return nil
}
