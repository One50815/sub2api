# Omnio 完整使用教程

<p class="guide-lede">从注册账号、创建 API Key，到完成第一次模型调用、查看用量和排查错误，这一页会带你走完整个使用流程。</p>

<div class="guide-quick-path">
  <span>最快上手路径</span>
  <strong>注册登录 → 充值或订阅 → 创建 API Key → 获取模型名 → 发起第一次请求</strong>
</div>

> 本教程中的 <code>YOUR_API_KEY</code> 和 <code>YOUR_MODEL</code> 都是占位符。请替换为你自己的 API Key，以及该 Key 实际可用的模型名。不要把真实 API Key 发给他人、写入公开仓库或放进网页前端代码。

## 一、使用前需要知道的地址

Omnio 主站和 API 使用同一个域名：

| 用途 | 地址 |
|---|---|
| 主站与控制台 | <code>https://omni0.top</code> |
| OpenAI / Anthropic 兼容 API 根地址 | <code>https://omni0.top/v1</code> |
| Gemini 原生兼容 API 根地址 | <code>https://omni0.top/v1beta</code> |
| 文档中心 | <code>https://omni0.top/docs/</code> |

不同 SDK 对“根地址”的要求可能不同：

- OpenAI SDK 的 <code>base_url</code> / <code>baseURL</code> 填 <code>https://omni0.top/v1</code>。
- Anthropic SDK 的 <code>base_url</code> / <code>baseURL</code> 填 <code>https://omni0.top</code>，SDK 会继续请求 <code>/v1/messages</code>。
- 直接发送 HTTP 请求时，使用本教程列出的完整接口地址。

## 二、注册、登录与准备余额

