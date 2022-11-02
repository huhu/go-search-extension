package data

type Pkg struct {
	Libtype    string
	FullPath   string // 全名 container/heap
	SimpleSpec string // 第一句话 Package heap provides heap operations for any type that implements heap.Interface
	Stars      int    // 被star的次数
	Version    string // 版本
}

// type func variable const
type Datum struct {
	Name       string // 数据名	Truncate
	Label      string // 锚点	Buffer.Truncate
	Definition string // 完整定义	func (b *Buffer) Truncate(n int)
	SimpleSpec string // 第一句话	Truncate discards all but the first n unread bytes from the bufferbut continues to use the same allocated storage
	Datatype   string // 类型	type func
	Pkg        string // 包名	bytes
}

type PkgVersions struct {
	ID         int64
	VersionNum int
	FullPath   string
}
