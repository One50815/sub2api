# Sub2API Frontend Design QA

## Comparison Target

- Source visual truth: `E:\最终版\前端`
- Rendered implementation: `http://127.0.0.1:4173`
- Desktop viewport: `1280x720`, dark and light themes
- Mobile viewport: `390x844`, dark theme
- States: login, register, forgot password, authenticated user API keys, authenticated admin users, empty tables

## Evidence

Full-view comparisons:

- Login, dark desktop: `E:\最终版\sub2api\docs\design-qa\comparison-login-dark.png`
- Admin users, dark desktop: `E:\最终版\sub2api\docs\design-qa\comparison-admin-users-dark.png`

Implementation views:

- Register, dark desktop: `E:\最终版\sub2api\docs\design-qa\implementation-register-dark.png`
- Forgot password, dark desktop: `E:\最终版\sub2api\docs\design-qa\implementation-forgot-password-dark.png`
- Login, mobile: `E:\最终版\sub2api\docs\design-qa\implementation-login-mobile.png`
- Admin users, mobile: `E:\最终版\sub2api\docs\design-qa\implementation-admin-users-mobile.png`
- User API keys, light desktop: `E:\最终版\sub2api\docs\design-qa\implementation-user-keys-light.png`
- User API keys, dark desktop: `E:\最终版\sub2api\docs\design-qa\implementation-user-keys-dark.png`

Focused region comparisons:

- Login form controls and card: `E:\最终版\sub2api\docs\design-qa\comparison-login-form-focused.png`
- Admin heading, filters, actions, and table header: `E:\最终版\sub2api\docs\design-qa\comparison-admin-toolbar-focused.png`

## Findings

- No actionable P0, P1, or P2 visual mismatches remain.
- Fonts and typography: Public Sans, optical weights, line heights, non-negative letter spacing, heading hierarchy, wrapping, and compact UI labels follow the source. Mobile text remains inside its controls.
- Spacing and layout rhythm: the 48px header, inset 192px sidebar, rounded main workspace, content padding, toolbar density, card geometry, table rhythm, and auth split composition match the source system. The 390px checks have no horizontal overflow.
- Colors and visual tokens: charcoal dark surfaces, white light surfaces, blue primary actions, restrained borders, status colors, focus rings, shadows, and auth purple/blue/cyan lighting are consistently mapped across shared components.
- Image quality and asset fidelity: the gateway background, brand mark, platform marks, and icons render sharply. Existing source assets and the repository icon system are used; there are no placeholder images or emoji substitutions.
- Copy and content: authentication copy follows the reference tone while retaining Sub2API naming. Admin and user screens preserve Sub2API-specific fields, permissions, and routes where the reference product has different business content.

Expected product differences:

- `Omnio` is replaced by the configured Sub2API site name and logo.
- Reference-only navigation and New API fields are not copied. Sub2API routes, role checks, filters, columns, balance data, and API contracts remain authoritative.
- Password reset, OAuth, agreement, Turnstile, and registration controls remain feature-gated by existing public settings.

## Patches Made Since Previous QA Pass

- Moved page titles and primary actions into the content workspace.
- Matched the reference 48px header, inset sidebar, rounded workspace, and compact navigation.
- Centered wide-table empty states in the visible viewport.
- Removed the duplicate empty-state create action on admin users.
- Displayed role and status filters by default for first-time users.
- Kept zero-result pagination visible and disabled correctly.
- Tightened login typography and site-name-aware copy.
- Added consistent light/dark/mobile handling across authentication, user, and admin views.

## Implementation Checklist

- [x] Desktop dark comparison
- [x] Desktop light comparison
- [x] Mobile responsive comparison
- [x] Authentication pages
- [x] User application shell and tables
- [x] Admin application shell and tables
- [x] Empty states and pagination
- [x] Typography, spacing, colors, assets, and copy
- [x] Lint, typecheck, tests, and production build

## Follow-up Polish

No blocking follow-up. Future upstream-only fields should be placed into the existing design system instead of restoring legacy Vue styling.

final result: passed
