export type EstadoCotizacion = 'abierta' | 'cerrada'

export type RemitenteCotizacion = 'cliente' | 'ia'

export interface VehiculoCotizado {
  id: number
  marca: string
  modelo: string
  anio: number
  precio: number
  condicion: string
  tipo: string
  imagen: string
}

export interface MensajeCotizacion {
  id: number
  remitente: RemitenteCotizacion
  contenido: string
  createdAt: string
}

export interface CotizacionResumen {
  id: number
  vehiculo: VehiculoCotizado
  estado: EstadoCotizacion
  ultimoMensaje?: {
    contenido: string
    createdAt: string
  }
  createdAt: string
  updatedAt: string
}

export interface Cotizacion {
  id: number
  vehiculo: VehiculoCotizado
  estado: EstadoCotizacion
  mensajes: MensajeCotizacion[]
  createdAt: string
  updatedAt: string
}

export interface CrearCotizacion {
  vehiculoId: number
  mensaje?: string
}

export interface RespuestaMensajeCotizacion {
  respuesta: string
  enviado: boolean
}