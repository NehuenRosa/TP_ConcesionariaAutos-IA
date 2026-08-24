import type {
  FiltrosVehiculos,
  PaginaVehiculos,
  PaginaVehiculosGestion,
  Vehiculo,
  VehiculoEntrada,
} from '../types/vehiculo'
import type { DatosLogin, DatosLoginGoogle, DatosRegistro, DatosUsuarioAdmin, ProveedoresAuth, RespuestaLogin, Usuario } from '../types/usuario'
import type { ConsultaResumen, CrearConsulta, Mensaje } from '../types/consulta'
import type { FranjaHoraria, SolicitarTestDrive, TurnoTestDrive } from '../types/testDrive'
import type { CrearReserva, DatosTransferencia, Reserva } from '../types/reserva'
import type { Metricas } from '../types/metricas'
import type {
  Cotizacion,
  CotizacionResumen,
  CrearCotizacion,
  RespuestaMensajeCotizacion,
} from '../types/cotizacion'
import type {
  ConfirmarTasacion,
  PeticionChatbot,
  RespuestaChatbot,
  RespuestaConfirmarTasacion,
  RespuestaTasacion,
} from '../types/chatbot'

const urlBase = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api'
const CLAVE_TOKEN = 'token_concesionaria'
// El backend espera como máximo 120s en el chatbot; 140s evita que la interfaz
// quede "pensando" para siempre si el servidor se cuelga.
const TIEMPO_MAXIMO_MILISEGUNDOS = 140_000

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
  const controlador = new AbortController()
  const temporizador = setTimeout(() => controlador.abort(), TIEMPO_MAXIMO_MILISEGUNDOS)
  try {
    const token = obtenerToken()
    const { headers: encabezadosOpcionales, ...restoOpciones } = opciones ?? {}
    const respuesta = await fetch(`${urlBase}${ruta}`, {
      ...restoOpciones,
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...encabezadosOpcionales,
      },
      signal: controlador.signal,
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

    // El backend responde 204 No Content en eliminaciones; un cuerpo vacío
    // rompería respuesta.json(), así que se corta antes.
    if (respuesta.status === 204) {
      return undefined as T
    }

    return (await respuesta.json()) as T
  } finally {
    clearTimeout(temporizador)
  }
}

