import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../services/api'

export function ContactSeller() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.post('/consultations', { vehicle_id: Number(id), message })
      setSuccess(true)
      setTimeout(() => navigate('/catalogo'), 2000)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al enviar la consulta')
    }
  }

  return (
    <div className="max-w-lg mx-auto mt-20 p-6 bg-white rounded-lg shadow">
      <h1 className="text-2xl font-bold mb-6">Consultar sobre este vehículo</h1>
      {success && (
        <div className="bg-green-100 text-green-700 p-3 rounded mb-4">
          Consulta enviada correctamente. Te responderemos a la brevedad.
        </div>
      )}
      {error && <div className="bg-red-100 text-red-700 p-3 rounded mb-4">{error}</div>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <textarea
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Escribe tu consulta o solicitud de cotización..."
          required
          rows={5}
          className="w-full border rounded px-3 py-2"
        />
        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700">
          Enviar consulta
        </button>
      </form>
    </div>
  )
}
