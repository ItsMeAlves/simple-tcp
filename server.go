package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

func fixResource(seed string) string {
	var value string
	if seed == "/" {
		value = "index.html"
	} else {
		value = seed[1:len(seed)]
	}

	if !strings.Contains(value, ".") {
		value = strings.Join([]string{value, ".html"}, "")
	}

	return strings.Join([]string{"resources/", value}, "")
}

func handle(conn net.Conn) {
	line, err := bufio.NewReader(conn).ReadBytes('\n')
	seed := strings.Split(string(line), " ")
	if err == nil {
		fmt.Println("no errors...")

		if seed[0] == "GET" {
			go get(conn, fixResource(seed[1]))
		} else if seed[0] == "POST" {
			go post(conn, fixResource(seed[1]))
		} else if seed[0] == "PATCH" {
			go patch(conn, fixResource(seed[1]))
		} else if seed[0] == "DELETE" {
			go delete(conn, fixResource(seed[1]))
		} else {
			fmt.Println(strings.Join([]string{"Not Accepted:", seed[0]}, " "))
			conn.Close()
		}
	}
}

func get(conn net.Conn, source string) {
	defer conn.Close()
	var file []byte
	var err error

	file, err = ioutil.ReadFile(source)

	if err != nil {
		fmt.Println(strings.Join([]string{"cannot read", source}, " "))
		page, _ := ioutil.ReadFile("resources/404.html")
		conn.Write(page)
	} else {
		conn.Write(file)
	}
}

func post(conn net.Conn, source string) {
	response := strings.Join([]string{"Alguém está dando um POST aqui:", source}, " ")
	conn.Write([]byte(response))
	conn.Close()
}

func patch(conn net.Conn, source string) {
	response := strings.Join([]string{"Alguém está dando um PATCH aqui:", source}, " ")
	conn.Write([]byte(response))
	conn.Close()
}

func delete(conn net.Conn, source string) {
	response := strings.Join([]string{"Alguém está dando um DELETE aqui:", source}, " ")
	conn.Write([]byte(response))
	conn.Close()
}

func main() {
	server, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("algo deu errado")
		return
	}
	fmt.Println("Listening...")

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("algo deu errado")
			continue
		}
		go handle(conn)
	}
}
