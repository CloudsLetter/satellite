package misc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log"
	"net"
	"satellite/models/minecraft"
	"satellite/models/respon"
	"strconv"
	"strings"
	"time"
)

// OfflineMessageDataId 宏定义
var OfflineMessageDataId = []byte{0x00, 0xFF, 0xFF, 0x00, 0xFE, 0xFE, 0xFE, 0xFE, 0xFD, 0xFD, 0xFD, 0xFD, 0x12, 0x34, 0x56, 0x78}

func ExtraRawRecursion(extra []minecraft.ExtraObjectOrString, descriptionRaw *[]respon.DescriptionRaw, depth int) {
	for _, v := range extra {
		result := respon.DescriptionRaw{}
		if v.ExtraString != "" {
			result.Text += v.ExtraString
		}
		if v.ExtraObject.Text != "" {
			result.Text += v.ExtraObject.Text
		}
		if v.ExtraObject.Color != "" {
			result.Color = v.ExtraObject.Color
		}
		if v.ExtraObject.Bold {
			result.Bold = v.ExtraObject.Bold
		}
		result.Depth = depth
		*descriptionRaw = append(*descriptionRaw, result)
		if v.ExtraObject.Extra != nil {
			ExtraRawRecursion(v.ExtraObject.Extra, descriptionRaw, depth+1)
		}
	}
}

func ExtraRecursion(extra []minecraft.ExtraObjectOrString) string {
	result := ""
	for _, v := range extra {
		if v.ExtraString != "" {
			result += v.ExtraString
		}
		if v.ExtraObject.Text != "" {
			result += v.ExtraObject.Text
		}
		if v.ExtraObject.Extra != nil {
			result += ExtraRecursion(v.ExtraObject.Extra)
		}
	}
	return result
}

// GetMotdFromJavaEdition17 (string, int) error
// 获取java版1.7及以上版本的motd
func GetMotdFromJavaEdition17(motd *minecraft.JavaMotdJson17, host string, port int) error {

	buf := bytes.Buffer{}
	resultBuf := bytes.Buffer{}
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
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
	startTime := time.Now().UnixMilli()

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
	endTime := time.Now().UnixMilli()
	motd.Latency = int(endTime - startTime)
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
	err = json.Unmarshal(recBytes, &motd)
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}
	return nil
}

// GetMotdFromBedrockEdition (string, int) error
// 获取基岩版的motd
func GetMotdFromBedrockEdition(motd *minecraft.BedrockMotd, host string, port int) error {
	originBuf := bytes.Buffer{}
	conn, err := net.DialTimeout("udp", host+":"+strconv.Itoa(port), 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
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
		return err
	}
	startTime := time.Now().UnixMilli()
	resBuf := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return err
	}
	firstByte, _ := resBuf.ReadByte()
	if bytes.Compare([]byte{firstByte}, []byte{0x1C}) != 0 {
		return errors.New("error response packet type")
	}
	endTime := time.Now().UnixMilli()
	motd.Latency = int(endTime - startTime)
	//丢 18 byte数据
	_, err = resBuf.Discard(16)
	if err != nil {
		return err
	}
	offlineMsg := make([]byte, 16)
	_, err = resBuf.Read(offlineMsg)
	if err != nil {
		return err
	}
	if bytes.Compare(offlineMsg, OfflineMessageDataId) != 0 {
		return errors.New("error response packet type")
	}
	length := make([]byte, 2)
	_, err = resBuf.Read(length)
	if err != nil {
		return err
	}
	lengthR := binary.BigEndian.Uint16(length)

	data := make([]byte, lengthR)

	for i := uint16(0); i < lengthR; {
		n, err := resBuf.Read(data[i:lengthR])
		if err != nil {
			return err
		}
		i = i + uint16(n)
	}
	for i, v := range strings.Split(string(data), ";") {
		switch i {
		case 0:
			motd.GameName = v
			continue
		case 1:
			motd.HostName = v
			continue
		case 2:
			motd.Protocol = v
			continue
		case 3:
			motd.Version = v
			continue
		case 4:
			motd.Players, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
			continue
		case 5:
			motd.MaxPlayers, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
			continue
		case 6:
			motd.ServerId, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}
			continue
		case 7:
			motd.Map = v
			continue
		case 8:
			motd.GameMode = v
			continue
		case 9:
			motd.NintendoLimited, err = strconv.ParseBool(v)
			if err != nil {
				return err
			}
			continue
		case 10:
			motd.IPv4Port, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
			continue
		case 11:
			motd.IPv6Port, err = strconv.Atoi(v)
			if err != nil {
				return err
			}
			continue
		case 12:
			motd.Extra = v
			continue
		}
	}
	return nil
}

