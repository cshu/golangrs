package golangrs

import (
	"encoding/binary"
	"io"
	"net"
)

func ReadBytesFromConn(conn *net.TCPConn, leng uint32) []byte {
	buf := make([]byte, leng)
	_, err := io.ReadFull(conn, buf)
	CheckErr(err)
	return buf
}
func ReadBytesWithVarLenPrefixFromConn(conn *net.TCPConn) []byte {
	leng := ReadVarLenUint32(conn)
	return ReadBytesFromConn(conn, leng)
}
func ReadVarLenUint32(conn *net.TCPConn) uint32 {
	leng := ReadOneByteFromConn(conn)
	if 255 == leng {
		return ReadUint32FromConn(conn)
	}
	return uint32(leng)
}
func ReadOneByteFromConn(conn *net.TCPConn) byte {
	return ReadBytesFromConn(conn, 1)[0]
}
func ReadUint32FromConn(conn *net.TCPConn) uint32 {
	return binary.LittleEndian.Uint32(ReadBytesFromConn(conn, 4))
}
func ReadUint64FromConn(conn *net.TCPConn) uint64 {
	return binary.LittleEndian.Uint64(ReadBytesFromConn(conn, 8))
}
func ReadStrFromConn(conn *net.TCPConn) string {
	return string(ReadBytesWithVarLenPrefixFromConn(conn))
}
func ReadStrWithUint32Prefix(conn *net.TCPConn) string {
	leng := ReadUint32FromConn(conn)
	return string(ReadBytesFromConn(conn, leng))
}
