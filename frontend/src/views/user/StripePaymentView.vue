<template>
  <component :is="isPopup ? 'div' : AppLayout" :class="isPopup ? 'min-h-screen bg-gray-50 dark:bg-dark-950' : ''">
    <div class="mx-auto w-full max-w-xl space-y-4 py-8" :class="isPopup ? 'px-4' : ''">
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
      </div>
      <div v-else-if="initError" class="rounded-xl border border-gray-200 bg-white p-8 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900">
        <div class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-red-50 ring-1 ring-red-100 dark:bg-red-500/10 dark:ring-red-500/20">
          <Icon name="exclamationCircle" size="xl" class="text-red-500" />
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('payment.stripeLoadFailed') }}</h3>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ initError }}</p>
        <button class="btn btn-primary mt-6" @click="router.push('/purchase')">{{ t('payment.result.backToRecharge') }}</button>
      </div>
      <template v-else>
        <!-- 金额头部 -->
        <div v-if="order" class="rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between gap-5 px-5 py-5 sm:px-6">
            <div class="flex items-center gap-3">
              <span class="flex h-10 w-10 items-center justify-center rounded-lg bg-primary-50 text-primary-600 ring-1 ring-primary-100 dark:bg-primary-500/10 dark:text-primary-400 dark:ring-primary-500/20">
                <Icon name="creditCard" size="md" />
              </span>
              <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('payment.actualPay') }}</p>
            </div>
            <p class="font-mono text-2xl font-bold tabular-nums text-gray-950 dark:text-white">{{ formatGatewayAmount(order.pay_amount) }}</p>
          </div>
        </div>

        <!-- 微信二维码展示 -->
        <template v-if="wechatQrUrl">
          <div class="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <div class="flex flex-col items-center space-y-4 px-5 py-6 sm:px-6">
              <p class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('payment.qr.scanWxpay') }}</p>
              <div class="relative rounded-xl border border-gray-200 bg-white p-3 ring-4 ring-gray-50 dark:border-dark-600 dark:ring-dark-800">
                <img :src="wechatQrUrl" alt="WeChat Pay QR" class="h-56 w-56 rounded" />
                <div class="pointer-events-none absolute inset-0 flex items-center justify-center">
                  <span class="rounded-full bg-white p-2 shadow-sm ring-1 ring-gray-200">
                    <img :src="wxpayIcon" alt="" class="h-6 w-6" />
                  </span>
                </div>
              </div>
              <p class="text-center text-sm text-gray-500 dark:text-gray-400">{{ t('payment.qr.scanWxpayHint') }}</p>
            </div>
            <div class="border-t border-gray-200 bg-gray-50 px-5 py-3 text-center dark:border-dark-700 dark:bg-dark-800/60">
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('payment.qr.waitingPayment') }}</p>
            </div>
          </div>
        </template>

        <!-- 支付宝跳转状态 -->
        <template v-else-if="redirecting">
          <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <div class="flex flex-col items-center space-y-4 py-4">
              <div class="h-10 w-10 animate-spin rounded-full border-4 border-[#00AEEF] border-t-transparent"></div>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('payment.qr.payInNewWindowHint') }}</p>
            </div>
          </div>
        </template>

        <!-- 成功状态 -->
        <template v-else-if="stripeSuccess">
          <div class="rounded-xl border border-gray-200 bg-white p-6 text-center shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <div class="flex flex-col items-center gap-3 py-4">
              <div class="flex h-14 w-14 items-center justify-center rounded-full bg-emerald-50 ring-1 ring-emerald-100 dark:bg-emerald-500/10 dark:ring-emerald-500/20">
                <Icon name="checkCircle" size="xl" class="text-emerald-600 dark:text-emerald-400" />
              </div>
              <p class="text-lg font-bold text-gray-900 dark:text-white">{{ t('payment.result.success') }}</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('payment.stripeSuccessProcessing') }}</p>
            </div>
          </div>
        </template>

        <!-- 无指定方式或未知方式时展示完整 Payment Element -->
        <template v-else-if="showPaymentElement">
          <div class="rounded-xl border border-gray-200 bg-white p-6 shadow-sm dark:border-dark-700 dark:bg-dark-900">
            <div id="stripe-payment-element" class="min-h-[200px]"></div>
            <p v-if="stripeError" class="mt-4 text-sm text-red-600 dark:text-red-400">{{ stripeError }}</p>
            <button class="btn btn-stripe mt-6 w-full py-3 text-base" :disabled="stripeSubmitting || !stripeReady" @click="handleGenericPay">
              <span v-if="stripeSubmitting" class="flex items-center justify-center gap-2">
                <span class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></span>
                {{ t('common.processing') }}
              </span>
              <span v-else>{{ t('payment.stripePay') }}</span>
            </button>
          </div>
          <div class="text-center">
            <button class="btn btn-secondary" @click="router.push('/purchase')">{{ t('payment.result.backToRecharge') }}</button>
          </div>
        </template>

        <!-- 错误状态 -->
        <div v-if="stripeError && !showPaymentElement" class="rounded-xl border border-red-200 bg-white p-4 shadow-sm dark:border-red-900/60 dark:bg-dark-900">
          <p class="text-sm text-red-600 dark:text-red-400">{{ stripeError }}</p>
          <button class="btn btn-secondary mt-3 w-full" @click="router.push('/purchase')">{{ t('payment.result.backToRecharge') }}</button>
        </div>
      </template>
    </div>
  </component>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { usePaymentStore } from '@/stores/payment'
