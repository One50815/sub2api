<template>
  <AppLayout>
    <div class="mx-auto max-w-7xl space-y-5 p-4 sm:p-6">
      <header class="flex flex-wrap items-end justify-between gap-4">
        <div>
          <p class="text-xs font-semibold uppercase text-primary-600 dark:text-primary-400">Omnio Pro Control</p>
          <h1 class="mt-1 text-2xl font-bold text-gray-950 dark:text-white">Omnio Pro 管理</h1>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">统一管理 Pro 等级、独立分组倍率、专属分组、售卖方案和订阅赠送关系。</p>
        </div>
        <button class="btn btn-secondary" :disabled="loading" @click="load">
          <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />刷新
        </button>
      </header>

      <div class="grid gap-5 xl:grid-cols-[1.1fr_0.9fr]">
        <section class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-dark-700">
            <h2 class="font-semibold text-gray-900 dark:text-white">Pro 等级与分组权益</h2>
            <button class="btn btn-primary" @click="newLevel">新增等级</button>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-dark-700">
            <article v-for="level in catalog.levels" :key="level.id" class="space-y-3 px-5 py-4">
              <div class="flex flex-wrap items-start justify-between gap-3">
                <div class="flex items-start gap-3">
                  <span class="mt-1.5 h-3 w-3 rounded-full" :style="{ backgroundColor: level.badge_color }"></span>
                  <div>
                    <p class="font-semibold text-gray-900 dark:text-white">{{ level.name }} <span class="ml-2 text-xs text-gray-400">{{ level.slug }} · Rank {{ level.rank }}</span></p>
                    <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ level.description || '未填写描述' }}</p>
                  </div>
                </div>
                <div class="flex gap-2"><button class="btn btn-secondary px-3" @click="newBenefit(level.id)">分组权益</button><button class="btn btn-secondary px-3" @click="editLevel(level)">编辑</button><button class="btn btn-danger px-3" @click="removeLevel(level)">删除</button></div>
              </div>
              <div class="flex flex-wrap gap-2 text-xs">
                <span class="rounded-full bg-gray-100 px-2.5 py-1 text-gray-600 dark:bg-dark-800 dark:text-gray-300">并发 +{{ level.concurrency_bonus }}</span>
                <span class="rounded-full bg-gray-100 px-2.5 py-1 text-gray-600 dark:bg-dark-800 dark:text-gray-300">{{ level.priority_support ? '优先支持' : '标准支持' }}</span>
                <button v-for="benefit in level.group_benefits" :key="benefit.id" class="rounded-full border border-gray-200 px-2.5 py-1 text-gray-500 hover:border-primary-300 hover:text-primary-600 dark:border-dark-700 dark:text-gray-400" @click="editBenefit(benefit)">
                  {{ benefit.group_name || `分组 #${benefit.group_id}` }}
                  <template v-if="benefit.pro_only"> · Pro 专属</template>
                  <template v-if="benefit.rate_multiplier != null"> · {{ benefit.rate_multiplier }}x</template>
                  <template v-if="benefit.rpm_limit != null"> · {{ benefit.rpm_limit }} RPM</template>
                </button>
              </div>
            </article>
            <p v-if="!catalog.levels.length" class="p-8 text-center text-sm text-gray-500">还没有 Omnio Pro 等级</p>
          </div>
        </section>

        <section class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-dark-700">
            <h2 class="font-semibold text-gray-900 dark:text-white">Omnio Pro 方案</h2>
            <button class="btn btn-primary" :disabled="!catalog.levels.length" @click="newOffer">新增方案</button>
          </div>
          <div class="divide-y divide-gray-100 dark:divide-dark-700">
            <article v-for="offer in catalog.offers" :key="offer.id" class="flex items-start justify-between gap-3 px-5 py-4">
              <div>
                <p class="font-semibold text-gray-900 dark:text-white">{{ offer.name }} <span class="ml-2 text-xs text-gray-400">{{ offer.level_name }}</span></p>
                <p class="mt-1 text-sm text-gray-600 dark:text-gray-300">{{ offer.currency }} {{ offer.price.toFixed(2) }} · {{ offer.duration_days }} 天</p>
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ offer.for_sale ? '销售中' : '已下架' }} · {{ offer.description || '未填写描述' }}</p>
              </div>
              <div class="flex gap-2"><button class="btn btn-secondary px-3" @click="editOffer(offer)">编辑</button><button class="btn btn-danger px-3" @click="removeOffer(offer.id)">删除</button></div>
            </article>
            <p v-if="!catalog.offers.length" class="p-8 text-center text-sm text-gray-500">还没有 Omnio Pro 方案</p>
          </div>
        </section>
      </div>

      <div class="grid gap-5 xl:grid-cols-[1.15fr_0.85fr]">
        <section v-if="false" class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
          <div class="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-dark-700">
            <div><h2 class="font-semibold text-gray-900 dark:text-white">订阅赠送 Omnio Pro</h2><p class="mt-1 text-xs text-gray-500 dark:text-gray-400">保留订阅兼容：购买指定套餐后按套餐期限或指定天数发放 Pro。</p></div>
            <button class="btn btn-primary" :disabled="!catalog.levels.length" @click="newPlanBenefit">新增关系</button>
          </div>
          <div class="overflow-x-auto">
            <table class="min-w-full text-left text-sm">
              <thead class="bg-gray-50 text-xs text-gray-500 dark:bg-dark-800 dark:text-gray-400"><tr><th class="px-5 py-3">订阅套餐</th><th class="px-5 py-3">Pro 等级</th><th class="px-5 py-3">天数</th><th class="px-5 py-3">操作</th></tr></thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-for="item in catalog.plan_benefits" :key="item.plan_id"><td class="px-5 py-3 text-gray-900 dark:text-white">{{ item.plan_name || `套餐 #${item.plan_id}` }}</td><td class="px-5 py-3 text-gray-600 dark:text-gray-300">{{ item.level_name }}</td><td class="px-5 py-3 text-gray-600 dark:text-gray-300">{{ item.duration_days || '跟随订阅' }}</td><td class="whitespace-nowrap px-5 py-3"><button class="mr-3 text-xs font-semibold text-primary-600" @click="editPlanBenefit(item)">编辑</button><button class="text-xs font-semibold text-red-600" @click="removePlanBenefit(item.plan_id)">删除</button></td></tr>
                <tr v-if="!catalog.plan_benefits.length"><td colspan="4" class="px-5 py-8 text-center text-sm text-gray-500">还没有订阅赠送关系</td></tr>
              </tbody>
            </table>
          </div>
        </section>

        <section class="rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
          <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700"><h2 class="font-semibold text-gray-900 dark:text-white">手动发放 Omnio Pro</h2><p class="mt-1 text-xs text-gray-500 dark:text-gray-400">用于补偿、活动和人工续期，操作会写入审计日志。</p></div>
          <form class="grid gap-4 p-5 sm:grid-cols-2" @submit.prevent="grantMembership">
            <label class="field"><span>用户 ID</span><input v-model.number="grantForm.user_id" type="number" min="1" class="input" required /></label>
            <label class="field"><span>Pro 等级</span><select v-model.number="grantForm.level_id" class="input" required><option v-for="level in catalog.levels" :key="level.id" :value="level.id">{{ level.name }}</option></select></label>
            <label class="field"><span>天数</span><input v-model.number="grantForm.days" type="number" min="1" class="input" required /></label>
            <label class="field"><span>备注</span><input v-model="grantForm.notes" class="input" /></label>
            <button class="btn btn-primary sm:col-span-2" :disabled="granting || !catalog.levels.length">{{ granting ? '发放中...' : '确认发放' }}</button>
          </form>
        </section>
      </div>

      <section class="overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-900">
        <div class="border-b border-gray-100 px-5 py-4 dark:border-dark-700"><h2 class="font-semibold text-gray-900 dark:text-white">Omnio Pro 审计记录</h2></div>
        <div class="overflow-x-auto"><table class="min-w-full text-left text-sm"><thead class="bg-gray-50 text-xs text-gray-500 dark:bg-dark-800 dark:text-gray-400"><tr><th class="px-5 py-3">时间</th><th class="px-5 py-3">动作</th><th class="px-5 py-3">用户 / 等级</th><th class="px-5 py-3">来源</th><th class="px-5 py-3">操作者</th></tr></thead><tbody class="divide-y divide-gray-100 dark:divide-dark-700"><tr v-for="log in auditLogs" :key="log.id"><td class="whitespace-nowrap px-5 py-3 text-gray-500">{{ formatDate(log.created_at) }}</td><td class="px-5 py-3 font-medium text-gray-900 dark:text-white">{{ log.action }}</td><td class="px-5 py-3 text-gray-600 dark:text-gray-300">{{ log.user_id ? `用户 #${log.user_id}` : '-' }} / {{ log.level_id ? `等级 #${log.level_id}` : '-' }}</td><td class="px-5 py-3 text-gray-500">{{ log.source_type }} {{ log.source_id }}</td><td class="px-5 py-3 text-gray-500">{{ log.operator }}</td></tr><tr v-if="!auditLogs.length"><td colspan="5" class="px-5 py-8 text-center text-sm text-gray-500">暂无审计记录</td></tr></tbody></table></div>
      </section>

      <BaseDialog :show="dialogOpen" :title="dialogTitle" width="wide" @close="dialogOpen = false">
        <form class="space-y-4" @submit.prevent="submitDialog">
          <div v-if="dialogType === 'level'" class="grid gap-4 sm:grid-cols-2">
            <label class="field"><span>名称</span><input v-model="levelForm.name" class="input" required /></label><label class="field"><span>Slug</span><input v-model="levelForm.slug" class="input" required /></label>
            <label class="field"><span>Rank</span><input v-model.number="levelForm.rank" type="number" min="0" class="input" required /></label><label class="field"><span>并发加成</span><input v-model.number="levelForm.concurrency_bonus" type="number" min="0" class="input" required /></label>
            <label class="field"><span>徽标颜色</span><input v-model="levelForm.badge_color" type="color" class="input h-10" /></label><label class="field"><span>排序</span><input v-model.number="levelForm.sort_order" type="number" class="input" /></label>
            <label class="field sm:col-span-2"><span>描述</span><textarea v-model="levelForm.description" class="input min-h-20" /></label><label class="flex items-center gap-2 text-sm"><input v-model="levelForm.priority_support" type="checkbox" />优先支持</label><label class="flex items-center gap-2 text-sm"><input v-model="levelForm.active" type="checkbox" />启用</label>
          </div>
          <div v-else-if="dialogType === 'offer'" class="grid gap-4 sm:grid-cols-2">
            <label class="field"><span>所属等级</span><select v-model.number="offerForm.level_id" class="input" required><option v-for="level in catalog.levels" :key="level.id" :value="level.id">{{ level.name }}</option></select></label><label class="field"><span>方案名称</span><input v-model="offerForm.name" class="input" required /></label>
            <label class="field"><span>价格</span><input v-model.number="offerForm.price" type="number" min="0" step="0.01" class="input" required /></label><label class="field"><span>原价</span><input v-model.number="offerForm.original_price" type="number" min="0" step="0.01" class="input" /></label>
            <label class="field"><span>币种</span><input v-model="offerForm.currency" class="input" maxlength="3" /></label><label class="field"><span>有效天数</span><input v-model.number="offerForm.duration_days" type="number" min="1" class="input" required /></label>
            <label class="field"><span>排序</span><input v-model.number="offerForm.sort_order" type="number" class="input" /></label><label class="flex items-center gap-2 pt-6 text-sm"><input v-model="offerForm.for_sale" type="checkbox" />前台销售</label><label class="field sm:col-span-2"><span>描述</span><textarea v-model="offerForm.description" class="input min-h-20" /></label>
          </div>
          <div v-else-if="dialogType === 'benefit'" class="grid gap-4 sm:grid-cols-2">
            <label class="field"><span>Pro 等级</span><select v-model.number="benefitForm.level_id" class="input" required><option v-for="level in catalog.levels" :key="level.id" :value="level.id">{{ level.name }}</option></select></label><label class="field"><span>分组 ID</span><input v-model.number="benefitForm.group_id" type="number" min="1" class="input" required /></label>
            <label class="field"><span>Pro 最终倍率（留空使用分组基础倍率）</span><input v-model.number="benefitForm.rate_multiplier" type="number" min="0" step="0.01" class="input" /><small>只覆盖该 Pro 等级用户的最终倍率，不修改分组全局倍率。</small></label><label class="field"><span>RPM 上限（留空不覆盖）</span><input v-model.number="benefitForm.rpm_limit" type="number" min="0" class="input" /></label>
            <label class="flex items-center gap-2 text-sm sm:col-span-2"><input v-model="benefitForm.allow_access" type="checkbox" />允许该等级访问此分组</label>
            <label class="flex items-center gap-2 text-sm sm:col-span-2"><input v-model="benefitForm.pro_only" type="checkbox" @change="onProOnlyChange" />仅 Omnio Pro 可见和绑定</label>
            <button v-if="benefitEditing" type="button" class="btn btn-danger sm:col-span-2" @click="removeBenefit">删除该分组权益</button>
          </div>
          <div v-else class="grid gap-4 sm:grid-cols-2">
            <label class="field"><span>订阅套餐 ID</span><input v-model.number="planBenefitForm.plan_id" type="number" min="1" class="input" required /></label><label class="field"><span>Pro 等级</span><select v-model.number="planBenefitForm.level_id" class="input" required><option v-for="level in catalog.levels" :key="level.id" :value="level.id">{{ level.name }}</option></select></label>
            <label class="field sm:col-span-2"><span>赠送天数（留空跟随订阅）</span><input v-model.number="planBenefitForm.duration_days" type="number" min="1" class="input" /></label>
          </div>
          <div class="flex justify-end gap-3"><button type="button" class="btn btn-secondary" @click="dialogOpen = false">取消</button><button type="submit" class="btn btn-primary">保存</button></div>
        </form>
      </BaseDialog>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import membershipAPI from '@/api/membership'
