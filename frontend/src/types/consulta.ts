export type EstadoConsulta = 'pendiente' | 'en_conversacion' | 'cerrada'

export interface UsuarioResumen {
  id: number
  nombre: string
}

export interface MensajeResumen {
  contenido: string
  createdAt: string
}

export interface Mensaje {
  id: number
  consultaId: number
  remitente: UsuarioResumen
  contenido: string
  leido: boolean
  createdAt: string
}

export interface ConsultaResumen {
  id: number
  vehiculo: {
    id: number
    marca: string
    modelo: string
    anio: number
    precio: number
    condicion: string
    tipo: string
    imagen: string
  }
  cliente: UsuarioResumen
  vendedor?: UsuarioResumen
  estado: EstadoConsulta
  ultimoMensaje?: MensajeResumen
  mensajesNuevos: number
  createdAt: string
  updatedAt: string
}

export interface Consulta {
  id: number
  vehiculo: {
    id: number
    marca: string
    modelo: string
    anio: number
    kilometraje: number
    combustible: string
    transmision: string
    tipo: string
    precio: number
    condicion: string
    estado: string
    imagenes: Array<{ id: number; url: string }>
  }
  cliente: UsuarioResumen
  vendedor?: UsuarioResumen
  estado: EstadoConsulta
  mensajes: Mensaje[]
  createdAt: string
  updatedAt: string
}

export interface CrearConsulta {
  vehiculoId: number
  mensaje: string
}
