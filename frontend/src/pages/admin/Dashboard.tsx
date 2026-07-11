import { useState, useEffect } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, PieChart, Pie, Cell, Legend } from 'recharts'
import { Link } from 'react-router-dom'
import api from '../../services/api'
import type { DashboardMetrics } from '../../types'

const COLORS = ['#22c55e', '#eab308', '#ef4444']
const PIE_COLORS = ['#6366f1', '#f59e0b', '#ef4444']

export function AdminDashboard() {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get('/admin/dashboard').then(({ data }) => {
      setMetrics(data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [])

  if (loading) return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-6">
        <div className="h-8 bg-gray-200 rounded-xl w-64" />
        <div className="grid grid-cols-4 gap-6">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-28 bg-gray-200 rounded-2xl" />
          ))}
        </div>
      </div>
    </div>
  )

  if (!metrics) return (
    <div className="text-center py-20">
      <p className="text-gray-500 font-medium">Error al cargar métricas</p>
    </div>
  )

  const vehicleData = [
    { name: 'Disponibles', value: metrics.vehiculos.disponible },
    { name: 'Reservados', value: metrics.vehiculos.reservado },
    { name: 'Vendidos', value: metrics.vehiculos.vendido },
  ]

  const cards = [
    { label: 'Vehículos totales', value: metrics.vehiculos.total, color: 'text-gray-900', icon: 'M9 17a2 2 0 11-4 0 2 2 0 014 0zM19 17a2 2 0 11-4 0 2 2 0 014 0z', bg: 'bg-gray-100', iconBg: 'text-gray-600' },
    { label: 'Disponibles', value: metrics.vehiculos.disponible, color: 'text-emerald-600', icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z', bg: 'bg-emerald-50', iconBg: 'text-emerald-600' },
    { label: 'Reservados', value: metrics.vehiculos.reservado, color: 'text-amber-600', icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z', bg: 'bg-amber-50', iconBg: 'text-amber-600' },
    { label: 'Vendidos', value: metrics.vehiculos.vendido, color: 'text-red-600', icon: 'M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6', bg: 'bg-red-50', iconBg: 'text-red-600' },
  ]

  const secondaryCards = [
    { label: 'Consultas totales', value: metrics.consultas.total, icon: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z', bg: 'bg-blue-50', iconBg: 'text-blue-600' },
    { label: 'Test Drives pendientes', value: metrics.test_drives_pendientes, sub: `${metrics.test_drives_totales} totales`, icon: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z', bg: 'bg-purple-50', iconBg: 'text-purple-600' },
    { label: 'Reservas activas', value: metrics.reservas_activas, icon: 'M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z', bg: 'bg-amber-50', iconBg: 'text-amber-600' },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-fade-in">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="text-gray-500 mt-1">Panel de control administrativo</p>
        </div>
        <Link to="/admin/vehiculos" className="btn-primary text-sm flex items-center gap-2">
          <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
          </svg>
          Gestionar vehículos
        </Link>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5 mb-8">
        {cards.map((card) => (
          <div key={card.label} className="card p-6 animate-slide-up">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-gray-500 mb-1">{card.label}</p>
                <p className={`text-3xl font-bold ${card.color}`}>{card.value}</p>
              </div>
              <div className={`w-12 h-12 ${card.bg} rounded-xl flex items-center justify-center`}>
                <svg className={`w-6 h-6 ${card.iconBg}`} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d={card.icon} />
                </svg>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-5 mb-8">
        {secondaryCards.map((card) => (
          <div key={card.label} className="card p-6 animate-slide-up">
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-gray-500 mb-1">{card.label}</p>
                <p className="text-3xl font-bold text-gray-900">{card.value}</p>
                {card.sub && <p className="text-sm text-gray-400 mt-1">{card.sub}</p>}
              </div>
              <div className={`w-12 h-12 ${card.bg} rounded-xl flex items-center justify-center`}>
                <svg className={`w-6 h-6 ${card.iconBg}`} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d={card.icon} />
                </svg>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card p-6 animate-slide-up">
          <h3 className="text-sm font-medium text-gray-500 mb-6">Vehículos por estado</h3>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={vehicleData}>
              <XAxis dataKey="name" tick={{ fontSize: 12 }} axisLine={false} tickLine={false} />
              <YAxis allowDecimals={false} axisLine={false} tickLine={false} tick={{ fontSize: 12 }} />
              <Tooltip
                contentStyle={{ borderRadius: '12px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
              />
              <Bar dataKey="value" radius={[8, 8, 0, 0]} barSize={50}>
                {vehicleData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>

        <div className="card p-6 animate-slide-up">
          <h3 className="text-sm font-medium text-gray-500 mb-6">Distribución de vehículos</h3>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie data={vehicleData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={100} label={({ name, percent }: { name?: string; percent?: number }) => `${name} ${((percent ?? 0) * 100).toFixed(0)}%`}>
                {vehicleData.map((_, i) => (
                  <Cell key={i} fill={PIE_COLORS[i]} />
                ))}
              </Pie>
              <Tooltip
                contentStyle={{ borderRadius: '12px', border: '1px solid #e5e7eb', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
              />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  )
}
