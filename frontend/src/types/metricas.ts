export type EstadoVehiculoMetrica = 'disponible' | 'reservado' | 'vendido' | 'dado_de_baja'

export interface VehiculoPorEstado {
  estado: EstadoVehiculoMetrica
  cantidad: number
}

export interface ConsultaPorDia {
  fecha: string
  cantidad: number
}

export interface ConteoPorMarca {
  marca: string
  cantidad: number
}

export interface VehiculoEnStock {
  id: number
  marca: string
  modelo: string
  anio: number
  diasEnStock: number
}

export interface Metricas {
  vehiculosPorEstado: VehiculoPorEstado[]
  consultasPorPeriodo: ConsultaPorDia[]
  ventasPorPeriodo: ConsultaPorDia[]
  ingresoPorPeriodo: number
  ventasPorMarca: ConteoPorMarca[]
  testDrivesPorPeriodo: ConsultaPorDia[]
  vehiculosEnStock: VehiculoEnStock[]
  reservasActivas: number
  reservasVendidas: number
  testDrivesAgendados: number
  testDrivesCompletados: number
  consultasAbiertas: number
  totalUsuarios: number
}
