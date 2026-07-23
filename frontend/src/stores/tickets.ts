import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { ticketsAPI } from '@/api/tickets'
import type { TicketConfig, TicketSummary } from '@/types/ticket'

const emptySummary = (): TicketSummary => ({
  total: 0,
  pending_admin: 0,
  pending_user: 0,
  resolved: 0,
  closed: 0,
  unread_count: 0,
  pending_admin_count: 0,
  open_count: 0,
})

export const useTicketStore = defineStore('tickets', () => {
  const config = ref<TicketConfig>({ user_center_enabled: true, accept_new_tickets: true })
  const configLoaded = ref(false)
  const userSummary = ref<TicketSummary>(emptySummary())
  const adminSummary = ref<TicketSummary>(emptySummary())
  const loading = ref(false)
  let pollingTimer: ReturnType<typeof setInterval> | null = null

  const userUnread = computed(() => userSummary.value.unread_count)
  const adminUnread = computed(() => adminSummary.value.unread_count)
  const adminPending = computed(() => adminSummary.value.pending_admin_count)

  async function fetchConfig(admin = false, force = false) {
    if (configLoaded.value && !force) return config.value
    config.value = admin ? await ticketsAPI.admin.getConfig() : await ticketsAPI.getConfig()
    configLoaded.value = true
    return config.value
  }

  async function refresh(admin = false) {
    if (loading.value) return
    loading.value = true
    try {
      const current = await fetchConfig(admin, true)
      if (admin) {
        adminSummary.value = await ticketsAPI.admin.summary()
      } else if (current.user_center_enabled) {
        userSummary.value = await ticketsAPI.summary()
      } else {
        userSummary.value = emptySummary()
      }
    } finally {
      loading.value = false
    }
  }

  async function updateConfig(next: TicketConfig) {
    config.value = await ticketsAPI.admin.updateConfig(next)
    configLoaded.value = true
    return config.value
  }

  function startPolling(admin = false) {
    stopPolling()
    void refresh(admin).catch(() => undefined)
    pollingTimer = setInterval(() => {
      void refresh(admin).catch(() => undefined)
    }, 30_000)
  }

  function stopPolling() {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  return {
    config,
    configLoaded,
    userSummary,
    adminSummary,
    userUnread,
    adminUnread,
    adminPending,
    loading,
    fetchConfig,
    refresh,
    updateConfig,
    startPolling,
    stopPolling,
  }
})

