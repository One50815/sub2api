package service

import "testing"

func TestAPIKeyServiceCanUserBindGroupInternalOmnioPro(t *testing.T) {
	service := &APIKeyService{}
	user := &User{ID: 1}
	standard := &Group{ID: 10, SubscriptionType: SubscriptionTypeStandard}

	tests := []struct {
		name        string
		group       Group
		subscribed  map[int64]bool
		benefits    map[int64]EffectiveMembershipGroupBenefit
		proOnly     map[int64]struct{}
		wantAllowed bool
	}{
		{
			name:        "public group remains available",
			group:       *standard,
			wantAllowed: true,
		},
		{
			name:    "pro only group is hidden without active benefit",
			group:   *standard,
			proOnly: map[int64]struct{}{standard.ID: {}},
		},
		{
			name:    "pro only group is allowed by active level benefit",
			group:   *standard,
			proOnly: map[int64]struct{}{standard.ID: {}},
			benefits: map[int64]EffectiveMembershipGroupBenefit{
				standard.ID: {MembershipGroupBenefit: MembershipGroupBenefit{AllowAccess: true}},
			},
			wantAllowed: true,
		},
		{
			name:  "pro benefit can grant an otherwise exclusive group",
			group: Group{ID: 20, IsExclusive: true, SubscriptionType: SubscriptionTypeStandard},
			benefits: map[int64]EffectiveMembershipGroupBenefit{
				20: {MembershipGroupBenefit: MembershipGroupBenefit{AllowAccess: true}},
			},
			wantAllowed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := service.canUserBindGroupInternal(user, &tt.group, tt.subscribed, tt.benefits, tt.proOnly)
			if allowed != tt.wantAllowed {
				t.Fatalf("allowed = %v, want %v", allowed, tt.wantAllowed)
			}
		})
	}
}
