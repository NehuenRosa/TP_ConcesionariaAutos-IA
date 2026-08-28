export type EstadoTurnoTestDrive = 'solicitado' | 'confirmado' | 'cancelado' | 'completado'

export interface FranjaHoraria {
  id: string
  inicio: string
  fin: string
  ocupada?: boolean
}

export interface ClienteResumen {
  id: number
  nombre: string
}

export interface TurnoTestDrive {
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
  fecha: string
  franja: string
  estado: EstadoTurnoTestDrive
}

export interface SolicitarTestDrive {
  vehiculoId: number
  fecha: string
  franja: string
}
