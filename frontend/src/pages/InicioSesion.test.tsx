import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, useLocation } from 'react-router'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { Usuario } from '../types/usuario'
import { ProveedorAutenticacion } from '../hooks/useAuth'
import { api, ErrorApi } from '../services/api'
import { InicioSesion } from './InicioSesion'

const mocks = vi.hoisted(() => ({
  iniciarSesion: vi.fn(),
  iniciarSesionGoogle: vi.fn(),
  obtenerProveedoresAuth: vi.fn(),
  registrar: vi.fn(),
  obtenerPerfil: vi.fn(),
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

const usuario: Usuario = { id: 1, nombre: 'Ana', email: 'ana@ejemplo.com', rol: 'cliente' }

function ProbadorUbicacion() {
  const ubicacion = useLocation()
  return <div data-testid="ubicacion">{ubicacion.pathname + ubicacion.search + ubicacion.hash}</div>
}

function renderizar(estado?: { desde?: { pathname?: string; search?: string; hash?: string } }) {
  return render(
    <MemoryRouter initialEntries={[{ pathname: '/ingreso', state: estado }]}>
      <ProveedorAutenticacion>
        <InicioSesion />
        <ProbadorUbicacion />
      </ProveedorAutenticacion>
    </MemoryRouter>,
  )
}

async function completarLogin() {
  const usuarioInteraccion = userEvent.setup()
  await usuarioInteraccion.type(screen.getByLabelText('Email'), 'ana@ejemplo.com')
  await usuarioInteraccion.type(screen.getByLabelText('Contraseña'), 'secreto123')
  await usuarioInteraccion.click(screen.getByRole('button', { name: 'Ingresar' }))
}

describe('InicioSesion', () => {
  beforeEach(() => {
    vi.mocked(api.iniciarSesion).mockReset()
    vi.mocked(api.iniciarSesion).mockResolvedValue({ token: 'token-abc', usuario })
    vi.mocked(api.obtenerProveedoresAuth).mockResolvedValue({ google: false })
  })

  it('redirige a "/" por defecto cuando no hay destino previo', async () => {
    renderizar()

    await completarLogin()

    expect(await screen.findByTestId('ubicacion')).toHaveTextContent('/')
  })

  it('conserva search y hash del destino previo', async () => {
    renderizar({ desde: { pathname: '/vehiculos/12', search: '?destacado=true', hash: '#galeria' } })

    await completarLogin()

    expect(await screen.findByTestId('ubicacion')).toHaveTextContent('/vehiculos/12?destacado=true#galeria')
  })

  it('mantiene el destino previo cuando solo tiene pathname', async () => {
    renderizar({ desde: { pathname: '/vendedor/bandeja' } })

    await completarLogin()

    expect(await screen.findByTestId('ubicacion')).toHaveTextContent('/vendedor/bandeja')
  })

  it('muestra el mensaje de error del backend cuando el login falla', async () => {
    vi.mocked(api.iniciarSesion).mockRejectedValue(new ErrorApi('credenciales inválidas', 401))
    renderizar()

    await completarLogin()

    expect(await screen.findByText('credenciales inválidas')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Ingresar' })).toBeEnabled()
  })
})
