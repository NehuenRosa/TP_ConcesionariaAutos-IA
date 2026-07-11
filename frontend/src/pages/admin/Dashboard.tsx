import { useState, useEffect } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts'
import api from '../../services/api'
import type { DashboardMetrics } from '../../types'

const COLORS = ['#22c55e', '#eab308', '#ef4444']

export function AdminDashboard() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/admin/dashboard').then(({ data }) => {
      setMetrics(data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-center py-20 text-gray-500">Cargando...</div>
  if (!metrics) return <div className="text-center py-20 text-gray-500">Error al cargar métricas</div>

  const vehicleData = [
    { name: 'Disponibles', value: metrics.vehiculos.disponible },
    { name: 'Reservados', value: metrics.vehiculos.reservado },
    { name: 'Vendidos', value: metrics.vehiculos.vendido },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Dashboard Administrativo</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Vehículos totales</h3>
          <p className="text-3xl font-bold">{metrics.vehiculos.total}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Disponibles</h3>
          <p className="text-3xl font-bold text-green-600">{metrics.vehiculos.disponible}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Reservados</h3>
          <p className="text-3xl font-bold text-yellow-600">{metrics.vehiculos.reservado}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Vendidos</h3>
          <p className="text-3xl font-bold text-red-600">{metrics.vehiculos.vendido}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Consultas totales</h3>
          <p className="text-3xl font-bold">{metrics.consultas.total}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Test Drives agendados</h3>
          <p className="text-3xl font-bold">{metrics.test_drives_pendientes}</p>
          <p className="text-sm text-gray-400">({metrics.test_drives_totales} totales)</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm">Reservas activas</h3>
          <p className="text-3xl font-bold text-blue-600">{metrics.reservas_activas}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm mb-4">Vehículos por estado</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={vehicleData}>
              <XAxis dataKey="name" />
              <YAxis allowDecimals={false} />
              <Tooltip />
              <Bar dataKey="value" fill="#3b82f6" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm mb-4">Distribución de vehículos</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie data={vehicleData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={100} label>
                {vehicleData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i]} />
                ))}
              </Pie>
              <Tooltip />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>

      <div className="mt-8 flex gap-4">
        <a href="/admin/vehiculos" className="bg-blue-600 text-white px-6 py-3 rounded hover:bg-blue-700">
          Gestionar vehículos
        </a>
      </div>
    </div>
  )
}
