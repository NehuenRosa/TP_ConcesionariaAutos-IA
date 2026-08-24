export type RolTurno = 'usuario' | 'asistente'

export interface TurnoChat {
  rol: RolTurno
  contenido: string
}

export interface PeticionChatbot {
  mensaje: string
  historial: TurnoChat[]
}

export interface RespuestaChatbot {
  respuesta: string
  cotizacionId?: number
  /** Ids de vehículos del stock que la IA mencionó, para enlazar sus fichas. */
  vehiculosMencionados?: number[]
}

export interface RespuestaTasacion {
  respuesta: string
  sesionId?: string
}

export interface RespuestaConfirmarTasacion {
  respuesta: string
  confirmada: boolean
}

export interface ConfirmarTasacion {
  sesionId: string
  mensaje: string
}
