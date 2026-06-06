local short_code = os.getenv("SHORT_CODE")

if not short_code or short_code == "" then
  error("SHORT_CODE is required")
end

wrk.method = "GET"
wrk.path = "/" .. short_code
