# 模型官方定价与 Omnio 分组倍率

<div class="pricing-notice">
  <strong>重要：以下价格是模型厂商公布的官方上游定价，仅作为 Omnio 计费基准。Omnio 实际价格必须以你所使用分组的倍率为准。</strong>
  <p><b>1.0 倍</b> = 官方原价；<b>0.5 倍</b> = 官方原价的一半；<b>2.0 倍</b> = 官方原价的 2 倍；其他倍率以相同方式计算。</p>
</div>

页面价格最后核对日期：**2026 年 7 月 21 日**。除特别说明外，价格单位均为 **美元 / 100 万 tokens（USD / MTok）**，展示的是厂商标准实时 API 价格，不包含 Batch、Flex、Priority、区域端点、联网搜索、工具调用、图片、音频和视频等附加费用。

## 如何计算 Omnio 实际价格

公式很简单：

> **Omnio 实际单价 = 官方上游单价 × 当前分组倍率**

例如，某模型官方输入价格为 `$2.50 / MTok`、输出价格为 `$15.00 / MTok`：

| 分组倍率 | 实际输入价格 | 实际输出价格 |
|---:|---:|---:|
| 0.5 | $1.25 / MTok | $7.50 / MTok |
| 1.0 | $2.50 / MTok | $15.00 / MTok |
| 2.0 | $5.00 / MTok | $30.00 / MTok |

分组倍率可能因平台、线路、可用性和服务配置不同而变化。最终倍率、可用模型及实际扣费请以 Omnio 控制台的分组说明、可用渠道和用量记录为准。

## OpenAI · GPT 与 Codex

下表为 OpenAI 标准模式、短上下文的官方价格。`缓存输入`表示官方可识别的 cached input；带长上下文或数据驻留要求的请求可能采用更高价格。

| 模型 | 输入 | 缓存输入 | 输出 |
|---|---:|---:|---:|
| `gpt-5.6-sol` | $5.00 | $0.50 | $30.00 |
| `gpt-5.6-terra` | $2.50 | $0.25 | $15.00 |
| `gpt-5.6-luna` | $1.00 | $0.10 | $6.00 |
| `gpt-5.5` | $5.00 | $0.50 | $30.00 |
| `gpt-5.5-pro` | $30.00 | — | $180.00 |
| `gpt-5.4` | $2.50 | $0.25 | $15.00 |
| `gpt-5.4-mini` | $0.75 | $0.075 | $4.50 |
| `gpt-5.4-nano` | $0.20 | $0.02 | $1.25 |
| `gpt-5.4-pro` | $30.00 | — | $180.00 |
| `gpt-5.3-codex` | $1.75 | $0.175 | $14.00 |

