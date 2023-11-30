package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

var _ Codec = (*JsonCodec)(nil)

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

func (g *JsonCodec) Close() error {
	return g.conn.Close()
}

func (g *JsonCodec) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g *JsonCodec) ReadBody(a any) error {
	return g.dec.Decode(a)
}

func (g *JsonCodec) Write(header *Header, body any) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()

	if err := g.enc.Encode(header); err != nil {
		log.Panicln("rpc codec: json error encoding header:", err)
		return err
	}

	if err := g.enc.Encode(body); err != nil {
		log.Panicln("rpc codec: json error encoding body:", err)
		return err
	}
	return nil
}
