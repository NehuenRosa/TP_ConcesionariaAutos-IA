import { useCallback, useEffect, useRef, useState } from 'react'
import { api } from '../services/api'

export function useNotificaciones(autenticado: boolean) {
  const [cantidadConsultas, setCantidadConsultas] = useState(0)
  const [cantidadCotizaciones, setCantidadCotizaciones] = useState(0)
  const [nuevoAviso, setNuevoAviso] = useState(false)
  const [canalAviso, setCanalAviso] = useState<'consultas' | 'cotizaciones' | null>(null)
  const anterioresRef = useRef({ consultas: 0, cotizaciones: 0 })
  const inicializadoRef = useRef(false)

  const verificar = useCallback(async () => {
    try {
      const resultado = await api.obtenerContadorNotificaciones()
      if (inicializadoRef.current) {
        const subioConsultas = resultado.consultas > anterioresRef.current.consultas
        const subioCotizaciones = resultado.cotizaciones > anterioresRef.current.cotizaciones
        const subio = subioConsultas || subioCotizaciones
        if (subio) {
          setNuevoAviso(true)
          // Si subió solo el canal de cotizaciones, el aviso apunta ahí para
          // que el usuario aterrice en la bandeja correcta al tocar "Ver".
          setCanalAviso(subioCotizaciones && !subioConsultas ? 'cotizaciones' : 'consultas')
        }
      }
      if (resultado.consultas === 0 && resultado.cotizaciones === 0) {
        setNuevoAviso(false)
        setCanalAviso(null)
      }
      inicializadoRef.current = true
      anterioresRef.current = { consultas: resultado.consultas, cotizaciones: resultado.cotizaciones }
      setCantidadConsultas(resultado.consultas)
      setCantidadCotizaciones(resultado.cotizaciones)
    } catch {
      // Ignorar errores de polling
    }
  }, [])

  const descartarAviso = useCallback(() => setNuevoAviso(false), [])

  useEffect(() => {
    if (!autenticado) return

    verificar()

    // Re-verificar cuando algún chat marca mensajes como leídos
    const handler = () => verificar()
    window.addEventListener('mensajes-leidos', handler)

    // Polling liviano: solo consulta el contador, no los hilos completos
    const intervalo = setInterval(verificar, 3000)

    return () => {
      window.removeEventListener('mensajes-leidos', handler)
      clearInterval(intervalo)
    }
  }, [autenticado, verificar])

  return { cantidadConsultas, cantidadCotizaciones, nuevoAviso, canalAviso, descartarAviso }
}
