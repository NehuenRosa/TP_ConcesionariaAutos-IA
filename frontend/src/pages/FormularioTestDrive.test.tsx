import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Route, Routes } from 'react-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { Vehiculo } from '../types/vehiculo'
import type { FranjaHoraria } from '../types/testDrive'
import { api, ErrorApi } from '../services/api'
import { FormularioTestDrive } from './FormularioTestDrive'

const mocks = vi.hoisted(() => ({
  obtenerVehiculo: vi.fn(),
  obtenerFranjas: vi.fn(),
  solicitarTestDrive: vi.fn(),
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

const vehiculo: Vehiculo = {
  id: 5,
  marca: 'Toyota',
  modelo: 'Corolla',
  anio: 2020,
  kilometraje: 45000,
  combustible: 'nafta',
  transmision: 'manual',
  tipo: 'sedan',
  precio: 20000,
  condicion: 'usado',
  estado: 'disponible',
  imagenes: [{ id: 1, url: 'https://ejemplo.com/foto.jpg' }],
}

const franjas: FranjaHoraria[] = [
  { id: '09:00', inicio: '09:00', fin: '10:00' },
  { id: '10:00', inicio: '10:00', fin: '11:00' },
  { id: '14:00', inicio: '14:00', fin: '15:00' },
]

// fechaMananaISO devuelve mañana en formato local: evita que las franjas de
// hoy ya hayan comenzado y queden deshabilitadas cuando corre el test.
function fechaMananaISO(): string {
  const ahora = new Date()
  const manana = new Date(ahora.getFullYear(), ahora.getMonth(), ahora.getDate() + 1)
  const mes = String(manana.getMonth() + 1).padStart(2, '0')
  const dia = String(manana.getDate()).padStart(2, '0')
  return `${manana.getFullYear()}-${mes}-${dia}`
}

function renderizar() {
  return render(
    <MemoryRouter initialEntries={['/catalogo/5/test-drive']}>
      <Routes>
        <Route path="/catalogo/:id/test-drive" element={<FormularioTestDrive />} />
      </Routes>
    </MemoryRouter>,
  )
}

async function elegirFechaManana() {
  const entrada = await screen.findByLabelText('Fecha')
  fireEvent.change(entrada, { target: { value: fechaMananaISO() } })
}

describe('FormularioTestDrive', () => {
  beforeEach(() => {
    vi.mocked(api.obtenerVehiculo).mockReset()
    vi.mocked(api.obtenerVehiculo).mockResolvedValue(vehiculo)
    vi.mocked(api.obtenerFranjas).mockReset()
    vi.mocked(api.obtenerFranjas).mockResolvedValue(franjas)
    vi.mocked(api.solicitarTestDrive).mockReset()
  })

  it('muestra las horas disponibles como botones', async () => {
    renderizar()

    expect(await screen.findByRole('button', { name: /^09:00/ })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /^10:00/ })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /^14:00/ })).toBeInTheDocument()
  })

  it('el botón de solicitar queda deshabilitado hasta elegir una hora', async () => {
    renderizar()

    const solicitar = await screen.findByRole('button', { name: 'Solicitar test drive' })
    expect(solicitar).toBeDisabled()

    await elegirFechaManana()
    const usuario = userEvent.setup()
    await usuario.click(screen.getByRole('button', { name: /^10:00/ }))

    await waitFor(() => expect(solicitar).toBeEnabled())
  })

  it('envía la hora elegida al backend y muestra la confirmación', async () => {
    const manana = fechaMananaISO()
    vi.mocked(api.solicitarTestDrive).mockResolvedValue({
      id: 1,
      vehiculo: { id: 5, marca: 'Toyota', modelo: 'Corolla', anio: 2020, precio: 20000, condicion: 'usado', tipo: 'sedan', imagen: '' },
      cliente: { id: 7, nombre: 'Ana' },
      fecha: manana,
      franja: '10:00',
      estado: 'solicitado',
    })

    renderizar()

    await elegirFechaManana()
    const usuario = userEvent.setup()
    await usuario.click(await screen.findByRole('button', { name: /^10:00/ }))
    await usuario.click(screen.getByRole('button', { name: 'Solicitar test drive' }))

    expect(await screen.findByText('Test drive solicitado')).toBeInTheDocument()
    expect(api.solicitarTestDrive).toHaveBeenCalledWith({
      vehiculoId: 5,
      fecha: manana,
      franja: '10:00',
    })
  })

  it('muestra el error del backend si la solicitud falla', async () => {
    vi.mocked(api.solicitarTestDrive).mockRejectedValue(new ErrorApi('esa hora ya pasó para hoy', 400))

    renderizar()

    await elegirFechaManana()
    const usuario = userEvent.setup()
    await usuario.click(await screen.findByRole('button', { name: /^10:00/ }))
    await usuario.click(screen.getByRole('button', { name: 'Solicitar test drive' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('esa hora ya pasó para hoy')
  })
})
