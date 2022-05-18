package cache

//只读数据结构 ByteView 用来表示缓存值
type ByteView struct {
	b []byte //b 将会存储真实的缓存,选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片
}

// Len returns the view's length
func (bv ByteView) Len() int {
	return len(bv.b) //返回其所占的内存大小
}

//String returns the data as a string, making a copy if necessary.
func (bv ByteView) String() string {
	return string(bv.b)
}

// ByteSlice returns a copy of the data as a byte slice.
func (bv ByteView) ByteSlice() []byte {
	return colneBytes(bv.b) //b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改
}

func colneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
