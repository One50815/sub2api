# Sub2API 定制版长期项目记忆

> 状态：生效中  
> 最后核对日期：2026-07-24
> 基线：`main`，`3f97d49e`（版本 `0.1.164`；官方提交 `cb24522dd53f8f363d008e3afdc8e4baf9788cab`）
> 决策：采用“方案一”，保留 Sub2API 的 Vue 业务实现，高保真、全站复刻独立 React 前端的视觉与交互  
> 前端重做状态：已完成，进入长期维护与上游同步阶段  
> 维护要求：每次改变架构、接口、合并策略或完成一个重做阶段时，同步更新本文档

## 1. 文档用途与优先级

本文档是本 Fork 的长期上下文和维护手册，供开发者、代码审查者及 AI 助手在以下场景中使用：

- 继续当前前端重做；
- 给后端增加本 Fork 独有功能；
- 将新后端功能接入前端；
- 合并上游 `Wei-Shaw/sub2api` 的更新；
- 排查一次更新是否破坏了定制界面或既有功能；
- 新参与者快速恢复项目背景。

发生冲突时，事实来源优先级如下：

1. 当前代码、数据库迁移、测试和构建脚本；
2. 本文档记录的定制决策与维护边界；
3. `DEV_GUIDE.md` 中的通用开发环境经验；
4. 历史讨论、临时说明和个人记忆。

如果本文档与当前代码不一致，不得静默选择其中一方。先核对 Git 历史和实际行为，再修正文档或代码。

## 2. 新任务开始前的快速恢复

开始任何相关任务前，至少阅读：

- `docs/LONG_TERM_PROJECT_MEMORY_CN.md`；
- `DEV_GUIDE.md`；
- `frontend/package.json` 和 `frontend/vite.config.ts`；
- `frontend/src/api/client.ts` 和 `frontend/src/api/url.ts`；
- 涉及页面对应的 `frontend/src/api/`、`frontend/src/types/`、路由和测试；
- 涉及后端功能对应的 `backend/internal/server/routes/`、handler、service、repository、domain/model 和迁移。

然后执行只读检查：

```powershell
git status --short --branch
git remote -v
git log -1 --oneline
```

不要假设工作树干净，不要覆盖未知的本地修改，也不要根据旧文档猜测当前远端地址。

## 3. 当前项目身份

### 3.1 仓库与远端

本文档创建时的状态：

| 项目 | 当前值 |
|---|---|
| 本地仓库 | 工作区中的 `sub2api/` |
| 当前分支 | `main` |
| Fork 远端 `origin` | `https://github.com/One50815/sub2api.git` |
| 官方上游 | `https://github.com/Wei-Shaw/sub2api.git` |
| `upstream` 远端 | `https://github.com/Wei-Shaw/sub2api.git` |

`DEV_GUIDE.md` 中出现过旧 Fork 地址。远端身份以 `git remote -v` 的实时结果和上表为准；远端变化后必须更新本节。

### 3.2 两套前端的角色

| 目录 | 角色 | 是否进入生产构建 |
|---|---|---|
| `sub2api/frontend/` | 正式前端、业务事实源 | 是 |
| 工作区同级的 `../前端/` | React 视觉和交互的验收基准 | 否 |

参考前端当前使用 React 19、Rsbuild、TanStack Router、React Query、Zustand 和另一套 API 契约。正式前端使用 Vue 3、Vite、Vue Router、Pinia 以及 Sub2API API。

参考目录是全站视觉和交互的验收基准，但不是运行时依赖。不得让正式构建依赖 `../前端/` 的存在；CI、Docker 和发布必须能够只检出 `sub2api` 仓库后独立完成。

### 3.3 构建与运行事实

- 正式前端的 API 默认前缀为 `/api/v1`；
- 正式前端使用 Bearer access token 和 refresh token；
- 标准响应由前端按 `{ code, message, data }` 处理，`code === 0` 表示成功；
- `frontend/vite.config.ts` 将产物输出到 `backend/internal/web/dist`；
- 发布构建通过 `backend/internal/web/embed_on.go` 的 `go:embed` 将前端打入 Go 可执行文件；
- 后端向 HTML 注入 `window.__APP_CONFIG__`，用于站点名称、Logo 和公开设置；
- Docker 构建使用 pnpm，并依赖 `frontend/pnpm-lock.yaml`。

上述链路属于受保护的发布契约。前端视觉重做不得绕开它。

## 4. 已确认的架构决策

### ADR-001：采用 Vue 原生视觉重做

**状态**：已接受。  
**日期**：2026-07-20。

目标是在完整保留 Sub2API 功能、接口和发布方式的前提下，用 Vue 高保真复刻参考 React 前端的完整样式与交互。目标不是“相似风格”或“局部借鉴”，而是以参考前端作为明确的视觉验收标准。

复刻范围覆盖整个正式站点：

- 登录、注册、忘记密码、重置密码、邮箱验证、2FA 和 OAuth 过程页；
- 首页、公开页、定价、法律文档、错误页和首次安装页；
- 普通用户的 Dashboard、Keys、Usage、Redeem、Profile、Subscriptions、Purchase、Orders 等全部页面；
- 管理后台的 Dashboard、Users、Groups、Channels、Accounts、Proxies、Ops、Audit、Risk Control、Settings、Payments 等全部页面；
- 应用外壳、导航、表格、表单、弹窗、抽屉、提示、图表、空状态、错误状态和加载状态；
- 桌面端、移动端、浅色主题和深色主题。

这项决策的含义：

- Vue 3 继续作为正式前端框架；
- `frontend/src/api/`、鉴权、权限、路由行为和 Sub2API 类型是业务事实源；
- 参考前端是设计令牌、布局、组件外观、信息密度、响应式和交互行为的验收基准；
- React 组件需要按 Vue 语义重做，不能直接复制 TSX 后强行混用；
- 页面可以按模块渐进实施，但最终必须完成全站复刻，不能长期保留新旧两套视觉；
- 任何实施阶段都保持主分支可构建、可登录、可回滚；
- 视觉提交和业务功能提交尽量分开，降低上游合并冲突。

