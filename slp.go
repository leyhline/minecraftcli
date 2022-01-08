/*
slp.go - Server List Ping
Copyright (C) 2021  Thomas Leyh <thomas.leyh@mailbox.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/binary"
	"bufio"
	"io"
	"encoding/json"
)

const PROTOCOL_VERSION = 0

type SLPResponse struct {
	Description struct {
		Text string
	}
	Players struct {
		Max int
		Online int
		Sample []struct {
			Id string
			Name string
		}
	}
	Version struct {
		Name string
		Protocol int
	}
}

func handshake(conn MCConn, host string, port uint16) {
	// TODO Error when host is too long (>1024)
	b := make([]byte, 1024)
	b[0] = 0x00
	n := 1
	n += binary.PutUvarint(b[n:], PROTOCOL_VERSION)
	n += binary.PutUvarint(b[n:], uint64(len(host)))
	n += copy(b[n:], host)
	binary.BigEndian.PutUint16(b[n:], port)
	n += 2
	b[n] = 0x01
	n++
	conn.Write(b[:n])
}

func request(conn MCConn) {
	b := []byte{0x00}
	conn.Write(b)
}

func response(conn MCConn) ([]byte, SLPResponse) {
	r := bufio.NewReader(conn)
	r.ReadByte() // Read the Packet ID which should be 0x00 (unchecked)
	dataLen, _ := binary.ReadUvarint(r)
	data := make([]byte, dataLen)
	_, err := io.ReadFull(r, data)
	if err != nil {
		// TODO
	}
	var slpRes SLPResponse
	err = json.Unmarshal(data, &slpRes)
	return data, slpRes
}
