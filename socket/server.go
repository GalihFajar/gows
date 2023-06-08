package socket

import (
	"bufio"
	"fmt"
	"github.com/GalihFajar/gows/constant"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Method string
	Body   string
}

func StartServer() {
	fmt.Println("Server Running...")
	server, err := net.Listen(constant.SERVER_TYPE, constant.SERVER_HOST+":"+constant.SERVER_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer server.Close()

	fmt.Println("Listening on " + constant.SERVER_HOST + ":" + constant.SERVER_PORT)
	fmt.Println("Waiting for client...")

	for {
		connection, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("client connected")
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024*8)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	req, err := parseRequest(buffer[:mLen])

	switch req.Method {
	case "POST":
		handlePost()
	case "GET":
		handleGet()
	default:
		fmt.Println("Unknown request")
	}

	_, err = connection.Write([]byte("pong"))
	if err != nil {
		fmt.Println("Error sending response:", err.Error())
	}
	connection.Close()
}

func parseRequest(buffer []byte) (*Request, error) {
	r := strings.NewReader(string(buffer))
	buf := bufio.NewReader(r)
	req, err := http.ReadRequest(buf)

	if err != nil {
		return nil, fmt.Errorf("error reading request: %s", err)
	}

	reqBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading request: %s", err)
	}

	return &Request{
		Method: req.Method,
		Body:   string(reqBytes),
	}, nil
}

func handlePost() {
	fmt.Println("Received a POST request")
}

func handleGet() {
	fmt.Println("Received a GET request")
}