如果 Sub2API 页面在参考前端中有直接对应页面，必须按 1:1 目标复刻并验收布局、尺寸、间距、字体、颜色、边框、圆角、阴影、图标、信息密度、交互反馈和响应式行为。如果 Sub2API 独有页面没有直接对应页面，则必须使用已经复刻的同一套设计令牌、应用外壳和组件模式进行扩展，使其看起来属于同一个产品，而不是保留旧 Vue 风格。

### ADR-002：不建立 New API 兼容层

参考前端使用 `/api/...`、`New-Api-User` 和 `{ success, data }` 等 New API 契约。正式前端不得为了减少页面改造而要求 Sub2API 后端批量模拟这些接口。

原因：兼容层会产生第二套 API、权限和数据模型，长期成本高于按现有 Sub2API API 重做 Vue 界面。

### ADR-003：Fork 主分支保留合并历史

本 Fork 的 `main` 需要长期接收上游更新。同步上游时优先使用合并分支和 merge commit，不在已经共享的 `main` 上重写历史。每次上游同步单独建立 `sync/upstream-YYYYMMDD` 分支并完整验证。

## 5. 不可破坏的边界

以下规则属于硬性约束，改变它们必须先新增 ADR：

1. 不将 React、React Router、Zustand 或 React Query 引入正式 Vue 应用，仅为复用参考组件。
2. 不把参考前端的 `/api/...` 请求代码复制进正式前端。
3. 不用页面本地请求绕过 `frontend/src/api/client.ts` 的鉴权、刷新和错误处理。
4. 不改变 `/api/v1` 契约来迁就纯视觉改动。
5. 不删除暂时未完成视觉重做的 Sub2API 功能或路由。
6. 不在视觉重做中顺手修改计费、余额、权限、支付或调度语义。
7. 不破坏 `backend/internal/web/dist`、`go:embed` 和 `window.__APP_CONFIG__` 链路。
8. 不把密钥、数据库密码、OAuth secret、支付 secret 或真实 Token 写入代码、文档、截图和测试夹具。
9. 不修改已经发布的历史数据库迁移；新变更使用新的迁移。
10. 不在上游冲突中对整个目录使用无差别的 `ours` 或 `theirs`。
11. 不做无目的或重复的检查；根据改动范围执行最小但充分的验证，同一工作树和环境下已经通过的同一检查不重复运行。

## 6. 源码事实源矩阵

| 关注点 | 事实源 | 参考前端可提供什么 |
|---|---|---|
| API 路径与载荷 | `frontend/src/api/`、后端 routes/handler DTO | 不使用其 API 实现 |
| 用户、权限和余额 | Sub2API 后端模型及 `frontend/src/types/` | 展示方式 |
| 登录与 Token 刷新 | `frontend/src/api/client.ts`、auth API/store | 登录页视觉和交互 |
| 路由与访问控制 | `frontend/src/router/` | 导航结构参考 |
| 支付与 OAuth 回调 | Sub2API 现有页面、API 和测试 | 表单与状态展示参考 |
| 视觉令牌 | 重做后写入 `frontend/src/styles/` | 验收基准 |
| 通用组件 | 重做后写入 `frontend/src/components/` | 视觉、状态和交互验收基准 |
| 文案与语言 | `frontend/src/i18n/` | 语气参考，不直接覆盖翻译键 |
| 构建和部署 | Vite、Dockerfile、Go embed | 不使用其 Rsbuild 发布方式 |

## 7. 前端重做长期路线

### 7.1 范围

重做包含全站高保真复刻：

- 颜色、字体、间距、圆角、阴影、边框和状态色；
- 应用外壳、顶部栏、侧边栏、面包屑和移动端导航；
- 表格、筛选、分页、表单、弹窗、抽屉、标签、提示和空状态；
- 用户端和管理端页面的信息层级；
- 加载、错误、空数据、无权限、禁用和危险操作状态；
- 深色模式、响应式、键盘操作和基本可访问性。

“全站”明确包含登录页、注册页、所有认证辅助页面、用户前台和完整管理后台，不允许只完成用户端或只替换应用外壳后宣告完成。

重做不自动包含：

- 参考前端存在但 Sub2API 后端没有的业务；
- 对现有计费、调度或权限规则的重新定义；
- New API 兼容接口；
- React 与 Vue 混合运行；
- 未经核对的品牌素材、付费字体或许可证不明确的资源。

### 7.2 阶段与门禁

| 阶段 | 内容 | 完成门禁 | 状态 |
|---|---|---|---|
| P0 基线 | 记录现有页面、关键流程、构建结果和测试基线 | 能登录、构建、运行关键测试；保存基线截图清单 | 已完成 |
| P1 设计基础 | 逐项复刻设计令牌、基础样式、通用 UI 原语 | 与参考前端完成组件级视觉对照，深浅主题、响应式和通用状态稳定 | 已完成 |
| P2 应用外壳与认证 | 复刻布局、导航、公开页、登录、注册及全部认证辅助页 | 视觉对照通过；路由守卫、OAuth 回调和移动导航无回归 | 已完成 |
| P3 用户端 | Dashboard、Keys、Usage、Redeem、Profile、Subscriptions、Purchase 等 | 用户关键流程逐页验收 | 已完成 |
| P4 管理端 | Dashboard、Users、Groups、Channels、Accounts、Proxies、Ops、Settings 等 | 权限、批量操作和危险操作逐页验收 | 已完成 |
| P5 收尾 | 一致性、性能、可访问性、文档和旧样式清理 | 全量检查通过，无未登记的旧组件 | 已完成 |

状态只能使用：`待开始`、`进行中`、`阻塞`、`已完成`。阶段状态变化时更新本表和第 13 节维护日志。

### 7.3 推荐实施顺序

