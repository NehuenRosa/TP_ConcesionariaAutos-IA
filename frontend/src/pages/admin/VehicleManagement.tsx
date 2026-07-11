import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { vehicleService } from '../../services/vehicleService'
import type { Vehicle } from '../../types'

export function VehicleManagement() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [loading, setLoading] = useState(true)

  const loadVehicles = () => {
    setLoading(true)
    vehicleService.list({ page_size: 100 }).then((res) => {
      setVehicles(res.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { loadVehicles() }, [])

  const handleDelete = async (id: number) => {
    if (!window.confirm('¿Estás seguro de eliminar este vehículo?')) return
    try {
      await vehicleService.delete(id)
      loadVehicles()
    } catch {
      alert('Error al eliminar el vehículo')
    }
  }

  const formatPrice = (price: number) =>
    new Intl.NumberFormat('es-AR', { style: 'currency', currency: 'ARS' }).format(price)

  const statusStyles: Record<string, string> = {
    disponible: 'badge-green',
    reservado: 'badge-yellow',
    vendido: 'badge-red',
  }

  if (loading) return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        <div className="flex justify-between">
          <div className="h-8 bg-gray-200 rounded-xl w-64" />
          <div className="h-10 bg-gray-200 rounded-xl w-40" />
        </div>
        <div className="bg-gray-200 rounded-2xl h-64" />
      </div>
    </div>
  )

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-fade-in">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Gestión de Vehículos</h1>
          <p className="text-gray-500 mt-1">{vehicles.length} vehículo{vehicles.length !== 1 ? 's' : ''} registrado{vehicles.length !== 1 ? 's' : ''}</p>
        </div>
        <Link to="/admin/vehiculos/nuevo" className="btn-primary text-sm flex items-center gap-2">
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          Nuevo vehículo
        </Link>
      </div>

      <div className="card overflow-hidden">
        {vehicles.length === 0 ? (
          <div className="text-center py-16">
            <p className="text-gray-500 font-medium">No hay vehículos registrados</p>
            <Link to="/admin/vehiculos/nuevo" className="text-brand-600 text-sm font-medium mt-2 inline-block hover:text-brand-700">
              Crear el primer vehículo
            </Link>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50/50">
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Vehículo</th>
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Año</th>
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Precio</th>
                  <th className="px-5 py-3.5 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">Estado</th>
                  <th className="px-5 py-3.5 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider">Acciones</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {vehicles.map((v) => (
                  <tr key={v.id} className="hover:bg-brand-50/50 transition-colors duration-150">
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-gray-100 overflow-hidden flex-shrink-0">
                          <img
                            src={v.images?.[0] || 'https://placehold.co/40x40'}
                            alt=""
                            className="w-full h-full object-cover"
                          />
                        </div>
                        <div>
                          <p className="font-medium text-gray-900">{v.brand} {v.model}</p>
                          <p className="text-xs text-gray-400 capitalize">{v.vehicle_type} • {v.fuel}</p>
                        </div>
                      </div>
                    </td>
                    <td className="px-5 py-4 text-sm text-gray-600">{v.year}</td>
                    <td className="px-5 py-4 text-sm font-medium text-gray-900">{formatPrice(v.price)}</td>
                    <td className="px-5 py-4">
                      <span className={statusStyles[v.status] || 'badge-gray'}>{v.status}</span>
                    </td>
                    <td className="px-5 py-4 text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Link
                          to={`/admin/vehiculos/${v.id}`}
                          className="px-3 py-1.5 text-sm font-medium text-brand-600 hover:bg-brand-50 rounded-lg transition-colors"
                        >
                          Editar
                        </Link>
                        <button
                          onClick={() => handleDelete(v.id)}
                          className="px-3 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        >
                          Eliminar
                        </button>
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
