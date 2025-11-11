-- 1. 定义 multipart 分隔符（自定义，确保唯一，避免与文件内容冲突）
local boundary = "----WebKitFormBoundary7MA4YWxkTrZu0gW"

-- 2. 读取本地文件内容（注意：路径是 WSL 中的路径，不是 Windows 路径）
local function read_file(path)
    local file, err = io.open(path, "rb")
    if not file then
        error("无法读取文件: " .. (err or "未知错误"))
    end
    local content = file:read("*a")
    file:close()
    return content
end

-- 待上传的文件路径（WSL 路径，例如 ~/test.txt 需替换为绝对路径 /home/your_username/test.txt）
local file_path = "/home/zkimleung/codes/test-git/role.json"
local file_content = read_file(file_path)
local file_name = "role.json"
local field_name = "file"

-- 3. 构造 multipart/form-data 格式的请求体
local body = string.format(
    "--%s\r\n" ..
    'Content-Disposition: form-data; name="%s"; filename="%s"\r\n' ..
    "Content-Type: application/octet-stream\r\n" ..
    "\r\n" ..  
    "%s\r\n" .. 
    "--%s--\r\n", 
    boundary, field_name, file_name, file_content, boundary
)

-- 4. 设置请求参数
wrk.method = "POST" 
wrk.body = body  
wrk.headers["Content-Type"] = "multipart/form-data; boundary=" .. boundary
wrk.headers["User-Agent"] = "wrk-benchmark/1.0"
wrk.headers["X-WX-OPENID"] = "test"
-- 可选：添加认证头（如需要）
-- wrk.headers["Authorization"] = "Bearer your_token_here"