1. 建立视觉基线和页面清单，不先改业务。
2. 从参考前端提取语义化设计令牌，不逐页面复制颜色值。
3. 先重做 Button、Input、Select、Dialog、Table、Badge、Toast、Skeleton 等通用原语。
4. 重做 App Shell，再处理登录、注册和公开页面。
5. 完成用户端高频路径：登录 -> Dashboard -> Keys -> Usage -> Purchase/Profile。
6. 完成管理端高风险路径：Users -> Groups -> Channels -> Accounts -> Settings。
7. 最后处理 Ops、审计、支付管理和低频页面。
8. 删除旧样式前用 `rg` 确认没有调用方，并在同一提交中完成验证。
9. 每个阶段用相同视口、相同主题和相同数据状态对参考前端与 Vue 实现进行并排截图对照。

### 7.4 单页面重做模板

每个页面开始前记录：

- 路由、访问角色和功能开关；
- 使用的 API、请求参数、分页和响应类型；
- loading、empty、error、disabled、permission-denied 状态；
- 所有创建、更新、删除和批量操作；
- 桌面与移动端关键布局；
- 已有测试和必须新增的回归测试。

实施顺序：

1. 先运行或补充行为测试，锁定现有业务；
2. 在 Vue 中组合现有 API/store/composable，保持数据流不变；
3. 使用统一设计令牌和通用组件实现新视觉；
4. 逐状态检查，不只检查有数据的成功页面；
5. 执行类型检查、相关测试、构建和浏览器人工验收；
6. 在提交说明中标明“仅视觉”或列出真实业务变化。

### 7.5 单页面完成定义

- 功能与重做前一致，新增行为有明确需求和测试；
- 普通用户、管理员和未登录用户的权限行为正确；
- loading、empty、error、disabled 和成功反馈齐全；
- 最窄支持视口没有横向溢出、文本遮挡或不可点击控件；
- 深色和浅色主题均可读；
- 与参考前端对应页面完成桌面和移动端并排视觉对照，主要布局、间距、尺寸、字体、颜色、边框、圆角、阴影、图标和信息密度一致；
- 表单键盘可操作，输入有 label，焦点状态清晰；
- 用户文案进入 i18n，不新增散落的硬编码文本；
- 不绕过统一 API client；
- `pnpm run lint:check`、`pnpm run typecheck`、相关 Vitest 和 `pnpm run build` 通过；
- 对关键页面保留桌面和移动端验收截图，截图不得包含真实凭据。
- 页面没有可见的旧版 Vue 视觉残留；Sub2API 独有页面使用同一套设计系统并通过一致性审查。

## 8. 后续后端功能增加规则

### 8.1 先登记，再实现

本 Fork 独有后端功能必须先在下表登记。一个功能使用稳定 ID，例如 `FORK-001`，并在相关迁移、文档或提交说明中引用。

| ID | 功能 | 状态 | 后端入口 | 前端入口 | 迁移 | 上游重叠风险 |
|---|---|---|---|---|---|---|
| FORK-001 | 持久化工单中心与管理员工单管理 | 已发布 | `backend/internal/server/routes/user.go` 的 `/api/v1/tickets`；`backend/internal/server/routes/admin.go` 的 `/api/v1/admin/tickets` | `/tickets`、`/tickets/:id`、`/admin/tickets`、`/admin/tickets/:id`；系统设置 → 站点设置 → 侧边栏模块 | `backend/migrations/185_tickets.sql` | 上游若增加工单模块，重点核对状态机、逐管理员阅读位置、权限与设置键，禁止整目录覆盖 |
| FORK-002 | 订阅权益快照与独立 Omnio Pro 体系 | 已发布 | `/api/v1/membership`、`/api/v1/admin/membership`、`/api/v1/groups/entitlements`；现有支付和统一计费服务的小型扩展点 | `/omnio-pro`（兼容 `/membership`）、`/purchase?tab=membership`、`/admin/omnio-pro`（兼容 `/admin/membership`）、后台分组管理；每个分组可添加多个会员等级并分别配置倍率 | `backend/migrations/186_membership_and_subscription_entitlements.sql`、`backend/migrations/187_omnio_pro_group_visibility.sql`、`backend/migrations/189_omnio_pro_group_settings.sql`、`backend/migrations/190_omnio_pro_free_quota.sql`、`backend/migrations/191_omnio_pro_level_benefits_and_quota.sql` | 上游若增加会员、套餐版本、订阅超额或自动续费功能，重点核对“人工用户倍率 > 等级与分组共同确定的 Pro 倍率 > 分组普通倍率”、按等级独立的 Pro 专属分组、每日额度、每月额度及用量计数、免费额度原子拆分、订单幂等、退款反向发放和计费缓存；会员优惠倍率的管理入口必须留在分组管理，且修改或删除只作用于当前分组与当前等级；禁止把 Pro 倍率写回分组普通倍率，禁止退回全局 Pro 配置，也禁止把独立 Pro 免费额度重新耦合到旧订阅 |
| FORK-003 | Omnio Pro 统一商品并退役公开订阅 | 已完成 | 新订单仅允许 `balance`、`membership`；历史 subscription API、已支付订单履约和退款保留兼容 | `/omnio-pro` 为唯一会员入口；`/subscriptions`、`/admin/subscriptions` 自动跳转 Omnio Pro；购买页隐藏订阅标签 | `backend/migrations/188_omnio_pro_retire_subscriptions.sql` | 不删除历史订阅表和已支付权益；新建订阅订单返回 `SUBSCRIPTIONS_RETIRED`，所有新商品通过 Omnio Pro levels/offers/grants 管理；上游合并时优先保留该边界 |

新增功能时填写实际路径，不要只写模块名称。状态使用：`设计中`、`开发中`、`已发布`、`已废弃`。

### 8.2 端到端实现顺序

1. **需求与权限**：明确角色、租户/用户边界、开关、审计和失败语义。
2. **契约**：先定义路由、DTO、错误码、分页、幂等性和兼容策略。
3. **数据**：设计 Ent schema 或 SQL 迁移、索引、约束、回填和回滚/前滚策略。
4. **后端分层**：按现有结构落到 domain/model、repository、service、handler 和 `internal/server/routes/`。
5. **安全**：校验输入、权限、敏感字段脱敏、日志脱敏、SSRF/文件大小/速率限制等风险。
6. **前端契约**：在 `frontend/src/types/` 和 `frontend/src/api/` 增加类型化调用。
7. **前端页面**：复用已经重做的 Vue 组件与设计令牌，不创建孤立的第二套样式。
8. **测试**：覆盖 service/handler、权限失败、迁移回归、API client 和关键页面行为。
9. **文档与登记**：更新本节功能表、配置说明、部署注意事项和用户文档。

