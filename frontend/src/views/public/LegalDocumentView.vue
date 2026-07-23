<template>
  <div class="legal-page">
    <header class="legal-header">
      <nav class="legal-nav">
        <RouterLink to="/home" class="legal-brand">
          <template v-if="settings">
            <span class="legal-brand-logo">
              <img :src="siteLogo || '/assets/brand/omnio-mark.svg?v=3'" alt="" />
            </span>
            <span>{{ siteName }}</span>
          </template>
          <template v-else>
            <span class="legal-skeleton legal-brand-logo-skeleton" aria-hidden="true"></span>
            <span class="legal-skeleton legal-brand-name-skeleton" aria-hidden="true"></span>
          </template>
        </RouterLink>

        <div class="legal-header-actions">
          <LocaleSwitcher />
          <button
            type="button"
            class="legal-tool-button"
            :aria-label="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon :name="isDark ? 'sun' : 'moon'" size="sm" />
          </button>
          <RouterLink to="/login" class="btn btn-primary legal-sign-in">
            {{ t('home.login') }}
          </RouterLink>
        </div>
      </nav>
    </header>

    <main class="legal-main">
      <div v-if="loading" class="legal-loading" aria-live="polite">
        <span class="legal-skeleton legal-loading-title"></span>
        <span class="legal-skeleton"></span>
        <span class="legal-skeleton"></span>
        <span class="legal-skeleton legal-loading-short"></span>
      </div>

      <section v-else-if="loadError" class="legal-state-card legal-state-error" role="alert">
        <span class="legal-state-icon"><Icon name="exclamationCircle" size="md" /></span>
        <div>
          <h1>{{ t('legal.loadFailed') }}</h1>
          <p>{{ t('legal.retryLater') }}</p>
        </div>
      </section>

      <section v-else-if="!currentDocument" class="legal-state-card">
        <span class="legal-state-icon"><Icon name="document" size="md" /></span>
        <div>
          <h1>{{ t('legal.notFound') }}</h1>
          <p>{{ t('legal.notFoundDescription') }}</p>
        </div>
      </section>

      <article v-else class="legal-article">
        <header class="legal-document-header">
          <p class="legal-document-type">{{ documentTypeLabel }}</p>
          <h1>{{ currentDocument.title }}</h1>
          <p v-if="updatedAt" class="legal-updated-at">
            {{ t('legal.updatedAt', { date: updatedAt }) }}
          </p>
        </header>

        <div v-if="hasContent" class="legal-document-content" v-html="renderedHtml"></div>
        <div v-else class="legal-empty">
          <Icon name="document" size="lg" />
          <span>{{ t('legal.empty') }}</span>
        </div>
      </article>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import { getLocale } from '@/i18n'
import { sanitizeUrl } from '@/utils/url'
import { useAppStore } from '@/stores/app'
import type { LoginAgreementDocument } from '@/types'
import zhAdminCompliance from '../../../../docs/legal/admin-compliance.zh.md?raw'
import enAdminCompliance from '../../../../docs/legal/admin-compliance.en.md?raw'

const route = useRoute()
const { t } = useI18n()
const appStore = useAppStore()
const settings = computed(() => appStore.cachedPublicSettings)
const loading = ref(!settings.value)
const loadError = ref(false)
const isDark = ref(document.documentElement.classList.contains('dark'))

marked.setOptions({
  breaks: true,
  gfm: true
})

const documentId = computed(() => String(route.params.documentId || ''))
const isAdminComplianceDocument = computed(() => documentId.value === 'admin-compliance')
const documents = computed(() => settings.value?.login_agreement_documents ?? [])
const siteName = computed(() => settings.value?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() =>
  sanitizeUrl(settings.value?.site_logo || appStore.siteLogo || '', {
    allowRelative: true,
    allowDataUrl: true
  })
)
const updatedAt = computed(() =>
  isAdminComplianceDocument.value ? '' : settings.value?.login_agreement_updated_at || ''
)
const documentTypeLabel = computed(() =>
  isAdminComplianceDocument.value ? t('legal.adminCompliance') : t('legal.loginAgreement')
)

const currentDocument = computed<LoginAgreementDocument | null>(() => {
  if (isAdminComplianceDocument.value) {
    return {
      id: 'admin-compliance',
      title: t('adminCompliance.title'),
      content_md: getLocale() === 'zh' ? zhAdminCompliance : enAdminCompliance
    }
  }
  const id = documentId.value
  if (!id) return null
  return documents.value.find((document) => document.id === id) ?? null
})