import { paymentAPI } from '@/api/payment'
import { extractI18nErrorMessage } from '@/utils/apiError'
import { isMobileDevice } from '@/utils/device'
import { formatPaymentAmount, normalizePaymentCurrency } from '@/components/payment/currency'
import { PAYMENT_RECOVERY_STORAGE_KEY, readPaymentRecoverySnapshot } from '@/components/payment/paymentFlow'
import type { PaymentOrder } from '@/types/payment'
import type { Stripe, StripeElements } from '@stripe/stripe-js'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import wxpayIcon from '@/assets/icons/wxpay.svg'

const i18n = useI18n()
const { t } = i18n
const route = useRoute()
const router = useRouter()
const paymentStore = usePaymentStore()

// 弹窗模式：指定支付宝或微信方式时跳过 AppLayout
const isPopup = computed(() => !!route.query.method)

const loading = ref(true)
const initError = ref('')
const stripeError = ref('')
const stripeSubmitting = ref(false)
const stripeSuccess = ref(false)
const stripeReady = ref(false)
const order = ref<PaymentOrder | null>(null)
const currency = ref('CNY')
const wechatQrUrl = ref('')
const redirecting = ref(false)
const showPaymentElement = ref(false)

let stripeInstance: Stripe | null = null
let elementsInstance: StripeElements | null = null
let redirectTimer: ReturnType<typeof setTimeout> | null = null

onMounted(async () => {
  const orderId = Number(route.query.order_id)
  const clientSecret = String(route.query.client_secret || '')
  const method = String(route.query.method || '')
  const resumeToken = typeof route.query.resume_token === 'string' ? route.query.resume_token : undefined

  if (!orderId || !clientSecret) {
    loading.value = false
    initError.value = t('payment.stripeMissingParams')
    return
  }

  try {
    if (typeof window !== 'undefined') {
      const restored = readPaymentRecoverySnapshot(
        window.localStorage.getItem(PAYMENT_RECOVERY_STORAGE_KEY),
        { resumeToken },
      )
      if (restored?.orderId === orderId) {
        currency.value = normalizePaymentCurrency(restored.currency)
      }
    }
    const res = await paymentAPI.getOrder(orderId)
    order.value = res.data
    if (res.data.currency) {
      currency.value = normalizePaymentCurrency(res.data.currency)
    }

    await paymentStore.fetchConfig()
    const publishableKey = paymentStore.config?.stripe_publishable_key
    if (!publishableKey) { initError.value = t('payment.stripeNotConfigured'); return }

    const { loadStripe } = await import('@stripe/stripe-js/pure')
    const stripe = await loadStripe(publishableKey)
    if (!stripe) { initError.value = t('payment.stripeLoadFailed'); return }

    stripeInstance = stripe
    loading.value = false

    // 指定方式直接确认，无需渲染完整 Payment Element
    if (method === 'alipay') {
      await confirmAlipay(stripe, clientSecret, orderId)
    } else if (method === 'wechat_pay') {
      await confirmWechatPay(stripe, clientSecret)
    } else {
      // 未指定方式时渲染完整 Payment Element
      showPaymentElement.value = true
      await nextTick()
      mountPaymentElement(stripe, clientSecret)
    }
  } catch (err: unknown) {
    initError.value = extractI18nErrorMessage(err, t, 'payment.errors', t('payment.stripeLoadFailed'))
  } finally {
    loading.value = false
  }
})

