import { useState } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import api from '../services/api'

export function ReserveVehicle() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [notes, setNotes] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await api.post('/reservations', { vehicle_id: Number(id), notes })
      setSuccess(true)
      setTimeout(() => navigate('/catalogo'), 2000)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al reservar el vehículo')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-[80vh] flex items-center justify-center px-4">
      <div className="w-full max-w-lg animate-slide-up">
        <Link to={`/vehiculos/${id}`} className="inline-flex items-center gap-1.5 text-sm font-medium text-gray-500 hover:text-brand-600 mb-6 transition-colors">
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
          Volver al vehículo
        </Link>

        {success ? (
          <div className="card p-8 text-center">
            <div className="w-16 h-16 bg-emerald-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <h2 className="text-xl font-bold text-gray-900 mb-2">¡Vehículo reservado!</h2>
            <p className="text-gray-500">El vendedor se comunicará para coordinar la venta. Redirigiendo al catálogo...</p>
          </div>
        ) : (
          <div className="card p-8">
            <div className="flex items-center gap-3 mb-6">
              <div className="w-12 h-12 bg-amber-100 rounded-xl flex items-center justify-center">
                <svg className="w-6 h-6 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                </svg>
              </div>
              <div>
                <h1 className="text-xl font-bold text-gray-900">Reservar Vehículo</h1>
                <p className="text-sm text-gray-500">Asegurá tu vehículo con una reserva</p>
              </div>
            </div>

            {error && (
              <div className="flex items-center gap-2 bg-red-50 text-red-700 p-3 rounded-xl mb-6 text-sm border border-red-100">
                <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1.5">Comentarios (opcional)</label>
                <textarea
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="Contanos si tenés alguna pregunta antes de la compra..."
                  rows={3}
                  className="input-field resize-none"
                />
              </div>
              <button
                type="submit"
                disabled={loading}
                className="btn-warning w-full flex items-center justify-center gap-2"
              >
                {loading ? (
                  <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                ) : (
                  <>
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                    </svg>
                    Reservar ahora
                  </>
                )}
              </button>
            </form>
          </div>
        )}
      </div>
    </div>
  )
}