const hasContent = computed(() => Boolean(currentDocument.value?.content_md?.trim()))

const renderedHtml = computed(() => {
  const content = currentDocument.value?.content_md?.trim() || ''
  if (!content) return ''
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
})

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(async () => {
  loadError.value = false
  const loadedSettings = await appStore.fetchPublicSettings()
  if (!loadedSettings) loadError.value = true
  loading.value = false
})
</script>

<style scoped>
.legal-page {
  min-height: 100svh;
  color: var(--omnio-foreground);
  background: var(--omnio-bg);
}

.legal-header {
  position: fixed;
  inset: 0 0 auto;
  z-index: 50;
  height: 4rem;
  pointer-events: none;
}

.legal-nav {
  display: flex;
  width: min(100%, 80rem);
  height: 4rem;
  margin: 0 auto;
  padding: 0 1.5rem;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  pointer-events: auto;
  border-bottom: 1px solid color-mix(in srgb, var(--omnio-border) 68%, transparent);
  background: color-mix(in srgb, var(--omnio-bg) 88%, transparent);
  backdrop-filter: blur(18px) saturate(1.2);
}

.legal-brand {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 0.625rem;
  color: var(--omnio-foreground);
  font-size: 0.875rem;
  font-weight: 630;
  text-decoration: none;
}

.legal-brand-logo {
  display: grid;
  width: 1.75rem;
  height: 1.75rem;
  flex: 0 0 1.75rem;
  place-items: center;
  overflow: hidden;
  border-radius: 0.45rem;
}

.legal-brand-logo img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.legal-header-actions {
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  gap: 0.35rem;
}

.legal-tool-button {
  display: grid;
  width: 2rem;
  height: 2rem;
  place-items: center;
  border: 0;
  border-radius: 0.5rem;
  color: var(--omnio-muted);
  background: transparent;
  transition: color 140ms ease, background 140ms ease;
}

.legal-tool-button:hover {
  color: var(--omnio-foreground);
  background: var(--omnio-surface-subtle);
}

.legal-sign-in {
  min-height: 2rem;
  margin-left: 0.35rem;
  padding: 0.35rem 0.875rem;
  border-radius: 0.5rem;
  font-size: 0.75rem;
}

.legal-main {
  width: min(100%, 56rem);
  margin: 0 auto;
  padding: 7rem 1.5rem 4rem;
}

.legal-loading {
  display: grid;
  gap: 1rem;
  padding-top: 1rem;
}

.legal-skeleton {
  display: block;
  width: 100%;
  height: 0.875rem;
  border-radius: 0.4rem;
  background: var(--omnio-surface-subtle);
  animation: legal-pulse 1.8s ease-in-out infinite;
}

.legal-brand-logo-skeleton {
  width: 1.75rem;
  height: 1.75rem;
}

.legal-brand-name-skeleton {
  width: 5rem;
  height: 0.85rem;
}

.legal-loading-title {
  width: 45%;
  height: 2rem;
  margin-bottom: 0.5rem;
}

.legal-loading-short {
  width: 78%;
}

@keyframes legal-pulse {
  50% { opacity: 0.48; }
}

.legal-state-card {
  display: flex;
  max-width: 42rem;
  align-items: flex-start;
  gap: 1rem;
  margin: 2rem auto 0;
  padding: 1.5rem;
  border: 1px dashed var(--omnio-border-strong);
  border-radius: 0.75rem;
  background: var(--omnio-surface);
}

.legal-state-card h1,
.legal-state-card p {
  margin: 0;
}

.legal-state-card h1 {
  font-size: 1rem;
  font-weight: 630;
}

.legal-state-card p {
  margin-top: 0.35rem;
  color: var(--omnio-muted);
  font-size: 0.85rem;
  line-height: 1.55;
}

.legal-state-icon {
  display: grid;
  width: 2.5rem;
  height: 2.5rem;
  flex: 0 0 2.5rem;
  place-items: center;
  border-radius: 0.6rem;
  color: var(--omnio-muted);
  background: var(--omnio-surface-subtle);
}

.legal-state-error {
  border-style: solid;
  border-color: rgb(239 68 68 / 25%);
}

.legal-state-error .legal-state-icon {
  color: #ef4444;
  background: rgb(239 68 68 / 8%);
}

