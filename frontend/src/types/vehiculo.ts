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
  tipo: string
  precio: number
  condicion: Condicion
  estado: EstadoVehiculo
  imagenes: Imagen[]
}

export interface ResumenVehiculo {
  id: number
  marca: string
  modelo: string
  anio: number
  precio: number
  condicion: Condicion
  tipo: string
  imagen: string
}

export interface PaginaVehiculos {
  datos: ResumenVehiculo[]
  pagina: number
  tamano: number
  total: number
}

export interface VehiculoEntrada {
  marca: string
  modelo: string
  anio: number
  kilometraje: number
  combustible: string
  transmision: string
  tipo: string
  precio: number
  condicion: Condicion
  estado: EstadoVehiculo
  imagenes: string[]
}

export type OrdenPorVehiculos = 'precio' | 'anio'

export type OrdenDireccion = 'asc' | 'desc'

export interface FiltrosVehiculos {
  busqueda?: string
  marca?: string
  modelo?: string
  anioMin?: number
  anioMax?: number
  precioMin?: number
  precioMax?: number
  tipo?: string
  combustible?: string
  condicion?: Condicion
  ordenPor?: OrdenPorVehiculos
  ordenDireccion?: OrdenDireccion
}

export interface ResumenVehiculoGestion {
  id: number
  marca: string
  modelo: string
  anio: number
  kilometraje: number
  precio: number
  condicion: Condicion
  estado: EstadoVehiculo
  imagen: string
}

export interface PaginaVehiculosGestion {
  datos: ResumenVehiculoGestion[]
  pagina: number
  tamano: number
  total: number
}
