import { useState, useEffect } from 'react'
import api from '../../services/api'
import type { TestDrive } from '../../types'

const statusStyles: Record<string, string> = {
  pendiente: 'badge-yellow',
  confirmado: 'badge-green',
  cancelado: 'badge-red',
  completado: 'badge-blue',
}

export function TestDriveManagement() {
  const [testDrives, setTestDrives] = useState<TestDrive[]>([])
  const [loading, setLoading] = useState(true)

  const load = () => {
    setLoading(true)
    api.get('/test-drives').then(({ data }) => {
      setTestDrives(data.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const handleStatus = async (id: number, status: string) => {
    try {
      await api.patch(`/test-drives/${id}/status`, { status })
      load()
    } catch {
      alert('Error al actualizar el estado')
    }
  }

  const formatDate = (dateStr: string) =>
    new Intl.DateTimeFormat('es-AR', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(dateStr))

  if (loading) return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        <div className="h-8 bg-gray-200 rounded-xl w-64" />
        <div className="bg-gray-200 rounded-2xl h-64" />
      </div>
    </div>
  )

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-fade-in">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Turnos de Test Drive</h1>
        <p className="text-gray-500 mt-1">{testDrives.length} turno{testDrives.length !== 1 ? 's' : ''}</p>
      </div>

      <div className="card overflow-hidden">
        {testDrives.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-gray-500 font-medium">No hay turnos registrados</p>
          </div>
        ) : (
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
                {testDrives.map((td) => (
                  <tr key={td.id} className="hover:bg-purple-50/50 transition-colors duration-150">
                    <td className="px-5 py-4">
                      <p className="font-medium text-gray-900">{td.client?.name || `Cliente #${td.client_id}`}</p>
                      <p className="text-xs text-gray-400">{td.client?.email}</p>
                    </td>
                    <td className="px-5 py-4">
                      <p className="text-sm text-gray-900">{td.vehicle?.brand} {td.vehicle?.model}</p>
                      <p className="text-xs text-gray-400">{td.vehicle?.year}</p>
                    </td>
                    <td className="px-5 py-4 text-sm text-gray-600">{formatDate(td.scheduled_at)}</td>
                    <td className="px-5 py-4">
                      <span className={statusStyles[td.status] || 'badge-gray'}>{td.status}</span>
                    </td>
                    <td className="px-5 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        {td.status === 'pendiente' && (
                          <>
                            <button onClick={() => handleStatus(td.id, 'confirmado')} className="px-3 py-1.5 text-sm font-medium text-emerald-600 hover:bg-emerald-50 rounded-lg transition-colors">Confirmar</button>
                            <button onClick={() => handleStatus(td.id, 'cancelado')} className="px-3 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors">Cancelar</button>
                          </>
                        )}
                        {td.status === 'confirmado' && (
                          <>
                            <button onClick={() => handleStatus(td.id, 'completado')} className="px-3 py-1.5 text-sm font-medium text-blue-600 hover:bg-blue-50 rounded-lg transition-colors">Completar</button>
                            <button onClick={() => handleStatus(td.id, 'cancelado')} className="px-3 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors">Cancelar</button>
                          </>
                        )}
                        {td.status === 'completado' && <span className="text-xs text-gray-400">Finalizado</span>}
                        {td.status === 'cancelado' && <span className="text-xs text-gray-400">Cancelado</span>}
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}