### 8.3 API 兼容原则

- 优先做向后兼容的字段新增，不随意重命名或改变字段语义；
- 新增可选字段时明确默认值和旧客户端行为；
- 删除或收紧行为必须有迁移期、版本说明和调用方搜索结果；
- handler 不直接承载复杂业务，业务规则进入 service/domain；
- 前端不得依赖未公开、未测试的临时响应字段；
- 所有管理员接口继续经过现有管理员权限与合规中间件；
- 影响余额、计费、支付和配额的接口必须考虑并发、幂等和审计。

### 8.4 数据库迁移原则

- 永不编辑已经进入共享分支或已部署环境的历史迁移；
- 新迁移编号基于合并上游后的最新状态选择，允许仓库既有的同号命名模式，但文件名必须唯一且语义明确；
- 大表变更评估锁表时间，索引优先考虑非事务/在线策略及仓库现有约定；
- 数据回填必须可重入、可观察，并为中断后继续执行设计；
- 删除列或数据至少分为“停止写入 -> 回填/观察 -> 删除”多个版本；
- 修改 Ent schema 后运行 `go generate ./ent`，检查并提交必要生成物；
- 上游合并出现迁移冲突时保留双方迁移，重新评估顺序，不用覆盖文件解决。

### 8.5 定制功能的可合并性

为了降低未来上游冲突：

- 新功能尽量放在独立文件和独立路由组中；
- 对上游核心函数优先增加小型扩展点，不复制整段实现；
- 配置项使用明确的 Fork 命名和默认关闭策略，除非它是必要的兼容修复；
- 定制逻辑必须有测试，使上游合并后能快速判断是否仍然生效；
- 每项定制记录它与上游模块的接触点以及未来删除条件。

## 9. 上游 Sub2API 更新合并策略

### 9.1 一次性配置官方上游

先检查远端：

```powershell
git remote -v
```

仅当 `upstream` 不存在时添加：

```powershell
git remote add upstream https://github.com/Wei-Shaw/sub2api.git
git fetch upstream --prune --tags
```

如果 `upstream` 已存在但 URL 不同，先确认原因，不要直接覆盖。

### 9.2 分支模型

| 分支 | 用途 |
|---|---|
| `main` | 始终保持可发布的 Fork 主线 |
| `feature/ui-<module>` | 单个前端模块重做 |
| `feature/FORK-<id>-<name>` | 本 Fork 独有功能 |
| `sync/upstream-YYYYMMDD` | 一次上游同步和冲突解决 |
| `hotfix/<name>` | 已发布版本的紧急修复 |

不要在共享的 `main` 上执行历史重写。功能提交保持小而清楚，避免把格式化、视觉改动、后端功能和上游同步混在一个提交中。

### 9.3 每次同步流程

同步前：

1. 确认工作树干净或已有明确提交；
2. 记录当前版本、提交、数据库迁移最高状态和部署配置；
3. 阅读上游 release notes、迁移和安全修复；
4. 确认当前定制功能表和前端阶段表已经更新；
5. 对生产数据执行独立备份，验证恢复方式。

建议命令：

```powershell
git fetch origin --prune --tags
git fetch upstream --prune --tags
git switch main
git pull --ff-only origin main
git switch -c sync/upstream-YYYYMMDD
git merge --no-ff --no-commit upstream/main
```

解决冲突和验证完成后再提交 merge。不要因为冲突很多就把上游整体 squash 成无法追踪的补丁。

### 9.4 冲突归属矩阵

| 冲突区域 | 默认处理原则 |
|---|---|
| `frontend/src/styles/`、纯视觉组件 | 保留 Fork 视觉，再人工吸收上游行为和可访问性修复 |
| 页面 `.vue` | 先识别上游业务变化，将变化移植到定制布局，不能只保留外观 |
| `frontend/src/api/`、types、auth store | 以上游契约和安全修复为基线，再恢复必要的 Fork 扩展 |
| `frontend/src/router/` | 保留双方新路由与守卫，逐条核对角色和功能开关 |
| `frontend/vite.config.ts` | 必须保留 dist 输出、配置注入、代理和 CSP 兼容 |
| `frontend/package.json`、lockfile | 先合并 package 声明，再用 pnpm 重新生成并验证 lockfile |
| 后端 routes/handler/service | 先理解调用链，保留上游修复与 Fork 功能，不整文件选边 |
| Ent schema、生成代码 | 先合并 schema，再重新生成，禁止手工拼凑生成文件 |
| SQL migrations | 保留双方文件，检查编号、顺序、约束、索引和回填交互 |
| Docker、CI、Makefile | 合并上游工具链变更，同时保留正式前端构建和完整测试门禁 |
| 文档 | 去除过期信息，保留 Fork 身份和本项目决策记录 |

特别注意：视觉定制文件也可能包含上游新增的权限判断、字段处理或错误状态。不能因为它“看起来是 UI 文件”就直接选择 Fork 版本。

### 9.5 合并后的检查顺序

1. 检查 `git status` 和所有未合并标记；
2. 搜索冲突残留：`rg -n "^(<<<<<<<|=======|>>>>>>>)"`；
3. 检查 migration、Ent schema、Go module、package 和 lockfile；
4. 运行后端生成、格式化、测试和 lint；
5. 运行前端安装、lint、类型检查、测试和构建；
6. 构建 Docker 镜像，确认 Go binary 内含新前端；
7. 按第 10 节执行人工回归；
8. 更新本文档的基线、阶段状态、功能表和维护日志；
9. 通过代码审查后合并同步分支，再标记 Fork 发布版本。

### 9.6 合并频率

优先小步、定期合并，不累计多个大型上游版本。遇到以下情况立即同步评估：

