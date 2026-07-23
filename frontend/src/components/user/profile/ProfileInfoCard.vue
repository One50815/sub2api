<template>
  <div class="space-y-4 sm:space-y-5">
    <section
      data-testid="profile-overview-hero"
      class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
    >
      <div class="p-3 sm:p-5">
        <div class="flex items-center gap-3 text-left sm:gap-4">
          <div
            class="flex h-12 w-12 shrink-0 items-center justify-center overflow-hidden rounded-xl bg-gray-900 text-sm font-semibold text-white ring-2 ring-white dark:bg-gray-100 dark:text-gray-900 dark:ring-dark-800 sm:h-16 sm:w-16 sm:rounded-2xl sm:text-lg sm:ring-4"
          >
            <img
              v-if="avatarUrl"
              :src="avatarUrl"
              :alt="displayName"
              class="h-full w-full object-cover"
            >
            <span v-else>{{ avatarInitial }}</span>
          </div>

          <div class="min-w-0 flex-1 space-y-1.5 sm:space-y-3">
              <div class="flex min-w-0 items-center gap-2">
                <h2 class="truncate text-xl font-semibold tracking-tight text-gray-950 dark:text-white sm:text-2xl">
                  {{ displayName }}
                </h2>
                <span :class="['badge', user?.role === 'admin' ? 'badge-primary' : 'badge-gray']">
                  {{ user?.role === 'admin' ? t('profile.administrator') : t('profile.user') }}
                </span>
                <span
                  :class="['badge', user?.status === 'active' ? 'badge-success' : 'badge-danger']"
                >
                  {{
                    user?.status === 'active'
                      ? t('common.active')
                      : t('common.disabled')
                  }}
                </span>
              </div>

              <div class="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-gray-500 dark:text-gray-400 sm:gap-x-3 sm:text-sm">
                <p v-if="primaryEmailDisplay" class="truncate">
                  {{ primaryEmailDisplay }}
                </p>
                <template
                  v-if="sourceHints.length"
                >
                  <span
                    v-for="hint in sourceHints"
                    :key="hint.key"
                    class="inline-flex min-w-0 items-center gap-1 rounded-md bg-gray-50 px-2 py-0.5 text-[11px] ring-1 ring-gray-200 dark:bg-dark-900/60 dark:ring-dark-700"
                  >
                    <Icon name="link" size="sm" />
                    <span class="truncate">{{ hint.text }}</span>
                  </span>
                </template>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-3 divide-x divide-gray-200 border-t border-gray-200 dark:divide-dark-700 dark:border-dark-700">
              <div
                data-testid="profile-overview-metric-balance"
                class="min-w-0 px-3 py-3 sm:px-5 sm:py-4"
              >
                <div class="flex items-center gap-2">
                  <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400">
                    <Icon name="dollar" size="sm" />
                  </span>
                <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">
                  {{ t('profile.accountBalance') }}
                </p>
                </div>
                <p class="mt-1.5 truncate font-mono text-base font-bold tabular-nums text-gray-950 dark:text-white sm:mt-2 sm:text-2xl">
                  {{ formatCurrency(user?.balance || 0) }}
                </p>
              </div>
              <div
                data-testid="profile-overview-metric-concurrency"
                class="min-w-0 px-3 py-3 sm:px-5 sm:py-4"
              >
                <div class="flex items-center gap-2">
                  <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-blue-50 text-blue-600 dark:bg-blue-500/10 dark:text-blue-400">
                    <Icon name="bolt" size="sm" />
                  </span>
                <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">
                  {{ t('profile.concurrencyLimit') }}
                </p>
                </div>
                <p class="mt-1.5 truncate font-mono text-base font-bold tabular-nums text-gray-950 dark:text-white sm:mt-2 sm:text-2xl">
                  {{ user?.concurrency || 0 }}
                </p>
              </div>
              <div
                data-testid="profile-overview-metric-member-since"
                class="min-w-0 px-3 py-3 sm:px-5 sm:py-4"
              >
                <div class="flex items-center gap-2">
                  <span class="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-violet-50 text-violet-600 dark:bg-violet-500/10 dark:text-violet-400">
                    <Icon name="calendar" size="sm" />
                  </span>
                <p class="truncate text-[11px] font-medium uppercase tracking-wider text-gray-500 dark:text-gray-400 sm:text-xs">
                  {{ t('profile.memberSince') }}
                </p>
                </div>
                <p class="mt-1.5 truncate font-mono text-base font-bold tabular-nums text-gray-950 dark:text-white sm:mt-2 sm:text-2xl">
                  {{ memberSinceLabel }}
                </p>
              </div>
      </div>
    </section>

    <div class="grid gap-4 sm:gap-5 xl:grid-cols-[minmax(0,1fr)_minmax(320px,0.46fr)] xl:items-start">
      <div data-testid="profile-main-column" class="space-y-4 sm:space-y-5">
        <section
          data-testid="profile-basics-panel"
          class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-800"
        >
          <div class="border-b border-gray-200 px-4 py-3.5 dark:border-dark-700 sm:px-5">
            <div>
              <h3 class="text-base font-semibold text-gray-950 dark:text-white">
                {{ t('profile.basicsTitle') }}
              </h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                {{ t('profile.basicsDescription') }}
              </p>
            </div>
          </div>

          <div class="grid gap-3 p-4 md:grid-cols-2 sm:p-5">
            <div class="rounded-lg border border-gray-200 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-900/30">
              <ProfileAvatarCard
                :user="user"
                embedded
              />
            </div>

            <div class="rounded-lg border border-gray-200 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-900/30">
              <ProfileEditForm
                :initial-username="user?.username || ''"
                embedded
              />
            </div>
          </div>
        </section>

        <section
          data-testid="profile-auth-bindings-panel"
          class="rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800 sm:p-5"
        >
          <ProfileIdentityBindingsSection
            :user="user"
            :linuxdo-enabled="linuxdoEnabled"
            :dingtalk-enabled="dingtalkEnabled"
            :oidc-enabled="oidcEnabled"
            :oidc-provider-name="oidcProviderName"
            :wechat-enabled="wechatEnabled"
            :wechat-open-enabled="wechatOpenEnabled"
            :wechat-mp-enabled="wechatMpEnabled"
            embedded
            compact
          />
        </section>
      </div>

      <div data-testid="profile-side-column" class="space-y-4 sm:space-y-5 xl:sticky xl:top-4">
        <section
          v-if="sourceHints.length"
          class="rounded-lg border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800 sm:p-5"
        >
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ t('profile.linkedProfileSources') }}
          </h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t('profile.linkedProfileSourcesDescription') }}
          </p>

          <div class="mt-4 grid gap-2">
            <div
              v-for="hint in sourceHints"
              :key="hint.key"
              class="flex items-start gap-3 rounded-lg border border-gray-200 bg-gray-50/70 px-3 py-2.5 text-sm text-gray-600 dark:border-dark-700 dark:bg-dark-900/30 dark:text-gray-300"
            >
              <Icon name="link" size="sm" class="mt-0.5 text-gray-400 dark:text-gray-500" />
              <span>{{ hint.text }}</span>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import ProfileAvatarCard from '@/components/user/profile/ProfileAvatarCard.vue'
