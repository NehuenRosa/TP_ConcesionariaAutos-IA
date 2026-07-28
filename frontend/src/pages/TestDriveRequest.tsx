import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import api from '../services/api'
import { vehicleService } from '../services/vehicleService'
import type { Vehicle } from '../types'

export function TestDriveRequest() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [vehicle, setVehicle] = useState<Vehicle | null>(null)
  const [date, setDate] = useState('')
  const [time, setTime] = useState('')
  const [notes, setNotes] = useState('')
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (id) {
      vehicleService.getById(Number(id)).then(setVehicle).catch(() => navigate('/catalogo'))
    }
  }, [id])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!date || !time) {
      setError('Seleccioná fecha y horario')
      return
    }
    setSubmitting(true)
    setError('')
    try {
      await api.post('/test-drives', {
        vehicle_id: Number(id),
        scheduled_at: `${date}T${time}`,
        notes,
      })
      navigate('/catalogo', { state: { message: 'Turno solicitado con éxito' } })
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al solicitar el turno')
    } finally {
      setSubmitting(false)
    }
  }

  if (!vehicle) return (
    <div className="max-w-lg mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        <div className="h-6 bg-gray-200 rounded-lg w-32" />
        <div className="h-12 bg-gray-200 rounded-xl" />
        <div className="h-12 bg-gray-200 rounded-xl" />
      </div>
    </div>
  )

  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  const minDate = tomorrow.toISOString().split('T')[0]

  return (
    <div className="max-w-lg mx-auto px-4 sm:px-6 py-8 animate-fade-in">
      <Link to={`/vehiculos/${vehicle.id}`} className="inline-flex items-center gap-1.5 text-sm font-medium text-gray-500 hover:text-brand-600 mb-6 transition-colors">
        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        Volver al vehículo
      </Link>

      <div className="card p-8">
        <div className="flex items-center gap-3 mb-8">
          <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
            <svg className="w-6 h-6 text-purple-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Solicitar Test Drive</h1>
            <p className="text-sm text-gray-500">{vehicle.brand} {vehicle.model} ({vehicle.year})</p>
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

        <form onSubmit={handleSubmit} className="space-y-5">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Fecha</label>
            <input type="date" value={date} onChange={(e) => setDate(e.target.value)} min={minDate} required className="input-field" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Horario</label>
            <input type="time" value={time} onChange={(e) => setTime(e.target.value)} required className="input-field" />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Notas (opcional)</label>
            <textarea placeholder="Comentarios adicionales..." value={notes} onChange={(e) => setNotes(e.target.value)} rows={3} className="input-field resize-none" />
          </div>
          <div className="flex items-center gap-3 pt-2">
            <button type="submit" disabled={submitting} className="btn-primary flex items-center gap-2">
              {submitting ? (
                <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
              ) : (
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              )}
              Solicitar turno
            </button>
            <Link to={`/vehiculos/${vehicle.id}`} className="btn-secondary text-sm">Cancelar</Link>
          </div>
        </form>
      </div>
    </div>
  )
}
