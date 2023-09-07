package golangrs

import (
	"encoding/binary"
	"io"
	"net"
)

func ReadBytesFromUnixConn(conn *net.UnixConn, leng uint32) []byte {
	buf := make([]byte, leng)
	_, err := io.ReadFull(conn, buf)
	CheckErr(err)
	return buf
}
func ReadBytesWithVarLenPrefixFromUnixConn(conn *net.UnixConn) []byte {
	leng := ReadVarLenUint32FromUnixConn(conn)
	return ReadBytesFromUnixConn(conn, leng)
}
func ReadVarLenUint32FromUnixConn(conn *net.UnixConn) uint32 {
	leng := ReadOneByteFromUnixConn(conn)
	if 255 == leng {
		return ReadUint32FromUnixConn(conn)
	}
	return uint32(leng)
}
func ReadOneByteFromUnixConn(conn *net.UnixConn) byte {
	return ReadBytesFromUnixConn(conn, 1)[0]
}
func ReadUint32FromUnixConn(conn *net.UnixConn) uint32 {
	return binary.LittleEndian.Uint32(ReadBytesFromUnixConn(conn, 4))
}
func ReadUint64FromUnixConn(conn *net.UnixConn) uint64 {
	return binary.LittleEndian.Uint64(ReadBytesFromUnixConn(conn, 8))
}
func ReadStrFromUnixConn(conn *net.UnixConn) string {
	return string(ReadBytesWithVarLenPrefixFromUnixConn(conn))
}
func ReadStrWithUint32PrefixFromUnixConn(conn *net.UnixConn) string {
	leng := ReadUint32FromUnixConn(conn)
	return string(ReadBytesFromUnixConn(conn, leng))
}