长上下文、Batch、Flex、Priority、Realtime、图片、视频与工具价格不在此表中。详情请查看 [OpenAI 官方 API Pricing](https://developers.openai.com/api/docs/pricing)。

## Anthropic · Claude

下表展示 Claude Platform 的基础输入、缓存命中和输出价格。缓存写入、Batch、数据驻留及第三方云区域端点存在独立规则。

| 模型 | 基础输入 | 缓存命中 | 输出 | 说明 |
|---|---:|---:|---:|---|
| `claude-fable-5` | $10.00 | $1.00 | $50.00 | 标准价格 |
| `claude-mythos-5` | $10.00 | $1.00 | $50.00 | 限定可用 |
| `claude-opus-4.8` | $5.00 | $0.50 | $25.00 | 标准价格 |
| `claude-opus-4.7` | $5.00 | $0.50 | $25.00 | 标准价格 |
| `claude-opus-4.6` | $5.00 | $0.50 | $25.00 | 标准价格 |
| `claude-opus-4.5` | $5.00 | $0.50 | $25.00 | 标准价格 |
| `claude-sonnet-5` | $2.00 | $0.20 | $10.00 | 官方优惠价至 2026-08-31；之后为 $3 / $0.30 / $15 |
| `claude-sonnet-4.6` | $3.00 | $0.30 | $15.00 | 标准价格 |
| `claude-sonnet-4.5` | $3.00 | $0.30 | $15.00 | 标准价格 |
| `claude-haiku-4.5` | $1.00 | $0.10 | $5.00 | 标准价格 |

详情请查看 [Anthropic 官方 Claude Pricing](https://platform.claude.com/docs/en/about-claude/pricing)。

## Google · Gemini

下表为 Gemini Developer API Paid Tier 的标准价格，主要展示文本或常规多模态输入、上下文缓存和文本输出。音频、图片生成、搜索增强和不同服务等级可能使用其他价格。

| 模型 | 输入 | 上下文缓存 | 输出 | 条件 |
|---|---:|---:|---:|---|
| `gemini-3.5-flash` | $1.50 | $0.15 | $9.00 | Standard |
| `gemini-3.1-pro-preview` | $2.00 | $0.20 | $12.00 | 提示词 ≤ 200K |
| `gemini-3.1-pro-preview` | $4.00 | $0.40 | $18.00 | 提示词 > 200K |
| `gemini-3.1-flash-lite` | $0.25 | $0.025 | $1.50 | 文本/图片/视频输入 |
| `gemini-3-flash-preview` | $0.50 | $0.05 | $3.00 | 文本/图片/视频输入 |
| `gemini-2.5-pro` | $1.25 | $0.125 | $10.00 | 提示词 ≤ 200K |
| `gemini-2.5-pro` | $2.50 | $0.25 | $15.00 | 提示词 > 200K |
| `gemini-2.5-flash` | $0.30 | $0.03 | $2.50 | 文本/图片/视频输入 |
| `gemini-2.5-flash-lite` | $0.10 | $0.01 | $0.40 | 文本/图片/视频输入 |

详情请查看 [Google 官方 Gemini Developer API Pricing](https://ai.google.dev/gemini-api/docs/pricing)。

## DeepSeek · V4

DeepSeek 将输入区分为缓存命中和缓存未命中，输出单独计费：

| 模型 | 输入：缓存命中 | 输入：缓存未命中 | 输出 |
|---|---:|---:|---:|
| `deepseek-v4-flash` | $0.0028 | $0.14 | $0.28 |
| `deepseek-v4-pro` | $0.003625 | $0.435 | $0.87 |

详情请查看 [DeepSeek 官方 Models & Pricing](https://api-docs.deepseek.com/quick_start/pricing/)。

## Alibaba Cloud · Qwen

Qwen 的官方价格会因部署区域、上下文长度、思考模式和限时活动而变化。下表使用官方列表价，不把临时折扣作为长期基准。

| 模型 | 输入 | 输出 | 范围与条件 |
|---|---:|---:|---|
| `qwen3.7-max` | $1.65 | $4.951 | Global / 中国大陆列表价，≤ 1M tokens |
| `qwen3.7-plus` | $0.40 | $1.60 | International，≤ 256K tokens |
| `qwen3.7-plus` | $1.20 | $4.80 | International，256K–1M tokens |
| `qwen-max` | $1.60 | $6.40 | International，无阶梯价格 |

详情请查看 [Alibaba Cloud Model Studio 官方模型定价](https://www.alibabacloud.com/help/en/model-studio/model-pricing)。

## 使用定价表时请注意

1. **模型可用性与定价是两件事。** 官方列出某个模型，不代表该模型一定已在你的 Omnio 分组中开放。
2. **先确认分组倍率。** 同一模型在不同分组中的最终价格可能不同。
3. **确认缓存类型。** 缓存命中、缓存写入和普通输入可能使用不同单价。
4. **确认长上下文。** 部分模型在输入超过指定长度后会提升输入和输出单价。
5. **工具费用可能另计。** 联网搜索、文件检索、代码执行、图片、音频、视频和其他工具不一定包含在 token 价格中。
6. **以实际账单为准。** 厂商可能随时调整公开价格；Omnio 最终扣费以请求发生时的平台配置和用量记录为准。

需要确认自己当前可用的模型时，请查看[可用渠道](https://omni0.top/available-channels)；需要查看已产生的费用时，请打开[用量查询](https://omni0.top/usage)。
