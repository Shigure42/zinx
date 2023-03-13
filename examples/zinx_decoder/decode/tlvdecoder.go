// TLV，即Tag(Type)—Length—Value，是一种简单实用的数据传输方案。
//在TLV的定义中，可以知道它包括三个域，分别为：标签域（Tag），长度域（Length），内容域（Value）。这里的长度域的值实际上就是内容域的长度。
//
//解码前 (20 bytes)                                   解码后 (20 bytes)
//+------------+------------+-----------------+      +------------+------------+-----------------+
//|     Tag    |   Length   |     Value       |----->|     Tag    |   Length   |     Value       |
//| 0x00000001 | 0x0000000C | "HELLO, WORLD"  |      | 0x00000001 | 0x0000000C | "HELLO, WORLD"  |
//+------------+------------+-----------------+      +------------+------------+-----------------+
// Tag：   uint32类型，占4字节，Tag作为MsgId，暂定为1
// Length：uint32类型，占4字节，Length标记Value长度12(hex:0x0000000C)
// Value： 共12个字符，占12字节
//
//   说明：
//   lengthFieldOffset   = 4            (Length的字节位索引下标是4) 长度字段的偏差
//   lengthFieldLength   = 4            (Length是4个byte) 长度字段占的字节数
//   lengthAdjustment    = 0            (Length只表示Value长度，程序只会读取Length个字节就结束，后面没有来，故为0，若Value后面还有crc占2字节的话，那么此处就是2。若Length标记的是Tag+Length+Value总长度，那么此处是-8)
//   initialBytesToStrip = 0            (这个0表示返回完整的协议内容Tag+Length+Value，如果只想返回Value内容，去掉Tag的4字节和Length的4字节，此处就是8) 从解码帧中第一次去除的字节数
//   maxFrameLength      = 2^32 + 4 + 4 (Length为uint类型，故2^32次方表示Value最大长度，此外Tag和Length各占4字节)

package decode

import (
	"encoding/binary"
	"fmt"
	"github.com/aceld/zinx/ziface"
)

const TLV_HEADER_SIZE = 8 //表示TLV空包长度

type TlvData struct {
	Tag    uint32
	Length uint32
	Value  string
}

type TLVDecoder struct {
}

func (this *TLVDecoder) Intercept(chain ziface.Chain) ziface.Response {
	request := chain.Request()
	if request != nil {
		switch request.(type) {
		case ziface.IRequest:
			iRequest := request.(ziface.IRequest)
			iMessage := iRequest.GetMessage()
			if iMessage != nil {
				data := iMessage.GetData()
				fmt.Println("1-TLV", len(data), data)
				datasize := len(data)
				_data := TlvData{}
				if datasize >= TLV_HEADER_SIZE {
					_data.Tag = binary.BigEndian.Uint32(data[0:4])
					_data.Length = binary.BigEndian.Uint32(data[4:8])
					_data.Value = string(data[8 : 8+_data.Length])
					iMessage.SetMsgID(_data.Tag)
					iRequest.SetResponse(_data)
					fmt.Println("2-TLV", _data)
				}
			}
		}
	}
	return chain.Proceed(chain.Request())
}