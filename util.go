package ftdc

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
)

func flattenBSON(d *bson.Document) (o []Metric) {
	iter := d.Iterator()
	for iter.Next() {
		e := iter.Element()
		val := e.Value()
		key := e.Key()

		switch val.Type() {
		case bson.TypeObjectID:
			// pass
		case bson.TypeString:
			// pass
		case bson.TypeArray:
			// pass
		case bson.TypeEmbeddedDocument:
			n := flattenBSON(val.MutableDocument())
			for _, ne := range n {
				o = append(o, Metric{
					Key:   key + "." + ne.Key,
					Value: ne.Value,
				})
			}
		case bson.TypeBoolean:
			if val.Boolean() {
				o = append(o, Metric{
					Key:   key,
					Value: 1,
				})
			} else {
				o = append(o, Metric{
					Key:   key,
					Value: 0,
				})
			}
		case bson.TypeDouble:
			o = append(o, Metric{
				Key:   key,
				Value: int(val.Double()),
			})
		case bson.TypeInt32:
			o = append(o, Metric{
				Key:   key,
				Value: int(val.Int32()),
			})
		case bson.TypeInt64:
			o = append(o, Metric{
				Key:   key,
				Value: int(val.Int64()),
			})
		case bson.TypeDateTime:
			o = append(o, Metric{
				Key:   key,
				Value: int(val.DateTime().Unix()) * 1000,
			})
		case bson.TypeTimestamp:
			t, _ := val.Timestamp()
			o = append(o, Metric{
				Key:   key,
				Value: int(t) * 1000,
			})
		default:
			fmt.Println("skipping:", key, val.Type())
		}
	}

	return
}

func unpackDelta(buf *bufio.Reader) (delta int, err error) {
	var res uint64
	var shift uint
	for {
		var b byte
		b, err = buf.ReadByte()
		if err != nil {
			return
		}
		bb := uint64(b)
		res |= (bb & 0x7F) << shift
		if bb&0x80 == 0 {
			// read as int64 (handle negatives)
			var n int64
			tmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmp, res)
			binary.Read(bytes.NewBuffer(tmp), binary.LittleEndian, &n)
			delta = int(n)
			return
		}
		shift += 7
	}
}

func unpackInt(bl []byte) int {
	return int(int32((uint32(bl[0]) << 0) |
		(uint32(bl[1]) << 8) |
		(uint32(bl[2]) << 16) |
		(uint32(bl[3]) << 24)))
}

func sum(l ...int) (s int) {
	for _, v := range l {
		s += v
	}
	return
}

func square(n int) int {
	return n * n
}
