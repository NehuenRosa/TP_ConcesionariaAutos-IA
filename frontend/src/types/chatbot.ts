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
}
