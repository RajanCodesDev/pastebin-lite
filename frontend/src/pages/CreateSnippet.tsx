import { useState } from 'react'
import { createSnippet } from '../api/snippets'

type State = 'idle' | 'loading' | 'success' | 'error'

export function CreateSnippet() {
  const [content, setContent] = useState('')
  const [state, setState] = useState<State>('idle')
  const [result, setResult] = useState<{ slug: string; url: string } | null>(null)
  const [errorMsg, setErrorMsg] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim()) {
      setErrorMsg('Content cannot be empty')
      setState('error')
      return
    }
    setState('loading')
    setErrorMsg('')
    try {
      const snippet = await createSnippet({ content })
      const url = `${window.location.origin}/s/${snippet.slug}`
      setResult({ slug: snippet.slug, url })
      setContent('')
      setState('success')
    } catch (err) {
      setErrorMsg(err instanceof Error ? err.message : 'Failed to create snippet')
      setState('error')
    }
  }

  return (
    <div className="container">
      <h1>Pastebin Lite</h1>
      <form onSubmit={handleSubmit}>
        <textarea
          value={content}
          onChange={(e) => { setContent(e.target.value); if (state === 'error') setState('idle') }}
          placeholder="Paste your code or text here..."
          rows={10}
          disabled={state === 'loading'}
        />
        <br />
        <button type="submit" disabled={state === 'loading'}>
          {state === 'loading' ? 'Creating...' : 'Create Snippet'}
        </button>
      </form>

      {state === 'success' && result && (
        <div className="result">
          <p>Snippet created!</p>
          <p>
            <a href={result.url}>{result.url}</a>
          </p>
        </div>
      )}

      {state === 'error' && errorMsg && (
        <div className="error">
          <p>{errorMsg}</p>
        </div>
      )}
    </div>
  )
}
