import { useParams, useNavigate } from 'react-router'
import { ChatConsulta } from '../components/ChatConsulta'

export function ChatVendedor() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  if (!id) {
    return (
      <div className="space-y-4">
        <h1 className="text-2xl font-bold text-gray-900">Consulta no encontrada</h1>
        <p className="text-gray-700">La consulta solicitada no existe.</p>
        <button
          type="button"
          onClick={() => navigate('/vendedor/bandeja')}
          className="inline-block rounded-md bg-gray-900 px-4 py-2 text-white hover:bg-gray-700"
        >
          Volver a la bandeja
        </button>
      </div>
    )
  }

  const estado: 'en_conversacion' | 'cerrada' = 'en_conversacion'

  return (
    <div className="flex h-[calc(100vh-200px)] flex-col">
      <div className="border-b border-gray-200 p-4">
        <button
          type="button"
          onClick={() => navigate('/vendedor/bandeja')}
          className="mb-2 text-sm text-gray-500 hover:text-gray-700"
        >
          ← Volver a la bandeja
        </button>
        <h1 className="text-lg font-semibold text-gray-900">Consulta #{id}</h1>
      </div>

      <div className="flex-1">
        <ChatConsulta
          consultaId={Number(id)}
          estado={estado}
        />
      </div>
    </div>
  )
}
