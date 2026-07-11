import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import api from '../services/api'

export function ReserveVehicle() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [notes, setNotes] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.post('/reservations', { vehicle_id: Number(id), notes })
      setSuccess(true)
      setTimeout(() => navigate('/catalogo'), 2000)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al reservar el vehículo')
    }
  }

  return (
    <div className="max-w-lg mx-auto mt-20 p-6 bg-white rounded-lg shadow">
      <h1 className="text-2xl font-bold mb-6">Reservar Vehículo</h1>
      {success && (
        <div className="bg-green-100 text-green-700 p-3 rounded mb-4">
          Vehículo reservado correctamente. El vendedor se comunicará para coordinar la venta.
        </div>
      )}
      {error && <div className="bg-red-100 text-red-700 p-3 rounded mb-4">{error}</div>}
      <form onSubmit={handleSubmit} className="space-y-4">
        <textarea
          value={notes}
          onChange={(e) => setNotes(e.target.value)}
          placeholder="Comentarios adicionales (opcional)"
          rows={3}
          className="w-full border rounded px-3 py-2"
        />
        <button type="submit" className="w-full bg-yellow-600 text-white py-2 rounded hover:bg-yellow-700">
          Reservar ahora
        </button>
      </form>
    </div>
  )
}
