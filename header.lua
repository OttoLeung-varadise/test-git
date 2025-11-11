-- headers.lua
-- 设置请求头（键为头名称，值为头内容）
wrk.headers["User-Agent"] = "wrk-benchmark/1.0"
wrk.headers["X-WX-OPENID"] = "test"