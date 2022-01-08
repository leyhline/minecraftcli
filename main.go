/*
minecraftcli - simple communication with Minecraft servers
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
	"fmt"
	"os"
	"time"
)

const TIMEOUT = 2

func main() {
	var address string
	address = os.Args[1]
	conn, err := DialTimeout(address, TIMEOUT * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	host, port, _ := SplitHostPort(address)
	handshake(conn, host, port)
	request(conn)
	_, slpRes := response(conn)
	fmt.Printf("%+v\n", slpRes)
}
