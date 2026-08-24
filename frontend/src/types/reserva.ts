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
  /** Monto de la seña (5 % del precio), calculado por el backend. */
  montoSenia: number
  /** Límite para subir el comprobante (RFC3339); vacío en reservas viejas. */
  vencimientoComprobante?: string
  /** Cuándo se subió el comprobante; ausente = pendiente de envío. */
  comprobanteEnviadoAt?: string
  /** Explicación del vendedor al cancelar; vacío si la anuló el cliente. */
  motivoCancelacion?: string
  createdAt: string
}

export interface CrearReserva {
  vehiculoId: number
}

/** Datos bancarios y monto para transferir la seña (los calcula el backend). */
export interface DatosTransferencia {
  cbu: string
  alias: string
  monto: number
}
