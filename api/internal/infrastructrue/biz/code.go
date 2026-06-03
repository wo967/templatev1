package biz

const OK = 200

// 通用的错误信息
var (
	DBError    = NewError(10000, "数据库错误")
	TokenError = NewError(10102, "token错误")
	RedisError = NewError(10103, "redis错误")
)