- 鉴权、Token、OAuth 或权限安全修复；
- 计费、余额、支付或数据库迁移修复；
- 上游 API 契约变化；
- Go、Node、Vue、Vite 或 pnpm 的安全更新；
- 上游修改了正在重做的同一页面。

## 10. 验证矩阵

### 10.1 最小充分验证原则

检查应与风险和改动范围匹配，不以重复运行命令代替分析：

- 开发过程中优先运行受影响文件、模块或关键路径的定向检查；
- 完整 lint、类型检查、全量测试和生产构建通常只在阶段合并、上游同步完成或发布前执行一次；
- 如果代码、依赖、配置、生成物和运行环境都没有变化，不重复运行已经通过的同一检查；
- 只有相关文件再次变化、上次检查失败、环境发生变化、结果疑似不稳定或进入新的发布门禁时才重新运行；
- 鉴权、权限、计费、支付、数据库迁移和上游合并等高风险改动仍须执行对应的必要回归；
- 记录已执行的命令和结果，后续人员先复用有效结果，再决定是否需要补充检查。

### 10.2 前端标准检查

在 `frontend/` 中执行：

```powershell
pnpm install --frozen-lockfile
pnpm run lint:check
pnpm run typecheck
pnpm run test:run
pnpm run build
```

开发期间可以只运行相关 Vitest，但合并阶段必须运行完整门禁。依赖变化时确认 `pnpm-lock.yaml` 同步。

### 10.3 后端标准检查

在 `backend/` 中执行：

```powershell
go generate ./ent
go generate ./cmd/server
go test ./...
golangci-lint run ./...
```

涉及集成环境时补充：

```powershell
go test -tags=integration ./...
go test -tags=e2e -v -timeout=300s ./internal/integration/...
```

不要为了通过检查随意更新生成物；先确认生成变化是否来自本次 schema、路由或版本信息变更。

### 10.4 人工关键流程

每次上游同步、认证改动、应用外壳改动或正式发布至少检查：

- 首次安装与 `/setup`；
- 登录、注册、退出、Token 过期刷新、2FA；
- 已启用的 OAuth 登录和回调；
- 普通用户与管理员路由隔离；
- Dashboard、Keys、Usage、Redeem、Profile；
- 订阅购买、支付发起、支付结果和订单恢复；
- 管理员 Users、Groups、Channels、Accounts、Proxies、Settings；
- Ops、审计、风控和危险操作确认；
- 深色/浅色、桌面/移动端、空数据/错误状态；
- 刷新任意 SPA 深层路由仍能返回正确 `index.html`；
- 生产构建中站点名称、Logo 和 `window.__APP_CONFIG__` 注入正常。

## 11. 发布与回退

- 发布前为 Fork 使用独立版本或标签，记录对应上游提交；
- 数据库升级前备份并验证恢复，不把“能生成备份”当作“能恢复”；
- 应用回退与数据库回退分开设计，迁移优先采用向前修复；
- 上游 merge 引入问题时优先用新的 revert/修复提交，不重写已共享历史；
- 视觉回退不得连带回退安全、鉴权、支付或数据库修复；
- 发布说明分别列出：上游变化、Fork 功能、前端重做进度、迁移和配置变化。

## 12. 许可证与素材

- Sub2API 根仓库当前为 LGPLv3；
- 参考 React 前端多个源码文件声明 AGPLv3，但参考目录创建本文档时没有根 `LICENSE` 文件；
- 复制代码、文案、Logo、图片、字体或其他资源前必须确认来源、许可证和署名要求；
- 默认只借鉴视觉与交互思想，在 Vue 中独立实现；
- 如果确需复制 AGPL 源码，必须保留版权头并在发布前完成许可证义务评估；
- 本节不是法律意见，许可证不清楚时停止复制具体资产或源码。

## 13. 维护记录

