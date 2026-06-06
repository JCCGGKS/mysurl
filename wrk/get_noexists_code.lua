local code_file = os.getenv("CODE_FILE") or "code.txt"
local codes = {}
local index = 0

local function load_codes(path)
  local file = io.open(path, "r")
  if not file then
    error("failed to open code file: " .. path)
  end

  for line in file:lines() do
    local code = line:match("^%s*([A-Za-z0-9]+)%s*$")
    if code then
      codes[#codes + 1] = code
    end
  end

  file:close()

  if #codes == 0 then
    error("no codes found in file: " .. path)
  end
end

load_codes(code_file)

wrk.method = "GET"

request = function()
  index = index + 1
  if index > #codes then
    index = 1
  end

  wrk.path = "/" .. codes[index]
  return wrk.format(nil, wrk.path)
end
