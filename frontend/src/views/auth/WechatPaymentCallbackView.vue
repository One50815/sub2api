<template>
  <AuthLayout>
    <div class="wechat-callback-state">
      <div
        class="wechat-callback-icon"
        :class="{ 'is-error': Boolean(errorMessage) }"
      >
        <Icon :name="errorMessage ? 'exclamationCircle' : 'creditCard'" size="xl" :stroke-width="1.6" />
      </div>

      <div class="wechat-callback-heading">
        <h2>{{ callbackTitleText }}</h2>
        <p>{{ errorMessage || callbackProcessingText }}</p>
      </div>

      <template v-if="!errorMessage">
        <div class="wechat-callback-progress">
          <Icon name="refresh" size="sm" class="animate-spin" :stroke-width="2" />
          <span>{{ callbackProcessingText }}</span>
        </div>
        <p class="wechat-callback-note">{{ t('auth.wechatPayment.callbackHint') }}</p>
      </template>

      <div v-else class="wechat-callback-error" role="alert">
        <p>{{ errorMessage }}</p>
        <button class="btn btn-primary" type="button" @click="goBackToPayment">
          {{ backToPaymentText }}
        </button>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores'
import AuthLayout from '@/components/layout/AuthLayout.vue'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const appStore = useAppStore()

const errorMessage = ref('')

watch(errorMessage, (message) => {
  if (message) {
    appStore.showError(message)
  }
})

const callbackProcessingText = computed(() => t('auth.wechatPayment.callbackProcessing'))
const callbackTitleText = computed(() => t('auth.wechatPayment.callbackTitle'))
const backToPaymentText = computed(() => t('auth.wechatPayment.backToPayment'))

function readQueryString(key: string): string {
  const value = route.query[key]
  if (Array.isArray(value)) {
    return typeof value[0] === 'string' ? value[0] : ''
  }
  return typeof value === 'string' ? value : ''
}

function parseFragmentParams(): URLSearchParams {
  const raw = typeof window !== 'undefined' ? window.location.hash : ''
  const hash = raw.startsWith('#') ? raw.slice(1) : raw
  return new URLSearchParams(hash)
}

function normalizeRedirectPath(path: string | null | undefined): string {
  const value = (path || '').trim()
  if (!value) return '/purchase'
  if (!value.startsWith('/')) return '/purchase'
  if (value.startsWith('//') || value.includes('://')) return '/purchase'
  if (value === '/payment') return '/purchase'
  if (value.startsWith('/payment?')) return '/purchase' + value.slice('/payment'.length)
  return value
}

function appendQueryParam(query: Record<string, string>, key: string, value: string) {
  if (value) {
    query[key] = value
  }
}

function goBackToPayment() {
  void router.replace('/purchase')
}

onMounted(async () => {
  const fragment = parseFragmentParams()
  const readParam = (key: string) => fragment.get(key) || readQueryString(key)

  const error = readParam('error') || readParam('err_msg') || readParam('errmsg')
  const errorDescription = readParam('error_description') || readParam('message')

  if (error) {
    errorMessage.value = errorDescription || error
    return
  }

  const resumeToken = readParam('wechat_resume_token')
  const openid = readParam('openid')
  const state = readParam('state')
  const scope = readParam('scope')
  const paymentType = readParam('payment_type')
  const amount = readParam('amount')
  const orderType = readParam('order_type')
  const planId = readParam('plan_id')
  const redirectURL = new URL(
    normalizeRedirectPath(readParam('redirect')),
    window.location.origin,
  )

  if (!resumeToken && !openid) {
    errorMessage.value = t('auth.wechatPayment.callbackMissingResumeToken')
    return
  }

  const query: Record<string, string> = {
    ...Object.fromEntries(redirectURL.searchParams.entries()),
    wechat_resume: '1',
  }

  if (resumeToken) {
    query.wechat_resume_token = resumeToken
  } else {
    query.openid = openid
    appendQueryParam(query, 'state', state)
    appendQueryParam(query, 'scope', scope)
    appendQueryParam(query, 'payment_type', paymentType)
    appendQueryParam(query, 'amount', amount)
    appendQueryParam(query, 'order_type', orderType)
    appendQueryParam(query, 'plan_id', planId)
  }

  await router.replace({
    path: redirectURL.pathname,
    query,
  })
})
</script>

<style scoped>
.wechat-callback-state {
  display: grid;
  justify-items: center;
  gap: 1.5rem;
  width: 100%;
  padding: 0.25rem 0;
  text-align: center;
}

.wechat-callback-icon {
  display: grid;
  width: 4rem;
  height: 4rem;
  place-items: center;
  border-radius: 1rem;
  color: var(--omnio-foreground);
  background: color-mix(in srgb, var(--omnio-foreground) 6%, transparent);
}

.wechat-callback-icon.is-error {
  color: #ef4444;
  background: rgb(239 68 68 / 10%);
}

.wechat-callback-heading {
  display: grid;
  gap: 0.5rem;
}

.wechat-callback-heading h2 {
  margin: 0;
  color: var(--omnio-foreground);
  font-size: 1.5rem;
  font-weight: 650;
  line-height: 1.25;
  letter-spacing: -0.025em;
}

.wechat-callback-heading p,
.wechat-callback-note {
  margin: 0;
  color: var(--omnio-muted);
  font-size: 0.875rem;
  line-height: 1.55;
}

.wechat-callback-progress {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--omnio-foreground);
  font-size: 0.875rem;
  font-weight: 580;
}

.wechat-callback-note {
  max-width: 23rem;
  font-size: 0.78rem;
}

.wechat-callback-error {
  display: grid;
  justify-items: center;
  gap: 1rem;
  width: 100%;
  padding: 0.875rem;
  border: 1px solid rgb(239 68 68 / 24%);
  border-radius: 0.75rem;
  color: #ef4444;
  background: rgb(239 68 68 / 8%);
}

.wechat-callback-error p {
  margin: 0;
  font-size: 0.82rem;
  line-height: 1.5;
}
</style>
