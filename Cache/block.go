package Cache

//队列中数据块的基本单位，每发布/订阅一个Block是一个原子操作
type Block struct{
	data []byte
	offset int
}
func (b *Block)Get() []byte{
	return b.data
}
func (b *Block) Set(x []byte) {
	b.data = x
}
func (b *Block) GetOffset() int{
	return b.offset
}
func (b *Block) SetOffset(offset int){
	b.offset = offset
}