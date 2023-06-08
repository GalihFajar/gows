package socket

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"testing"
)

type MockConn struct {
	mock.Mock
	net.Conn
}

func (m *MockConn) Read(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Write(b []byte) (int, error) {
	_ = m.Called(b)
	return 0, nil
}

func (m *MockConn) Close() error {
	m.Called()
	return nil
}

func Test_parseRequest(test *testing.T) {
	body := "hello"
	parsedRequest, err := parseRequest([]byte(mockPostRequest(body)))

	assert.Equal(test, body, parsedRequest.Body, "parsed body should be equal to given argument to mock request")
	assert.Nil(test, err, "error should be nil")

	parsedRequest, err = parseRequest([]byte(""))

	assert.Nil(test, parsedRequest, "parsed request should be nil")
	assert.NotNil(test, err, "error should not be nil")
}

func Test_processClient(test *testing.T) {

	testObj := new(MockConn)
	testObj.On("Read", mock.AnythingOfType("[]uint8")).Return(len(mockPostRequest("hello")), nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).([]byte)
		for i, val := range []byte(mockPostRequest("hello")) {
			arg[i] = val
		}
	})

	testObj.On("Write", mock.Anything)
	testObj.On("Close", mock.Anything)

	processClient(testObj)

	testObj.AssertExpectations(test)
}

func mockPostRequest(body string) string {
	header := fmt.Sprintf("POST / HTTP/1.1\nContent-Type: plain\nContent-Length: %d\n\n", len(body))
	return fmt.Sprintf(header + body)
}
