export const mockCreateSnippet = vi.fn()
export const mockGetSnippet = vi.fn()
export const mockDeleteSnippet = vi.fn()

vi.mock('../api/snippets', () => ({
  createSnippet: (...args: unknown[]) => mockCreateSnippet(...args),
  getSnippet: (...args: unknown[]) => mockGetSnippet(...args),
  deleteSnippet: (...args: unknown[]) => mockDeleteSnippet(...args),
}))
