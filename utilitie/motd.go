package utilitie

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"satellite/model/motd"
	"strconv"
)

// (string, int, error)
func GetMotdFromJavaVersion17(javaJson *motd.MotdJavaJson, host string, port int) error {

	buf := bytes.Buffer{}
	resultBuf := bytes.Buffer{}
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	// write packet id
	buf.Write([]byte{0x00})
	// write protocol
	buf.Write([]byte{0x01})
	// write host length
	buf.Write([]byte{byte(len(host))})
	// write host
	buf.Write([]byte(host))
	// write port
	portBuf := make([]byte, 2)
	// encode port to big endian array
	binary.BigEndian.PutUint16(portBuf, uint16(port))
	// write port
	buf.Write(portBuf)
	buf.Write([]byte{0x01})
	// write nextState
	packageLength := []byte{uint8(buf.Len())}
	resultBuf.Write(append(packageLength, buf.Bytes()...))
	conn.Write(resultBuf.Bytes())
	conn.Write([]byte{0x01, 0x00})
	resBuf := bufio.NewReader(conn)
	// drop less data
	_, err = binary.ReadUvarint(resBuf)
	if err != nil {
		return err
	}

	packetType, _ := resBuf.ReadByte()
	if bytes.Compare([]byte{packetType}, []byte{0x00}) != 0 {
		panic("error response packet type")
	}
	length, err := binary.ReadUvarint(resBuf)
	if length < 10 {
		return errors.New("error to small response")
	} else if length > 700000 {
		return errors.New("error to big response")
	}
	recBytes := make([]byte, length)
	// read buffer to recBytes
	for i := uint64(0); i < length; {
		n, _ := resBuf.Read(recBytes[i:length])
		i = i + uint64(n)
	}
	json.Unmarshal(recBytes, &javaJson)
	return nil
}

//func GetMotdFromJavaVersion16(host string, port int) {
//	conn, _ := net.Dial("tcp", host+":"+strconv.Itoa(port))
//	buf := bytes.Buffer{}
//	newbuf := bytes.Buffer{}
//	buf.Write([]byte{0xFE, 0x01, 0xFA})
//	newbuf.Write(append([]byte{uint8(buf.Len())}, buf.Bytes()...))
//	conn.Write(newbuf.Bytes())
//	resBuf := bufio.NewReader(conn)
//	resbuf := make([]byte, 512)
//	resBuf.Read(resbuf)
//	log.Fatalln(resbuf)
//}
//
//func GetMotdFromJavaVersion14(host string, port int) {
//	conn, _ := net.Dial("tcp", host+":"+strconv.Itoa(port))
//	buf := bytes.Buffer{}
//	newbuf := bytes.Buffer{}
//	buf.Write([]byte{0xFE, 0x01})
//	newbuf.Write(append([]byte{uint8(buf.Len())}, buf.Bytes()...))
//	conn.Write(newbuf.Bytes())
//	resBuf := bufio.NewReader(conn)
//	resbuf := make([]byte, 512)
//	resBuf.Read(resbuf)
//	log.Fatalln(resbuf)
//}
//
//func GetMotdFromJavaVersion13(host string, port int) {
//	conn, _ := net.Dial("tcp", host+":"+strconv.Itoa(port))
//	buf := bytes.Buffer{}
//	newbuf := bytes.Buffer{}
//	buf.Write([]byte{0xFE})
//	newbuf.Write(append([]byte{uint8(buf.Len())}, buf.Bytes()...))
//	conn.Write(newbuf.Bytes())
//	resBuf := bufio.NewReader(conn)
//	resbuf := make([]byte, 512)
//	resBuf.Read(resbuf)
//	log.Fatalln(resbuf)
//}
