import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import api from '../services/api'
import type { Consultation } from '../types'

export function MyConsultations() {
  const [consultations, setConsultations] = useState<Consultation[]>([])
  const [selected, setSelected] = useState<Consultation | null>(null)
  const [responseText, setResponseText] = useState('')
  const [loading, setLoading] = useState(true)

  const load = () => {
    api.get('/consultations/mine').then(({ data }) => {
      setConsultations(data.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const sorted = [...consultations].sort((a, b) => {
    if (a.has_unread_for_client && !b.has_unread_for_client) return -1
    if (!a.has_unread_for_client && b.has_unread_for_client) return 1
    if (a.status === 'pendiente' && b.status !== 'pendiente') return -1
    if (a.status !== 'pendiente' && b.status === 'pendiente') return 1
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })

  const handleSendResponse = async () => {
    if (!selected || !responseText.trim()) return
    await api.post(`/consultations/${selected.id}/responses`, { message: responseText })
    setResponseText('')
    const { data } = await api.get(`/consultations/${selected.id}`)
    setSelected(data)
    load()
  }

  const handleDelete = async (id: number) => {
    if (!confirm('¿Eliminar esta consulta?')) return
    await api.delete(`/consultations/${id}`)
    if (selected?.id === id) setSelected(null)
    load()
  }

  const statusStyles: Record<string, string> = {
    pendiente: 'badge-pendiente',
    en_conversacion: 'badge-conversacion',
    cerrada: 'badge-cerrada',
  }

  const statusLabels: Record<string, string> = {
    pendiente: 'Pendiente',
    en_conversacion: 'En conversación',
    cerrada: 'Cerrada',
  }

  if (loading) return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="animate-pulse space-y-4">
        <div className="h-8 bg-gray-200 rounded-xl w-64" />
        <div className="grid grid-cols-2 gap-6">
          <div className="h-96 bg-gray-200 rounded-2xl" />
          <div className="h-96 bg-gray-200 rounded-2xl" />
        </div>
      </div>
    </div>
  )

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 animate-fade-in">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Mis Consultas</h1>
        <p className="text-gray-500 mt-1">{consultations.length} consulta{consultations.length !== 1 ? 's' : ''}</p>
      </div>

      <div className="grid lg:grid-cols-5 gap-6">
        <div className="lg:col-span-2 card overflow-hidden">
          {consultations.length === 0 ? (
            <div className="p-8 text-center">
              <div className="w-14 h-14 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-3">
                <svg className="w-7 h-7 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                </svg>
              </div>
              <p className="text-gray-500 font-medium">No tenés consultas</p>
              <p className="text-gray-400 text-sm mt-1">
                <Link to="/catalogo" className="text-brand-600 hover:text-brand-700">Visitá el catálogo</Link> para consultar sobre un vehículo
              </p>
            </div>
          ) : (
            <div className="divide-y divide-gray-50 max-h-[600px] overflow-y-auto">
              {sorted.map((c) => (
                <button
                  key={c.id}
                  onClick={() => {
                    setSelected(c)
                    api.get(`/consultations/${c.id}`).then(({ data }) => setSelected(data))
                  }}
                  className={`w-full text-left p-4 transition-colors duration-150 ${
                    selected?.id === c.id ? 'bg-brand-50' : 'hover:bg-gray-50'
                  }`}
                >
                  <div className="flex items-start justify-between gap-3 mb-2">
                    <div className="flex items-center gap-2 min-w-0">
                      <div className="relative">
                        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-gray-400 to-gray-600 flex items-center justify-center flex-shrink-0">
                          <span className="text-xs font-bold text-white">
                            {c.vehicle?.brand?.charAt(0).toUpperCase() || '?'}
                          </span>
                        </div>
                        {c.has_unread_for_client && (
                          <span className="absolute -top-0.5 -right-0.5 w-3 h-3 bg-red-500 border-2 border-white rounded-full" />
                        )}
                      </div>
                      <div className="min-w-0">
                        <p className="text-sm font-medium text-gray-900 truncate">{c.vehicle?.brand} {c.vehicle?.model}</p>
                        <p className="text-xs text-gray-400 truncate">{c.vehicle?.year}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className={statusStyles[c.status] || 'badge-gray'}>{statusLabels[c.status] || c.status}</span>
                      <button
                        onClick={(e) => { e.stopPropagation(); handleDelete(c.id) }}
                        className="p-1 text-gray-300 hover:text-red-500 transition-colors"
                        title="Eliminar"
                      >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                          <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>
                  </div>
                  <p className="text-sm text-gray-500 line-clamp-2 pl-10">{c.message}</p>
                </button>
              ))}
            </div>
          )}
        </div>

        <div className="lg:col-span-3">
          {selected ? (
            <div className="card p-6 animate-fade-in">
              <div className="flex items-start justify-between mb-6">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-gray-400 to-gray-600 flex items-center justify-center">
                    <span className="text-lg font-bold text-white">
                      {selected.vehicle?.brand?.charAt(0).toUpperCase() || '?'}
                    </span>
                  </div>
                  <div>
                    <h2 className="text-lg font-semibold text-gray-900">{selected.vehicle?.brand} {selected.vehicle?.model}</h2>
                    <p className="text-sm text-gray-400">{selected.vehicle?.year}</p>
                  </div>
                </div>
                <span className={statusStyles[selected.status] || 'badge-gray'}>{statusLabels[selected.status] || selected.status}</span>
              </div>

              <div className="mb-6">
                <p className="text-xs text-gray-400 uppercase tracking-wider font-medium mb-2">Tu consulta inicial</p>
                <div className="bg-brand-50 rounded-xl p-4">
                  <p className="text-sm text-gray-700">{selected.message}</p>
                </div>
              </div>

              {selected.responses && selected.responses.length > 0 && (
                <div className="mb-6">
                  <p className="text-xs text-gray-400 uppercase tracking-wider font-medium mb-3">Conversación</p>
                  <div className="space-y-3">
                    {selected.responses.map((r) => (
                      <div key={r.id} className="bg-gray-50 rounded-xl p-4">
                        <div className="flex items-center gap-2 mb-1.5">
                          <span className="text-xs font-medium text-gray-500">{r.user?.name}</span>
                          <span className="text-xs text-gray-300">•</span>
                          <span className="text-xs text-gray-400">{new Date(r.created_at).toLocaleString('es-AR')}</span>
                        </div>
                        <p className="text-sm text-gray-700">{r.message}</p>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {selected.status !== 'cerrada' && (
                <div>
                  <p className="text-xs text-gray-400 uppercase tracking-wider font-medium mb-2">Responder</p>
                  <div className="flex gap-2">
                    <textarea
                      value={responseText}
                      onChange={(e) => setResponseText(e.target.value)}
                      placeholder="Escribí tu mensaje..."
                      rows={3}
                      className="input-field resize-none flex-1"
                    />
                    <button
                      onClick={handleSendResponse}
                      disabled={!responseText.trim()}
                      className="btn-primary btn-sm self-end flex items-center gap-1.5 disabled:opacity-50"
                    >
                      <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 19V5m0 0l-7 7m7-7l7 7" />
                      </svg>
                      Enviar
                    </button>
                  </div>
                </div>
              )}

              {selected.status === 'cerrada' && (
                <div className="bg-gray-50 rounded-xl p-4 text-center">
                  <p className="text-sm text-gray-500">Esta consulta está cerrada.</p>
                </div>
              )}
            </div>
          ) : (
            <div className="card p-8 text-center h-full flex items-center justify-center">
              <div>
                <div className="w-16 h-16 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
                  <svg className="w-8 h-8 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
                  </svg>
                </div>
                <p className="text-gray-500 font-medium">Seleccioná una consulta</p>
                <p className="text-gray-400 text-sm mt-1">Elegí una del listado para ver su detalle</p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
