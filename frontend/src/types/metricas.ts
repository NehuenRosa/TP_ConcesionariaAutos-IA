export type EstadoVehiculoMetrica = 'disponible' | 'reservado' | 'vendido' | 'dado_de_baja'

export interface VehiculoPorEstado {
  estado: EstadoVehiculoMetrica
  cantidad: number
}

export interface ConsultaPorDia {
  fecha: string
  cantidad: number
}

export interface Metricas {
  vehiculosPorEstado: VehiculoPorEstado[]
  consultasPorPeriodo: ConsultaPorDia[]
  reservasActivas: number
  reservasVendidas: number
  testDrivesAgendados: number
  testDrivesCompletados: number
  consultasAbiertas: number
  totalUsuarios: number
}
