import { useState, useEffect } from 'react'
import api from '../../services/api'
import type { Consultation } from '../../types'

export function ConsultationInbox() {
  const [consultations, setConsultations] = useState<Consultation[]>([])
  const [selected, setSelected] = useState<Consultation | null>(null)
  const [responseText, setResponseText] = useState('')
  const [loading, setLoading] = useState(true)

  const load = () => {
    api.get('/consultations').then(({ data }) => {
      setConsultations(data.data)
      setLoading(false)
    }).catch(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const handleStatusChange = async (id: number, status: string) => {
    await api.patch(`/consultations/${id}/status`, { status })
    load()
  }

  const handleSendResponse = async () => {
    if (!selected || !responseText.trim()) return
    await api.post(`/consultations/${selected.id}/responses`, { message: responseText })
    setResponseText('')
    const { data } = await api.get(`/consultations/${selected.id}`)
    setSelected(data)
    load()
  }

  const statusStyles: Record<string, string> = {
    pendiente: 'badge-yellow',
    en_conversacion: 'badge-blue',
    cerrada: 'badge-gray',
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
        <h1 className="text-3xl font-bold text-gray-900">Bandeja de Consultas</h1>
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
              <p className="text-gray-500 font-medium">No hay consultas</p>
              <p className="text-gray-400 text-sm mt-1">Las consultas de los clientes aparecerán aquí</p>
            </div>
          ) : (
            <div className="divide-y divide-gray-50 max-h-[600px] overflow-y-auto">
              {consultations.map((c) => (
                <button
                  key={c.id}
                  onClick={() => setSelected(c)}
                  className={`w-full text-left p-4 transition-colors duration-150 ${
                    selected?.id === c.id ? 'bg-brand-50' : 'hover:bg-gray-50'
                  }`}
                >
                  <div className="flex items-start justify-between gap-3 mb-2">
                    <div className="flex items-center gap-2 min-w-0">
                      <div className="w-8 h-8 rounded-full bg-gradient-to-br from-brand-400 to-brand-600 flex items-center justify-center flex-shrink-0">
                        <span className="text-xs font-bold text-white">
                          {c.client?.name?.charAt(0).toUpperCase() || '?'}
                        </span>
                      </div>
                      <div className="min-w-0">
                        <p className="text-sm font-medium text-gray-900 truncate">{c.client?.name || 'Cliente'}</p>
                        <p className="text-xs text-gray-400 truncate">{c.vehicle?.brand} {c.vehicle?.model}</p>
                      </div>
                    </div>
                    <span className={statusStyles[c.status] || 'badge-gray'}>{c.status.replace('_', ' ')}</span>
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
                  <div className="w-12 h-12 rounded-full bg-gradient-to-br from-brand-400 to-brand-600 flex items-center justify-center">
                    <span className="text-lg font-bold text-white">
                      {selected.client?.name?.charAt(0).toUpperCase() || '?'}
                    </span>
                  </div>
                  <div>
                    <h2 className="text-lg font-semibold text-gray-900">{selected.client?.name}</h2>
                    <p className="text-sm text-gray-400">{selected.client?.email}</p>
                  </div>
                </div>
                <span className={statusStyles[selected.status] || 'badge-gray'}>{selected.status.replace('_', ' ')}</span>
              </div>

              <div className="bg-gray-50 rounded-xl p-4 mb-6">
                <p className="text-xs text-gray-400 uppercase tracking-wider font-medium mb-1">Vehículo</p>
                <p className="text-sm font-medium text-gray-800">
                  {selected.vehicle?.brand} {selected.vehicle?.model} ({selected.vehicle?.year})
                </p>
              </div>

              <div className="mb-6">
                <p className="text-xs text-gray-400 uppercase tracking-wider font-medium mb-2">Mensaje del cliente</p>
                <div className="bg-brand-50 rounded-xl p-4">
                  <p className="text-sm text-gray-700">{selected.message}</p>
                </div>
              </div>

              {selected.status !== 'cerrada' && (
                <div className="flex gap-2 mb-6">
                  {selected.status === 'pendiente' && (
                    <button onClick={() => handleStatusChange(selected.id, 'en_conversacion')} className="btn-primary btn-sm flex items-center gap-1.5">
                      <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                      </svg>
                      Tomar consulta
                    </button>
                  )}
                  <button onClick={() => handleStatusChange(selected.id, 'cerrada')} className="btn-secondary btn-sm flex items-center gap-1.5">
                    <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    Cerrar consulta
                  </button>
                </div>
              )}

              {selected.responses && selected.responses.length > 0 && (
                <div className="mb-6">
                  <p className="text-xs text-gray-400 uppercase tracking-wider font-medium mb-3">Historial de respuestas</p>
                  <div className="space-y-3">
                    {selected.responses.map((r) => (
                      <div key={r.id} className="bg-gray-50 rounded-xl p-4">
                        <div className="flex items-center gap-2 mb-1.5">
                          <span className="text-xs font-medium text-gray-500">{r.user?.name}</span>
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
                      placeholder="Escribí tu respuesta..."
                      rows={3}
                      className="input-field resize-none flex-1"
                    />
                    <button
                      onClick={handleSendResponse}
                      disabled={!responseText.trim()}
                      className="btn-success btn-sm self-end flex items-center gap-1.5 disabled:opacity-50"
                    >
                      <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M12 19V5m0 0l-7 7m7-7l7 7" />
                      </svg>
                      Enviar
                    </button>
                  </div>
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
