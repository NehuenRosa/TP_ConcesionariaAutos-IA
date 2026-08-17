import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router'
import { api, ErrorApi } from '../services/api'
import type { TurnoChat } from '../types/chatbot'
import { Boton } from './ui/Boton'

const MAXIMO_FOTOS = 5
const TAMANO_MAXIMO_MB = 5
const FORMATOS_PERMITIDOS = ['image/jpeg', 'image/png', 'image/webp']
const EXTENSIONES_PERMITIDAS = ['jpg', 'jpeg', 'png', 'webp']

const esImagenAceptada = (archivo: File): boolean => {
  if (FORMATOS_PERMITIDOS.includes(archivo.type)) return true
  const extension = archivo.name.split('.').pop()?.toLowerCase() ?? ''
  return EXTENSIONES_PERMITIDAS.includes(extension)
}

const generarSesionId = (): string => {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `${Date.now()}-${Math.random().toString(36).slice(2)}`
}

type FotoConUrl = { archivo: File; url: string }

export function Chatbot() {
  const navigate = useNavigate()
  const [abierto, setAbierto] = useState(false)
  const [modo, setModo] = useState<'chat' | 'tasacion'>('chat')
  const [mensajesConsulta, setMensajesConsulta] = useState<TurnoChat[]>([])
  const [mensajesTasacion, setMensajesTasacion] = useState<TurnoChat[]>([])
  const [entrada, setEntrada] = useState('')
  const [cargando, setCargando] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [fotos, setFotos] = useState<FotoConUrl[]>([])
  const [descripcion, setDescripcion] = useState('')
  const [sesionTasacion, setSesionTasacion] = useState<string | null>(null)
  const [entradaVisita, setEntradaVisita] = useState('')
  const mensajesRef = useRef<HTMLDivElement>(null)
  const entradaRef = useRef<HTMLTextAreaElement>(null)
  const descripcionRef = useRef<HTMLTextAreaElement>(null)
  const visitaRef = useRef<HTMLTextAreaElement>(null)
  const urlsCreadasRef = useRef<string[]>([])
  const enviandoRef = useRef(false)

  const crearUrl = (archivo: File): string => {
    const url = URL.createObjectURL(archivo)
    urlsCreadasRef.current.push(url)
    return url
  }

  const mensajes = modo === 'tasacion' ? mensajesTasacion : mensajesConsulta

  const revocarUrl = (url: string) => {
    const indice = urlsCreadasRef.current.indexOf(url)
    if (indice !== -1) {
      urlsCreadasRef.current.splice(indice, 1)
    }
    URL.revokeObjectURL(url)
  }

  useEffect(() => {
    return () => {
      urlsCreadasRef.current.forEach((url) => URL.revokeObjectURL(url))
      urlsCreadasRef.current = []
    }
  }, [])

  useEffect(() => {
    if (mensajesRef.current) {
      mensajesRef.current.scrollTop = mensajesRef.current.scrollHeight
    }
  }, [mensajes, cargando, abierto])

  useEffect(() => {
    if (abierto) {
      if (modo === 'tasacion' && sesionTasacion) {
        visitaRef.current?.focus()
      } else if (modo === 'tasacion') {
        descripcionRef.current?.focus()
      } else {
        entradaRef.current?.focus()
      }
    }
  }, [abierto, modo, sesionTasacion])

  const enviarMensaje = async () => {
    const mensaje = entrada.trim()
    if (!mensaje || enviandoRef.current) return

    enviandoRef.current = true
    setCargando(true)
    setError(null)
    setEntrada('')
    setMensajesConsulta((m) => [...m, { rol: 'usuario', contenido: mensaje }])

    try {
      const historial = mensajesConsulta.map((turno) => ({
        rol: turno.rol,
        contenido: turno.contenido,
      }))
      const respuesta = await api.enviarMensajeChatbot({ mensaje, historial })
      setMensajesConsulta((m) => [...m, { rol: 'asistente', contenido: respuesta.respuesta }])
      // Cuando la IA creó una cotización, se redirige al panel para seguir la
      // conversación sobre ese vehículo.
      if (respuesta.cotizacionId) {
        setTimeout(() => {
          setAbierto(false)
          navigate(`/mis-cotizaciones/${respuesta.cotizacionId}`)
        }, 350)
      }
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo enviar el mensaje')
    } finally {
      enviandoRef.current = false
      setCargando(false)
    }
  }

  const enviarTasacion = async () => {
    const detalle = descripcion.trim()
    if ((fotos.length === 0 && detalle === '') || enviandoRef.current) return

    enviandoRef.current = true
    setCargando(true)
    setError(null)

    try {
      const respuesta = await api.enviarTasacion(
        fotos.map((foto) => foto.archivo),
        descripcion,
        sesionTasacion ?? generarSesionId(),
      )
      const partesUsuario = [`Tasación de mi auto (${fotos.length} foto${fotos.length > 1 ? 's' : ''})`]
      if (detalle !== '') {
        partesUsuario.push(detalle)
      }
      setMensajesTasacion((m) => [
        ...m,
        { rol: 'usuario', contenido: partesUsuario.join('\n') },
        { rol: 'asistente', contenido: respuesta.respuesta },
      ])
      if (respuesta.sesionId) {
        setSesionTasacion(respuesta.sesionId)
      }
      fotos.forEach((foto) => revocarUrl(foto.url))
      setFotos([])
      setDescripcion('')
    } catch (e: unknown) {
      if (e instanceof DOMException && e.name === 'AbortError') {
        setError('La tasación tardó demasiado. Probá de nuevo con una foto más liviana o con menos fotos.')
      } else {
        setError(e instanceof ErrorApi ? e.message : 'No se pudo realizar la tasación')
      }
    } finally {
      enviandoRef.current = false
      setCargando(false)
    }
  }

  const confirmarVisita = async () => {
    const mensaje = entradaVisita.trim()
    if (!mensaje || !sesionTasacion || enviandoRef.current) return

    enviandoRef.current = true
    setCargando(true)
    setError(null)
    setEntradaVisita('')
    setMensajesTasacion((m) => [...m, { rol: 'usuario', contenido: mensaje }])

    try {
      const respuesta = await api.confirmarTasacion({ sesionId: sesionTasacion, mensaje })
      setMensajesTasacion((m) => [...m, { rol: 'asistente', contenido: respuesta.respuesta }])
      if (respuesta.confirmada) {
        setSesionTasacion(null)
      }
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo confirmar la visita')
      setMensajesTasacion((m) => m.slice(0, -1))
    } finally {
      enviandoRef.current = false
      setCargando(false)
    }
  }

  const agregarFotos = (archivos: FileList | null) => {
    if (!archivos) return
    setError(null)
    const nuevas: FotoConUrl[] = []
    for (const archivo of Array.from(archivos)) {
      if (!esImagenAceptada(archivo)) {
        setError('Formato no soportado: se aceptan JPG, PNG o WebP')
        continue
      }
      if (archivo.size > TAMANO_MAXIMO_MB * 1024 * 1024) {
        setError(`Cada foto debe pesar menos de ${TAMANO_MAXIMO_MB} MB`)
        continue
      }
      nuevas.push({ archivo, url: crearUrl(archivo) })
    }
    setFotos((actuales) => {
      const combinadas = [...actuales, ...nuevas]
      const sobrantes = combinadas.slice(MAXIMO_FOTOS)
      sobrantes.forEach((foto) => revocarUrl(foto.url))
      return combinadas.slice(0, MAXIMO_FOTOS)
    })
  }

  const quitarFoto = (indice: number) => {
    setFotos((actuales) => {
      const objetivo = actuales[indice]
      if (objetivo) revocarUrl(objetivo.url)
      return actuales.filter((_, i) => i !== indice)
    })
  }

  const seleccionarModo = (nuevoModo: 'chat' | 'tasacion') => {
    setModo(nuevoModo)
    setError(null)
  }

  return (
    <>
      <button
        type="button"
        onClick={() => setAbierto((a) => !a)}
        className="fixed right-4 bottom-4 z-40 flex h-14 w-14 items-center justify-center rounded-full bg-acento-500 text-carbono-900 shadow-luz transition-transform hover:scale-105 focus-visible:ring-2 focus-visible:ring-acento-400 focus-visible:ring-offset-2 focus-visible:ring-offset-carbono-950 focus-visible:outline-none sm:right-6 sm:bottom-6"
        aria-label={abierto ? 'Cerrar asistente' : 'Abrir asistente'}
      >
        {abierto ? (
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="h-5 w-5">
            <path strokeLinecap="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        ) : (
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
            <path
              strokeLinejoin="round"
              d="M12 3a9 9 0 0 0-7.8 13.5L3 21l4.6-1.2A9 9 0 1 0 12 3z"
            />
            <path strokeLinecap="round" d="M8.5 12h.01M12 12h.01M15.5 12h.01" strokeWidth="2.5" />
          </svg>
        )}
      </button>

      {abierto && (
        <div className="fixed right-4 bottom-20 z-40 flex h-[min(600px,calc(100vh-6rem))] w-[calc(100%-2rem)] max-w-sm flex-col overflow-hidden rounded-2xl border border-white/10 bg-carbono-900/95 shadow-luz backdrop-blur-xl sm:right-6 sm:bottom-24">
          <div className="flex items-center justify-between gap-2 border-b border-white/8 bg-carbono-850/80 px-4 py-3">
            <div className="flex items-center gap-2.5">
              <span className="flex h-9 w-9 items-center justify-center rounded-full bg-acento-500/15 text-acento-400">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-5 w-5">
                  <path strokeLinejoin="round" d="M12 3a9 9 0 0 0-7.8 13.5L3 21l4.6-1.2A9 9 0 1 0 12 3z" />
                </svg>
              </span>
              <div>
                <p className="font-display text-sm font-semibold text-plata-100">Asistente Aurum</p>
                <p className="text-xs text-plata-500">
                  {modo === 'tasacion' ? 'Tasación de tu auto por fotos' : 'Respondo sobre nuestro stock'}
                </p>
              </div>
            </div>
            <button
              type="button"
              onClick={() => setAbierto(false)}
              className="rounded-md p-1.5 text-plata-500 transition-colors hover:bg-white/5 hover:text-plata-200"
              aria-label="Cerrar asistente"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-4 w-4">
                <path strokeLinecap="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div className="flex gap-1 border-b border-white/8 px-4 py-2">
            <Boton
              tamano="sm"
              variante={modo === 'chat' ? 'acento' : 'fantasma'}
              className="rounded-full"
              onClick={() => seleccionarModo('chat')}
            >
              Consultas
            </Boton>
            <Boton
              tamano="sm"
              variante={modo === 'tasacion' ? 'acento' : 'fantasma'}
              className="rounded-full"
              onClick={() => seleccionarModo('tasacion')}
            >
              Tasá tu auto
            </Boton>
          </div>

          <div ref={mensajesRef} className="flex-1 space-y-3 overflow-y-auto px-4 py-4">
            {mensajes.length === 0 && (
              <p className="text-sm leading-relaxed text-plata-400">
                {modo === 'tasacion'
                  ? 'Subí fotos de tu auto y/o contame sus detalles (marca, modelo, año, versión, estado, rayones…). Te doy una tasación con valores reales de la guía oficial: venta normal y otras dos opciones de cobro. Después coordinamos tu visita a la concesionaria.'
                  : '¡Hola! Preguntame sobre los vehículos disponibles, compará modelos o pedí un test drive.'}
              </p>
            )}
            {mensajes.map((turno, indice) => (
              <div
                key={indice}
                className={`flex ${turno.rol === 'usuario' ? 'justify-end' : 'justify-start'}`}
              >
                <div
                  className={`max-w-[85%] rounded-2xl px-3.5 py-2.5 text-sm leading-relaxed whitespace-pre-wrap ${
                    turno.rol === 'usuario'
                      ? 'rounded-br-md bg-acento-500 text-carbono-900'
                      : 'rounded-bl-md border border-white/8 bg-carbono-800/70 text-plata-100'
                  }`}
                >
                  {turno.contenido}
                </div>
              </div>
            ))}
            {cargando && (
              <div className="flex justify-start">
                <div className="flex items-center gap-2 rounded-2xl rounded-bl-md border border-white/8 bg-carbono-800/70 px-3.5 py-3">
                  <span className="flex items-center gap-1.5">
                    <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-plata-400" />
                    <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-plata-400 [animation-delay:0.15s]" />
                    <span className="h-1.5 w-1.5 animate-bounce rounded-full bg-plata-400 [animation-delay:0.3s]" />
                  </span>
                  {modo === 'tasacion' && (
                    <span className="text-xs text-plata-400">
                      Analizando tus fotos, puede tardar hasta un minuto…
                    </span>
                  )}
                </div>
              </div>
            )}
          </div>

          {error && (
            <div className="border-t border-red-500/30 bg-red-500/10 px-4 py-2.5 text-xs text-red-300">
              {error}
            </div>
          )}

          <div className="border-t border-white/8 bg-carbono-850/60 p-3">
            {modo === 'tasacion' && (
              <div className="mb-3 space-y-2">
                {sesionTasacion ? (
                  <>
                    <p className="text-xs text-plata-400">
                      Coordinemos tu visita: escribí qué día y entre qué franja horaria te podés acercar
                      (ej. "el jueves a las 15").
                    </p>
                    <textarea
                      ref={visitaRef}
                      value={entradaVisita}
                      onChange={(e) => setEntradaVisita(e.target.value)}
                      onKeyDown={(e) => {
                        if (e.key === 'Enter' && !e.shiftKey) {
                          e.preventDefault()
                          void confirmarVisita()
                        }
                      }}
                      rows={2}
                      placeholder="Día y franja horaria…"
                      className="max-h-28 min-h-[52px] w-full resize-none rounded-lg border border-white/10 bg-carbono-900 px-3 py-2 text-sm text-plata-100 placeholder:text-plata-600 focus:border-acento-400 focus:outline-none"
                    />
                    <Boton
                      tamano="sm"
                      variante="acento"
                      disabled={!entradaVisita.trim() || cargando}
                      onClick={confirmarVisita}
                    >
                      Confirmar visita
                    </Boton>
                  </>
                ) : (
                  <>
                    <div className="flex flex-wrap gap-2">
                      {fotos.map((foto, indice) => (
                        <div key={`${foto.archivo.name}-${indice}`} className="relative">
                          <img
                            src={foto.url}
                            alt={foto.archivo.name}
                            className="h-14 w-14 rounded-lg border border-white/10 object-cover"
                          />
                          <button
                            type="button"
                            onClick={() => quitarFoto(indice)}
                            className="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-white"
                            aria-label="Quitar foto"
                          >
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="h-3 w-3">
                              <path strokeLinecap="round" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        </div>
                      ))}
                    </div>
                    <textarea
                      ref={descripcionRef}
                      value={descripcion}
                      onChange={(e) => setDescripcion(e.target.value)}
                      rows={2}
                      placeholder="Detallá tu auto: marca, modelo, año, versión, estado, kilometraje, rayones, abolladuras…"
                      className="max-h-28 min-h-[52px] w-full resize-none rounded-lg border border-white/10 bg-carbono-900 px-3 py-2 text-sm text-plata-100 placeholder:text-plata-600 focus:border-acento-400 focus:outline-none"
                    />
                    <div className="flex items-center gap-2">
                      <label className="cursor-pointer rounded-full border border-plata-400/25 bg-carbono-800/50 px-3 py-1.5 text-xs font-medium text-plata-200 transition-colors hover:border-plata-300/50 hover:bg-carbono-700/60">
                        Subir fotos ({fotos.length}/{MAXIMO_FOTOS})
                        <input
                          type="file"
                          accept="image/jpeg,image/png,image/webp"
                          multiple
                          className="hidden"
                          onChange={(e) => agregarFotos(e.target.files)}
                        />
                      </label>
                      <Boton
                        tamano="sm"
                        variante="acento"
                        disabled={(fotos.length === 0 && descripcion.trim() === '') || cargando}
                        onClick={enviarTasacion}
                      >
                        Tasá
                      </Boton>
                    </div>
                  </>
                )}
              </div>
            )}

            {modo === 'chat' && (
              <div className="flex items-end gap-2">
                <textarea
                  ref={entradaRef}
                  value={entrada}
                  onChange={(e) => setEntrada(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' && !e.shiftKey) {
                      e.preventDefault()
                      void enviarMensaje()
                    }
                  }}
                  rows={1}
                  placeholder="Escribí tu consulta…"
                  className="max-h-28 min-h-[38px] flex-1 resize-none rounded-lg border border-white/10 bg-carbono-900 px-3 py-2 text-sm text-plata-100 placeholder:text-plata-600 focus:border-acento-400 focus:outline-none"
                />
                <Boton
                  tamano="sm"
                  variante="acento"
                  className="h-[38px] px-3"
                  disabled={!entrada.trim() || cargando}
                  onClick={enviarMensaje}
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className="h-4 w-4">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12h15m0 0l-6-6m6 6l-6 6" />
                  </svg>
                </Boton>
              </div>
            )}
          </div>
        </div>
      )}
    </>
  )
}
