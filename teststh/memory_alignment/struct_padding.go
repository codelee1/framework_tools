package main

import (
	"fmt"
	"sync"
	"unsafe"
)

type A1 struct {
	//字段名  类型  所占字节 起始偏移量   往前字段补齐量  当前内存结束地址
	a bool  // 1        0           0               0+1=1
	b int32 // 4        4           3               4+4=8
	c byte  // 1        8           0               1+8=9
	d int64 // 8        16          7               8+16=24
	e int8  // 1        24          0               24+1=25
}
type A2 struct {
	arr [2]int8          // 2
	sl  []int32          // 24 切片结构体字段为：uintptr的Data，int类型的Len，int类型的Cap，所占字节为24
	m   map[string]int32 // 8 map结构有很多字段的，但是只使用map的的指针值
	ptr *int64           // 8
	st  struct {         // 16，和结构体内的字段类型有关
		str string // string类型字段为：uintptr的Data，int类型的Len，所占字节为16
	}
	i interface{} //16
}
type Ag2 struct {
	sl []int16
	i  interface{}
	st struct {
		str string
	}
	ptr *int64
	m   map[string]int32
	arr [2]int8
	bl  bool
	wg  sync.WaitGroup
}

type A1Rearrange struct {
	d int64 // 8        0           0                8
	b int32 // 4        8           0               8+4=12
	a bool  // 1        12          0               12+1=13
	c byte  // 1        13          0               13+1=14
	e int8  // 1        14          0               14+1=15
}

//type A3 struct {
//	a bool // 1        0           0               0+1=1
//	c byte // 1        1           0               1+1=2
//	e int8 // 1        2           1               2+1=3
//	// padding    1	  3 +1 =1
//	b int32 // 4        4           0               4+4=8
//	d int64 // 8        8           0               8+8=16
//}
//
func main() {
	// 获取数据类型对齐系数：unsafe.Alignof(type)
	// 如果字段类型大于当前系统位数/8，则取系统位数/8，32位取4，64位取8
	// 否则，则去当前字段类型的字节数。

	//获取字段的实际大小，如果是常用的int，byte基本上一看就知道了，但是如map、struct和切片等数据类型的实际大小呢？
	//可以用
	a := A2{}
	fmt.Printf("A2.arr size(array): %d\n", unsafe.Sizeof(a.arr))
	fmt.Printf("A2.sl size(slice): %d\n", unsafe.Sizeof(a.sl))
	fmt.Printf("A2.ptr size(pointer): %d\n", unsafe.Sizeof(a.ptr))
	fmt.Printf("A2.st size(struct): %d\n", unsafe.Sizeof(a.st))
	fmt.Printf("A2.m size(map): %d\n", unsafe.Sizeof(a.m))
	fmt.Printf("A2.i size(interface): %d\n", unsafe.Sizeof(a.i))
	fmt.Printf("A2 size(struct): %d\n", unsafe.Sizeof(a))

	a1 := A1{}
	fmt.Printf("A1.arr aligno(a): %d\n", unsafe.Alignof(a1.a))
	fmt.Printf("A1.arr aligno(b): %d\n", unsafe.Alignof(a1.b))
	fmt.Printf("A1.arr aligno(c): %d\n", unsafe.Alignof(a1.c))
	fmt.Printf("A1.arr aligno(d): %d\n", unsafe.Alignof(a1.d))
	fmt.Printf("A1.arr aligno(e): %d\n", unsafe.Alignof(a1.e))

	fmt.Printf("A1.arr Offsetof(a): %d\n", unsafe.Offsetof(a1.a))
	fmt.Printf("A1.arr Offsetof(b): %d\n", unsafe.Offsetof(a1.b))
	fmt.Printf("A1.arr Offsetof(c): %d\n", unsafe.Offsetof(a1.c))
	fmt.Printf("A1.arr Offsetof(d): %d\n", unsafe.Offsetof(a1.d))
	fmt.Printf("A1.arr Offsetof(e): %d\n", unsafe.Offsetof(a1.e))

	// 除了知道了结构体内部的字段的各自占用多少字节，如果按照各个字段所占用的字节来统计，那么结构体A2的大小应该是91，但是实际却是79，
	// 为什么是这样呢，这就需要我们知道结构体中实际的分配是怎么样的，各个字段从什么地方开始，占用多少字节，这里需要用到unsafe.Offsetof(x ArbitraryType)方法

}

/**
分析工具
layout
go get github.com/ajstarks/svgo/structlayout-svg
go get -u honnef.co/go/tools
go install honnef.co/go/tools/cmd/structlayout
go install honnef.co/go/tools/cmd/structlayout-pretty



建议工具
optimize
go install honnef.co/go/tools/cmd/structlayout-optimize

生成svg图片分析
structlayout -json framework_tools/main Ag|structlayout-svg -t "ag-padding" >ag.svg

字段重排建议：
structlayout -json framework_tools/main Ag|structlayout-optimize -r

1、介绍内存对齐是什么
2、为什么需要内存对齐
https://www.jianshu.com/p/49f7e6f56568

*/
