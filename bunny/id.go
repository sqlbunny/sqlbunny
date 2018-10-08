package bunny

const (
	IDEncodedLen = 20
	IDRawLen     = 12

	idEncoding = "123456789abcdefghjklmnpqrstuvxyz"
)

var IDNullBytes = []byte("null")

var idDec [256]byte

func init() {
	for i := 0; i < len(idDec); i++ {
		idDec[i] = 0xFF
	}
	for i := 0; i < len(idEncoding); i++ {
		idDec[idEncoding[i]] = byte(i)
	}
}

type ID interface {
	String() string
	IDBytes() []byte
}

type IDData [IDRawLen]byte

func (id IDData) Encode(dst []byte) {
	dst[0] = idEncoding[id[0]>>3]
	dst[1] = idEncoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[2] = idEncoding[(id[1]>>1)&0x1F]
	dst[3] = idEncoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[4] = idEncoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[5] = idEncoding[(id[3]>>2)&0x1F]
	dst[6] = idEncoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[7] = idEncoding[id[4]&0x1F]
	dst[8] = idEncoding[id[5]>>3]
	dst[9] = idEncoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[10] = idEncoding[(id[6]>>1)&0x1F]
	dst[11] = idEncoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[12] = idEncoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[13] = idEncoding[(id[8]>>2)&0x1F]
	dst[14] = idEncoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[15] = idEncoding[id[9]&0x1F]
	dst[16] = idEncoding[id[10]>>3]
	dst[17] = idEncoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[18] = idEncoding[(id[11]>>1)&0x1F]
	dst[19] = idEncoding[(id[11]<<4)&0x1F]
}

func (id *IDData) Decode(src []byte) bool {
	for _, c := range src {
		if idDec[c] == 0xFF {
			return false
		}
	}
	if idDec[src[19]]&0x0F != 0 {
		return false
	}
	id[0] = idDec[src[0]]<<3 | idDec[src[1]]>>2
	id[1] = idDec[src[1]]<<6 | idDec[src[2]]<<1 | idDec[src[3]]>>4
	id[2] = idDec[src[3]]<<4 | idDec[src[4]]>>1
	id[3] = idDec[src[4]]<<7 | idDec[src[5]]<<2 | idDec[src[6]]>>3
	id[4] = idDec[src[6]]<<5 | idDec[src[7]]
	id[5] = idDec[src[8]]<<3 | idDec[src[9]]>>2
	id[6] = idDec[src[9]]<<6 | idDec[src[10]]<<1 | idDec[src[11]]>>4
	id[7] = idDec[src[11]]<<4 | idDec[src[12]]>>1
	id[8] = idDec[src[12]]<<7 | idDec[src[13]]<<2 | idDec[src[14]]>>3
	id[9] = idDec[src[14]]<<5 | idDec[src[15]]
	id[10] = idDec[src[16]]<<3 | idDec[src[17]]>>2
	id[11] = idDec[src[17]]<<6 | idDec[src[18]]<<1 | idDec[src[19]]>>4
	return true
}