1. 打开 [Omnio 注册页](https://omni0.top/register) 创建账号；已有账号可直接 [登录](https://omni0.top/login)。
2. 登录后进入控制台，先确认账号状态正常。
3. 如当前没有可用余额或订阅，进入 [购买与充值](https://omni0.top/purchase)。
4. 已购买的订阅可在 [我的订阅](https://omni0.top/subscriptions) 查看，支付记录可在 [我的订单](https://omni0.top/orders) 查看。
5. 如果你获得了兑换码，可前往 [兑换中心](https://omni0.top/redeem) 使用。

余额、订阅和具体扣减规则以控制台实时展示为准。涉及退款或未消费余额处理时，请查看[《数字充值服务政策》](https://omni0.top/docs/recharge-policy/)。

## 三、创建你的第一个 API Key

打开控制台的 [API Keys](https://omni0.top/keys) 页面，点击创建并完成以下设置：

1. **名称**：建议写清用途，例如“个人开发环境”“生产服务 A”。
2. **分组**：分组决定 Key 可访问的平台和模型。优先选择与你要使用的模型相符的分组。
3. **额度上限**：可为单个 Key 设置独立额度，降低误调用造成的损失。
4. **有效期**：临时测试 Key 建议设置到期时间，长期服务按实际维护周期设置。
5. **IP 白名单 / 黑名单**：固定服务器调用时建议配置白名单；动态网络环境下请谨慎开启。
6. **时间窗口限制**：如页面提供 5 小时、每日或 7 日额度限制，可按项目预算设置。

创建成功后立即复制并妥善保存 Key。建议一个项目、一个环境使用一把独立 Key，不要让开发、测试和生产环境共用。

## 四、先获取当前可用模型

模型会随分组、上游能力和平台配置变化，因此不要仅凭旧教程猜测模型名。你可以在控制台的[可用渠道](https://omni0.top/available-channels)页面查看，也可以直接请求模型列表：

~~~bash
export OMNIO_API_KEY="YOUR_API_KEY"

curl "https://omni0.top/v1/models" \
  -H "Authorization: Bearer $OMNIO_API_KEY"
~~~

Windows PowerShell：

~~~powershell
$env:OMNIO_API_KEY = "YOUR_API_KEY"
$headers = @{ Authorization = "Bearer $env:OMNIO_API_KEY" }
Invoke-RestMethod -Uri "https://omni0.top/v1/models" -Headers $headers
~~~

从返回结果中复制你需要的模型 <code>id</code>，并在后面的示例中替换 <code>YOUR_MODEL</code>。如果模型列表为空或返回权限错误，请检查 Key 绑定的分组、余额和状态。

## 五、完成第一次 OpenAI 兼容调用

### Chat Completions

这是兼容范围最广的调用方式：

~~~bash
curl "https://omni0.top/v1/chat/completions" \
  -H "Authorization: Bearer $OMNIO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "YOUR_MODEL",
    "messages": [
      {
        "role": "user",
        "content": "请用一句话介绍 Omnio。"
      }
    ],
    "stream": false
  }'
~~~

成功时，文本通常位于 <code>choices[0].message.content</code>。

### Responses API

支持 Responses API 的模型可以使用：

~~~bash
curl "https://omni0.top/v1/responses" \
  -H "Authorization: Bearer $OMNIO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "YOUR_MODEL",
    "input": "给我三个提高代码可读性的建议。"
  }'
~~~

是否支持 Responses、工具调用、图片等能力，取决于 Key 所属分组和所选模型。接口返回“不支持”时，不要反复重试；先确认模型和分组能力。

## 六、使用 OpenAI SDK

### Python

安装 SDK：

~~~bash
pip install openai
~~~

调用 Chat Completions：

~~~python
import os
from openai import OpenAI

client = OpenAI(
    api_key=os.environ["OMNIO_API_KEY"],
    base_url="https://omni0.top/v1",
)

response = client.chat.completions.create(
    model="YOUR_MODEL",
    messages=[
        {"role": "user", "content": "写一个 Python 快速排序示例。"}
    ],
)

print(response.choices[0].message.content)
~~~

调用 Responses API：

~~~python
response = client.responses.create(
    model="YOUR_MODEL",
    input="解释什么是 API 网关。",
)

print(response.output_text)
~~~

### Node.js

安装 SDK：

~~~bash
npm install openai
~~~

~~~javascript
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: process.env.OMNIO_API_KEY,
  baseURL: "https://omni0.top/v1",
});

const response = await client.chat.completions.create({
  model: "YOUR_MODEL",
  messages: [
    { role: "user", content: "用 JavaScript 写一个防抖函数。" },
  ],
});

console.log(response.choices[0].message.content);
~~~

## 七、使用 Anthropic Messages 接口

支持 Anthropic 兼容协议的分组可以请求 <code>/v1/messages</code>。API Key 既可使用 Bearer 认证，也兼容 <code>x-api-key</code> 请求头。

~~~bash
curl "https://omni0.top/v1/messages" \
  -H "x-api-key: $OMNIO_API_KEY" \
  -H "anthropic-version: 2023-06-01" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "YOUR_MODEL",
    "max_tokens": 1024,
    "messages": [
      {
        "role": "user",
        "content": "请解释这段代码可能出现的并发问题。"
      }
    ]
  }'
~~~

Python Anthropic SDK：

~~~bash
pip install anthropic
~~~

~~~python
import os
from anthropic import Anthropic

client = Anthropic(
    api_key=os.environ["OMNIO_API_KEY"],
    base_url="https://omni0.top",
)

message = client.messages.create(
    model="YOUR_MODEL",
    max_tokens=1024,
    messages=[
        {"role": "user", "content": "给这个函数补充单元测试。"}
    ],
)

print(message.content[0].text)
~~~

## 八、使用 Gemini 原生接口

Gemini 原生兼容接口位于 <code>/v1beta</code>。先查询该 Key 在 Gemini 分组下可用的模型：

~~~bash
curl "https://omni0.top/v1beta/models" \
  -H "x-goog-api-key: $OMNIO_API_KEY"
~~~

发起内容生成请求：

~~~bash
curl "https://omni0.top/v1beta/models/YOUR_MODEL:generateContent" \
  -H "x-goog-api-key: $OMNIO_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [
      {
        "role": "user",
        "parts": [
          {"text": "概括一下这段需求的核心目标。"}
        ]
      }
    ]
  }'
