package main

// Imports
import (
	"bufio"         // Read buffer strings
	"fmt"           // Format strings and prints
	"io/ioutil"     // File reading
	"net"           // Baasic TCP handler
	"strings"       // String helper functions
)

// Apply patch to HTTP resource
func fixResource(seed string) string {
	var value string

    // If the homepage ("/") is requested, returns the "index.html" resource

    if seed == "/" {
		value = "index.html"
	} else {
		value = seed[1:len(seed)]   // Remove the trailing bar to handle its data
	}

    // If the resource requested does not have ".", it means it requested a url without file extension
    // In this case, lets assume its a html file, so append to the resource the extension ".html"
	if !strings.Contains(value, ".") {
		value = strings.Join([]string{value, ".html"}, "")
	}

    // Since all files are inside resources directory, insert "resources/" first and return that url
    // Example: if user requested "/someUrl", it should return "resources/someUrl.html"
	return strings.Join([]string{"resources/", value}, "")
}

// Handling all incoming connections (received as a parameter)
func handle(conn net.Conn) {

    // Read first line from the HTTP request as an array of bytes
	line, err := bufio.NewReader(conn).ReadBytes('\n')

    // Split by spaces
    // Example, if the request "GET / HTTP/1.1", seed should be ["GET", "/", "HTTP/1.1"]
    seed := strings.Split(string(line), " ")

    // Check if everything is ok
	if err == nil {
		fmt.Println("no errors...")

        // Check HTTP method requested and use the adequate function
        // Example: If HTTP method is GET, execute the get function and so on
		if seed[0] == "GET" {
			go get(conn, fixResource(seed[1]))
		} else if seed[0] == "POST" {
			go post(conn, fixResource(seed[1]))
		} else if seed[0] == "PATCH" {
			go patch(conn, fixResource(seed[1]))
		} else if seed[0] == "DELETE" {
			go delete(conn, fixResource(seed[1]))
		} else {
            // In case the user requested an unsupported method, tell it and closes connection
			fmt.Println(strings.Join([]string{"Not Accepted:", seed[0]}, " "))
			conn.Close()
		}
	}
}

// Response for a GET request
func get(conn net.Conn, source string) {
	defer conn.Close()  // Close connection at the end of function
	var file []byte
	var err error

    // Read the file from the source requested
	file, err = ioutil.ReadFile(source)
    
	if err != nil {
        // In case the file was not found
		fmt.Println(strings.Join([]string{"cannot read", source}, " "))
		page, _ := ioutil.ReadFile("resources/404.html")
		conn.Write(page)
	} else {
        // File found, send it to the user
		conn.Write(file)
	}
}

// Response for a POST request
func post(conn net.Conn, source string) {
    // Just send a string back
	response := strings.Join([]string{"Alguém está dando um POST aqui:", source}, " ")
	conn.Write([]byte(response))
	conn.Close()
}

// Response for a PATCH reqeust
func patch(conn net.Conn, source string) {
    // Just send a string back
	response := strings.Join([]string{"Alguém está dando um PATCH aqui:", source}, " ")
	conn.Write([]byte(response))
	conn.Close()
}

// Response for a DELELE request
func delete(conn net.Conn, source string) {
    // Just send a string back
	response := strings.Join([]string{"Alguém está dando um DELETE aqui:", source}, " ")
	conn.Write([]byte(response))
	conn.Close()
}

// Program entry point
func main() {
    // Open server at the port 80 using tcp
	server, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("algo deu errado")
		return
	}
    // Everything went good
	fmt.Println("Listening...")

    // Infinite loop
	for {
        // Accept every incoming request
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("algo deu errado")
			continue
		}
        // Handle the request
		go handle(conn)
	}
}
