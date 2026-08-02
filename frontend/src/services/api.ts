import type {
  PaginaVehiculos,
  PaginaVehiculosGestion,
  Vehiculo,
  VehiculoEntrada,
} from '../types/vehiculo'

const urlBase = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api'

export class ErrorApi extends Error {
  constructor(
    mensaje: string,
    public estado: number,
  ) {
    super(mensaje)
    this.name = 'ErrorApi'
  }
}

async function peticion<T>(ruta: string, opciones?: RequestInit): Promise<T> {
  const respuesta = await fetch(`${urlBase}${ruta}`, {
    headers: {
      'Content-Type': 'application/json',
      ...opciones?.headers,
    },
    ...opciones,
  })

  if (!respuesta.ok) {
    let mensaje = 'Ocurrió un error inesperado. Intente nuevamente.'
    try {
      const cuerpo = (await respuesta.json()) as { error?: string }
      if (cuerpo.error) {
        mensaje = cuerpo.error
      }
    } catch {
      // Se ignora: se mantiene el mensaje por defecto.
    }
    throw new ErrorApi(mensaje, respuesta.status)
  }

  return (await respuesta.json()) as T
}

export const api = {
  obtenerEstado: () => peticion<{ estado: string }>('/health'),
  listarVehiculos: (pagina: number, tamano: number) =>
    peticion<PaginaVehiculos>(`/vehiculos?pagina=${pagina}&tamano=${tamano}`),
  obtenerVehiculo: (id: number) => peticion<Vehiculo>(`/vehiculos/${id}`),
  listarVehiculosGestion: (pagina: number, tamano: number, estado?: string) =>
    peticion<PaginaVehiculosGestion>(
      `/admin/vehiculos?pagina=${pagina}&tamano=${tamano}${estado ? `&estado=${estado}` : ''}`,
    ),
  obtenerVehiculoGestion: (id: number) => peticion<Vehiculo>(`/admin/vehiculos/${id}`),
  crearVehiculo: (datos: VehiculoEntrada) =>
    peticion<Vehiculo>('/admin/vehiculos', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  actualizarVehiculo: (id: number, datos: VehiculoEntrada) =>
    peticion<Vehiculo>(`/admin/vehiculos/${id}`, {
      method: 'PUT',
      body: JSON.stringify(datos),
    }),
  darDeBajaVehiculo: (id: number) =>
    peticion<Vehiculo>(`/admin/vehiculos/${id}`, { method: 'DELETE' }),
}