~~~

不要把 Key 放在 URL 查询参数中。Omnio 支持 <code>x-goog-api-key</code>、Bearer 或 <code>x-api-key</code> 请求头，Gemini 客户端优先使用 <code>x-goog-api-key</code>。

## 九、流式输出

需要边生成边显示内容时，在支持的请求中设置：

~~~json
{
  "stream": true
}
~~~

OpenAI Python SDK 示例：

~~~python
stream = client.chat.completions.create(
    model="YOUR_MODEL",
    messages=[{"role": "user", "content": "写一篇短文。"}],
    stream=True,
)

for chunk in stream:
    delta = chunk.choices[0].delta.content
    if delta:
        print(delta, end="", flush=True)
~~~

流式请求中断后，先确认客户端超时、网络和模型状态。不要在不知道上一次请求是否已被执行的情况下无限自动重试，以免产生重复调用。

## 十、常用 API 端点

| 方法 | 端点 | 用途 |
|---|---|---|
| GET | <code>/v1/models</code> | 查询当前 Key 可见的模型 |
| POST | <code>/v1/chat/completions</code> | OpenAI Chat Completions 兼容调用 |
| POST | <code>/v1/responses</code> | OpenAI Responses 兼容调用 |
| POST | <code>/v1/messages</code> | Anthropic Messages 兼容调用 |
| POST | <code>/v1/messages/count_tokens</code> | 支持的分组中预估消息 Token |
| GET | <code>/v1/usage</code> | 查询当前 Key 的额度和用量 |
| POST | <code>/v1/embeddings</code> | 支持的 OpenAI 分组中生成向量 |
| POST | <code>/v1/images/generations</code> | 支持的模型中生成图片 |
| POST | <code>/v1/images/edits</code> | 支持的模型中编辑图片 |
| GET | <code>/v1beta/models</code> | 查询 Gemini 原生接口模型 |
| POST | <code>/v1beta/models/{model}:generateContent</code> | Gemini 原生内容生成 |

端点存在不代表每个模型都支持对应能力。实际能力以当前 Key 的分组、模型列表和接口响应为准。

## 十一、查询余额与用量

### 在控制台查看

