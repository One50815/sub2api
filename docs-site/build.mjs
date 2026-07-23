import { cp, mkdir, readFile, rm, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { marked } from '../frontend/node_modules/marked/lib/marked.esm.js'

const root = path.dirname(fileURLToPath(import.meta.url))
const dist = path.join(root, 'dist')

const documents = [
  {
    slug: 'about-omnio',
    title: '关于 Omnio',
    description: '了解 Omnio 是什么、为谁服务、解决哪些问题，以及我们重视的服务体验。',
    file: 'Omnio关于我们.md',
    category: 'brand',
    meta: '了解 Omnio',
    status: '持续更新',
    featured: true,
  },
  {
    slug: 'brand-story',
    title: '品牌故事',
    description: '从词元101到 Omnio，了解名字背后的含义、坚持与方向。',
    file: 'Omnio品牌故事.md',
    category: 'brand',
    meta: 'Omnio 品牌',
    status: '持续书写',
    featured: true,
  },
  {
    slug: 'usage-guide',
    title: '完整使用教程',
    description: '从注册、创建 API Key 到 SDK 调用、用量查询、错误排查和工单支持。',
    file: 'Omnio使用教程.md',
    category: 'guide',
    meta: 'Omnio 使用指南',
    status: '随产品更新',
    featured: true,
  },
  {
    slug: 'pricing',
    title: '模型定价',
    description: '按模型厂商分组查看官方上游定价，并依据 Omnio 分组倍率计算实际价格。',
    file: 'Omnio定价说明.md',
    category: 'guide',
    meta: 'Omnio 定价说明',
    status: '随官方价格更新',
    featured: true,
  },
  {
    slug: 'user-agreement',
    title: '用户服务协议',
    description: '了解使用 Omnio 服务时双方的权利、责任与使用规范。',
    file: 'Omnio用户服务协议.md',
    category: 'legal',
    meta: 'Omnio 法律与政策',
    status: '持续生效',
  },
  {
    slug: 'recharge-policy',
    title: '数字充值服务政策',
    description: '了解按量付费充值、余额、优惠、退款与争议处理规则。',
    file: 'Omnio数字充值服务政策.md',
    category: 'legal',
    meta: 'Omnio 法律与政策',
    status: '持续生效',
  },
]

function escapeHtml(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
}

function stripTags(value) {
  return String(value).replace(/<[^>]*>/g, '').replaceAll('&amp;', '&').trim()
}

function addHeadingAnchors(html) {
  const used = new Map()
  const toc = []
  let fallback = 0

  const body = html.replace(/<h([23])>([\s\S]*?)<\/h\1>/g, (_, depth, inner) => {
    fallback += 1
    const label = stripTags(inner)
    const base = label
      .toLowerCase()
      .replace(/[^\p{Letter}\p{Number}\s-]/gu, '')
      .trim()
      .replace(/\s+/g, '-') || `section-${fallback}`
    const count = used.get(base) || 0
    used.set(base, count + 1)
    const id = count === 0 ? base : `${base}-${count + 1}`
    toc.push({ depth: Number(depth), id, label })
    return `<h${depth} id="${escapeHtml(id)}"><a class="heading-anchor" href="#${escapeHtml(id)}" aria-label="链接到本节">${inner}</a></h${depth}>`
  })

  return { body, toc }
}

function shell({ title, description, activeSlug = '', content, toc = '' }) {
  const renderNavigation = (category) => documents
    .filter((document) => document.category === category)
    .map((document) => {
      const current = document.slug === activeSlug
      return `<a class="doc-nav-link${current ? ' is-active' : ''}" href="/docs/${document.slug}/"${current ? ' aria-current="page"' : ''}>${escapeHtml(document.title)}</a>`
    })
    .join('')
  const brandNavigation = renderNavigation('brand')
  const guideNavigation = renderNavigation('guide')
  const legalNavigation = renderNavigation('legal')

  return `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="description" content="${escapeHtml(description)}">
  <meta name="theme-color" content="#746df5">
  <title>${escapeHtml(title)} · Omnio 文档</title>
  <link rel="icon" href="/docs/assets/omnio-mark.svg" type="image/svg+xml">
  <link rel="stylesheet" href="/docs/assets/styles.css">
</head>
<body id="top">
  <header class="site-header">
    <a class="brand" href="/docs/" aria-label="Omnio 文档首页">
      <img src="/docs/assets/omnio-mark.svg" alt="">
      <span>Omnio</span><span class="brand-divider"></span><span class="brand-section">文档中心</span>
    </a>
    <a class="back-link" href="/">返回主站 <span aria-hidden="true">↗</span></a>
  </header>
  <div class="page-shell">
    <aside class="sidebar" aria-label="文档导航">
      <div class="sidebar-groups">
        <section class="sidebar-group">
          <p class="sidebar-label">了解 Omnio</p>
          <nav>${brandNavigation}</nav>
        </section>
        <section class="sidebar-group">
          <p class="sidebar-label">使用教程</p>
          <nav>${guideNavigation}</nav>
        </section>
        <section class="sidebar-group">
          <p class="sidebar-label">法律与政策</p>
          <nav>${legalNavigation}</nav>
        </section>
      </div>
      <div class="support-card">
        <span>需要帮助？</span>
        <a href="mailto:3290800970@qq.com">联系客服</a>
      </div>
    </aside>
    <main class="main-content">${content}</main>
    ${toc}
  </div>
  <footer class="site-footer">
    <span>© 2026 Omnio</span>
    <a href="mailto:3290800970@qq.com">3290800970@qq.com</a>
  </footer>
</body>
</html>`
}

function renderToc(items) {
  if (items.length === 0) return ''
  const links = items
    .map((item) => `<a class="toc-link depth-${item.depth}" href="#${escapeHtml(item.id)}">${escapeHtml(item.label)}</a>`)
    .join('')
  return `<aside class="toc" aria-label="本文目录"><p>本文目录</p><nav>${links}</nav></aside>`
}

function renderIndex() {
  const sections = [
    { category: 'brand', title: '了解品牌', description: '先了解 Omnio 的产品定位，再阅读名字与起点背后的故事。' },
    { category: 'guide', title: '使用教程', description: '从第一次登录到正式接入，一页完成全部配置。' },
    { category: 'legal', title: '法律与政策', description: '了解使用服务时的权利、责任、充值与退款规则。' },
  ]

  const documentSections = sections.map((section) => {
    const categoryDocuments = documents.filter((document) => document.category === section.category)
    const sectionCards = categoryDocuments.map((document) => {
      const index = documents.indexOf(document)
      return `<a class="document-card${document.featured ? ' is-featured' : ''}" href="/docs/${document.slug}/">
        <span class="card-index">${String(index + 1).padStart(2, '0')}</span>
        <span class="card-title">${escapeHtml(document.title)}</span>
        <span class="card-description">${escapeHtml(document.description)}</span>
        <span class="card-action">阅读文档 <span aria-hidden="true">→</span></span>
      </a>`
    }).join('')

    return `<section class="home-document-section">
      <div class="home-section-heading">
        <h2>${escapeHtml(section.title)}</h2>
        <p>${escapeHtml(section.description)}</p>
      </div>
      <div class="document-grid${categoryDocuments.length === 1 ? ' is-single' : ''}">${sectionCards}</div>
    </section>`
  }).join('')

  const content = `<section class="docs-hero">
      <span class="eyebrow">Omnio Documentation</span>
      <h1>了解 Omnio，也清晰了解服务规则</h1>
      <p>先认识品牌，再完成接入；需要时随时查阅法律与政策。</p>
    </section>
    <div class="home-document-sections" aria-label="全部文档">${documentSections}</div>
    <section class="contact-strip"><div><strong>没有找到需要的信息？</strong><span>发送邮件，我们会尽快协助。</span></div><a href="mailto:3290800970@qq.com">联系 Omnio</a></section>`

  return shell({
    title: '文档中心',
    description: 'Omnio 用户协议、充值政策与服务文档。',
    content,
  })
}

await rm(dist, { recursive: true, force: true })
await mkdir(path.join(dist, 'assets'), { recursive: true })
await cp(path.join(root, 'src', 'styles.css'), path.join(dist, 'assets', 'styles.css'))
await cp(path.join(root, 'src', 'omnio-mark.svg'), path.join(dist, 'assets', 'omnio-mark.svg'))
await writeFile(path.join(dist, 'index.html'), renderIndex(), 'utf8')

for (const document of documents) {
  const source = await readFile(path.join(root, 'content', document.file), 'utf8')
  const rendered = marked.parse(source, { gfm: true })
  const { body, toc } = addHeadingAnchors(rendered)
  const isBrandStory = document.slug === 'brand-story'
  const isAboutOmnio = document.slug === 'about-omnio'
  const content = `<article class="legal-document${isBrandStory ? ' brand-story' : ''}${isAboutOmnio ? ' about-omnio' : ''}">
    <div class="document-breadcrumb"><a href="/docs/">文档中心</a><span>/</span><span>${escapeHtml(document.title)}</span></div>
    <div class="document-meta"><span>${escapeHtml(document.meta)}</span><span>${escapeHtml(document.status)}</span></div>
    <div class="markdown-body">${body}</div>
    <div class="document-end"><span>${isBrandStory ? '故事未完，持续书写' : isAboutOmnio ? '认识 Omnio，从这里开始' : '文档结束'}</span><a href="#top">返回顶部 ↑</a></div>
  </article>`
  const target = path.join(dist, document.slug)
  await mkdir(target, { recursive: true })
  await writeFile(path.join(target, 'index.html'), shell({
    title: document.title,
    description: document.description,
    activeSlug: document.slug,
    content,
    toc: renderToc(toc),
  }), 'utf8')
}

console.log(`Built ${documents.length + 1} pages in ${dist}`)