// peticionFormulario envía un FormData (multipart) con token si existe, sin
// fijar Content-Type para que el navegador agregue el límite del multipart.
async function peticionFormulario<T>(ruta: string, formulario: FormData): Promise<T> {
  const controlador = new AbortController()
  const temporizador = setTimeout(() => controlador.abort(), TIEMPO_MAXIMO_MILISEGUNDOS)
  try {
    const token = obtenerToken()
    const respuesta = await fetch(`${urlBase}${ruta}`, {
      method: 'POST',
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
      body: formulario,
      signal: controlador.signal,
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

    if (respuesta.status === 204) {
      return undefined as T
    }

    return (await respuesta.json()) as T
  } finally {
    clearTimeout(temporizador)
  }
}

// peticionMultipart arma el formulario de la tasación del chatbot y delega en
// el envío genérico de multipart.
async function peticionMultipart<T>(ruta: string, fotos: File[], descripcion: string, sesionId?: string): Promise<T> {
  const formulario = new FormData()
  for (const foto of fotos) {
    formulario.append('fotos', foto)
  }
  formulario.append('descripcion', descripcion)
  if (sesionId) {
    formulario.append('sesion_id', sesionId)
  }
  return peticionFormulario<T>(ruta, formulario)
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
  iniciarSesionGoogle: (datos: DatosLoginGoogle) =>
    peticion<RespuestaLogin>('/auth/google', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  obtenerProveedoresAuth: () => peticion<ProveedoresAuth>('/auth/proveedores'),
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
    peticion<{ contador: number; consultas: number; cotizaciones: number }>('/notificaciones/contador'),

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

  // Reservas
  crearReserva: (datos: CrearReserva) =>
    peticion<Reserva>('/reservas', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  obtenerDatosTransferencia: (vehiculoId: number) =>
    peticion<DatosTransferencia>(`/reservas/datos-transferencia?vehiculoId=${vehiculoId}`),
  subirComprobanteReserva: (id: number, archivo: File) => {
    const formulario = new FormData()
    formulario.append('comprobante', archivo)
    return peticionFormulario<Reserva>(`/reservas/${id}/comprobante`, formulario)
  },
  // Devuelve la imagen del comprobante como Blob para poder mostrarla con
  // URL.createObjectURL (la petición necesita el token en el encabezado).
  obtenerComprobanteReserva: async (id: number): Promise<Blob> => {
    const token = obtenerToken()
    const respuesta = await fetch(`${urlBase}/reservas/${id}/comprobante`, {
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
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
    return respuesta.blob()
  },
  listarMisReservas: () => peticion<Reserva[]>('/reservas/mis-reservas'),
  cancelarReserva: (id: number) =>
    peticion<Reserva>(`/reservas/${id}`, { method: 'DELETE' }),
  listarReservas: (estado?: string) =>
    peticion<Reserva[]>(`/reservas${estado ? `?estado=${estado}` : ''}`),
  confirmarReservaVenta: (id: number) =>
    peticion<Reserva>(`/reservas/${id}/confirmar`, { method: 'PUT' }),
  cancelarReservaVendedor: (id: number, motivo: string) =>
    peticion<Reserva>(`/reservas/${id}/cancelar`, { method: 'PUT', body: JSON.stringify({ motivo }) }),

  // Métricas del panel de administración
  obtenerMetricas: (periodo?: number) =>
    peticion<Metricas>(`/admin/metricas${periodo ? `?periodo=${periodo}` : ''}`),

  // Cotizaciones del cliente (atendidas por la IA)
  crearCotizacion: (datos: CrearCotizacion) =>
    peticion<Cotizacion>('/cotizaciones', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  listarMisCotizaciones: () => peticion<CotizacionResumen[]>('/cotizaciones/mis-cotizaciones'),
  obtenerCotizacion: (id: number) => peticion<Cotizacion>(`/cotizaciones/${id}`),
  enviarMensajeCotizacion: (id: number, mensaje: string) =>
    peticion<RespuestaMensajeCotizacion>(`/cotizaciones/${id}/mensajes`, {
      method: 'POST',
      body: JSON.stringify({ mensaje }),
    }),
  cerrarCotizacion: (id: number) =>
    peticion<Cotizacion>(`/cotizaciones/${id}/cerrar`, { method: 'PUT' }),

  // Atención personal de cotizaciones (vendedor)
  listarBandejaCotizaciones: () => peticion<CotizacionResumen[]>('/cotizaciones/bandeja'),
  obtenerCotizacionPersonal: (id: number) => peticion<Cotizacion>(`/cotizaciones/${id}/personal`),
  tomarCotizacion: (id: number) =>
    peticion<Cotizacion>(`/cotizaciones/${id}/tomar`, { method: 'PUT' }),
  responderCotizacionVendedor: (id: number, mensaje: string) =>
    peticion<Cotizacion>(`/cotizaciones/${id}/mensajes-vendedor`, {
      method: 'POST',
      body: JSON.stringify({ mensaje }),
    }),
  cerrarCotizacionPersonal: (id: number) =>
    peticion<Cotizacion>(`/cotizaciones/${id}/cerrar-personal`, { method: 'PUT' }),

  // Chatbot
  enviarMensajeChatbot: (datos: PeticionChatbot) =>
    peticion<RespuestaChatbot>('/chatbot/mensajes', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
  enviarTasacion: (fotos: File[], descripcion: string, sesionId?: string) =>
    peticionMultipart<RespuestaTasacion>('/chatbot/tasacion', fotos, descripcion, sesionId),
  confirmarTasacion: (datos: ConfirmarTasacion) =>
    peticion<RespuestaConfirmarTasacion>('/chatbot/tasacion/confirmar', {
      method: 'POST',
      body: JSON.stringify(datos),
    }),
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
