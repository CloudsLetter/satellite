package utilitie

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"satellite/model/misc"
	"strconv"
	"strings"
	"time"
)

// OfflineMessageDataId 宏定义
var OfflineMessageDataId = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}

// GetMotdFromJavaEdition17 (string, int) error
// 获取java版1.7及以上版本的motd
func GetMotdFromJavaEdition17(javaJson *misc.MotdJson, host string, port int) error {

	buf := bytes.Buffer{}
	resultBuf := bytes.Buffer{}
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), 1*time.Second)

	if err != nil {
		return err
	}
	// 写入包id
	buf.Write([]byte{0x00})
	// 写入请求协议版本
	buf.Write([]byte{0x01})
	// 写入host长度
	buf.Write([]byte{byte(len(host))})
	// 写入host
	buf.Write([]byte(host))
	// 写入port
	portBuf := make([]byte, 2)
	// short，大字节序列
	binary.BigEndian.PutUint16(portBuf, uint16(port))
	// 写入端口
	buf.Write(portBuf)
	// 写入状态
	buf.Write([]byte{0x01})
	// 写入包长度
	packageLength := []byte{uint8(buf.Len())}
	resultBuf.Write(append(packageLength, buf.Bytes()...))
	_, err = conn.Write(resultBuf.Bytes())
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte{0x01, 0x00})
	if err != nil {
		return err
	}
	resBuf := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return err
	}

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
	str := string(recBytes)
	str = strings.ReplaceAll(str, "\"\\n\",", "")
	str = strings.ReplaceAll(str, "\"\\r\",", "")
	str = strings.ReplaceAll(str, ",\" \",", ",")
	err = json.Unmarshal([]byte(str), &javaJson)
	if err != nil {
		return err
	}
	return nil
}

// GetMotdFromBedrockEdition (string, int) error
// 获取基岩版的motd
func GetMotdFromBedrockEdition(host string, port int) ([]string, error) {
	originBuf := bytes.Buffer{}
	conn, err := net.DialTimeout("udp", host+":"+strconv.Itoa(port), 1*time.Second)

	if err != nil {
		return nil, err
	}
	// 写入包id
	originBuf.Write([]byte{0x01})
	timeBuf := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBuf, uint64(time.Now().Unix()))
	// 写入时间戳
	originBuf.Write(timeBuf)
	guidBuf := make([]byte, 8)
	// 写入 magic
	originBuf.Write(OfflineMessageDataId)
	binary.BigEndian.PutUint64(guidBuf, 2)
	// 写入guid
	originBuf.Write(guidBuf)
	_, err = conn.Write(originBuf.Bytes())
	if err != nil {
		return nil, err
	}
	resBuf := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return nil, err
	}
	firstByte, _ := resBuf.ReadByte()
	if bytes.Compare([]byte{firstByte}, []byte{0x1C}) != 0 {
		return nil, errors.New("error response packet type")
	}
	//丢 18 byte数据
	_, err = resBuf.Discard(16)
	if err != nil {
		return nil, err
	}
	offliemsg := make([]byte, 16)
	_, err = resBuf.Read(offliemsg)
	if err != nil {
		return nil, err
	}
	if bytes.Compare(offliemsg, OfflineMessageDataId) != 0 {
		return nil, errors.New("error response packet type")
	}
	length := make([]byte, 2)
	_, err = resBuf.Read(length)
	if err != nil {
		return nil, err
	}

	data := make([]byte, binary.BigEndian.Uint16(length))

	for i := 0; i < int(binary.BigEndian.Uint16(length)); {
		n, err := resBuf.Read(data[i:int(binary.BigEndian.Uint16(length))])
		if err != nil {
			return nil, err
		}
		i = i + n
	}
	return strings.Split(string(data), ";"), nil
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