| 日期 | 基线/版本 | 变更 | 决策或后续动作 |
|---|---|---|---|
| 2026-07-20 | `d4b9797f` / `0.1.161` | 创建长期项目记忆；确认方案一及全站高保真复刻目标；明确包含登录、注册、用户端和完整管理后台；记录当前 origin 和官方 upstream | 下一步执行 P0 基线与页面清单 |
| 2026-07-20 | `d4b9797f` / `0.1.161` | 完成 P0-P5：以参考 React 前端为唯一视觉基准，统一认证页、公开页、用户端、管理端、应用外壳、通用组件、深浅主题和移动端；保留 Vue、Sub2API API、鉴权、权限与业务语义 | 进入维护模式；后续新增功能沿用本设计系统，上游同步按第 9 节逐项移植业务变化 |
| 2026-07-20 | `d4b9797f` / `0.1.161` | 完成全站收尾：接入参考端现成首页/认证素材，补齐搜索面板、用户订阅与兑换、支付/Profile/Usage、管理 Settings/Ops/订单、低频辅助页和复杂弹窗；统一设计令牌、通用组件及深浅主题，并同步认证回调测试夹具 | 全量 lint、类型检查、测试和生产构建通过；本次不重复既有浏览器视觉验收，最终主观验收由项目负责人执行 |
| 2026-07-21 | `d4b9797f` / `0.1.161` | 登记并实现 FORK-001 工单第一版：工单/消息/逐账户阅读位置独立存储；用户新建、对话、回复、解决、关闭；管理员筛选、回复、领取/转派、优先级与状态管理；增加“显示用户工单中心”和“接受新工单”两个独立开关；侧边栏 30 秒轮询未读数字；首版不支持附件 | 工单定向单测、后端必要编译、前端 typecheck/lint/生产构建已通过；未来附件必须单独设计安全对象存储、鉴权下载、恶意文件检测与多节点一致性 |
| 2026-07-21 | `sub2api:custom-20260721` | FORK-001 与全站定制前端正式部署至 `omni0.top` 和 `www.omni0.top`；应用仅监听 `127.0.0.1:8080`，由宝塔 Nginx 反向代理；PostgreSQL、Redis 与应用容器均健康，工单迁移表已建立 | 双域名 Let's Encrypt HTTPS 已启用，HTTP 统一跳转主域名，证书自动续期及 Nginx 重载演练通过；管理员凭据仅保存在服务器 `/root/sub2api-admin-credentials.txt` |
| 2026-07-21 | `docs-site` | 在 `/docs/` 发布独立静态 Omnio 文档中心，首批收录《用户服务协议》和《数字充值服务政策》；由 Nginx 直接提供，不改动 Sub2API 应用与端口；生产环境 `doc_url` 已设为 `https://omni0.top/docs/`，主站顶部、首页、页脚及用量查询页的“文档”入口统一指向该地址 | 源码位于 `docs-site/`，后续通过 `content/` 与 `build.mjs` 增加或更新文档；生产 HTML 注入存在进程内缓存，绕过设置服务直接改数据库后必须重载应用或主动失效缓存；本次已重载并验证 HTML 注入地址及文档页；法律协议入口继续使用原有 `/legal/:documentId`，不与文档入口混用 |
| 2026-07-21 | `docs-site/brand-story` | 发布 Omnio 品牌故事页，正式记录从“词元101”升级为 Omnio 的品牌演进，并以 `Omni + I/O` 解释名称方向，明确简单、稳定、透明、同行四项品牌坚持；入口已加入文档首页和侧边栏 | 页面位于 `/docs/brand-story/`，内容事实边界保持为品牌表达，不虚构公司历史、客户规模或绝对服务承诺；后续品牌定位变化时同步更新该页 |
| 2026-07-21 | `docs-site/usage-guide` | 文档中心按“了解品牌 → 使用教程 → 法律与政策”重组，并发布完整使用教程，覆盖账号、API Key、模型查询、OpenAI/Anthropic/Gemini 调用、SDK、流式响应、用量、计费、安全、错误排查与工单支持 | 教程不写死具体模型名，要求通过 `/v1/models` 获取；接口与当前 `gateway.go` 对齐，后续网关端点或控制台流程变化时同步更新；按项目要求仅执行构建、内容与在线可用性检查，视觉验收由项目负责人完成 |
| 2026-07-21 | `docs-site/about-omnio` | 在“了解 Omnio”中新增独立《关于 Omnio》页面并置于《品牌故事》之前；内容聚焦产品定位、适用对象、解决的问题、使用流程、服务体验与边界，避免重复品牌更名历史；主站桌面端和移动端“关于”源码入口改为 `/docs/about-omnio/` | 生产 Nginx 同时将旧 `/legal/about` 永久重定向到新页面，保证当前已部署前端与历史链接立即生效；页面位于 `/docs/about-omnio/`，后续产品定位或服务边界变化时同步更新 |
| 2026-07-21 | `docs-site/pricing` | 在“使用教程”下方新增“模型定价”，按 OpenAI、Anthropic、Google、DeepSeek、Alibaba Cloud 五个官方上游分组列出常用文本模型公开价格，并用醒目提示说明 Omnio 最终价格按分组倍率计算：1.0 为原价、0.5 为半价、2.0 为两倍 | 页面位于 `/docs/pricing/`；官方价格最后核对日期为 2026-07-21，表格明确标准实时 API、上下文、区域和优惠边界；厂商调价、模型更名或平台分组规则变化时必须重新核对官方来源并同步更新，实际扣费始终以 Omnio 控制台倍率和用量记录为准 |
| 2026-07-22 | `sub2api:custom-20260722-nav` | 主站顶部导航将“模型广场、排行榜”替换为“定价、教程”，桌面端和移动端分别直接链接 `/docs/pricing/` 与 `/docs/usage-guide/`，中英文文案同步为 Pricing/Tutorial | 新镜像已发布，应用、PostgreSQL、Redis 与文档页保持健康；以后调整文档路径时必须同时核对主站顶部导航和文档中心路由 |
| 2026-07-22 | `sub2api:custom-20260722-membership` | 完成 FORK-002 首版：新增会员等级、独立 VIP 方案、分组访问/倍率/RPM 权益、手动发放、订阅赠送 VIP、审计日志、订单支付履约与退款撤销；订阅订单写入购买时权益快照，计费优先使用管理员倍率、VIP 倍率和快照倍率；前端新增会员页、VIP 购买页、订阅快照额度与超额策略开关、完整管理员会员控制页 | 通过 `go test ... -run '^$'` 后端定向编译、`pnpm run typecheck` 和 `pnpm run build`；全量 service 测试曾因既有 Redis/外部运行时测试噪声失败，未将其误记为本次功能回归；首版明确不包含自动续费、成长值、按天差价和附件 |
| 2026-07-22 | `FORK-002-omnio-pro` | 将独立 VIP 正式更名并扩展为 Omnio Pro：人工用户倍率、Pro 倍率和分组基础倍率分层存储，最终优先级固定为“人工 > Pro > 基础”；API Key 分组选择器和可用渠道同时展示基础倍率与 Pro/人工最终倍率；新增按等级配置的 Pro 专属分组，普通用户不可见、不可绑定，Pro 失效后已有 Key 也会在网关鉴权时被拒绝 | 前台主入口为 `/omnio-pro` 与 `/admin/omnio-pro`；本地后端受影响包编译、前端 typecheck 与生产构建通过，视觉验收按项目要求由负责人执行；已随 FORK-003 部署生产 |
| 2026-07-22 | `sub2api:custom-20260722-omnio-pro-fix1` / `FORK-003` | 公开订阅商品正式并入 Omnio Pro：旧订阅前端入口自动跳转，购买页只展示充值和 Omnio Pro；后端禁止新建 subscription 订单；迁移 187/188 已在生产执行，历史可售订阅方案仅在存在时迁移为 Omnio Pro offers 和基础访问权益 | 保留历史订阅订单履约、退款、表结构和内部 API，避免既有用户权益丢失；生产当前没有历史订阅方案，因此已创建 Omnio Pro 等级但没有自动生成售卖方案，管理员须在 `/admin/omnio-pro` 配置价格和期限后才能购买；本次不重复视觉验收 |
| 2026-07-23 | `sub2api:custom-20260723-omnio-pro-groups` / `FORK-002-omnio-pro-group-settings` | 修正 Omnio Pro 分组定价模型并部署生产：每个分组分别保存“普通用户倍率”和“Omnio Pro 倍率”，例如同一站点可让 plus 号池使用 `0.2 / 0.15`、pro 号池使用 `0.3 / 0.25`；计费优先级固定为“人工用户倍率 > 分组独立 Pro 倍率 > 分组普通倍率”；API Key 分组选择器对 Pro 用户展示普通倍率删除线和格式化后的 Pro 实际倍率 | 独立 Pro 倍率与“仅 Omnio Pro 可见和绑定”开关统一放在后台“分组管理”的新建/编辑表单；生产迁移 189 已执行并新增 `omnio_pro_group_settings` 作为分组级来源，兼容旧等级权益；Omnio Pro 等级列表增加安全删除，存在方案、发放或权益关联时拒绝删除 |
| 2026-07-23 | `sub2api:custom-20260723-omnio-pro-quota` / `FORK-002-omnio-pro-free-quota` | 为每个分组增加独立的 Omnio Pro 每日/月度免费美元额度并部署生产；有效 Pro 用户请求先按最终 Pro 倍率计算费用，再在统一账务事务中原子拆分为免费额度与钱包扣费，按北京时间自然日和自然月重置；每日或月度任一上限用完后，超出部分自动扣钱包 | 生产迁移 190 已执行，新增额度配置、当前窗口聚合和不可变逐请求拆分事件；余额不足但仍有免费额度时允许请求，跨额度边界的单次请求按剩余额度与钱包精确拆分；后台分组编辑页配置额度，用户 `/omnio-pro` 页面展示今日和本月用量与剩余 |
| 2026-07-23 | 官方 `0.1.164` / `cb24522dd53f8f363d008e3afdc8e4baf9788cab` | 从 Omnio 快照 `eb3087cc` 建立 `sync/upstream-20260723`，通过透明导入提交 `204c3bbd` 合并官方源码；吸收组合模型路由、分组 reasoning effort 策略、Ollama Cloud 用量、图片存储/鉴权缓存、支付宝移动端 deep link、运维移动端和可访问性修复 | 保留 Omnio Vue 视觉、Logo、首页、认证、导航、工单、Omnio Pro 及公开订阅退役边界；支付响应同时保留 `membership_offers` 与 `alipay_mobile_precreate_deep_link`，安装向导保留 Omnio 页面并加入 Redis ACL 用户名；新增迁移 `185_group_reasoning_effort_policy.sql`、`186_alipay_mobile_precreate_deep_link.sql`、`186_group_auth_cache_image_generation.sql`，与 Fork 既有同号迁移按完整文件名独立执行；本地验证已通过，Docker 镜像构建与生产发布等待可用服务器登录 |
| 2026-07-24 | `sub2api:custom-20260724-0.1.164` / `d9811ce5` | 在生产服务器完成备份、服务器端 Docker 多阶段构建和应用容器滚动替换；仅重建 `sub2api`，未重启 PostgreSQL/Redis | 备份位于服务器 `/opt/sub2api/backups/20260723T164802Z-pre-0.1.164/`，包含数据库 dump、源码配置归档和旧镜像；四个新增迁移已执行，应用/数据库/Redis 均 healthy，内网和 `https://omni0.top/health` 均返回 `{"status":"ok"}`；生产源码同步到 `0.1.164`，旧镜像仍保留可回退 |
| 2026-07-24 | `sub2api:custom-20260724-omnio-pro-levels` / `e96dc50e` | 完成 Omnio Pro 四档独立化并部署生产：`Omnio Pro`、`Omnio Pro Max`、`Omnio Pro Ultra`、`Omnio 内测` 的分组倍率、专属分组、每日免费额度、每月免费额度均按等级与分组分别保存；用量按 `(user_id, level_id, group_id)` 独立累计，调整任一档不再联动其他档 | 迁移 191 已执行，将旧全局配置和旧用量复制到现有四档，避免上线后权益或额度突变；运行时已切换到等级权益和等级用量表，旧表仅为旧镜像回退保留；备份位于 `/opt/sub2api/backups/20260723T182048Z-pre-omnio-pro-levels/`，仅替换 `sub2api`，PostgreSQL/Redis 未重启，旧镜像 `sub2api:custom-20260724-0.1.164` 仍可回退 |
| 2026-07-24 | `sub2api:custom-20260724-group-member-rates` / `3f97d49e` | 将各独立分组的会员优惠倍率统一放入分组管理并部署生产：新建或编辑分组时可添加多个不同会员等级，每个等级独立设置当前分组倍率；`0` 表示免费，留空继承分组基础倍率，重复等级不能添加 | 继续复用 `membership_level_group_benefits`，无数据库迁移；保存倍率时保留该等级的专属分组、RPM、每日和每月额度，删除只影响当前分组与当前等级。备份位于 `/opt/sub2api/backups/20260723T192600Z-pre-group-member-rates/`；仅替换 `sub2api`，PostgreSQL/Redis 容器 ID 未变，旧镜像 `sub2api:custom-20260724-omnio-pro-levels` 仍可回退 |

