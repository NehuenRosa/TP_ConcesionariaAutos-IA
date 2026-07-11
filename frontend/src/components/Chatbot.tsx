import { useState, useRef, useEffect } from 'react'
import api from '../services/api'

export function Chatbot() {
  const [isOpen, setIsOpen] = useState(false)
  const [messages, setMessages] = useState<{ role: string; text: string }[]>([
    { role: 'bot', text: '¡Hola! Soy el asistente de AutoPrime. ¿En qué puedo ayudarte?' },
  ])
  const [question, setQuestion] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!question.trim()) return

    setMessages((prev) => [...prev, { role: 'user', text: question }])
    setLoading(true)
    setQuestion('')

    try {
      const { data } = await api.post('/chatbot/ask', { question })
      setMessages((prev) => [...prev, { role: 'bot', text: data.answer }])
    } catch {
      setMessages((prev) => [
        ...prev,
        { role: 'bot', text: 'Lo siento, hubo un error al procesar tu consulta.' },
      ])
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed bottom-6 right-6 w-14 h-14 bg-gradient-to-br from-brand-500 to-brand-700 text-white rounded-2xl shadow-xl shadow-brand-500/30 hover:shadow-brand-500/50 hover:scale-105 active:scale-95 transition-all duration-200 z-50 flex items-center justify-center"
      >
        {isOpen ? (
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        ) : (
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
          </svg>
        )}
      </button>

      {isOpen && (
        <div className="fixed bottom-24 right-6 w-[380px] bg-white rounded-2xl shadow-2xl border border-gray-100 z-50 flex flex-col max-h-[600px] animate-scale-in overflow-hidden">
          <div className="bg-gradient-to-r from-brand-600 to-brand-500 p-4 flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-white/20 flex items-center justify-center backdrop-blur-sm">
              <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
              </svg>
            </div>
            <div>
              <h3 className="font-semibold text-white text-sm">Asistente AutoPrime</h3>
              <p className="text-white/70 text-xs">Online • Respuesta inmediata</p>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto p-4 space-y-3 min-h-[350px] max-h-[400px] bg-gray-50/50">
            {messages.map((msg, i) => (
              <div
                key={i}
                className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'} animate-slide-up`}
              >
                <div
                  className={`px-4 py-2.5 rounded-2xl max-w-[85%] text-sm leading-relaxed ${
                    msg.role === 'user'
                      ? 'bg-brand-600 text-white rounded-br-md'
                      : 'bg-white text-gray-700 shadow-sm border border-gray-100 rounded-bl-md'
                  }`}
                >
                  {msg.text}
                </div>
              </div>
            ))}
            {loading && (
              <div className="flex justify-start">
                <div className="bg-white px-4 py-3 rounded-2xl shadow-sm border border-gray-100 rounded-bl-md">
                  <div className="flex gap-1.5">
                    <div className="w-2 h-2 bg-gray-300 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                    <div className="w-2 h-2 bg-gray-300 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                    <div className="w-2 h-2 bg-gray-300 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                  </div>
                </div>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>

          <form onSubmit={handleSubmit} className="border-t border-gray-100 p-3 flex gap-2 bg-white">
            <input
              type="text"
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder="Escribí tu mensaje..."
              className="flex-1 bg-gray-50 border-0 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500/30 transition-all"
            />
            <button
              type="submit"
              disabled={loading || !question.trim()}
              className="w-10 h-10 bg-brand-600 text-white rounded-xl hover:bg-brand-700 disabled:opacity-40 disabled:cursor-not-allowed transition-all flex items-center justify-center flex-shrink-0"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 19V5m0 0l-7 7m7-7l7 7" />
              </svg>
            </button>
          </form>
        </div>
      )}
    </>
  )
}
