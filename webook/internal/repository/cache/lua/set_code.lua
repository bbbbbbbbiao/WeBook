-- 【业务说明】该脚本用于手机验证码存储与发送频率控制
-- 【返回值说明】
--  0 ：操作成功，已重新生成并存储验证码
-- -1 ：发送过于频繁，需等待（剩余有效期>9分钟）
-- -2 ：key存在但无过期时间（手动配置错误）
-- -3 ：参数不合法（KEYS[1]或ARGV[1]为空）
-- -4 ：未知错误（ttl转换失败等）

-- 1. 合法性校验：检查KEYS[1]和ARGV[1]是否非空
if not KEYS[1] or not ARGV[1] then
    return -3
end

-- Redis 中的 key (phone_code:login:130xxxxxxxx)
local key = KEYS[1]
-- 验证次数 key (phone_code:login:130xxxxxxxx:cnt)
local cntKey = key .. ":cnt"
-- 你的验证码
local val = ARGV[1]

-- 2. 获取key的剩余过期时间，并转换为数字（兜底处理转换失败）
local ttl_raw = redis.call("ttl", key)
local ttl = tonumber(ttl_raw)

-- 3. 处理ttl转换失败的情况
if not ttl then
    return -4
end

-- 4. 核心业务逻辑
if ttl == -1 then
    -- key 存在，没有过期时间 （手动设置错误，没有给过期时间）
    return -2
elseif ttl == -2 or ttl < 540 then
    -- key 不存在 or 过期时间小于9分钟了（已经过了一分钟了）
    -- 优化：单条set命令完成赋值+过期时间设置，更高效
    redis.call("set", key, val, "EX", 600)
    redis.call("set", cntKey, 3, "EX", 600)
    return 0
else
    -- 发送太频繁
    return -1
end