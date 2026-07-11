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

  if (loading) return <div className="text-center py-20 text-gray-500">Cargando...</div>

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Gestión de Reservas</h1>
      {reservations.length === 0 ? (
        <p className="text-gray-500">No hay reservas registradas.</p>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Cliente</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Vehículo</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Fecha</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Estado</th>
                <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Acciones</th>
              </tr>
            </thead>
            <tbody className="divide-y">
              {reservations.map((r) => (
                <tr key={r.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3">{r.client?.name}</td>
                  <td className="px-4 py-3">{r.vehicle?.brand} {r.vehicle?.model}</td>
                  <td className="px-4 py-3">{new Date(r.created_at).toLocaleDateString()}</td>
                  <td className="px-4 py-3">
                    <span className={`px-2 py-1 rounded text-xs ${
                      r.status === 'activa' ? 'bg-yellow-100 text-yellow-700' :
                      r.status === 'confirmada' ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                    }`}>{r.status}</span>
                  </td>
                  <td className="px-4 py-3 flex gap-2">
                    {r.status === 'activa' && (
                      <>
                        <button onClick={() => handleConfirm(r.id)} className="bg-green-600 text-white px-2 py-1 rounded text-xs">Confirmar venta</button>
                        <button onClick={() => handleCancel(r.id)} className="bg-red-600 text-white px-2 py-1 rounded text-xs">Cancelar</button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
