import { useState, useEffect } from 'react'
import api from '../../services/api'
import type { TestDrive } from '../../types'

export function TestDriveManagement() {
  const [testDrives, setTestDrives] = useState<TestDrive[]>([])
  const [loading, setLoading] = useState(true)

  const load = () => {
    api.get('/test-drives').then(({ data }) => {
      setTestDrives(data.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const handleStatusChange = async (id: number, status: string) => {
    await api.patch(`/test-drives/${id}/status`, { status })
    load()
  }

  if (loading) return <div className="text-center py-20 text-gray-500">Cargando...</div>

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Gestión de Test Drives</h1>
      {testDrives.length === 0 ? (
        <p className="text-gray-500">No hay turnos de test drive registrados.</p>
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
              {testDrives.map((td) => (
                <tr key={td.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3">{td.client?.name}</td>
                  <td className="px-4 py-3">{td.vehicle?.brand} {td.vehicle?.model}</td>
                  <td className="px-4 py-3">{new Date(td.scheduled_at).toLocaleString()}</td>
                  <td className="px-4 py-3">
                    <span className={`px-2 py-1 rounded text-xs ${
                      td.status === 'pendiente' ? 'bg-yellow-100 text-yellow-700' :
                      td.status === 'confirmado' ? 'bg-green-100 text-green-700' :
                      td.status === 'completado' ? 'bg-blue-100 text-blue-700' : 'bg-red-100 text-red-700'
                    }`}>{td.status}</span>
                  </td>
                  <td className="px-4 py-3 flex gap-2">
                    {td.status === 'pendiente' && (
                      <>
                        <button onClick={() => handleStatusChange(td.id, 'confirmado')} className="bg-green-600 text-white px-2 py-1 rounded text-xs">Confirmar</button>
                        <button onClick={() => handleStatusChange(td.id, 'cancelado')} className="bg-red-600 text-white px-2 py-1 rounded text-xs">Cancelar</button>
                      </>
                    )}
                    {td.status === 'confirmado' && (
                      <button onClick={() => handleStatusChange(td.id, 'completado')} className="bg-blue-600 text-white px-2 py-1 rounded text-xs">Completar</button>
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