import { useAppStore } from '@/stores'
import type { MembershipAuditLog, MembershipCatalog, MembershipGroupBenefit, MembershipLevel, MembershipOffer, MembershipPlanBenefit } from '@/types/membership'

const appStore = useAppStore()
const loading = ref(false)
const granting = ref(false)
const dialogOpen = ref(false)
const dialogType = ref<'level' | 'offer' | 'benefit' | 'plan'>('level')
const benefitEditing = ref(false)
const catalog = reactive<MembershipCatalog>({ levels: [], offers: [], plan_benefits: [] })
const auditLogs = ref<MembershipAuditLog[]>([])
const levelForm = reactive<Partial<MembershipLevel>>({ name: '', slug: '', description: '', rank: 0, badge_color: '#2563eb', concurrency_bonus: 0, priority_support: false, active: true, sort_order: 0 })
const offerForm = reactive<Partial<MembershipOffer>>({ level_id: 0, name: '', description: '', price: 0, original_price: null, currency: 'USD', duration_days: 30, for_sale: true, sort_order: 0 })
const benefitForm = reactive<{ level_id: number; group_id: number; allow_access: boolean; pro_only: boolean; rate_multiplier: number | null; rpm_limit: number | null }>({ level_id: 0, group_id: 0, allow_access: true, pro_only: false, rate_multiplier: null, rpm_limit: null })
const planBenefitForm = reactive<MembershipPlanBenefit>({ plan_id: 0, plan_name: '', level_id: 0, level_name: '', duration_days: null })
const grantForm = reactive({ user_id: 0, level_id: 0, days: 30, notes: '' })

