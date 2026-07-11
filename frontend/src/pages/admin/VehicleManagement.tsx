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
    if (!confirm('¿Estás seguro de eliminar este vehículo?')) return
    try {
      await vehicleService.delete(id)
      loadVehicles()
    } catch {
      alert('Error al eliminar el vehículo')
    }
  }

  if (loading) return <div className="text-center py-20 text-gray-500">Cargando...</div>

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Gestión de Vehículos</h1>
        <Link
          to="/admin/vehiculos/nuevo"
          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          + Nuevo vehículo
        </Link>
      </div>

      <div className="bg-white rounded-lg shadow overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Marca</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Modelo</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Año</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Precio</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Estado</th>
              <th className="px-4 py-3 text-left text-sm font-medium text-gray-500">Acciones</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {vehicles.map((v) => (
              <tr key={v.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">{v.brand}</td>
                <td className="px-4 py-3">{v.model}</td>
                <td className="px-4 py-3">{v.year}</td>
                <td className="px-4 py-3">${v.price.toLocaleString()}</td>
                <td className="px-4 py-3">
                  <span className={`px-2 py-1 rounded text-xs ${
                    v.status === 'disponible' ? 'bg-green-100 text-green-700' :
                    v.status === 'reservado' ? 'bg-yellow-100 text-yellow-700' : 'bg-red-100 text-red-700'
                  }`}>{v.status}</span>
                </td>
                <td className="px-4 py-3 flex gap-2">
                  <Link to={`/admin/vehiculos/${v.id}`} className="text-blue-600 hover:underline text-sm">
                    Editar
                  </Link>
                  <button onClick={() => handleDelete(v.id)} className="text-red-600 hover:underline text-sm">
                    Eliminar
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
