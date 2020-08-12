package model

type Err struct {
	Code int
	Err  error
}

const (
	ModelNoErr      = iota // 无错误
	ModelErr               // 未知错误
	ModelErrExisted        // 数据已存在
	ModelErrDel            // 数据删除失败
	ModelErrUpdate         // 数据更新失败
	ModelErrFind           // 查无此数据
)