const dialogTitle = computed(() => {
  if (dialogType.value === 'level') return levelForm.id ? '编辑 Pro 等级' : '新增 Pro 等级'
  if (dialogType.value === 'offer') return offerForm.id ? '编辑 Omnio Pro 方案' : '新增 Omnio Pro 方案'
  if (dialogType.value === 'benefit') return '配置分组权益'
  return '配置订阅赠送 Omnio Pro'
})

function showError(error: unknown) {
  appStore.showError(error instanceof Error ? error.message : '操作失败')
}

async function load() {
  loading.value = true
  try {
    const [catalogResponse, auditResponse] = await Promise.all([membershipAPI.adminCatalog(), membershipAPI.auditLogs(50)])
    Object.assign(catalog, catalogResponse.data)
    auditLogs.value = auditResponse.data
    if (!grantForm.level_id && catalog.levels.length) grantForm.level_id = catalog.levels[0].id
  } catch (error) {
    showError(error)
  } finally {
    loading.value = false
  }
}

function newLevel() { Object.assign(levelForm, { id: undefined, name: '', slug: '', description: '', rank: 0, badge_color: '#2563eb', concurrency_bonus: 0, priority_support: false, active: true, sort_order: 0 }); dialogType.value = 'level'; dialogOpen.value = true }
function editLevel(level: MembershipLevel) { Object.assign(levelForm, level); dialogType.value = 'level'; dialogOpen.value = true }
function newOffer() { Object.assign(offerForm, { id: undefined, level_id: catalog.levels[0]?.id || 0, name: '', description: '', price: 0, original_price: null, currency: 'USD', duration_days: 30, for_sale: true, sort_order: 0 }); dialogType.value = 'offer'; dialogOpen.value = true }
function editOffer(offer: MembershipOffer) { Object.assign(offerForm, offer); dialogType.value = 'offer'; dialogOpen.value = true }
function newBenefit(levelId: number) { Object.assign(benefitForm, { level_id: levelId, group_id: 0, allow_access: true, pro_only: false, rate_multiplier: null, rpm_limit: null }); benefitEditing.value = false; dialogType.value = 'benefit'; dialogOpen.value = true }
function editBenefit(benefit: MembershipGroupBenefit) { Object.assign(benefitForm, benefit); benefitEditing.value = true; dialogType.value = 'benefit'; dialogOpen.value = true }
function newPlanBenefit() { Object.assign(planBenefitForm, { plan_id: 0, plan_name: '', level_id: catalog.levels[0]?.id || 0, level_name: '', duration_days: null }); dialogType.value = 'plan'; dialogOpen.value = true }
function editPlanBenefit(item: MembershipPlanBenefit) { Object.assign(planBenefitForm, item); dialogType.value = 'plan'; dialogOpen.value = true }

