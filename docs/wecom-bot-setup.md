# 企业微信智能机器人配置指南

本文档介绍如何在 ZenOps 中配置和使用企业微信智能机器人。

## 功能特性

- 支持与企业微信用户进行智能对话
- 集成 LLM 大模型,提供智能问答功能
- 支持流式响应,实时返回AI生成的内容
- 自动处理消息加密解密
- 支持查询云资源(阿里云、腾讯云)和 CI/CD 信息(Jenkins)

## 前置条件

1. 拥有企业微信管理员权限
2. ZenOps 服务器可以通过公网访问(企业微信需要回调)
3. 已配置 LLM 大模型(推荐使用 DeepSeek 或其他兼容 OpenAI API 的模型)

## 配置步骤

### 1. 创建企业微信AI机器人

1. 登录 [企业微信管理后台](https://work.weixin.qq.com/)
2. 进入 **应用管理** > **应用** > **智能助手**
3. 创建一个新的智能助手应用
4. 进入应用配置页面,获取以下信息:
   - **Token**: 用于验证消息来源
   - **EncodingAESKey**: 用于消息加密解密(43位字符串)

### 2. 配置回调地址

在企业微信智能助手的配置页面:

1. 设置 **接收消息服务器配置**
2. 回调 URL 格式: `http://your-domain.com/api/wecom/callback`
   - 将 `your-domain.com` 替换为你的实际域名或IP
   - 确保该地址可以从公网访问
3. 填写 Token 和 EncodingAESKey
4. 点击保存

### 3. 修改 ZenOps 配置文件

编辑 `config.yaml` 文件,添加企业微信配置:

```yaml
# 企业微信配置
wecom:
  enabled: true  # 启用企业微信AI机器人
  token: "YOUR_WECOM_BOT_TOKEN"  # 从企业微信后台获取
  encoding_aes_key: "YOUR_ENCODING_AES_KEY"  # 43位字符串

# LLM 大模型配置(必须启用)
llm:
  enabled: true
  model: "DeepSeek-V3"  # 或其他模型
  api_key: "YOUR_LLM_API_KEY"
  base_url: "https://api.deepseek.com"  # 可选,自定义 API 端点

# HTTP 服务配置
server:
  http:
    enabled: true
    port: 8080
```

### 4. 启动 ZenOps 服务

```bash
# 启动服务
./zenops serve

# 或使用 Docker
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  your-registry/zenops:latest
```

### 5. 验证配置

1. 在企业微信智能助手配置页面点击 **保存** 按钮
2. 企业微信会向回调地址发送验证请求
3. 如果配置正确,会提示 "配置成功"
4. 检查 ZenOps 日志,应该能看到:
   ```
   Wecom verify request: signature=xxx, timestamp=xxx, nonce=xxx
   Wecom message handler initialized
   ```

## 使用方式

### 与机器人对话

1. 在企业微信中找到你创建的智能助手应用
2. 点击进入对话界面
3. 直接发送消息,机器人会自动回复

### 支持的命令

- **帮助** / **help**: 查看使用帮助
- 直接提问: 例如 "查询阿里云 ECS 列表"、"列出腾讯云的 CVM 实例"

### 示例对话

```
用户: 帮我查询一下阿里云杭州地区的 ECS 实例
机器人: 正在查询阿里云杭州地区的 ECS 实例...

       查询结果:
       - 实例ID: i-bp1234567890
         名称: web-server-01
         状态: Running
         IP: 192.168.1.10

       ✅ 回答完成 | 由 ZenOps 智能机器人提供
```

## 工作原理

### 消息流程

1. **用户发送消息** → 企业微信服务器
2. **企业微信服务器** → 加密消息 → ZenOps 回调接口
3. **ZenOps** → 解密消息 → LLM 处理
4. **LLM** → 流式返回结果 → ZenOps
5. **ZenOps** → 加密响应 → 企业微信服务器
6. **企业微信服务器** → 展示给用户

### 技术特点

- **消息加密**: 使用 AES-256-CBC 加密算法
- **签名验证**: 确保消息来源可靠
- **流式响应**: 支持 AI 生成内容的实时推送
- **状态管理**: 维护对话上下文,支持多轮对话

## 故障排查

### 1. 回调验证失败

**症状**: 企业微信提示 "URL验证失败"

**解决方法**:
- 检查回调 URL 是否可以从公网访问
- 确认 Token 和 EncodingAESKey 配置正确
- 查看 ZenOps 日志中的错误信息

### 2. 机器人不响应

**症状**: 发送消息后没有回复

**解决方法**:
- 检查 LLM 配置是否正确
- 确认 `llm.enabled = true`
- 查看日志中是否有 API 调用错误
- 验证 API Key 是否有效

### 3. 消息解密失败

**症状**: 日志显示 "Failed to decrypt Wecom message"

**解决方法**:
- 确认 EncodingAESKey 长度为 43 位字符
- 检查配置文件中的 EncodingAESKey 是否正确
- 确保没有多余的空格或换行符

### 4. 服务启动失败

**症状**: 启动时提示 "Failed to create Wecom message handler"

**解决方法**:
- 检查配置文件格式是否正确(YAML 语法)
- 确认所有必填字段都已配置
- 查看详细的错误日志

## 安全建议

1. **保护敏感信息**
   - 不要将 Token 和 EncodingAESKey 提交到版本控制系统
   - 使用环境变量或密钥管理服务存储敏感信息

2. **网络安全**
   - 建议使用 HTTPS 协议
   - 配置防火墙,仅允许企业微信服务器 IP 访问回调地址
   - 企业微信服务器 IP 段: [官方文档](https://developer.work.weixin.qq.com/document/path/90930)

3. **访问控制**
   - 限制机器人的可见范围
   - 定期更新 Token 和密钥

## 参考文档

- [企业微信智能助手开发文档](https://developer.work.weixin.qq.com/document/path/100719)
- [消息加密解密说明](https://developer.work.weixin.qq.com/document/path/90968)
- [ZenOps 项目主页](https://github.com/eryajf/zenops)

## 常见问题

**Q: 是否支持群聊?**
A: 目前仅支持单聊模式。群聊功能正在开发中。

**Q: 如何更换 LLM 模型?**
A: 修改配置文件中的 `llm.model` 和 `llm.base_url`,确保新模型兼容 OpenAI API 格式。

**Q: 消息有延迟怎么办?**
A: 检查网络连接和 LLM API 的响应时间。可以更换更快的模型或优化网络配置。

**Q: 如何查看详细日志?**
A: 启动服务时添加 `--debug` 参数,或在配置文件中启用 debug 模式。

## 更新日志

- **v1.0.0** (2024-01): 首次发布企业微信AI机器人功能
  - 支持基础对话功能
  - 集成 LLM 大模型
  - 支持流式响应
  - 支持云资源查询
