import { describe, expect, it } from "vitest";
import type { MembershipLevel } from "@/types/membership";
import {
  buildGroupMembershipBenefitChangeSet,
  cloneGroupMembershipBenefitDrafts,
  createGroupMembershipBenefitDraft,
  groupMembershipBenefitDrafts,
  validateGroupMembershipBenefitDrafts,
} from "../groupsMembershipBenefits";

const levels: MembershipLevel[] = [
  {
    id: 19,
    name: "Omnio Pro",
    slug: "pro",
    description: "",
    rank: 1,
    badge_color: "#2563eb",
    concurrency_bonus: 0,
    priority_support: false,
    active: true,
    sort_order: 1,
    group_benefits: [
      {
        id: 101,
        level_id: 19,
        group_id: 7,
        allow_access: true,
        pro_only: false,
        rate_multiplier: 0.8,
        rpm_limit: 120,
        daily_free_usd: 2,
        monthly_free_usd: 20,
      },
    ],
  },
  {
    id: 20,
    name: "Omnio Pro Max",
    slug: "pro-max",
    description: "",
    rank: 2,
    badge_color: "#16a34a",
    concurrency_bonus: 0,
    priority_support: false,
    active: true,
    sort_order: 2,
    group_benefits: [],
  },
];

describe("group membership benefit editor helpers", () => {
  it("loads only the selected group while preserving non-rate entitlements", () => {
    expect(groupMembershipBenefitDrafts(levels, 7)).toEqual([
      {
        level_id: 19,
        allow_access: true,
        pro_only: false,
        rate_multiplier: 0.8,
        rpm_limit: 120,
        daily_free_usd: 2,
        monthly_free_usd: 20,
      },
    ]);
  });

  it("builds minimal upserts and deletions for independent levels", () => {
    const original = groupMembershipBenefitDrafts(levels, 7);
    const current = cloneGroupMembershipBenefitDrafts(original);
    current[0].rate_multiplier = 0.6;
    current.push(createGroupMembershipBenefitDraft(20, 0.5));

    const changes = buildGroupMembershipBenefitChangeSet(7, original, current);

    expect(changes.deletedLevelIds).toEqual([]);
    expect(changes.upserts).toEqual([
      expect.objectContaining({
        level_id: 19,
        group_id: 7,
        rate_multiplier: 0.6,
        rpm_limit: 120,
        daily_free_usd: 2,
        monthly_free_usd: 20,
      }),
      expect.objectContaining({
        level_id: 20,
        group_id: 7,
        rate_multiplier: 0.5,
      }),
    ]);

    const removed = buildGroupMembershipBenefitChangeSet(7, current, [current[1]]);
    expect(removed.deletedLevelIds).toEqual([19]);
    expect(removed.upserts).toEqual([]);
  });

  it("rejects duplicate levels and negative multipliers", () => {
    const first = createGroupMembershipBenefitDraft(19, 0.8);
    expect(
      validateGroupMembershipBenefitDrafts([first, { ...first }], levels),
    ).toBe("duplicate-level");
    expect(
      validateGroupMembershipBenefitDrafts(
        [{ ...first, rate_multiplier: -0.1 }],
        levels,
      ),
    ).toBe("invalid-rate");
  });
});
