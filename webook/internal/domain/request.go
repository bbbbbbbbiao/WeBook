package domain

type Request struct {
	ID      int64
	Tpl     string
	Args    []string
	Numbers []string

	// 时间，毫秒级别
	Ctime int64
	Utime int64
	Dtime int64
}
