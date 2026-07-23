import type {
  MembershipGroupBenefit,
  MembershipLevel,
} from "@/types/membership";

export interface GroupMembershipBenefitDraft {
  level_id: number;
  allow_access: boolean;
  pro_only: boolean;
  rate_multiplier: number | null;
  rpm_limit: number | null;
  daily_free_usd: number | null;
  monthly_free_usd: number | null;
}

export interface GroupMembershipBenefitChangeSet {
  upserts: MembershipGroupBenefit[];
  deletedLevelIds: number[];
}

export type GroupMembershipBenefitValidationError =
  | "duplicate-level"
  | "invalid-level"
  | "invalid-rate";

export function createGroupMembershipBenefitDraft(
  levelId: number,
  baseRateMultiplier: number,
): GroupMembershipBenefitDraft {
  return {
    level_id: levelId,
    allow_access: true,
    pro_only: false,
    rate_multiplier: baseRateMultiplier,
    rpm_limit: null,
    daily_free_usd: null,
    monthly_free_usd: null,
  };
}

export function groupMembershipBenefitDrafts(
  levels: MembershipLevel[],
  groupId: number,
): GroupMembershipBenefitDraft[] {
  return levels.flatMap((level) => {
    const benefit = level.group_benefits?.find((item) => item.group_id === groupId);
    if (!benefit) return [];
    return [
      {
        level_id: level.id,
        allow_access: benefit.allow_access,
        pro_only: benefit.pro_only,
        rate_multiplier: benefit.rate_multiplier ?? null,
        rpm_limit: benefit.rpm_limit ?? null,
        daily_free_usd: benefit.daily_free_usd ?? null,
        monthly_free_usd: benefit.monthly_free_usd ?? null,
      },
    ];
  });
}

export function cloneGroupMembershipBenefitDrafts(
  rows: GroupMembershipBenefitDraft[],
): GroupMembershipBenefitDraft[] {
  return rows.map((row) => ({ ...row }));
}

export function validateGroupMembershipBenefitDrafts(
  rows: GroupMembershipBenefitDraft[],
  levels: MembershipLevel[],
): GroupMembershipBenefitValidationError | null {
  const validLevelIds = new Set(levels.map((level) => level.id));
  const selectedLevelIds = new Set<number>();

  for (const row of rows) {
    if (!validLevelIds.has(row.level_id)) return "invalid-level";
    if (selectedLevelIds.has(row.level_id)) return "duplicate-level";
    selectedLevelIds.add(row.level_id);
    if (
      row.rate_multiplier !== null &&
      (!Number.isFinite(row.rate_multiplier) || row.rate_multiplier < 0)
    ) {
      return "invalid-rate";
    }
  }

  return null;
}

function normalizedBenefit(
  groupId: number,
  row: GroupMembershipBenefitDraft,
): MembershipGroupBenefit {
  return {
    id: 0,
    level_id: row.level_id,
    group_id: groupId,
    allow_access: row.pro_only ? true : row.allow_access,
    pro_only: row.pro_only,
    rate_multiplier: row.rate_multiplier,
    rpm_limit: row.rpm_limit,
    daily_free_usd: row.daily_free_usd,
    monthly_free_usd: row.monthly_free_usd,
  };
}

function benefitsEqual(
  left: MembershipGroupBenefit,
  right: MembershipGroupBenefit,
): boolean {
  return (
    left.level_id === right.level_id &&
    left.group_id === right.group_id &&
    left.allow_access === right.allow_access &&
    left.pro_only === right.pro_only &&
    left.rate_multiplier === right.rate_multiplier &&
    left.rpm_limit === right.rpm_limit &&
    left.daily_free_usd === right.daily_free_usd &&
    left.monthly_free_usd === right.monthly_free_usd
  );
}

export function buildGroupMembershipBenefitChangeSet(
  groupId: number,
  originalRows: GroupMembershipBenefitDraft[],
  currentRows: GroupMembershipBenefitDraft[],
): GroupMembershipBenefitChangeSet {
  const originalByLevel = new Map(
    originalRows.map((row) => [row.level_id, normalizedBenefit(groupId, row)]),
  );
  const currentByLevel = new Map(
    currentRows.map((row) => [row.level_id, normalizedBenefit(groupId, row)]),
  );

  const upserts = [...currentByLevel.values()].filter((benefit) => {
    const original = originalByLevel.get(benefit.level_id);
    return !original || !benefitsEqual(original, benefit);
  });
  const deletedLevelIds = [...originalByLevel.keys()].filter(
    (levelId) => !currentByLevel.has(levelId),
  );

  return { upserts, deletedLevelIds };
}
