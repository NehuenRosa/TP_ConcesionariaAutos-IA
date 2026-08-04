import { useEffect, useState, useCallback } from 'react'
import { api } from '../services/api'

export function useNotificaciones(autenticado: boolean) {
  const [cantidad, setCantidad] = useState(0)

  const verificar = useCallback(async () => {
    try {
      const resultado = await api.obtenerContadorNotificaciones()
      setCantidad(resultado.contador)
    } catch {
      // Ignorar errores de polling
    }
  }, [])

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

  return cantidad > 0
}
