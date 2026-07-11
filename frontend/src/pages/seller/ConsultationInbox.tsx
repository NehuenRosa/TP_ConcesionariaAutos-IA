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

  if (loading) return <div className="text-center py-20 text-gray-500">Cargando...</div>

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Bandeja de Consultas</h1>
      <div className="grid md:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow">
          {consultations.length === 0 ? (
            <p className="p-6 text-gray-500">No hay consultas</p>
          ) : (
            <ul className="divide-y">
              {consultations.map((c) => (
                <li
                  key={c.id}
                  onClick={() => setSelected(c)}
                  className={`p-4 cursor-pointer hover:bg-gray-50 ${selected?.id === c.id ? 'bg-blue-50' : ''}`}
                >
                  <div className="flex justify-between items-start">
                    <div>
                      <p className="font-medium">{c.client?.name || 'Cliente'}</p>
                      <p className="text-sm text-gray-500">{c.vehicle?.brand} {c.vehicle?.model}</p>
                    </div>
                    <span className={`px-2 py-1 rounded text-xs ${
                      c.status === 'pendiente' ? 'bg-yellow-100 text-yellow-700' :
                      c.status === 'en_conversacion' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-700'
                    }`}>{c.status}</span>
                  </div>
                  <p className="text-sm text-gray-600 mt-2 line-clamp-2">{c.message}</p>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          {selected ? (
            <>
              <h2 className="font-semibold text-lg mb-2">Detalle de la consulta</h2>
              <p className="text-gray-600 mb-2">
                <strong>Cliente:</strong> {selected.client?.name} ({selected.client?.email})
              </p>
              <p className="text-gray-600 mb-2">
                <strong>Vehículo:</strong> {selected.vehicle?.brand} {selected.vehicle?.model} ({selected.vehicle?.year})
              </p>
              <p className="text-gray-600 mb-4"><strong>Mensaje:</strong> {selected.message}</p>

              <div className="flex gap-2 mb-4">
                {selected.status !== 'cerrada' && (
                  <>
                    <button onClick={() => handleStatusChange(selected.id, 'en_conversacion')} className="bg-blue-600 text-white px-3 py-1 rounded text-sm">Tomar</button>
                    <button onClick={() => handleStatusChange(selected.id, 'cerrada')} className="bg-gray-600 text-white px-3 py-1 rounded text-sm">Cerrar</button>
                  </>
                )}
              </div>

              {selected.responses && selected.responses.length > 0 && (
                <div className="mb-4 space-y-2">
                  <h3 className="font-medium">Historial de respuestas</h3>
                  {selected.responses.map((r) => (
                    <div key={r.id} className="bg-gray-50 p-3 rounded">
                      <p className="text-xs text-gray-400">{r.user?.name}</p>
                      <p className="text-sm">{r.message}</p>
                    </div>
                  ))}
                </div>
              )}

              {selected.status !== 'cerrada' && (
                <div className="flex gap-2">
                  <textarea
                    value={responseText}
                    onChange={(e) => setResponseText(e.target.value)}
                    placeholder="Escribí tu respuesta..."
                    rows={3}
                    className="flex-1 border rounded px-3 py-2 text-sm"
                  />
                  <button onClick={handleSendResponse} className="bg-green-600 text-white px-4 py-2 rounded text-sm self-end">
                    Enviar
                  </button>
                </div>
              )}
            </>
          ) : (
            <p className="text-gray-500">Seleccioná una consulta para ver el detalle</p>
          )}
        </div>
      </div>
    </div>
  )
}
