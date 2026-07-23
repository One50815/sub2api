<template>
  <div class="header-search">
    <button
      type="button"
      class="header-search-trigger"
      :aria-label="searchLabel"
      aria-haspopup="dialog"
      :aria-expanded="open"
      @click="openSearch"
    >
      <Icon name="search" size="sm" />
      <span class="header-search-label">{{ searchLabel }}</span>
      <kbd class="header-search-shortcut">{{ shortcutLabel }}</kbd>
    </button>

    <Teleport to="body">
      <Transition name="command-palette">
        <div
          v-if="open"
          class="header-search-overlay"
          role="presentation"
          @mousedown.self="closeSearch"
        >
          <section
            class="header-search-panel"
            role="dialog"
            aria-modal="true"
            :aria-label="searchLabel"
          >
            <div class="header-search-input-row">
              <Icon name="search" size="sm" />
              <input
                ref="inputRef"
                v-model="query"
                type="search"
                :placeholder="searchLabel"
                autocomplete="off"
                @keydown.esc.prevent="closeSearch"
                @keydown.down.prevent="moveSelection(1)"
                @keydown.up.prevent="moveSelection(-1)"
                @keydown.enter.prevent="openSelected"
              />
              <kbd>Esc</kbd>
            </div>

            <div class="header-search-results" role="listbox">
              <p class="header-search-group-label">{{ navigationLabel }}</p>
              <button
                v-for="(item, index) in filteredItems"
                :key="item.path"
                type="button"
                role="option"
                :aria-selected="selectedIndex === index"
                :class="['header-search-result', selectedIndex === index && 'is-selected']"
                @mousemove="selectedIndex = index"
                @click="navigate(item.path)"
              >
                <span class="header-search-result-icon">
                  <Icon :name="item.icon" size="sm" />
                </span>
                <span class="min-w-0 flex-1 truncate">{{ item.label }}</span>
                <span class="header-search-result-path">{{ item.path }}</span>
              </button>

              <div v-if="filteredItems.length === 0" class="header-search-empty">
                {{ emptyLabel }}
              </div>
            </div>
          </section>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import Icon from '@/components/icons/Icon.vue'
import { useAuthStore } from '@/stores/auth'

type IconName =
  | 'home'
  | 'key'
  | 'chartBar'
  | 'creditCard'
  | 'document'
  | 'user'
  | 'users'
  | 'cube'
  | 'server'
  | 'userCircle'
  | 'chart'
  | 'clipboard'
  | 'cog'

interface SearchItem {
  path: string
  label: string
  icon: IconName
}

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const open = ref(false)
const query = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)

const searchLabel = computed(() => {
  const value = t('common.search')
  return value === 'common.search' ? 'Search' : value
})
const navigationLabel = computed(() => {
  const value = t('nav.console')
  return value === 'nav.console' ? 'Navigation' : value
})
const emptyLabel = computed(() => {
  const value = t('common.noResults')
  return value === 'common.noResults' ? 'No results found' : value
})
const shortcutLabel = computed(() => {
  if (typeof navigator === 'undefined') return 'Ctrl K'
  return /Mac|iPhone|iPad/i.test(navigator.platform) ? '⌘K' : 'Ctrl K'
})

const items = computed<SearchItem[]>(() => {
  const personal: SearchItem[] = [
    { path: '/dashboard', label: t('dashboard.title'), icon: 'home' },
    { path: '/keys', label: t('keys.title'), icon: 'key' },
    { path: '/usage', label: t('usage.title'), icon: 'chartBar' },
    { path: '/purchase', label: t('nav.buySubscription'), icon: 'creditCard' },
    { path: '/orders', label: t('nav.myOrders'), icon: 'document' },
    { path: '/profile', label: t('profile.title'), icon: 'user' },
  ]

  if (!authStore.isAdmin) return personal

  return [
    { path: '/admin/dashboard', label: t('admin.dashboard.title'), icon: 'home' },
    { path: '/admin/users', label: t('admin.users.title'), icon: 'users' },
    { path: '/admin/groups', label: t('admin.groups.title'), icon: 'cube' },
    { path: '/admin/channels/pricing', label: t('admin.channels.title'), icon: 'server' },
    { path: '/admin/accounts', label: t('admin.accounts.title'), icon: 'userCircle' },
    { path: '/admin/usage', label: t('admin.usage.title'), icon: 'chartBar' },
    { path: '/admin/ops', label: t('admin.ops.title'), icon: 'chart' },
    { path: '/admin/audit-logs', label: t('admin.audit.title'), icon: 'clipboard' },
    { path: '/admin/settings', label: t('admin.settings.title'), icon: 'cog' },
    ...personal,
  ]
})

const filteredItems = computed(() => {
  const needle = query.value.trim().toLocaleLowerCase()
  if (!needle) return items.value
  return items.value.filter((item) =>
    `${item.label} ${item.path}`.toLocaleLowerCase().includes(needle)
  )
})

watch(filteredItems, () => {
  selectedIndex.value = 0
})

function openSearch(): void {
  open.value = true
  query.value = ''
  selectedIndex.value = 0
  nextTick(() => inputRef.value?.focus())
}

function closeSearch(): void {
  open.value = false
}

function moveSelection(direction: number): void {
  const count = filteredItems.value.length
  if (!count) return
  selectedIndex.value = (selectedIndex.value + direction + count) % count
}

function openSelected(): void {
  const item = filteredItems.value[selectedIndex.value]
  if (item) navigate(item.path)
}

async function navigate(path: string): Promise<void> {
  closeSearch()
  await router.push(path)
}

function handleShortcut(event: KeyboardEvent): void {
  if ((event.metaKey || event.ctrlKey) && event.key.toLocaleLowerCase() === 'k') {
    event.preventDefault()
    open.value ? closeSearch() : openSearch()
    return
  }
  if (event.key === 'Escape' && open.value) closeSearch()
}

onMounted(() => document.addEventListener('keydown', handleShortcut))
onBeforeUnmount(() => document.removeEventListener('keydown', handleShortcut))
</script>
