export default {
  batchImageGuide: {
    title: '图片批量生成',
    description: '一次提交多条提示词，任务完成后可统一下载图片结果'
  },
  // Home Page
  home: {
    viewOnGithub: '在 GitHub 上查看',
    viewDocs: '查看文档',
    docs: '文档',
    switchToLight: '切换到浅色模式',
    switchToDark: '切换到深色模式',
    dashboard: '控制台',
    login: '登录',
    getStarted: '立即开始',
    goToDashboard: '进入控制台',
    // 新增：面向用户的价值主张
    heroSubtitle: '一个密钥，畅用多个 AI 模型',
    heroDescription: '无需管理多个订阅账号，一站式接入 Claude、GPT、Gemini 等主流 AI 服务',
    tags: {
      subscriptionToApi: '订阅转 API',
      stickySession: '会话保持',
      realtimeBilling: '按量计费'
    },
    // 用户痛点区块
    painPoints: {
      title: '你是否也遇到这些问题？',
      items: {
        expensive: {
          title: '订阅费用高',
          desc: '每个 AI 服务都要单独订阅，每月支出越来越多'
        },
        complex: {
          title: '多账号难管理',
          desc: '不同平台的账号、密钥分散各处，管理起来很麻烦'
        },
        unstable: {
          title: '服务不稳定',
          desc: '单一账号容易触发限制，影响正常使用'
        },
        noControl: {
          title: '用量无法控制',
          desc: '不知道钱花在哪了，也无法限制团队成员的使用'
        }
      }
    },
    // 解决方案区块
    solutions: {
      title: '我们帮你解决',
      subtitle: '简单三步，开始省心使用 AI'
    },
    features: {
      unifiedGateway: '一键接入',
      unifiedGatewayDesc: '获取一个 API 密钥，即可调用所有已接入的 AI 模型，无需分别申请。',
      multiAccount: '稳定可靠',
      multiAccountDesc: '智能调度多个上游账号，自动切换和负载均衡，告别频繁报错。',
      balanceQuota: '用多少付多少',
      balanceQuotaDesc: '按实际使用量计费，支持设置配额上限，团队用量一目了然。'
    },
    // 优势对比
    comparison: {
      title: '为什么选择我们？',
      headers: {
        feature: '对比项',
        official: '官方订阅',
        us: '本平台'
      },
      items: {
        pricing: {
          feature: '付费方式',
          official: '固定月费，用不完也付',
          us: '按量付费，用多少付多少'
        },
        models: {
          feature: '模型选择',
          official: '单一服务商',
          us: '多模型随意切换'
        },
        management: {
          feature: '账号管理',
          official: '每个服务单独管理',
          us: '统一密钥，一站管理'
        },
        stability: {
          feature: '服务稳定性',
          official: '单账号易触发限制',
          us: '多账号池，自动切换'
        },
        control: {
          feature: '用量控制',
          official: '无法限制',
          us: '可设配额、查明细'
        }
      }
    },
    providers: {
      title: '已支持的 AI 模型',
      description: '一个 API，多种选择',
      supported: '已支持',
      soon: '即将推出',
      claude: 'Claude',
      gemini: 'Gemini',
      antigravity: 'Antigravity',
      more: '更多'
    },
    // CTA 区块
    cta: {
      title: '准备好开始了吗？',
      description: '注册即可获得免费试用额度，体验一站式 AI 服务',
      button: '免费注册'
    },
    footer: {
      allRightsReserved: '保留所有权利。'
    }
  },

  homeReplica: {
    primaryNavigation: '主导航',
    pricing: '定价',
    tutorial: '教程',
    about: '关于',
    notifications: '通知',
    toggleNavigation: '切换导航菜单',
    modelPricing: '模型定价',
    heroCopy: '从一个端点开始。继续使用现有 SDK，由 {brand} 负责路由、故障回退和可观测性。',
    quickstart: '快速开始',
    gatewayOnline: '网关在线',
    sdk: 'SDK 示例',
    readDocs: '阅读文档',
    getApiKey: '获取 API Key',
    everyModel: '所有模型，一个网关。',
    everyModelCopy: '为每项任务选择最合适的模型，无需重建现有技术栈。',
    everything: '覆盖应用与模型之间的一切。',
    averageRequest: '平均请求成本',
    flexiblePricing: '灵活，始于设计。',
    flexiblePricingCopy: '从按量计费开始，流量增长后选择订阅，或构建专属企业方案。',
    usageBilling: '按量计费',
    subscription: 'Omnio Pro',
    builtLike: '像基础设施一样构建。',
    operatedLike: '像产品一样运营。',
    intelligenceLayer: '为每一次 AI 请求提供快速、可观测的智能层。',
    protocolSurface: '统一协议接口',
    providers: '服务商',
    modelNetwork: '模型网络',
    firstRequest: '几分钟发出首个请求。',
    fullControl: '从第一天开始完全掌控。',
    exploreCopy: '检查每次路由请求的路径、延迟与最终结果。',
    faq: '常见问题',
    faqCopy: '开始通过 {brand} 承载生产级 AI 流量所需的一切。',
    createGateway: '创建你的网关',
    product: '产品',
    resources: '资源',
    userAgreement: '用户协议',
    builtOnSub2api: '基于 Sub2API 构建',
    headlineGateway: '通过一个网关路由所有 AI 模型',
    headlineVisible: '让每次请求都可观测',
    headlineControl: '控制成本、延迟和可靠性',
    capabilityReasoning: '前沿推理',
    capabilityContext: '超长上下文智能',
    capabilityMultimodal: '实时多模态',
    capabilityEfficient: '高效推理',
    capabilityMultilingual: '多语言规模化',
    modelNetworkCopy: '通过一套稳定的 API 契约访问广泛的模型网络。',
    requestVisible: '每次请求都可见',
    requestVisibleCopy: '在同一个运营视图中查看延迟、Token、成本与状态。',
    policySpeed: '策略实时执行',
    policySpeedCopy: '在请求到达服务商前应用访问、预算与路由策略。',
    failover: '自动故障切换',
    failoverCopy: '服务商、账号或模型不可用时，让流量持续运行。',
    costSignal: '实时成本信号',
    costSignalCopy: '按请求跟踪单位成本，在账单到来前完成优化。',
    protocolSurfaceCopy: '切换服务商、模型或路由策略时保留现有 SDK。',
    usageBased: '按量使用',
    usageBasedCopy: '只为成功的模型用量付费。',
    payAsYouGo: '随用随付',
    viewPricing: '查看模型定价',
    perModelPricing: '透明的模型费率',
    noCommitment: '无需月度承诺',
    usageAnalytics: '详细用量分析',
    automaticRouting: '自动服务商路由',
    subscriptionCopy: '为持续增长的流量提供可预测容量。',
    monthlyPlans: '月度方案',
    mostPopular: '最受欢迎',
    viewSubscriptions: '查看订阅',
    includedQuota: '包含月度额度',
    subscriptionBilling: '可预测的月度账单',
    teamVisibility: '团队共享可见性',
    upgradeTraffic: '随流量增长升级',
    enterprise: '企业版',
    enterpriseCopy: '为关键业务定制专属网关。',
    custom: '定制',
    tailored: '专属',
    contactUs: '联系我们',
    customQuota: '定制额度与条款',
    rolesPermissions: '基于角色的权限',
    onboardingSupport: '接入支持',
    customRouting: '定制路由策略',
    whatIs: '什么是 {brand}？',
    whatIsAnswer: '{brand} 是统一 AI 网关，通过一个 API 将应用连接到多个服务商，并内置路由、用量与成本控制。',
    existingSdk: '可以继续使用现有 SDK 吗？',
    existingSdkAnswer: '可以。将 OpenAI 兼容客户端指向 {brand} 的基础地址，其余集成方式保持熟悉。',
    routingQuestion: '自动路由如何工作？',
    routingAnswer: 'Sub2API 会评估可用渠道、账号健康状态与已配置策略，再为每次请求选择兼容路由。',
    trackQuestion: '可以跟踪用量与成本吗？',
    trackAnswer: '可以。控制台会根据用户和管理员权限提供请求级用量、延迟、状态与计费数据。'
  },

  // Key Usage Query Page
  keyUsage: {
    showKey: '显示 API 密钥',
    hideKey: '隐藏 API 密钥',
    title: 'API Key 用量查询',
    subtitle: '输入您的 API Key 以查看实时消费金额与使用状态',
    placeholder: 'sk-ant-mirror-xxxxxxxxxxxx',
    query: '查询',
    querying: '查询中...',
    privacyNote: '您的 Key 仅在浏览器本地处理，不会被存储',
    dateRange: '统计范围:',
    dateRangeToday: '今日',
    dateRange7d: '7 天',
    dateRange30d: '30 天',
    dateRange90d: '90 天',
    dateRangeCustom: '自定义',
    apply: '应用',
    used: '已使用',
    detailInfo: '详细信息',
    tokenStats: 'Token 统计',
    dailyDetail: '按日明细',
    modelStats: '模型用量统计',
    // Table headers
    date: '日期',
    model: '模型',
    requests: '请求数',
    inputTokens: '输入 Tokens',
    outputTokens: '输出 Tokens',
    cacheCreationTokens: '缓存创建',
    cacheReadTokens: '缓存读取',
    cacheWriteTokens: '缓存写入',
    totalTokens: '总 Tokens',
    cost: '费用',
    // Status
    quotaMode: 'Key 限额模式',
    walletBalance: '钱包余额',
    // Ring card titles
    totalQuota: '总额度',
    limit5h: '5 小时限额',
    limitDaily: '日限额',
    limit7d: '7 天限额',
    limitWeekly: '周限额',
    limitMonthly: '月限额',
    // Detail rows
    remainingQuota: '剩余额度',
    expiresAt: '过期时间',
    todayExpires: '(今日到期)',
    daysLeft: '({days} 天)',
    usedQuota: '已用额度',
    resetNow: '即将重置',
    subscriptionType: '订阅类型',
    subscriptionExpires: '订阅到期',
    // Usage stat cells
    todayRequests: '今日请求',
    todayInputTokens: '今日输入',
    todayOutputTokens: '今日输出',
    todayTokens: '今日 Tokens',
    todayCacheCreation: '今日缓存创建',
    todayCacheRead: '今日缓存读取',
    todayCost: '今日费用',
    rpmTpm: 'RPM / TPM',
    totalRequests: '累计请求',
    totalInputTokens: '累计输入',
    totalOutputTokens: '累计输出',
    totalTokensLabel: '累计 Tokens',
    totalCacheCreation: '累计缓存创建',
    totalCacheRead: '累计缓存读取',
    totalCost: '累计费用',
    avgDuration: '平均耗时',
    // Messages
    enterApiKey: '请输入 API Key',
    querySuccess: '查询成功',
    queryFailed: '查询失败',
    queryFailedRetry: '查询失败，请稍后重试',
    noDailyUsage: '暂无按日用量数据',
  },

  // Setup Wizard
  setup: {
    title: 'Sub2API 安装向导',
    initializeTitle: '初始化 {siteName}',
    description: '配置您的 Sub2API 实例',
    wizardTitle: '系统安装向导',
    wizardDescription: '完成以下步骤以完成首次安装。',
    progressLabel: '安装进度',
    logoAlt: '系统 Logo',
    database: {
      title: '数据库配置',
      description: '连接到您的 PostgreSQL 数据库',
      stepDescription: '验证数据库连接',
      host: '主机',
      port: '端口',
      username: '用户名',
      password: '密码',
      databaseName: '数据库名称',
      sslMode: 'SSL 模式',
      passwordPlaceholder: '密码',
      ssl: {
        disable: '禁用',
        require: '要求',
        verifyCa: '验证 CA',
        verifyFull: '完全验证'
      }
    },
    redis: {
      title: 'Redis 配置',
      description: '连接到您的 Redis 服务器',
      stepDescription: '验证缓存服务连接',
      host: '主机',
      port: '端口',
      password: '密码（可选）',
      database: '数据库',
      passwordPlaceholder: '密码',
      enableTls: '启用 TLS',
      enableTlsHint: '连接 Redis 时使用 TLS（公共 CA 证书）'
    },
    admin: {
      title: '管理员账户',
      description: '创建您的管理员账户',
      stepDescription: '创建根管理员凭据',
      email: '邮箱',
      password: '密码',
      confirmPassword: '确认密码',
      passwordPlaceholder: '至少 8 个字符',
      confirmPasswordPlaceholder: '确认密码',
      passwordMismatch: '密码不匹配'
    },
    ready: {
      title: '准备安装',
      description: '检查您的配置并完成安装',
      stepDescription: '确认设置并初始化',
      database: '数据库',
      redis: 'Redis',
      adminEmail: '管理员邮箱'
    },
    status: {
      testing: '测试中...',
      success: '连接成功',
      testConnection: '测试连接',
      installing: '安装中...',
      completeInstallation: '完成安装',
      completed: '安装完成！',
      redirecting: '正在跳转到登录页面...',
      restarting: '服务正在重启，请稍候...',
      timeout: '服务重启时间超出预期，请手动刷新页面。'
    }
  },

  // Common
}
