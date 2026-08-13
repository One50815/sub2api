# Homepage Design QA

## Scope

- Homepage only: `frontend/src/views/HomeView.vue`, `frontend/src/styles/home-v2.css`, and landing-page locale copy.
- Login, user console, admin backend, routes, APIs, authentication, billing, permissions, and theme infrastructure remain unchanged.
- Existing `home_content` URL and HTML overrides remain supported.

## Design Read

- Product: public AI gateway homepage for everyday AI users and developers.
- Direction: privacy-first, technically credible, restrained, and transparent about model choice and cost.
- Reference language: warm neutral surface, dark floating navigation, strong typography, and layered product interaction inspired by CollectiveOS.
- Dials: `DESIGN_VARIANCE 7`, `MOTION_INTENSITY 5`, `VISUAL_DENSITY 4`.
- Foundation: existing Vue 3 application, native CSS, existing icon components, and existing theme tokens.

## Visual Truth

- Reference URL: `https://collectiveos.vercel.app/#`
- Source visual truth: `artifacts/design/collectiveos-reference-top-1280x720.png`
- Primary implementation screenshot: `artifacts/design/omnio-home-desktop-light-1280x720.png`
- Full-view comparison: `artifacts/design/homepage-reference-comparison-1280x720.png`
- Focused hero comparison: `artifacts/design/homepage-reference-focused-comparison-1280x420.png`
- Local implementation: `http://127.0.0.1:4173/home`

## Viewports And States

- Chinese desktop light, `1280x720`: `artifacts/design/omnio-home-desktop-light-1280x720.png`
- English desktop light, `1280x720`: `artifacts/design/omnio-home-desktop-en-light-1280x720.png`
- English desktop dark, `1280x720`: `artifacts/design/omnio-home-desktop-dark-1280x720.png`
- Chinese mobile light, `382x844`: `artifacts/design/omnio-home-mobile-light-382x844.png`
- English mobile light, `382x844`: `artifacts/design/omnio-home-mobile-en-light-382x844.png`
- English mobile dark, `382x844`: `artifacts/design/omnio-home-mobile-en-dark-382x844.png`
- Mobile workflow: `artifacts/design/omnio-home-mobile-workflow-382x844.png`
- Mobile FAQ and ending: `artifacts/design/omnio-home-mobile-faq-382x844.png`

## Comparison Evidence

The implementation intentionally borrows the reference's floating navigation, neutral palette, bold type, restrained radii, and first-viewport product signal. It does not clone the centered CollectiveOS composition or its business content. Omnio uses an asymmetric hero so the privacy, model, and cost proof can remain visible beside the headline.

The focused comparison verifies headline scale, CTA hierarchy, line length, and proof-rail density. The full-view comparison verifies navigation proportion, first-viewport rhythm, and the visible handoff into the interactive task launcher.

## Findings

- No actionable P0, P1, or P2 visual findings remain.
- Typography: Public Sans keeps a consistent optical weight across Chinese and English. Hero copy remains two lines at desktop and two controlled lines at `382px` mobile.
- Spacing: sections use distinct layout families, consistent vertical rhythm, and explicit single-column mobile collapse rules.
- Color: one emerald accent is used across both themes. Light and dark modes preserve hierarchy without section-level theme flips.
- Assets: existing product marks and model icons render sharply in both themes. No placeholder image or custom decorative SVG replacement was introduced.
- Copy: privacy, exact model selection, request records, and cost transparency use one concrete product voice in both locales.
- Responsiveness: document width remains within the viewport at `1280x720` and `382x844`. The task selector is intentionally horizontally scrollable on mobile without creating document-level overflow.
- Interaction: task selection updates the launcher content; copy shows a visible success state; FAQ buttons expand and collapse with correct `aria-expanded`; mobile navigation opens and exposes all expected routes.
- Accessibility: focus-visible states, semantic headings, navigation labels, reduced-motion handling, and practical mobile tap targets remain present.

## Patches Made

- Rebuilt the workflow, model directory, transparency evidence, FAQ, and final CTA with distinct layout families.
- Rewrote Chinese and English copy around privacy, exact model selection, request traceability, and clear cost records.
- Removed the earlier fallback language that conflicted with the no-hidden-substitution promise.
- Tuned English and Chinese display type independently for desktop and mobile wrapping.
- Preserved the current hero while improving the rest of the page's information rhythm.
- Corrected the homepage language menu so normal, hover, and selected options retain readable contrast in light and dark modes.
- Preserved homepage URL sanitization, public overrides, authentication behavior, and existing routes.

## Verification

- Targeted ESLint: passed.
- Locale compile and collision tests: 8 passed.
- Production build: passed, 1026 modules transformed.
- Desktop and mobile visual checks: passed.
- Chinese and English layout checks: passed.
- Light and dark theme checks: passed.
- Task launcher, copy feedback, FAQ, and mobile navigation checks: passed.
- Page copy contains no em dash or en dash characters.
- Build warnings are existing project warnings for Browserslist age, mixed dynamic/static imports, chunk size, and a Node child-process deprecation.

final result: passed