.legal-document-header {
  padding-bottom: 2rem;
  border-bottom: 1px solid var(--omnio-border);
}

.legal-document-type {
  margin: 0;
  color: var(--omnio-primary);
  font-size: 0.78rem;
  font-weight: 620;
}

.legal-document-header h1 {
  margin: 0.65rem 0 0;
  color: var(--omnio-foreground);
  font-size: clamp(2rem, 5vw, 2.5rem);
  font-weight: 680;
  line-height: 1.15;
  letter-spacing: -0.035em;
  overflow-wrap: anywhere;
}

.legal-updated-at {
  margin: 0.75rem 0 0;
  color: var(--omnio-muted);
  font-size: 0.78rem;
}

.legal-document-content {
  padding-top: 1.5rem;
  color: var(--omnio-foreground);
  font-size: 0.925rem;
  line-height: 1.8;
  overflow-wrap: anywhere;
}

.legal-document-content :deep(h1),
.legal-document-content :deep(h2),
.legal-document-content :deep(h3),
.legal-document-content :deep(h4) {
  color: var(--omnio-foreground);
  font-weight: 650;
  letter-spacing: -0.02em;
}

.legal-document-content :deep(h1) {
  margin: 2.25rem 0 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--omnio-border);
  font-size: 1.75rem;
}

.legal-document-content :deep(h2) { margin: 2rem 0 0.75rem; font-size: 1.4rem; }
.legal-document-content :deep(h3) { margin: 1.75rem 0 0.6rem; font-size: 1.15rem; }
.legal-document-content :deep(h4) { margin: 1.5rem 0 0.5rem; font-size: 1rem; }

.legal-document-content :deep(p),
.legal-document-content :deep(ul),
.legal-document-content :deep(ol) {
  margin: 0 0 1rem;
  color: color-mix(in srgb, var(--omnio-foreground) 80%, transparent);
}

.legal-document-content :deep(ul),
.legal-document-content :deep(ol) { padding-left: 1.5rem; }
.legal-document-content :deep(li) { margin-bottom: 0.3rem; }

.legal-document-content :deep(a) {
  color: var(--omnio-primary);
  text-decoration: underline;
  text-underline-offset: 0.2em;
}

.legal-document-content :deep(blockquote) {
  margin: 1.25rem 0;
  padding-left: 1rem;
  border-left: 3px solid var(--omnio-border-strong);
  color: var(--omnio-muted);
}

.legal-document-content :deep(code) {
  padding: 0.15rem 0.35rem;
  border-radius: 0.3rem;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 0.84em;
  background: var(--omnio-surface-subtle);
}

.legal-document-content :deep(pre) {
  margin: 1.25rem 0;
  padding: 1rem;
  overflow-x: auto;
  border: 1px solid var(--omnio-border);
  border-radius: 0.65rem;
  color: var(--omnio-foreground);
  background: var(--omnio-surface-subtle);
}

.legal-document-content :deep(pre code) { padding: 0; background: transparent; }

.legal-document-content :deep(table) {
  display: block;
  width: 100%;
  margin: 1.25rem 0;
  overflow-x: auto;
  border-collapse: collapse;
}

.legal-document-content :deep(th),
.legal-document-content :deep(td) {
  padding: 0.6rem 0.75rem;
  border: 1px solid var(--omnio-border);
  text-align: left;
}

.legal-document-content :deep(th) { font-weight: 620; background: var(--omnio-surface-subtle); }
.legal-document-content :deep(img) { max-width: 100%; height: auto; margin: 1.25rem 0; border-radius: 0.65rem; }
.legal-document-content :deep(hr) { margin: 1.75rem 0; border: 0; border-top: 1px solid var(--omnio-border); }

.legal-empty {
  display: grid;
  min-height: 12rem;
  margin-top: 1.5rem;
  place-items: center;
  align-content: center;
  gap: 0.75rem;
  border: 1px dashed var(--omnio-border-strong);
  border-radius: 0.75rem;
  color: var(--omnio-muted);
  background: var(--omnio-surface);
  font-size: 0.85rem;
}

@media (max-width: 640px) {
  .legal-nav { padding: 0 1rem; }
  .legal-main { padding: 6rem 1rem 3rem; }
  .legal-brand > span:last-child { max-width: 7rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .legal-document-header { padding-bottom: 1.5rem; }
}
</style>