每次更新记录一行，内容应说明“为什么”，不能只写“更新文档”。当表格过长时，可按年度归档到 `docs/history/`，但保留最近一年和所有未完成决策。

## 14. 当前完成基线

### 14.1 视觉与组件事实源

- 全局视觉令牌和兼容覆盖：`frontend/src/style.css`、`frontend/tailwind.config.js`；
- 应用外壳：`frontend/src/components/layout/AppHeader.vue`、`AppSidebar.vue`、`AppLayout.vue`；
- 顶栏搜索：`frontend/src/components/common/HeaderSearch.vue`，支持 `Ctrl/Cmd+K`；
- 认证外壳：`frontend/src/components/layout/AuthLayout.vue`，覆盖登录、注册、忘记密码、重置密码、邮箱验证、2FA 和 OAuth 过程页；
- 表格页结构：`frontend/src/components/layout/TablePageLayout.vue`、`frontend/src/components/common/DataTable.vue`、`Pagination.vue`；
- 参考端现成素材已复制到 `frontend/public/assets/`，正式构建不依赖仓库外的 `../前端/`；
- 主题、图表和页面级差异继续通过现有 Vue 组件实现，不建立第二套 React 运行时；
- 视觉验收报告：`design-qa.md`；稳定截图：`docs/design-qa/`。

