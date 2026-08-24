import { useRef, useState } from 'react'
import { api, ErrorApi } from '../services/api'
import { Boton } from './ui/Boton'

interface PropiedadesSubirComprobante {
  reservaId: number
  yaEnviado?: boolean
  alEnviar?: () => void
}

// SubirComprobante permite al cliente adjuntar la imagen del comprobante de
// la seña (JPG/PNG/WebP de hasta 5 MB) de una reserva activa.
export function SubirComprobante({ reservaId, yaEnviado = false, alEnviar }: PropiedadesSubirComprobante) {
  const entradaArchivo = useRef<HTMLInputElement>(null)
  const [archivo, setArchivo] = useState<File | null>(null)
  const [enviando, setEnviando] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [enviado, setEnviado] = useState(yaEnviado)

  const handleSeleccionar = (e: React.ChangeEvent<HTMLInputElement>) => {
    setError(null)
    setArchivo(e.target.files?.[0] ?? null)
  }

  const handleEnviar = async () => {
    if (!archivo || enviando) return
    setEnviando(true)
    setError(null)
    try {
      await api.subirComprobanteReserva(reservaId, archivo)
      setEnviado(true)
      alEnviar?.()
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo subir el comprobante')
    } finally {
      setEnviando(false)
    }
  }

  if (enviado) {
    return (
      <p className="flex items-center gap-2 text-sm text-emerald-300">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="h-4 w-4 shrink-0">
          <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
        </svg>
        Comprobante enviado. Un vendedor lo va a revisar.
      </p>
    )
  }

  return (
    <div className="space-y-2">
      <div className="flex flex-wrap items-center gap-2">
        <input
          ref={entradaArchivo}
          type="file"
          accept="image/jpeg,image/png,image/webp"
          onChange={handleSeleccionar}
          className="block w-full max-w-xs cursor-pointer rounded-lg border border-white/10 bg-carbono-900/60 text-sm text-plata-300 file:mr-3 file:cursor-pointer file:rounded-l-lg file:border-0 file:bg-carbono-800 file:px-3 file:py-2 file:text-sm file:text-plata-200"
        />
        <Boton tamano="sm" onClick={handleEnviar} disabled={!archivo || enviando}>
          {enviando ? 'Subiendo…' : 'Subir comprobante'}
        </Boton>
      </div>
      {error && <p className="text-sm text-red-400">{error}</p>}
    </div>
  )
}
