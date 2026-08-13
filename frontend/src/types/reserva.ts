export type EstadoReserva = 'activa' | 'vendida' | 'cancelada'

export interface ClienteResumen {
  id: number
  nombre: string
}

export interface Reserva {
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
  cliente: ClienteResumen
  estado: EstadoReserva
  createdAt: string
}

export interface CrearReserva {
  vehiculoId: number
}