### 14.2 2026-07-20 验证结果

- `pnpm run lint:check`：全量通过；
- `pnpm run typecheck`：通过；
- `pnpm run test:run`：全量通过；认证外壳新增语言与主题控件后，OAuth/微信支付回调测试夹具已按稳定选择器和完整 store stub 同步；
- `pnpm run build`：通过，产物正确写入 `backend/internal/web/dist`；
- 浏览器验收：`1280x720` 深色/浅色、`390x844` 移动端、登录/注册/忘记密码、用户 API Keys、管理员 Users 通过；移动端文档宽度等于视口宽度，无横向溢出；
- 本次最终代码收尾按项目负责人要求不重复浏览器视觉验收，复用既有对照证据，最终主观视觉验收由项目负责人执行；
- `design-qa.md` 最终结果为 `passed`。
- 2026-07-22 FORK-002 发布门禁：后端 `cmd/server`、service、handler、repository、middleware 定向编译通过；前端 `pnpm run typecheck` 与 `pnpm run build` 通过，产物已写入 `backend/internal/web/dist`；本次按“只做必要检查”原则不重复全量视觉验收。
- 2026-07-22 Omnio Pro 扩展门禁：`go test ./cmd/server ./internal/service ./internal/handler ./internal/repository ./internal/server/middleware -run '^$'`、Pro 分组权限定向单测、`pnpm run typecheck`、`pnpm run build` 通过；没有重复运行全量前端测试或浏览器视觉验收。
- 2026-07-23 分组级 Omnio Pro 配置门禁：受影响的 service、repository、admin handler 和 routes 后端定向编译通过，前端 `pnpm run typecheck` 与 `pnpm run build` 通过，`git diff --check` 通过；按项目要求不重复构建和视觉验收。
- 2026-07-23 Omnio Pro 免费额度门禁：后端受影响包定向编译通过，额度跨日/月边界原子拆分单测通过；前端 `pnpm run typecheck` 与一次 Vite 生产构建通过；未重复执行全量测试或视觉验收。
- 2026-07-23 官方 `0.1.164` 同步门禁：`go generate ./ent`、`go generate ./cmd/server` 完成；`go test ./...` 全量通过；前端 `pnpm run lint:check`、`pnpm run typecheck`、`pnpm run test:run`（187 个测试文件、1288 个测试）和 `pnpm run build` 全量通过；`git diff --check` 通过且无冲突标记。当前 Windows 环境没有 Docker CLI，候选镜像须在生产服务器或 Docker 构建机上生成后再发布。
- 2026-07-24 生产发布门禁：服务器 Docker 构建完成，镜像内 `vue-tsc`、Vite 和 Go embed 构建通过；新容器 `sub2api:custom-20260724-0.1.164` healthy，PostgreSQL/Redis 保持 healthy；`schema_migrations` 已登记 `172_composite_model_routes.sql`、`185_group_reasoning_effort_policy.sql`、`186_alipay_mobile_precreate_deep_link.sql`、`186_group_auth_cache_image_generation.sql`；内网、公网 `/health` 及关键用户/管理员/Omnio Pro/API 请求均成功。
- 2026-07-24 Omnio Pro 四档独立化门禁：后端等级权益保存/删除定向测试与每日/月度额度跨窗口原子拆分 SQL mock 测试通过；前端 `pnpm run lint:check`、`pnpm run typecheck`、生产构建通过；`go test ./...` 首次仅在既有 `internal/service` 并发运行时用例抖动失败，单独重跑整个 `internal/service` 通过，未发现本次功能相关失败。生产镜像内前端和 Go 构建通过，容器 healthy；`schema_migrations` 已登记 `191_omnio_pro_level_benefits_and_quota.sql`，四档各有 4 条旧配置精确回填，等级用量/事件表已产生 12/7 条记录；PostgreSQL/Redis 保持原实例 healthy，内外网 `/health` 均返回 `{"status":"ok"}`，`/omnio-pro` 与 `/admin/omnio-pro` 均返回 200。
- 2026-07-24 分组会员倍率门禁：定向测试 3 个通过，分组页面回归 4 个文件、14 个测试通过，`pnpm run lint:check`、`pnpm run typecheck`、`pnpm run build` 通过；桌面端和 `390px` 移动端预览无横向溢出。生产镜像内 `vue-tsc`、Vite 和 Go embed 构建通过，新容器 `sub2api:custom-20260724-group-member-rates` healthy、重启次数为 0；内外网 `/health` 均返回 `{"status":"ok"}`，生产 `GroupsView` 静态资源返回 200 且包含新增会员等级、重复限制与继承提示，PostgreSQL/Redis 保持原实例 healthy。

### 14.3 后续维护入口

1. 开始任何新工作前先完整阅读本文档，只执行与改动风险匹配的必要检查，不重复已经有效的同一门禁。
2. 新后端功能先按第 8 节登记稳定 `FORK-xxx` ID，再接入现有 Vue API、类型、权限和设计系统。
3. 上游同步前核对并在需要时配置官方 `upstream`；使用独立 `sync/upstream-YYYYMMDD` 分支和 merge commit。
4. 上游页面冲突先提取 API、权限、状态和字段变化，再移植进本 Fork 的布局；禁止整文件覆盖定制视觉，也禁止为了保视觉而丢失上游业务修复。
5. 上游新增页面或本 Fork 新功能必须复用当前 Header、Sidebar、TablePageLayout、表单、弹窗、主题和响应式规则，不能恢复旧版 Vue 风格。
- 2026-07-22 production deployment: `sub2api:custom-20260722-omnio-pro-fix1` is active on `64.83.34.203`; migrations 187 and 188 are applied; `sub2api`, PostgreSQL, and Redis are healthy; `https://omni0.top/health` returned `{"status":"ok"}`. Pre-deployment database, source archive, and container image backups remain on the server.
