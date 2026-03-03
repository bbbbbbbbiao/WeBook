package domain

/**
 * @author: biao
 * @date: 2025/12/22 下午9:22
 * @description: 用户领域对象
 */

type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	OpenId   string
	UnionId  string

	//　这里可以进行组合，但是我们没有组合，因为后面如果有QQ、钉钉扫码的话，可以复用（OpenId 和 UnionId）
	//　WeChatInfo WeChatInfo
	NickName     string
	Birthday     string
	Introduction string

	Ctime int64
	Utime int64
}
