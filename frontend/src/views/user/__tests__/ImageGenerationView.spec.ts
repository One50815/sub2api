import { readFileSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'

const currentDirectory = dirname(fileURLToPath(import.meta.url))
const viewSource = readFileSync(resolve(currentDirectory, '../ImageGenerationView.vue'), 'utf8')
const routerSource = readFileSync(resolve(currentDirectory, '../../../router/index.ts'), 'utf8')
const sidebarSource = readFileSync(resolve(currentDirectory, '../../../components/layout/AppSidebar.vue'), 'utf8')

describe('ImageGenerationView integration contract', () => {
  it('is a native page backed by the existing image gateway and balance refresh', () => {
    expect(viewSource).toContain('submitImageGeneration')
    expect(viewSource).toContain("response_format: 'b64_json'")
    expect(viewSource).toContain('authStore.refreshUser()')
    expect(viewSource).not.toContain('<iframe')
  })

  it('is reachable from the authenticated router and sidebar', () => {
    expect(routerSource).toContain("path: '/image-generation'")
    expect(routerSource).toContain("component: () => import('@/views/user/ImageGenerationView.vue')")
    expect(sidebarSource).toContain("path: '/image-generation'")
    expect(sidebarSource).toContain("t('nav.imageGeneration')")
  })
})
