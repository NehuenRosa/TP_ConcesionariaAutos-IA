import { useState, useEffect } from 'react'
import api from '../../services/api'
import type { Reservation } from '../../types'

export function ReservationManagement() {
  const [reservations, setReservations] = useState<Reservation[]>([])
  const [loading, setLoading] = useState(true)

  const load = () => {
    api.get('/reservations').then(({ data }) => {
      setReservations(data.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const handleConfirm = async (id: number) => {
    await api.post(`/reservations/${id}/confirm`)
    load()
  }

  const handleCancel = async (id: number) => {
    await api.post(`/reservations/${id}/cancel`)
    load()
  }

  const statusStyles: Record<string, string> = {
    activa: 'badge-yellow',
    confirmada: 'badge-green',
    cancelada: 'badge-red',
  }

  if (loading) return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        <div className="h-8 bg-gray-200 rounded-xl w-64" />
        <div className="h-64 bg-gray-200 rounded-2xl" />
      </div>
    </div>
  )

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-fade-in">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Gestión de Reservas</h1>
        <p className="text-gray-500 mt-1">{reservations.length} reserva{reservations.length !== 1 ? 's' : ''}</p>
      </div>

      {reservations.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="w-16 h-16 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
            </svg>
          </div>
          <p className="text-gray-500 font-medium">No hay reservas registradas</p>
        </div>
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50/50">
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Cliente</th>
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Vehículo</th>
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Fecha</th>
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Estado</th>
                  <th className="px-5 py-3.5 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider">Acciones</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {reservations.map((r) => (
                  <tr key={r.id} className="hover:bg-brand-50/50 transition-colors duration-150">
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-2">
                        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-brand-400 to-brand-600 flex items-center justify-center">
                          <span className="text-xs font-bold text-white">{r.client?.name?.charAt(0).toUpperCase() || '?'}</span>
                        </div>
                        <span className="text-sm font-medium text-gray-900">{r.client?.name}</span>
                      </div>
                    </td>
                    <td className="px-5 py-4 text-sm text-gray-600">{r.vehicle?.brand} {r.vehicle?.model}</td>
                    <td className="px-5 py-4 text-sm text-gray-600">{new Date(r.created_at).toLocaleDateString('es-AR')}</td>
                    <td className="px-5 py-4">
                      <span className={statusStyles[r.status] || 'badge-gray'}>{r.status}</span>
                    </td>
                    <td className="px-5 py-4 text-right">
                      {r.status === 'activa' && (
                        <div className="flex items-center justify-end gap-2">
                          <button onClick={() => handleConfirm(r.id)} className="btn-success btn-sm flex items-center gap-1">
                            <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            Confirmar venta
                          </button>
                          <button onClick={() => handleCancel(r.id)} className="btn-danger btn-sm flex items-center gap-1">
                            <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                            Cancelar
                          </button>
                        </div>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}
