import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { CreateSnippet } from './CreateSnippet'
import * as api from '../api/snippets'

vi.mock('../api/snippets', () => ({
  createSnippet: vi.fn(),
}))

const mockCreateSnippet = vi.fn<typeof api.createSnippet>(api.createSnippet)
vi.mocked(api.createSnippet).mockImplementation(mockCreateSnippet)

describe('CreateSnippet', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    mockCreateSnippet.mockReset()
  })

  it('renders textarea and create button', () => {
    render(<CreateSnippet />)
    expect(screen.getByRole('textbox')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create snippet/i })).toBeInTheDocument()
  })

  it('shows validation error on empty submission', async () => {
    render(<CreateSnippet />)
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.click(btn)
    await waitFor(() => {
      expect(screen.getByText(/content cannot be empty/i)).toBeInTheDocument()
    })
  })

  it('shows loading state while submitting', async () => {
    mockCreateSnippet.mockImplementationOnce(() => new Promise((resolve) =>
      setTimeout(() => resolve({ id: 1, slug: 'abc123', content: 'hello', created_at: new Date().toISOString() }), 50)))

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'hello world' } })
    fireEvent.click(btn)

    expect(btn).toBeDisabled()
    expect(screen.getByText(/creating\.\.\./i)).toBeInTheDocument()

    await waitFor(() => {
      expect(mockCreateSnippet).toHaveBeenCalledWith({ content: 'hello world' })
    })
  })

  it('displays generated URL on success', async () => {
    mockCreateSnippet.mockResolvedValueOnce({ id: 1, slug: 'a8f31c', content: 'hello world', created_at: new Date().toISOString() })

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'hello world' } })
    fireEvent.click(btn)

    await waitFor(() => {
      const link = screen.getByRole('link', { name: /\/s\/a8f31c/i })
      expect(link).toBeInTheDocument()
      expect(link).toHaveAttribute('href', '/s/a8f31c')
    })
  })

  it('displays error on API failure', async () => {
    mockCreateSnippet.mockRejectedValueOnce(new Error('Network error'))

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'hello' } })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(screen.getByText('Network error')).toBeInTheDocument()
    })
  })

  it('clears textarea after successful creation', async () => {
    mockCreateSnippet.mockResolvedValueOnce({ id: 1, slug: 'xyz', content: 'test', created_at: new Date().toISOString() })

    render(<CreateSnippet />)
    const textarea = screen.getByRole('textbox')
    const btn = screen.getByRole('button', { name: /create snippet/i })
    fireEvent.change(textarea, { target: { value: 'my code' } })
    fireEvent.click(btn)

    await waitFor(() => {
      expect(textarea).toHaveValue('')
    })
  })
})