const localeCode = computed(() => {
  const raw = i18n.locale as unknown
  if (typeof raw === 'string') return raw
  if (raw && typeof raw === 'object' && 'value' in raw) {
    return String((raw as { value?: string }).value || '')
  }
  return undefined
})

function formatGatewayAmount(value: number): string {
  return formatPaymentAmount(value, currency.value, localeCode.value)
}

async function confirmAlipay(stripe: Stripe, clientSecret: string, orderId: number) {
  redirecting.value = true
  const returnUrl = window.location.origin + '/payment/result?order_id=' + orderId + '&status=success'
  const { error } = await stripe.confirmAlipayPayment(clientSecret, { return_url: returnUrl })
  if (error) {
    redirecting.value = false
    stripeError.value = error.message || t('payment.result.failed')
  }
  // 无错误时 Stripe 会自动跳转
}

async function confirmWechatPay(stripe: Stripe, clientSecret: string) {
  const { paymentIntent, error } = await (stripe as Stripe & {
    confirmWechatPayPayment: (cs: string, opts: Record<string, unknown>) => Promise<{ paymentIntent?: { status: string; next_action?: { wechat_pay_display_qr_code?: { image_data_url?: string } } }; error?: { message?: string } }>
  }).confirmWechatPayPayment(clientSecret, {
    payment_method_options: { wechat_pay: { client: isMobileDevice() ? 'mobile_web' : 'web' } },
  })

  if (error) {
    stripeError.value = error.message || t('payment.result.failed')
    return
  }

  // 从 next_action 中提取二维码
  const qrData = paymentIntent?.next_action?.wechat_pay_display_qr_code?.image_data_url
  if (qrData) {
    wechatQrUrl.value = qrData
    // 轮询支付完成状态
    startPolling()
  } else if (paymentIntent?.status === 'succeeded') {
    stripeSuccess.value = true
    scheduleClose()
  } else {
    stripeError.value = t('payment.result.failed')
  }
}

function mountPaymentElement(stripe: Stripe, clientSecret: string) {
  const isDark = document.documentElement.classList.contains('dark')
  const elements = stripe.elements({
    clientSecret,
    appearance: { theme: isDark ? 'night' : 'stripe', variables: { borderRadius: '8px' } },
  })
  elementsInstance = elements
  const paymentElement = elements.create('payment', {
    layout: 'tabs',
    paymentMethodOrder: ['alipay', 'wechat_pay', 'card', 'link'],
  } as Record<string, unknown>)
  paymentElement.mount('#stripe-payment-element')
  paymentElement.on('ready', () => { stripeReady.value = true })
}

async function handleGenericPay() {
  if (!stripeInstance || !elementsInstance || stripeSubmitting.value) return
  stripeSubmitting.value = true
  stripeError.value = ''
  try {
    const { error } = await stripeInstance.confirmPayment({
      elements: elementsInstance,
      confirmParams: {
        return_url: window.location.origin + '/payment/result?order_id=' + route.query.order_id + '&status=success',
      },
      redirect: 'if_required',
    })
    if (error) {
      stripeError.value = error.message || t('payment.result.failed')
    } else {
      stripeSuccess.value = true
      scheduleClose()
    }
  } catch (err: unknown) {
    stripeError.value = extractI18nErrorMessage(err, t, 'payment.errors', t('payment.result.failed'))
  } finally {
    stripeSubmitting.value = false
  }
}

let pollTimer: ReturnType<typeof setInterval> | null = null

function startPolling() {
  const orderId = Number(route.query.order_id)
  if (!orderId) return
  pollTimer = setInterval(async () => {
    const o = await paymentStore.pollOrderStatus(orderId)
    if (!o) return
    if (o.status === 'COMPLETED' || o.status === 'PAID') {
      if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
      stripeSuccess.value = true
      wechatQrUrl.value = ''
      scheduleClose()
    }
  }, 3000)
}

function scheduleClose() {
  if (window.opener) {
    redirectTimer = setTimeout(() => { window.close() }, 2000)
  } else {
    redirectTimer = setTimeout(() => {
      router.push({ path: '/payment/result', query: { order_id: String(route.query.order_id || ''), status: 'success' } })
    }, 2000)
  }
}

onUnmounted(() => {
  if (redirectTimer) clearTimeout(redirectTimer)
  if (pollTimer) clearInterval(pollTimer)
})
</script>