import ProfileEditForm from '@/components/user/profile/ProfileEditForm.vue'
import ProfileIdentityBindingsSection from '@/components/user/profile/ProfileIdentityBindingsSection.vue'
import type { User, UserAuthBindingStatus, UserAuthProvider, UserProfileSourceContext } from '@/types'

const props = withDefaults(defineProps<{
  user: User | null
  linuxdoEnabled?: boolean
  dingtalkEnabled?: boolean
  oidcEnabled?: boolean
  oidcProviderName?: string
  wechatEnabled?: boolean
  wechatOpenEnabled?: boolean
  wechatMpEnabled?: boolean
}>(), {
  linuxdoEnabled: false,
  dingtalkEnabled: false,
  oidcEnabled: false,
  oidcProviderName: 'OIDC',
  wechatEnabled: false,
  wechatOpenEnabled: undefined,
  wechatMpEnabled: undefined,
})

const { t } = useI18n()

function normalizeBindingStatus(binding: boolean | UserAuthBindingStatus | undefined): boolean | null {
  if (typeof binding === 'boolean') {
    return binding
  }
  if (!binding) {
    return null
  }
  if (typeof binding.bound === 'boolean') {
    return binding.bound
  }
  return Boolean(binding.provider_subject || binding.issuer || binding.provider_key)
}

