package config

/**
 * @author: biao
 * @date: 2026/1/8 下午10:49
 * @description:
 */

type config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}
