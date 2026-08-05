import { useParams, useNavigate } from 'react-router'
import { ChatConsulta } from '../components/ChatConsulta'
import { Boton } from '../components/ui/Boton'

export function ChatVendedor() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  if (!id) {
    return (
      <div className="mx-auto max-w-3xl px-4 py-24 text-center sm:px-6">
        <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">Error</p>
        <h1 className="font-display text-3xl font-bold text-plata-100">Consulta no encontrada</h1>
        <p className="mt-3 text-plata-400">La consulta solicitada no existe.</p>
        <div className="mt-8">
          <Boton onClick={() => navigate('/vendedor/bandeja')}>Volver a la bandeja</Boton>
        </div>
      </div>
    )
  }

  const estado: 'en_conversacion' | 'cerrada' = 'en_conversacion'

  return (
    <div className="mx-auto flex h-[calc(100vh-6.5rem)] max-w-4xl flex-col px-4 py-6 sm:px-6">
      <div className="mb-3 flex items-center justify-between rounded-2xl border border-white/8 bg-carbono-850/60 px-5 py-4 backdrop-blur-sm">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => navigate('/vendedor/bandeja')}
            className="inline-flex items-center gap-1.5 font-display text-sm font-medium text-plata-400 transition-colors hover:text-plata-100"
          >
            <span aria-hidden>←</span> Bandeja
          </button>
          <span className="h-4 w-px bg-white/10" />
          <h1 className="font-display text-lg font-semibold text-plata-100">Consulta #{id}</h1>
        </div>
        <span className="rounded-full border border-sky-400/40 bg-sky-400/10 px-3 py-1 font-display text-[11px] font-semibold tracking-wide text-sky-300">
          En conversación
        </span>
      </div>

      <div className="min-h-0 flex-1 overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/60 backdrop-blur-sm">
        <ChatConsulta consultaId={Number(id)} estado={estado} />
      </div>
    </div>
  )
}
