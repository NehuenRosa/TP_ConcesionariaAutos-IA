export type Condicion = 'nuevo' | 'usado'

export type EstadoVehiculo = 'disponible' | 'reservado' | 'vendido' | 'dado_de_baja'

export interface Imagen {
  id: number
  url: string
}

export interface Vehiculo {
  id: number
  marca: string
  modelo: string
  anio: number
  kilometraje: number
  combustible: string
  transmision: string
  precio: number
  condicion: Condicion
  estado: EstadoVehiculo
  imagenes: Imagen[]
}

export interface PaginaVehiculos {
  datos: Vehiculo[]
  pagina: number
  tamano: number
  total: number
}
