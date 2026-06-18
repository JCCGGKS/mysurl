local batch_size = tonumber(os.getenv("BATCH_SIZE") or "20")
local url_prefix = os.getenv("URL_PREFIX") or "https://github.com/JCCGGKS/mysurl-batch-"
local counter = 0

if not batch_size or batch_size <= 0 then
  error("BATCH_SIZE must be greater than zero")
end

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

request = function()
  local long_urls = {}

  for i = 1, batch_size do
    counter = counter + 1
    long_urls[i] = string.format("%q", url_prefix .. counter)
  end

  wrk.body = '{"long_urls":[' .. table.concat(long_urls, ",") .. ']}'
  return wrk.format(nil, "/api/v1/links/batch")
end
