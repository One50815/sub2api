<template>
  <section class="border-t border-gray-200 pt-5 dark:border-dark-700">
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div>
        <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
          会员等级优惠
        </h3>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          当前分组基础倍率 {{ formatRate(baseRateMultiplier) }}
        </p>
      </div>
      <button
        type="button"
        class="btn btn-secondary px-3"
        :disabled="loading || !nextAvailableLevel"
        @click="addLevel"
      >
        <Icon name="plus" size="sm" />
        添加会员等级
      </button>
    </div>

    <div
      v-if="loading"
      class="mt-4 border-y border-gray-100 py-5 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400"
    >
      正在加载会员等级...
    </div>
    <div
      v-else-if="!levels.length"
      class="mt-4 border-y border-gray-100 py-5 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400"
    >
      暂无可用会员等级
    </div>
    <div
      v-else-if="!modelValue.length"
      class="mt-4 border-y border-gray-100 py-5 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400"
    >
      未配置会员等级优惠
    </div>
    <div
      v-else
      class="mt-4 divide-y divide-gray-100 border-y border-gray-100 dark:divide-dark-700 dark:border-dark-700"
    >
      <div
        v-for="(row, index) in modelValue"
        :key="`${row.level_id}-${index}`"
        class="grid gap-3 py-4 sm:grid-cols-[minmax(0,1fr)_minmax(0,0.7fr)_2.5rem] sm:items-end"
      >
        <label>
          <span class="input-label">会员等级</span>
          <select
            class="input"
            :value="row.level_id"
            @change="updateLevel(index, $event)"
          >
            <option
              v-for="level in levels"
              :key="level.id"
              :value="level.id"
              :disabled="isLevelSelectedByAnother(level.id, index)"
            >
              {{ level.name }}{{ level.active ? "" : "（停用）" }}
            </option>
          </select>
        </label>

        <label>
          <span class="input-label">优惠倍率</span>
          <input
            class="input"
            type="number"
            min="0"
            step="0.001"
            :value="row.rate_multiplier ?? ''"
            :placeholder="`继承 ${formatRate(baseRateMultiplier)}`"
            @input="updateRate(index, $event)"
          />
        </label>

        <button
          type="button"
          class="flex h-10 w-10 justify-self-end items-center justify-center text-gray-400 transition-colors hover:text-red-600 dark:text-gray-500 dark:hover:text-red-400 sm:justify-self-auto"
          :title="`删除 ${levelName(row.level_id)} 的优惠`"
          @click="removeLevel(index)"
        >
          <Icon name="trash" size="sm" />
        </button>
      </div>
    </div>

    <p
      v-if="modelValue.length"
      class="mt-2 text-xs text-gray-500 dark:text-gray-400"
    >
      0 表示免费；留空时继承当前分组基础倍率。
    </p>
  </section>
</template>

<script setup lang="ts">
import { computed } from "vue";
import Icon from "@/components/icons/Icon.vue";
import type { MembershipLevel } from "@/types/membership";
import {
  createGroupMembershipBenefitDraft,
  type GroupMembershipBenefitDraft,
} from "@/views/admin/groupsMembershipBenefits";

const props = withDefaults(
  defineProps<{
    modelValue: GroupMembershipBenefitDraft[];
    levels: MembershipLevel[];
    baseRateMultiplier: number;
    loading?: boolean;
  }>(),
  { loading: false },
);

const emit = defineEmits<{
  "update:modelValue": [value: GroupMembershipBenefitDraft[]];
}>();

const selectedLevelIds = computed(
  () => new Set(props.modelValue.map((row) => row.level_id)),
);
const nextAvailableLevel = computed(() =>
  props.levels.find((level) => !selectedLevelIds.value.has(level.id)),
);

function formatRate(value: number): string {
  return `${Number(value || 0).toFixed(3).replace(/\.?0+$/, "")}x`;
}

function levelName(levelId: number): string {
  return props.levels.find((level) => level.id === levelId)?.name ?? "该等级";
}

function isLevelSelectedByAnother(levelId: number, currentIndex: number): boolean {
  return props.modelValue.some(
    (row, index) => index !== currentIndex && row.level_id === levelId,
  );
}

function replaceRow(
  index: number,
  patch: Partial<GroupMembershipBenefitDraft>,
): void {
  emit(
    "update:modelValue",
    props.modelValue.map((row, rowIndex) =>
      rowIndex === index ? { ...row, ...patch } : row,
    ),
  );
}

function addLevel(): void {
  if (!nextAvailableLevel.value) return;
  emit("update:modelValue", [
    ...props.modelValue,
    createGroupMembershipBenefitDraft(
      nextAvailableLevel.value.id,
      props.baseRateMultiplier,
    ),
  ]);
}

function removeLevel(index: number): void {
  emit(
    "update:modelValue",
    props.modelValue.filter((_, rowIndex) => rowIndex !== index),
  );
}

function updateLevel(index: number, event: Event): void {
  replaceRow(index, {
    level_id: Number((event.target as HTMLSelectElement).value),
  });
}

function updateRate(index: number, event: Event): void {
  const value = (event.target as HTMLInputElement).value;
  replaceRow(index, { rate_multiplier: value === "" ? null : Number(value) });
}
</script>
