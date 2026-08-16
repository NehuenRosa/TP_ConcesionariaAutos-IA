import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { ConsultaResumen } from '../types/consulta'
import { ProveedorAutenticacion } from '../hooks/useAuth'
import { api } from '../services/api'
import { ChatVendedor } from './ChatVendedor'

const mocks = vi.hoisted(() => ({
  listarBandeja: vi.fn(),
  tomarConsulta: vi.fn(),
  obtenerMensajes: vi.fn(),
  marcarComoLeidos: vi.fn(),
  obtenerMensajesNuevos: vi.fn(),
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
    id: 7,
    vehiculo: {
      id: 1,
      marca: 'Toyota',
      modelo: 'Corolla',
      anio: 2020,
      precio: 20000,
      condicion: 'usado',
      tipo: 'sedan',
      imagen: '',
    },
    cliente: { id: 5, nombre: 'Ana' },
    estado: 'pendiente',
    mensajesNuevos: 0,
    createdAt: '2026-08-15T10:00:00Z',
    updatedAt: '2026-08-15T10:00:00Z',
  }
}

function renderizar() {
  return render(
    <MemoryRouter initialEntries={['/vendedor/consultas/7']}>
      <ProveedorAutenticacion>
        <Routes>
          <Route path="/vendedor/consultas/:id" element={<ChatVendedor />} />
        </Routes>
      </ProveedorAutenticacion>
    </MemoryRouter>,
  )
}

describe('ChatVendedor', () => {
  beforeEach(() => {
    vi.mocked(api.listarBandeja).mockReset()
    vi.mocked(api.tomarConsulta).mockReset()
    vi.mocked(api.tomarConsulta).mockResolvedValue(consulta())
  })

  it('muestra el botón "Tomar consulta" para una consulta pendiente', async () => {
    vi.mocked(api.listarBandeja).mockResolvedValue([consulta()])

    renderizar()

    expect(await screen.findByRole('button', { name: 'Tomar consulta' })).toBeInTheDocument()
    expect(screen.getByText('Pendiente')).toBeInTheDocument()
  })

  it('toma la consulta y recarga la bandeja al hacer clic', async () => {
    vi.mocked(api.listarBandeja).mockResolvedValue([consulta()])

    renderizar()

    const usuario = userEvent.setup()
    await usuario.click(await screen.findByRole('button', { name: 'Tomar consulta' }))

    await waitFor(() => expect(api.tomarConsulta).toHaveBeenCalledWith(7))
    expect(api.listarBandeja).toHaveBeenCalledTimes(2)
  })

  it('muestra el chat cuando la consulta ya está en conversación', async () => {
    vi.mocked(api.listarBandeja).mockResolvedValue([{ ...consulta(), estado: 'en_conversacion' }])
    vi.mocked(api.obtenerMensajes).mockResolvedValue([])
    vi.mocked(api.marcarComoLeidos).mockResolvedValue(undefined)
    vi.mocked(api.obtenerMensajesNuevos).mockResolvedValue([])

    renderizar()

    expect(await screen.findByPlaceholderText(/escribí tu mensaje/i)).toBeInTheDocument()
    expect(api.obtenerMensajes).toHaveBeenCalledWith(7)
  })

  it('avisa cuando la consulta no está disponible en la bandeja', async () => {
    vi.mocked(api.listarBandeja).mockResolvedValue([])

    renderizar()

    expect(
      await screen.findByText('La consulta no existe o no está disponible en tu bandeja.'),
    ).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Volver a la bandeja' })).toBeInTheDocument()
  })
})
