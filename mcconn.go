/*
mcconn.go - methods for abstracting Minecraft server communication
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
	"net"
	"encoding/binary"
	"bufio"
	"io"
	"time"
	"strconv"
)

type MCConn struct {
	net.Conn
}

func (c MCConn) Write(b []byte) (int, error) {
	lenBuf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(lenBuf, uint64(len(b)))
	lenBuf = append(lenBuf[:n], b...)
	return c.Conn.Write(lenBuf)
}

func (c MCConn) Read(b []byte) (int, error) {
	r := bufio.NewReader(c.Conn)
	resLen, err := binary.ReadUvarint(r)
	if err != nil {
		// TODO: Adjust error message
		return 0, err
	}
	res := make([]byte, resLen)
	n, err := io.ReadFull(r, res)
	if err != nil {
		// TODO: Error message
		return n, err
	}
	copy(b, res)
	return n, nil
}

func DialTimeout(address string, timeout time.Duration) (MCConn, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		// TODO
	}
	return MCConn{conn}, err
}

func SplitHostPort(hostport string) (string, uint16, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", 0, err
	}
	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return "", 0, err
	}
	return host, uint16(portInt), nil
}
