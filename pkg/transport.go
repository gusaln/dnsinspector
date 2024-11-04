package pkg

import "net"

func SendQuery(server string, query *Query) (*Response, error) {
	conn, err := net.Dial("udp", server+":53")
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	msg := query.AsBytes()
	if _, err = conn.Write(msg); err != nil {
		return nil, err
	}

	raw := make([]byte, 0, 1024)
	for {
		buf := [512]byte{}
		read, err := conn.Read(buf[:])
		if err != nil {
			return nil, err
		}

		raw = append(raw, buf[:read]...)

		if read < 512 {
			break
		}
	}

	var response *Response
	if response, err = ParsedResponse(raw); err != nil {
		return nil, err
	}

	return response, nil
}
