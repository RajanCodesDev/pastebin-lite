import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { ViewSnippet } from './ViewSnippet'
import * as api from '../api/snippets'

const mockGetSnippet = vi.fn<typeof api.getSnippet>(api.getSnippet)
const mockDeleteSnippet = vi.fn<typeof api.deleteSnippet>(api.deleteSnippet)

vi.mock('../api/snippets', () => ({
  getSnippet: (...args: unknown[]) => mockGetSnippet(...args),
  deleteSnippet: (...args: unknown[]) => mockDeleteSnippet(...args),
}))

function renderWithRouter() {
  return render(
    <MemoryRouter initialEntries={['/s/test-slug']}>
      <Routes>
        <Route path="/s/:slug" element={<ViewSnippet />} />
        <Route path="/" element={<span>Home</span>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('ViewSnippet', () => {
  beforeEach(() => {
    vi.resetAllMocks()
  })

  it('shows loading state', () => {
    mockGetSnippet.mockImplementation(() => new Promise(() => {}))
    renderWithRouter()
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })

  it('displays snippet content on success', async () => {
    const snippet = { id: 1, slug: 'test-slug', content: 'console.log("hi")', created_at: new Date().toISOString() }
    mockGetSnippet.mockResolvedValueOnce(snippet)

    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByText('console.log("hi")')).toBeInTheDocument()
    })
    expect(screen.getByText(/created:/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
  })

  it('shows 404 when snippet not found', async () => {
    const err = Object.assign(new Error('Not found'), { status: 404 })
    mockGetSnippet.mockRejectedValueOnce(err)

    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByText('Snippet not found')).toBeInTheDocument()
    })
  })

  it('shows error on API failure', async () => {
    mockGetSnippet.mockRejectedValueOnce(new Error('Server error'))

    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByText('Server error')).toBeInTheDocument()
    })
  })

  it('calls DELETE endpoint on delete click', async () => {
    mockGetSnippet.mockResolvedValueOnce({ id: 1, slug: 'test-slug', content: 'code', created_at: new Date().toISOString() })
    mockDeleteSnippet.mockResolvedValueOnce(undefined)

    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /delete/i }))
    expect(mockDeleteSnippet).toHaveBeenCalledWith('test-slug')
  })

  it('redirects to / after successful delete', async () => {
    mockGetSnippet.mockResolvedValueOnce({ id: 1, slug: 'test-slug', content: 'code', created_at: new Date().toISOString() })
    mockDeleteSnippet.mockResolvedValueOnce(undefined)

    renderWithRouter()
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument()
    })
    fireEvent.click(screen.getByRole('button', { name: /delete/i }))
    await waitFor(() => {
      expect(screen.getByText('Home')).toBeInTheDocument()
    })
  })
})
