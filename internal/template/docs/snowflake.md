# Snowflake 方案说明

- 使用 `bwmarrin/snowflake` 生成本地唯一 ID
- 启动时通过 `Short.Snowflake.WorkerID` 指定节点号
- 生成后的 ID 再转 Base62，作为最终 `short_code`
