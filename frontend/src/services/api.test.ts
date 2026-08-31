import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { api, ErrorApi, eliminarToken, guardarToken } from './api'

const urlBase = 'http://localhost:8080/api'

function respuestaJSON<T>(cuerpo: T, estado = 200): Response {
  return {
    ok: estado >= 200 && estado < 300,
    status: estado,
    json: async () => cuerpo,
  } as Response
}

describe('cliente api', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn())
    eliminarToken()
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('hace GET a la ruta correcta con Content-Type JSON', async () => {
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ estado: 'ok' }))

    await api.obtenerEstado()

    expect(fetch).toHaveBeenCalledWith(
      `${urlBase}/health`,
      expect.objectContaining({
        headers: { 'Content-Type': 'application/json' },
      }),
    )
  })

  it('agrega Authorization cuando hay token guardado', async () => {
    guardarToken('token-123')
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ estado: 'ok' }))

    await api.obtenerEstado()

    expect(fetch).toHaveBeenCalledWith(
      expect.any(String),
      expect.objectContaining({
        headers: { 'Content-Type': 'application/json', Authorization: 'Bearer token-123' },
      }),
    )
  })

  it('devuelve undefined ante una respuesta 204', async () => {
    vi.mocked(fetch).mockResolvedValue(respuestaJSON(null, 204))

    const resultado = await api.eliminarConsulta(5)

    expect(resultado).toBeUndefined()
    expect(fetch).toHaveBeenCalledWith(
      `${urlBase}/consultas/5`,
      expect.objectContaining({ method: 'DELETE' }),
    )
  })

  it('elimina un test drive propio por la ruta de baja lógica', async () => {
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ id: 4, borradoPorCliente: true }))

    const resultado = await api.eliminarTestDrive(4)

    expect(fetch).toHaveBeenCalledWith(
      `${urlBase}/test-drives/4/eliminar`,
      expect.objectContaining({ method: 'DELETE' }),
    )
    expect(resultado).toMatchObject({ id: 4, borradoPorCliente: true })
  })

  it('normaliza el error del backend a ErrorApi', async () => {
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ error: 'el mensaje no puede estar vacío' }, 400))

    try {
      await api.crearConsulta({ vehiculoId: 1, mensaje: '' })
      expect.unreachable('debería lanzar ErrorApi')
    } catch (e) {
      expect(e).toBeInstanceOf(ErrorApi)
      expect(e).toMatchObject({ estado: 400, message: 'el mensaje no puede estar vacío' })
    }
  })

  it('usa el mensaje por defecto si el cuerpo del error no es parseable', async () => {
    vi.mocked(fetch).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => {
        throw new Error('cuerpo no JSON')
      },
    } as unknown as Response)

    try {
      await api.eliminarUsuario(3)
      expect.unreachable('debería lanzar ErrorApi')
    } catch (e) {
      expect(e).toBeInstanceOf(ErrorApi)
      expect(e).toMatchObject({ estado: 500, message: 'Ocurrió un error inesperado. Intente nuevamente.' })
    }
  })

  it('arma la query string del catálogo omitiendo filtros vacíos', async () => {
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ vehiculos: [], total: 0 }))

    await api.listarVehiculos(2, 10, { marca: 'Toyota', anioMin: 2018 })

    expect(fetch).toHaveBeenCalledWith(
      `${urlBase}/vehiculos?pagina=2&tamano=10&marca=Toyota&anio_min=2018`,
      expect.anything(),
    )
  })

  it('peticionMultipart envía FormData sin fijar Content-Type', async () => {
    guardarToken('token-123')
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ respuesta: 'ok' }))

    const foto = new File(['datos'], 'auto.jpg', { type: 'image/jpeg' })
    await api.enviarTasacion([foto], 'Gol 2019')

    expect(fetch).toHaveBeenCalledTimes(1)
    const opciones = vi.mocked(fetch).mock.calls[0]![1]!
    expect(opciones).toMatchObject({ method: 'POST' })
    expect(opciones.body).toBeInstanceOf(FormData)
    const headers = opciones.headers as Record<string, string>
    expect(headers['Content-Type']).toBeUndefined()
    expect(headers.Authorization).toBe('Bearer token-123')
  })

  it('peticionMultipart normaliza el error del backend', async () => {
    vi.mocked(fetch).mockResolvedValue(respuestaJSON({ error: 'no se identificó el vehículo' }, 422))

    try {
      await api.enviarTasacion([], '')
      expect.unreachable('debería lanzar ErrorApi')
    } catch (e) {
      expect(e).toBeInstanceOf(ErrorApi)
      expect(e).toMatchObject({ estado: 422, message: 'no se identificó el vehículo' })
    }
  })
})
