import { useState } from 'react'
import api from '../services/api'

export function Chatbot() {
  const [isOpen, setIsOpen] = useState(false)
  const [messages, setMessages] = useState<{ role: string; text: string }[]>([
    { role: 'bot', text: '¡Hola! Soy el asistente de la concesionaria. ¿En qué puedo ayudarte?' },
  ])
  const [question, setQuestion] = useState('')
  const [loading, setLoading] = useState(false)

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
        className="fixed bottom-4 right-4 bg-blue-600 text-white w-14 h-14 rounded-full shadow-lg hover:bg-blue-700 z-50 text-2xl"
      >
        {isOpen ? '✕' : '💬'}
      </button>

      {isOpen && (
        <div className="fixed bottom-20 right-4 w-96 bg-white rounded-lg shadow-xl border z-50 flex flex-col max-h-[500px]">
          <div className="bg-blue-600 text-white p-3 rounded-t-lg font-semibold">
            Asistente Concesionaria
          </div>
          <div className="flex-1 overflow-y-auto p-3 space-y-2 min-h-[300px]">
            {messages.map((msg, i) => (
              <div
                key={i}
                className={`p-2 rounded-lg max-w-[80%] ${
                  msg.role === 'user'
                    ? 'bg-blue-100 ml-auto'
                    : 'bg-gray-100'
                }`}
              >
                {msg.text}
              </div>
            ))}
            {loading && <div className="text-gray-400 italic">Escribiendo...</div>}
          </div>
          <form onSubmit={handleSubmit} className="border-t p-2 flex gap-2">
            <input
              type="text"
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              placeholder="Haz una pregunta..."
              className="flex-1 border rounded px-3 py-2 focus:outline-none"
            />
            <button
              type="submit"
              disabled={loading}
              className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
            >
              Enviar
            </button>
          </form>
        </div>
      )}
    </>
  )
}
