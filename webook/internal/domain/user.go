package domain

/**
 * @author: biao
 * @date: 2025/12/22 下午9:22
 * @description: 用户领域对象
 */

type User struct {
	Id           int64
	Email        string
	Password     string
	NickName     string
	Birthday     string
	Introduction string
}
