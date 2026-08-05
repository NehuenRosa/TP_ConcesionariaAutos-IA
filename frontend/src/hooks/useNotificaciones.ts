import { useCallback, useEffect, useRef, useState } from 'react'
import { api } from '../services/api'

export function useNotificaciones(autenticado: boolean) {
  const [cantidad, setCantidad] = useState(0)
  const [nuevoAviso, setNuevoAviso] = useState(false)
  const anteriorRef = useRef(0)
  const inicializadoRef = useRef(false)

  const verificar = useCallback(async () => {
    try {
      const resultado = await api.obtenerContadorNotificaciones()
      if (resultado.contador > anteriorRef.current && inicializadoRef.current) {
        setNuevoAviso(true)
      }
      if (resultado.contador === 0) {
        setNuevoAviso(false)
      }
      inicializadoRef.current = true
      anteriorRef.current = resultado.contador
      setCantidad(resultado.contador)
    } catch {
      // Ignorar errores de polling
    }
  }, [])

  const descartarAviso = useCallback(() => setNuevoAviso(false), [])

  useEffect(() => {
    if (!autenticado) return

    verificar()

    // Re-verificar cuando el chat marca mensajes como leídos
    const handler = () => verificar()
    window.addEventListener('mensajes-leidos', handler)

    // Polling liviano: solo consulta el contador, no las consultas completas
    const intervalo = setInterval(verificar, 3000)

    return () => {
      window.removeEventListener('mensajes-leidos', handler)
      clearInterval(intervalo)
    }
  }, [autenticado, verificar])

  return { cantidad, nuevoAviso, descartarAviso }
}