- [用量记录](https://omni0.top/usage)：查看请求、Token、费用、模型等记录。
- [API Key 用量查询](https://omni0.top/key-usage)：不登录控制台也可以用自己的 Key 查询余额和统计。
- [我的订阅](https://omni0.top/subscriptions)：查看订阅额度、周期与到期时间。

### 通过 API 查询

~~~bash
curl "https://omni0.top/v1/usage" \
  -H "Authorization: Bearer $OMNIO_API_KEY"
~~~

按日期查询时可使用 <code>start_date</code>、<code>end_date</code> 和 <code>timezone</code> 参数。页面展示和接口统计均以服务端实际账单记录为准。

## 十二、在第三方客户端中填写什么

大多数 OpenAI 兼容客户端只需要以下三项：

| 客户端字段 | 填写内容 |
|---|---|
| API Key | 你在 Omnio 创建的 Key |
| Base URL / API Address | <code>https://omni0.top/v1</code> |
| Model | 从 <code>/v1/models</code> 获取的模型 ID |

如果客户端会自动拼接 <code>/v1</code>，则 Base URL 只填 <code>https://omni0.top</code>。最终请求地址不应出现重复的 <code>/v1/v1</code>。

Claude Code 或 Anthropic 兼容客户端通常使用：

~~~bash
export ANTHROPIC_BASE_URL="https://omni0.top"
export ANTHROPIC_AUTH_TOKEN="YOUR_API_KEY"
~~~

不同版本的客户端配置项可能变化。如果客户端已经提供“自定义 OpenAI / Anthropic 服务商”界面，优先通过界面填写上述地址和 Key。

## 十三、API Key 安全建议

1. **不要放进浏览器前端**：网页源码和网络请求可能暴露 Key，应由自己的后端调用 Omnio。
2. **使用环境变量**：不要把 Key 直接写进代码、镜像或公开配置文件。
3. **按项目隔离**：开发、测试、生产分别创建 Key，便于停用、审计和统计。
4. **设置额度与有效期**：临时项目设置合理上限，减少泄露后的风险。
5. **固定服务配置 IP 白名单**：仅允许可信服务器调用。
6. **定期轮换**：发现泄露、异常调用或人员变动时，立即停用旧 Key 并创建新 Key。
7. **日志脱敏**：不要在报错日志、工单截图和聊天记录中展示完整 Key。

## 十四、常见错误与处理方式

| 状态或现象 | 常见原因 | 建议处理 |
|---|---|---|
| 400 请求错误 | JSON 格式、参数或模型能力不匹配 | 检查请求体，并先用最小示例测试 |
| 401 未认证 | Key 缺失、错误、停用或过期 | 检查认证请求头和 Key 状态 |
| 402 / 额度不足 | 余额、订阅或 Key 独立额度不足 | 查看用量、订阅和充值页面 |
| 403 无权限 | 分组、模型权限或 IP 规则不允许 | 检查 Key 分组及 IP 白名单/黑名单 |
| 404 不支持 | 路径错误，或当前平台不支持该端点 | 核对根地址、端点和模型能力 |
| 429 请求受限 | 并发、频率或时间窗口额度达到上限 | 降低并发，等待窗口重置或调整 Key 限制 |
| 5xx / 暂时不可用 | 上游模型、网络或网关出现临时故障 | 稍后重试；持续出现时提交工单 |
| 返回模型不可用 | 模型名过期或分组不匹配 | 重新请求 <code>/v1/models</code> 并选择有效模型 |

排查时建议先发送一个非流式、短文本、无工具调用的最小请求。最小请求成功后，再逐步恢复长上下文、流式输出、工具或多模态参数。

## 十五、提交工单获得帮助

如果问题仍未解决，登录后进入 [我的工单](https://omni0.top/tickets)：

1. 填写清晰的主题和问题描述。
2. 选择最接近的分类。
3. 提供发生时间、调用端点、模型名和 HTTP 状态码。
4. 如果响应中包含请求 ID，请填写到“关联请求 ID”，方便定位日志。
5. 不要提交完整 API Key、密码、访问令牌或其他敏感凭据。

用户可以在工单中继续回复、标记已解决或关闭。当前第一版工单不支持附件；如需展示错误内容，请先对 Key、Token、邮箱和个人信息进行脱敏后粘贴文字。

也可以通过客服邮箱 [3290800970@qq.com](mailto:3290800970@qq.com) 联系 Omnio。

## 十六、上线前检查清单

- [ ] API Key 来自独立的生产环境配置，不在代码仓库中。
- [ ] Base URL 没有重复或缺少 <code>/v1</code>。
- [ ] 模型名来自当前 Key 的 <code>/v1/models</code> 返回结果。
- [ ] Key 已设置合理的额度、有效期和 IP 规则。
- [ ] 客户端设置了合理超时，并区分可重试与不可重试错误。
- [ ] 流式中断不会触发无限重试或重复业务操作。
- [ ] 日志中不会记录完整 Key 或敏感输入。
- [ ] 已确认余额、订阅、用量统计和告警方式。
- [ ] 发生问题时能够保留请求 ID 并提交工单。

完成以上步骤后，你已经具备从测试调用走向正式接入 Omnio 的基本条件。
