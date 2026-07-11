import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../services/api'

export function TestDriveRequest() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [scheduledAt, setScheduledAt] = useState('')
  const [notes, setNotes] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.post('/test-drives', {
        vehicle_id: Number(id),
        scheduled_at: new Date(scheduledAt).toISOString(),
        notes,
      })
      setSuccess(true)
      setTimeout(() => navigate('/catalogo'), 2000)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al solicitar el turno')
    }
  }

  return (
    <div className="max-w-lg mx-auto mt-20 p-6 bg-white rounded-lg shadow">
      <h1 className="text-2xl font-bold mb-6">Solicitar Test Drive</h1>
      {success && (
        <div className="bg-green-100 text-green-700 p-3 rounded mb-4">
          Turno solicitado correctamente. Esperá la confirmación del vendedor.
        </div>
      )}
      {error && <div className="bg-red-100 text-red-700 p-3 rounded mb-4">{error}</div>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm text-gray-600 mb-1">Fecha y hora del test drive</label>
          <input
            type="datetime-local"
            value={scheduledAt}
            onChange={(e) => setScheduledAt(e.target.value)}
            required
            className="w-full border rounded px-3 py-2"
          />
        </div>
        <textarea
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Comentarios adicionales (opcional)"
          rows={3}
          className="w-full border rounded px-3 py-2"
        />
        <button type="submit" className="w-full bg-green-600 text-white py-2 rounded hover:bg-green-700">
          Solicitar turno
        </button>
      </form>
    </div>
  )
}
