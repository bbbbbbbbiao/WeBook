//go:build !k8s

package config

/**
 * @author: biao
 * @date: 2026/1/8 下午10:51
 * @description:
 */

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