async function submitDialog() {
  try {
    if (dialogType.value === 'level') await membershipAPI.saveLevel(levelForm)
    else if (dialogType.value === 'offer') await membershipAPI.saveOffer(offerForm)
    else if (dialogType.value === 'benefit') {
      if (benefitForm.pro_only) benefitForm.allow_access = true
      await membershipAPI.saveBenefit(benefitForm)
    }
    else await membershipAPI.savePlanBenefit(planBenefitForm)
    dialogOpen.value = false
    await load()
  } catch (error) { showError(error) }
}

function onProOnlyChange() { if (benefitForm.pro_only) benefitForm.allow_access = true }
async function removeOffer(id: number) { if (!window.confirm('删除该 Omnio Pro 方案？已有订单不会受影响。')) return; try { await membershipAPI.deleteOffer(id); await load() } catch (error) { showError(error) } }
async function removeLevel(level: MembershipLevel) { if (!window.confirm(`确定删除 Omnio Pro 等级“${level.name}”吗？仍有关联方案或用户权益时系统会拒绝删除。`)) return; try { await membershipAPI.deleteLevel(level.id); await load() } catch (error) { showError(error) } }
async function removeBenefit() { if (!window.confirm('删除该分组权益？')) return; try { await membershipAPI.deleteBenefit(benefitForm.level_id, benefitForm.group_id); dialogOpen.value = false; await load() } catch (error) { showError(error) } }
async function removePlanBenefit(planId: number) { if (!window.confirm('删除该订阅赠送关系？')) return; try { await membershipAPI.deletePlanBenefit(planId); await load() } catch (error) { showError(error) } }
async function grantMembership() { granting.value = true; try { await membershipAPI.grant(grantForm); Object.assign(grantForm, { user_id: 0, days: 30, notes: '' }); await load() } catch (error) { showError(error) } finally { granting.value = false } }
function formatDate(value: string) { return new Date(value).toLocaleString('zh-CN') }

onMounted(() => { void load() })
</script>
