import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { ConsultaResumen } from '../types/consulta'
import { ProveedorAutenticacion } from '../hooks/useAuth'
import { api } from '../services/api'
import { MisConsultas } from './MisConsultas'

const mocks = vi.hoisted(() => ({
  listarMisConsultas: vi.fn(),
  obtenerMensajes: vi.fn(),
  marcarComoLeidos: vi.fn(),
  obtenerMensajesNuevos: vi.fn(),
  enviarMensaje: vi.fn(),
}))

vi.mock('../services/api', () => ({
  ErrorApi: class ErrorApi extends Error {
    estado: number
    constructor(mensaje: string, estado: number) {
      super(mensaje)
      this.name = 'ErrorApi'
      this.estado = estado
    }
  },
  api: mocks,
  obtenerToken: vi.fn(() => null),
  guardarToken: vi.fn(),
  eliminarToken: vi.fn(),
}))

function consulta(): ConsultaResumen {
  return {
    id: 1,
    vehiculo: {
      id: 2,
      marca: 'Toyota',
      modelo: 'Corolla',
      anio: 2020,
      precio: 20000,
      condicion: 'usado',
      tipo: 'sedan',
      imagen: '',
    },
    cliente: { id: 5, nombre: 'Ana' },
    estado: 'en_conversacion',
    mensajesNuevos: 0,
    createdAt: '2026-08-15T10:00:00Z',
    updatedAt: '2026-08-15T10:00:00Z',
  }
}

function renderizar(ruta: string) {
  return render(
    <MemoryRouter initialEntries={[ruta]}>
      <ProveedorAutenticacion>
        <Routes>
          <Route path="/mis-consultas" element={<MisConsultas />} />
          <Route path="/mis-consultas/:id" element={<MisConsultas />} />
        </Routes>
      </ProveedorAutenticacion>
    </MemoryRouter>,
  )
}

describe('MisConsultas', () => {
  beforeEach(() => {
    vi.mocked(api.listarMisConsultas).mockReset()
    vi.mocked(api.listarMisConsultas).mockResolvedValue([])
  })

  it('muestra el estado vacío cuando no hay consultas', async () => {
    renderizar('/mis-consultas')

    expect(await screen.findByText('No tenés consultas todavía.')).toBeInTheDocument()
  })

  it('muestra "Consulta no encontrada" cuando el id no pertenece al cliente', async () => {
    renderizar('/mis-consultas/999')

    expect(await screen.findByText('Consulta no encontrada')).toBeInTheDocument()
    expect(screen.getByText('La consulta no existe o no te pertenece.')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Volver a mis consultas' })).toBeInTheDocument()
  })

  it('renderiza la conversación cuando el id coincide con una consulta', async () => {
    vi.mocked(api.listarMisConsultas).mockResolvedValue([consulta()])
    vi.mocked(api.obtenerMensajes).mockResolvedValue([])
    vi.mocked(api.marcarComoLeidos).mockResolvedValue(undefined)
    vi.mocked(api.obtenerMensajesNuevos).mockResolvedValue([])

    renderizar('/mis-consultas/1')

    expect(await screen.findByPlaceholderText(/escribí tu mensaje/i)).toBeInTheDocument()
    expect(api.obtenerMensajes).toHaveBeenCalledWith(1)
  })
})
