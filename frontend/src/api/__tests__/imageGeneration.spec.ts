import { afterEach, describe, expect, it, vi } from 'vitest'
import {
  ImageGenerationAPIError,
  imageOutputSource,
  isImageGenerationModel,
  mergeImageGenerationModels,
  submitImageGeneration,
} from '../imageGeneration'

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })
}

describe('imageGeneration API', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('recognizes the image model families supported by the gateway', () => {
    expect(isImageGenerationModel('gpt-image-2')).toBe(true)
    expect(isImageGenerationModel('grok-imagine')).toBe(true)
    expect(isImageGenerationModel('grok-imagine-image')).toBe(true)
    expect(isImageGenerationModel('gpt-5.2')).toBe(false)
    expect(isImageGenerationModel('grok-4')).toBe(false)
  })

  it('keeps preferred image models first while preserving discovered models', () => {
    expect(mergeImageGenerationModels(
      ['gpt-image-2', 'gpt-image-1.5', 'gpt-image-1'],
      ['gpt-image-1.5', 'gpt-image-2', 'gpt-image-custom'],
    )).toEqual([
      'gpt-image-2',
      'gpt-image-1.5',
      'gpt-image-1',
      'gpt-image-custom',
    ])
  })

  it('returns an async task when the task endpoint accepts the request', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({
      id: 'imgtask_123',
      task_id: 'imgtask_123',
      object: 'image.task',
      status: 'processing',
      created_at: 1,
      expires_at: 2,
    }, 202))
    vi.stubGlobal('fetch', fetchMock)

    const submission = await submitImageGeneration('sk-test', {
      model: 'gpt-image-2',
      prompt: 'A paper city',
    })

    expect(submission.mode).toBe('async')
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(String(fetchMock.mock.calls[0][0])).toContain('/v1/images/generations/async')
  })

  it('falls back to the synchronous endpoint only when async tasks are disabled', async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(jsonResponse({
        error: { type: 'not_found_error', message: 'async image tasks are not enabled' },
      }, 404))
      .mockResolvedValueOnce(jsonResponse({ data: [{ b64_json: 'YWJj' }] }))
    vi.stubGlobal('fetch', fetchMock)

    const submission = await submitImageGeneration('sk-test', {
      model: 'gpt-image-2',
      prompt: 'A paper city',
    })

    expect(submission.mode).toBe('sync')
    expect(fetchMock).toHaveBeenCalledTimes(2)
    expect(String(fetchMock.mock.calls[1][0])).toContain('/v1/images/generations')
    expect(String(fetchMock.mock.calls[1][0])).not.toContain('/async')
  })

  it('uses the synchronous route directly when async is not preferred', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse({ data: [{ url: 'https://example.com/image.png' }] }))
    vi.stubGlobal('fetch', fetchMock)

    const submission = await submitImageGeneration('sk-test', {
      model: 'art-route',
      prompt: 'A paper city',
    }, { preferAsync: false })

    expect(submission.mode).toBe('sync')
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(String(fetchMock.mock.calls[0][0])).toContain('/v1/images/generations')
    expect(String(fetchMock.mock.calls[0][0])).not.toContain('/async')
  })

  it('does not hide unrelated async endpoint errors', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(jsonResponse({
      error: { type: 'not_found_error', message: 'Images API is not supported for this platform' },
    }, 404)))

    await expect(submitImageGeneration('sk-test', {
      model: 'gpt-image-2',
      prompt: 'A paper city',
    })).rejects.toMatchObject<ImageGenerationAPIError>({
      status: 404,
      message: 'Images API is not supported for this platform',
    })
  })

  it('turns base64 output into a displayable data URL', () => {
    expect(imageOutputSource({ b64_json: 'YWJj' }, 'png')).toBe('data:image/png;base64,YWJj')
  })
})