function isEmailBound(user: User | null | undefined): boolean {
  if (typeof user?.email_bound === 'boolean') {
    return user.email_bound
  }

  const nested = user?.auth_bindings?.email ?? user?.identity_bindings?.email
  const normalized = normalizeBindingStatus(nested)
  return normalized ?? false
}

const avatarUrl = computed(() => props.user?.avatar_url?.trim() || '')
const displayName = computed(() => props.user?.username?.trim() || props.user?.email?.trim() || t('profile.user'))
const primaryEmailDisplay = computed(() => {
  const email = props.user?.email?.trim() || ''
  if (!email) {
    return ''
  }
  if (email.endsWith('.invalid') && !isEmailBound(props.user)) {
    return ''
  }
  return email
})
const avatarInitial = computed(() => displayName.value.charAt(0).toUpperCase() || 'U')
const memberSinceLabel = computed(() => {
  const raw = props.user?.created_at?.trim()
  if (!raw) {
    return '-'
  }

  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) {
    return '-'
  }

  return new Intl.DateTimeFormat(undefined, {
    year: 'numeric',
    month: 'short',
  }).format(date)
})

const providerLabels = computed<Record<UserAuthProvider, string>>(() => ({
  email: t('profile.authBindings.providers.email'),
  linuxdo: t('profile.authBindings.providers.linuxdo'),
  dingtalk: t('profile.authBindings.providers.dingtalk'),
  oidc: t('profile.authBindings.providers.oidc', { providerName: props.oidcProviderName }),
  wechat: t('profile.authBindings.providers.wechat'),
  github: 'GitHub',
  google: 'Google'
}))

function formatCurrency(value: number): string {
  return `$${value.toFixed(2)}`
}

function normalizeProvider(value: string): UserAuthProvider | null {
  const normalized = value.trim().toLowerCase()
  if (
    normalized === 'email' ||
    normalized === 'linuxdo' ||
    normalized === 'wechat' ||
    normalized === 'github' ||
    normalized === 'google'
  ) {
    return normalized
  }
  if (normalized === 'oidc' || normalized.startsWith('oidc:') || normalized.startsWith('oidc/')) {
    return 'oidc'
  }
  return null
}

function readObjectString(source: Record<string, unknown>, ...keys: string[]): string {
  for (const key of keys) {
    const value = source[key]
    if (typeof value === 'string' && value.trim()) {
      return value.trim()
    }
  }
  return ''
}

function resolveThirdPartySource(
  rawSource: string | UserProfileSourceContext | null | undefined
): { provider: UserAuthProvider; label: string } | null {
  if (!rawSource) {
    return null
  }

  if (typeof rawSource === 'string') {
    const provider = normalizeProvider(rawSource)
    if (!provider || provider === 'email') {
      return null
    }
    return {
      provider,
      label: providerLabels.value[provider]
    }
  }

  const sourceRecord = rawSource as Record<string, unknown>
  const provider = normalizeProvider(
    readObjectString(sourceRecord, 'provider', 'source', 'provider_type', 'auth_provider')
  )
  if (!provider || provider === 'email') {
    return null
  }

  const explicitLabel = readObjectString(
    sourceRecord,
    'provider_label',
    'label',
    'provider_name',
    'providerName'
  )

  return {
    provider,
    label: explicitLabel || providerLabels.value[provider]
  }
}

const sourceHints = computed(() => {
  const currentUser = props.user
  if (!currentUser) {
    return []
  }

  const hints: Array<{ key: string; text: string }> = []
  const avatarSource = resolveThirdPartySource(
    currentUser.profile_sources?.avatar ?? currentUser.avatar_source
  )
  const usernameSource = resolveThirdPartySource(
    currentUser.profile_sources?.username ??
      currentUser.profile_sources?.display_name ??
      currentUser.profile_sources?.nickname ??
      currentUser.display_name_source ??
      currentUser.username_source ??
      currentUser.nickname_source
  )

  if (avatarSource) {
    hints.push({
      key: 'avatar',
      text: t('profile.authBindings.source.avatar', { providerName: avatarSource.label })
    })
  }

  if (usernameSource) {
    hints.push({
      key: 'username',
      text: t('profile.authBindings.source.username', { providerName: usernameSource.label })
    })
  }

  return hints
})
</script>
