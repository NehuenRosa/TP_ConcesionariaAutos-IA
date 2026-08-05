import type {
  FiltrosVehiculos,
  PaginaVehiculos,
  PaginaVehiculosGestion,
  Vehiculo,
  VehiculoEntrada,
} from '../types/vehiculo'
import type { DatosLogin, DatosRegistro, DatosUsuarioAdmin, RespuestaLogin, Usuario } from '../types/usuario'
import type { ConsultaResumen, CrearConsulta, Mensaje } from '../types/consulta'
import type { FranjaHoraria, SolicitarTestDrive, TurnoTestDrive } from '../types/testDrive'

const urlBase = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api'
const CLAVE_TOKEN = 'token_concesionaria'

export class ErrorApi extends Error {
  constructor(
    mensaje: string,
    public estado: number,
  ) {
    super(mensaje)
    this.name = 'ErrorApi'
  }
}

export function obtenerToken(): string | null {
  return localStorage.getItem(CLAVE_TOKEN)
}

export function guardarToken(token: string): void {
  localStorage.setItem(CLAVE_TOKEN, token)
}

export function eliminarToken(): void {
  localStorage.removeItem(CLAVE_TOKEN)
}

async function peticion<T>(ruta: string, opciones?: RequestInit): Promise<T> {
  const token = obtenerToken()
  const respuesta = await fetch(`${urlBase}${ruta}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
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
  registrar: (datos: DatosRegistro) =>
    peticion<Usuario>('/auth/registro', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  iniciarSesion: (datos: DatosLogin) =>
    peticion<RespuestaLogin>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  obtenerPerfil: () => peticion<Usuario>('/auth/perfil'),

  // Usuarios (administrador)
  listarUsuarios: () => peticion<Usuario[]>('/admin/usuarios'),
  crearUsuario: (datos: DatosUsuarioAdmin) =>
    peticion<Usuario>('/admin/usuarios', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  actualizarUsuario: (id: number, datos: DatosUsuarioAdmin) =>
    peticion<Usuario>(`/admin/usuarios/${id}`, {
      method: 'PUT',
      body: JSON.stringify(datos),
    }),
  eliminarUsuario: (id: number) =>
    peticion<void>(`/admin/usuarios/${id}`, { method: 'DELETE' }),
  listarVehiculos: (pagina: number, tamano: number, filtros?: FiltrosVehiculos) =>
    peticion<PaginaVehiculos>(`/vehiculos?${construirConsultaVehiculos(pagina, tamano, filtros)}`),
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

  // Consultas
  crearConsulta: (datos: CrearConsulta) =>
    peticion<ConsultaResumen>('/consultas', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  listarMisConsultas: () => peticion<ConsultaResumen[]>('/consultas/mis-consultas'),
  listarBandeja: () => peticion<ConsultaResumen[]>('/consultas/bandeja'),
  tomarConsulta: (id: number) =>
    peticion<ConsultaResumen>(`/consultas/${id}/tomar`, { method: 'PUT' }),
  cerrarConsulta: (id: number) =>
    peticion<ConsultaResumen>(`/consultas/${id}/cerrar`, { method: 'PUT' }),
  eliminarConsulta: (id: number) =>
    peticion<void>(`/consultas/${id}`, { method: 'DELETE' }),

  // Mensajes
  enviarMensaje: (consultaId: number, contenido: string) =>
    peticion<Mensaje>(`/consultas/${consultaId}/mensajes`, {
      method: 'POST',
      body: JSON.stringify({ contenido }),
    }),
  obtenerMensajes: (consultaId: number) =>
    peticion<Mensaje[]>(`/consultas/${consultaId}/mensajes`),
  obtenerMensajesNuevos: (consultaId: number, desde: string) =>
    peticion<Mensaje[]>(`/consultas/${consultaId}/mensajes/nuevos?desde=${desde}`),
  marcarComoLeidos: (consultaId: number) =>
    peticion<void>(`/consultas/${consultaId}/mensajes/leidos`, { method: 'PUT' }),

  // Notificaciones
  obtenerContadorNotificaciones: () =>
    peticion<{ contador: number }>('/notificaciones/contador'),

  // Test drives
  obtenerFranjas: () => peticion<FranjaHoraria[]>('/test-drives/franjas'),
  solicitarTestDrive: (datos: SolicitarTestDrive) =>
    peticion<TurnoTestDrive>('/test-drives', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  listarMisTestDrives: () => peticion<TurnoTestDrive[]>('/test-drives/mis-turnos'),
  cancelarTestDrive: (id: number) =>
    peticion<TurnoTestDrive>(`/test-drives/${id}`, { method: 'DELETE' }),
  listarTestDrives: (estado?: string) =>
    peticion<TurnoTestDrive[]>(`/test-drives${estado ? `?estado=${estado}` : ''}`),
  confirmarTestDrive: (id: number) =>
    peticion<TurnoTestDrive>(`/test-drives/${id}/confirmar`, { method: 'PUT' }),
  cancelarTestDriveVendedor: (id: number) =>
    peticion<TurnoTestDrive>(`/test-drives/${id}/cancelar`, { method: 'PUT' }),
  completarTestDrive: (id: number) =>
    peticion<TurnoTestDrive>(`/test-drives/${id}/completar`, { method: 'PUT' }),
}

// construirConsultaVehiculos arma la query string del catálogo público con
// paginación y filtros opcionales, omitiendo los parámetros vacíos.
function construirConsultaVehiculos(
  pagina: number,
  tamano: number,
  filtros?: FiltrosVehiculos,
): string {
  const parametros = new URLSearchParams()
  parametros.set('pagina', String(pagina))
  parametros.set('tamano', String(tamano))

  if (filtros) {
    const mapeo: Array<[keyof FiltrosVehiculos, string]> = [
      ['busqueda', 'busqueda'],
      ['marca', 'marca'],
      ['modelo', 'modelo'],
      ['anioMin', 'anio_min'],
      ['anioMax', 'anio_max'],
      ['precioMin', 'precio_min'],
      ['precioMax', 'precio_max'],
      ['tipo', 'tipo'],
      ['combustible', 'combustible'],
      ['condicion', 'condicion'],
      ['ordenPor', 'orden_por'],
      ['ordenDireccion', 'orden_direccion'],
    ]
    for (const [campo, parametro] of mapeo) {
      const valor = filtros[campo]
      if (valor !== undefined && valor !== '' && valor !== null) {
        parametros.set(parametro, String(valor))
      }
    }
  }

  return parametros.toString()
}