func GetMotdFromJavaEdition16(motd *minecraft.JavaMotd16, host string, port int, gType int) error {

	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := bytes.Buffer{}
	if gType == 2 {
		// 写入包id
		buf.Write([]byte{0xFE, 0x01, 0xFA})
		// 写入包长度
		buf.Write([]byte{0x00, 0x0B})
		// 写入 MC|PingHost 字符串
		buf.Write([]byte{0x00, 0x4D, 0x00, 0x43, 0x00, 0x7C, 0x00, 0x50, 0x00, 0x69, 0x00, 0x6E, 0x00, 0x67, 0x00, 0x48, 0x00, 0x6F, 0x00, 0x73, 0x00, 0x74})
		hostLen := make([]byte, 2)
		binary.BigEndian.PutUint16(hostLen, 7+uint16(len(host)))
		// 写入其余数据长度
		buf.Write(hostLen)
		// 写入协议版本
		buf.Write([]byte{0x4A})
		hostLen = make([]byte, 2)
		binary.BigEndian.PutUint16(hostLen, uint16(len(host)))
		// 写入host长度
		buf.Write(hostLen)
		// 写入host
		buf.Write([]byte(host))
		portByte := make([]byte, 4)
		binary.BigEndian.PutUint32(portByte, uint32(port))
		buf.Write(portByte)
	}

	if gType == 3 {
		buf.Write([]byte{0xFE, 0x01})
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	startTime := time.Now().UnixMilli()

	resBuf := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	firstByte, err := resBuf.ReadByte()
	if err != nil {
		return err
	}
	endTime := time.Now().UnixMilli()
	motd.Latency = int(endTime - startTime)

	if bytes.Compare([]byte{firstByte}, []byte{0xFF}) != 0 {
		return errors.New("error response packet type")
	}
	lengthByte := make([]byte, 2)
	_, err = resBuf.Read(lengthByte)
	length := 2 * binary.BigEndian.Uint16(lengthByte)

	data := make([]byte, length)
	for i := uint16(0); i < length; {
		n, err := resBuf.Read(data[i:length])
		if err != nil {
			return err
		}
		i = i + uint16(n)
	}
	data = bytes.ReplaceAll(data, []byte{0x00, 0xA7, 0x00, 0x31, 0x00, 0x00}, []byte{})

	for i, v := range bytes.Split(data, []byte{0x00, 0x00}) {
		byteData := bytes.ReplaceAll(v, []byte{0x00}, []byte{})
		switch i {
		case 0:
			motd.Protocol = string(byteData)
			if err != nil {
				return err
			}
			continue
		case 1:
			motd.Version = string(byteData)
			continue
		case 2:
			motd.MessageDaily = string(byteData)
			continue
		case 3:
			motd.OnlinePlayers, err = strconv.Atoi(string(byteData))
			if err != nil {
				return err
			}
			continue
		case 4:
			motd.MaxPlayers, err = strconv.Atoi(string(byteData))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetMotdFromJavaVersion13(motd *minecraft.JavaMotd13, host string, port int) error {
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), 1*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := bytes.Buffer{}
	// 写入包id
	buf.Write([]byte{0xFE})
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	startTime := time.Now().UnixMilli()

	resBuf := bufio.NewReader(conn)
	err = conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	firstByte, err := resBuf.ReadByte()
	if err != nil {
		return err
	}
	endTime := time.Now().UnixMilli()
	motd.Latency = int(endTime - startTime)

	if bytes.Compare([]byte{firstByte}, []byte{0xFF}) != 0 {
		return errors.New("error response packet type")
	}
	_, err = resBuf.Discard(4)
	if err != nil {
		return err
	}
	lengthByte := make([]byte, 2)
	_, err = resBuf.Read(lengthByte)
	if err != nil {
		return err
	}
	strLength := binary.BigEndian.Uint16(lengthByte)
	data := make([]byte, strLength+10)

	for i := uint16(0); i < strLength+10; {
		n, err := resBuf.Read(data[i : strLength+10])
		if err != nil {
			return err
		}
		i = i + uint16(n)
	}
	motd.Motd = string(bytes.ReplaceAll(data[:strLength], []byte{0x00}, []byte{}))
	for i, v := range bytes.Split(data[strLength+2:], []byte{0x00, 0xA7}) {
		byteData := bytes.ReplaceAll(v, []byte{0x00}, []byte{})
		if i == 0 {
			motd.OnlinePlayers, err = strconv.Atoi(string(byteData))
			if err != nil {
				return err
			}
			continue
		}
		if i == 1 {
			motd.MaxPlayers, err = strconv.Atoi(string(byteData))
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}
