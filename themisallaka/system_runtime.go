package themisallaka

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/thaianhsoft/gothemis/themiscore/container"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

var baseAllaSystem *AllaSystem

type AllaSystem struct {
	*AllaKator
	containerAllakators *container.SplayTree
	globalProcessId     process
	dispatcher          *Dispatcher
	pool                *Pool
	possiblePort        int
	signalAddress       string
	containerRemoteHost map[net.Conn]bool
}

func NewAllaSystem() *AllaSystem {
	if baseAllaSystem != nil {
		return baseAllaSystem
	}
	baseAllaSystem = &AllaSystem{
		containerAllakators: &container.SplayTree{},
		globalProcessId:     1,
		dispatcher:          NewDispatcher(100),
		possiblePort:        8000,
		signalAddress:       ":8000",
		containerRemoteHost: map[net.Conn]bool{},
	}
	baseAllakator := newAllaKator(baseAllaSystem.globalProcessId, nil)
	baseAllaSystem.AllaKator = baseAllakator
	baseAllaSystem.containerAllakators.Insert(uint64(baseAllaSystem.globalProcessId), baseAllakator)
	return baseAllaSystem
}

func (as *AllaSystem) nextProcessId() process {
	v := as.globalProcessId + 1
	as.globalProcessId++
	return v
}

func (as *AllaSystem) WaitSignal() {
	signalAddress := "localhost:8000"
	as.signalAddress = signalAddress
	local, _ := net.ResolveTCPAddr("tcp", signalAddress)
	listener, _ := net.ListenTCP("tcp", local)
	for {
		//fmt.Println("waiting for signal")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		if _, ok := as.containerRemoteHost[conn]; !ok {
			as.containerRemoteHost[conn] = true
			fmt.Println("conn: ", conn, " is connected success !")
		}
		as.handle(conn)
	}
}

func (as *AllaSystem) ConnectRemoteServer(host string) {
	addr, _ := net.ResolveTCPAddr("tcp", host)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err == nil {
		fmt.Println("connect server success")
		as.handle(conn)
	} else {
		fmt.Println("connect failed: ", err)
	}
}

func (as *AllaSystem) serveWriter(rw io.ReadWriteCloser, chanMessage *chan []byte) {
	defer rw.Close()
	writer := bufio.NewWriter(rw)
	for {
		select {
		case message := <-*chanMessage:
			if string(message) == "closed" {
				fmt.Println("I am writer, exit thread handle from connection be closed!!")
				return
			}
		default:
			scanner := bufio.NewReader(os.Stdin)
			line, err := scanner.ReadString('\n')
			if err != nil {
				continue
			}
			message := strings.ReplaceAll(line, "\n", "")
			if _, err := writer.Write([]byte(message)); err != nil {
				fmt.Println("write message to connection failed: ", err)
			} else {
				writer.Flush()
			}
		}
	}
}

func (as *AllaSystem) serveReader(rw io.ReadWriteCloser, chanMessage *chan []byte) {
	//fmt.Println("serve reader running")
	defer rw.Close()
	buf := make([]byte, 1024)
	reader := bufio.NewReader(rw)
	connectionClosedSignal := "closed"
	for {
		n, err := reader.Read(buf)
		switch err {
		case io.EOF:
			// close connection from client is occured
			*chanMessage <- []byte(connectionClosedSignal)
			break
		case nil:
			t := time.Now()
			hour := t.Hour()
			minute := t.Minute()
			second := t.Second()
			c := fmt.Sprintf("%v:%v:%v", hour, minute, second)
			fmt.Printf("->>remote message [time=%v]: %v\n", c, string(buf[0:n]))
		}
	}
	fmt.Println("I am reader, exit thread handle from connection be closed !!")
}

type InternalReadWriteCloser struct {
	*bytes.Buffer
}

func (i *InternalReadWriteCloser) Close() error {
	return nil
}

func (as *AllaSystem) handle(rw io.ReadWriteCloser) {
	chanMessage := make(chan []byte)
	go as.serveWriter(rw, &chanMessage)
	go as.serveReader(rw, &chanMessage)
}
