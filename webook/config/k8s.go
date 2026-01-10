//go:build k8s

package config

/**
 * @author: biao
 * @date: 2026/1/8 下午10:54
 * @description:
 */

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:31000)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:32000",
	},
}
