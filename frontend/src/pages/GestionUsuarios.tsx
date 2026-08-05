import { useEffect, useState } from 'react'
import { api, ErrorApi } from '../services/api'
import type { DatosUsuarioAdmin, Rol, Usuario } from '../types/usuario'
import { ETIQUETAS_ROL, ROLES } from '../types/usuario'
import { Boton } from '../components/ui/Boton'
import { CampoSeleccion, CampoTexto } from '../components/ui/Campo'
import { ContenidoCargando } from '../components/ui/Spinner'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { EstadoVacio } from '../components/ui/EstadoVacio'

const ESTILOS_ROL: Record<Rol, string> = {
  cliente: 'border-plata-400/25 bg-plata-400/10 text-plata-300',
  vendedor: 'border-acento-500/40 bg-acento-500/10 text-acento-400',
  administrador: 'border-white/25 bg-white/10 text-plata-100',
}

interface EstadoFormulario {
  id: number | null
  nombre: string
  email: string
  password: string
  rol: Rol
}

const FORMULARIO_VACIO: EstadoFormulario = {
  id: null,
  nombre: '',
  email: '',
  password: '',
  rol: 'cliente',
}

export function GestionUsuarios() {
  const [usuarios, setUsuarios] = useState<Usuario[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formularioAbierto, setFormularioAbierto] = useState(false)
  const [formulario, setFormulario] = useState<EstadoFormulario>(FORMULARIO_VACIO)
  const [guardando, setGuardando] = useState(false)
  const [errorFormulario, setErrorFormulario] = useState<string | null>(null)

  const cargarUsuarios = () => {
    setCargando(true)
    setError(null)
    api
      .listarUsuarios()
      .then(setUsuarios)
      .catch((e: unknown) => {
        setError(e instanceof ErrorApi ? e.message : 'Ocurrió un error inesperado al cargar los usuarios.')
      })
      .finally(() => setCargando(false))
  }

  useEffect(() => {
    cargarUsuarios()
  }, [])

  function abrirNuevo() {
    setFormulario(FORMULARIO_VACIO)
    setErrorFormulario(null)
    setFormularioAbierto(true)
  }

  function abrirEdicion(usuario: Usuario) {
    setFormulario({
      id: usuario.id,
      nombre: usuario.nombre,
      email: usuario.email,
      password: '',
      rol: usuario.rol,
    })
    setErrorFormulario(null)
    setFormularioAbierto(true)
  }

  async function guardarUsuario(e: React.FormEvent) {
    e.preventDefault()
    setGuardando(true)
    setErrorFormulario(null)

    const datos: DatosUsuarioAdmin = {
      nombre: formulario.nombre,
      email: formulario.email,
      rol: formulario.rol,
      ...(formulario.password ? { password: formulario.password } : {}),
    }

    try {
      if (formulario.id === null) {
        await api.crearUsuario(datos)
      } else {
        await api.actualizarUsuario(formulario.id, datos)
      }
      setFormularioAbierto(false)
      cargarUsuarios()
    } catch (e: unknown) {
      setErrorFormulario(e instanceof ErrorApi ? e.message : 'No se pudo guardar el usuario.')
    } finally {
      setGuardando(false)
    }
  }

  async function eliminar(usuario: Usuario) {
    const confirmacion = window.confirm(
      `¿Eliminar la cuenta de ${usuario.nombre} (${usuario.email})? Esta acción no se puede deshacer.`,
    )
    if (!confirmacion) return

    try {
      await api.eliminarUsuario(usuario.id)
      setUsuarios((actuales) => actuales.filter((u) => u.id !== usuario.id))
    } catch (e: unknown) {
      setError(e instanceof ErrorApi ? e.message : 'No se pudo eliminar el usuario.')
    }
  }

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <EncabezadoPagina
        destacado="Administración"
        titulo="Gestión de usuarios"
        descripcion="Dá de alta, modificá o eliminá las cuentas del sistema."
        acciones={
          <Boton onClick={abrirNuevo}>+ Nuevo usuario</Boton>
        }
      />

      {cargando && <ContenidoCargando etiqueta="Cargando usuarios…" />}

      {!cargando && error && (
        <div className="mt-6">
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">{error}</div>
        </div>
      )}

      {!cargando && !error && usuarios.length === 0 && (
        <div className="mt-8">
          <EstadoVacio
            titulo="No hay usuarios para mostrar"
            descripcion="Creá la primera cuenta para comenzar a operar el sistema."
            accion={<Boton onClick={abrirNuevo}>+ Nuevo usuario</Boton>}
          />
        </div>
      )}

      {!cargando && !error && usuarios.length > 0 && (
        <div className="mt-6 overflow-hidden rounded-2xl border border-white/8 bg-carbono-850/50">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="border-b border-white/8 bg-carbono-900/70 text-plata-500">
                <tr>
                  <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Usuario</th>
                  <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Email</th>
                  <th className="px-5 py-3.5 font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Rol</th>
                  <th className="px-5 py-3.5 text-right font-display text-[11px] font-semibold tracking-[0.15em] uppercase">Acciones</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/6">
                {usuarios.map((usuario) => (
                  <tr key={usuario.id} className="transition-colors hover:bg-white/3">
                    <td className="px-5 py-3.5">
                      <p className="font-display font-semibold text-plata-100">{usuario.nombre}</p>
                      <p className="text-xs text-plata-500">ID {usuario.id}</p>
                    </td>
                    <td className="px-5 py-3.5 text-plata-300">{usuario.email}</td>
                    <td className="px-5 py-3.5">
                      <span
                        className={`inline-flex rounded-full border px-2.5 py-0.5 text-xs font-semibold ${ESTILOS_ROL[usuario.rol]}`}
                      >
                        {ETIQUETAS_ROL[usuario.rol]}
                      </span>
                    </td>
                    <td className="px-5 py-3.5">
                      <div className="flex items-center justify-end gap-2">
                        <Boton variante="secundario" tamano="sm" onClick={() => abrirEdicion(usuario)}>
                          Editar
                        </Boton>
                        <Boton variante="peligro" tamano="sm" onClick={() => eliminar(usuario)}>
                          Eliminar
                        </Boton>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {formularioAbierto && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-carbono-950/80 p-4 backdrop-blur-sm">
          <div className="w-full max-w-lg rounded-2xl border border-white/10 bg-carbono-850 p-6 shadow-luz">
            <h2 className="font-display text-xl font-bold text-plata-100">
              {formulario.id === null ? 'Nuevo usuario' : 'Editar usuario'}
            </h2>
            <p className="mt-1 text-sm text-plata-400">
              {formulario.id === null
                ? 'Completá los datos para crear la cuenta.'
                : 'Modificá los datos y guardá los cambios.'}
            </p>

            <form onSubmit={guardarUsuario} className="mt-6 space-y-4">
              {errorFormulario && (
                <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-3 text-sm text-red-300">
                  {errorFormulario}
                </div>
              )}

              <CampoTexto
                id="nombre-usuario"
                etiqueta="Nombre completo"
                required
                value={formulario.nombre}
                onChange={(e) => setFormulario((f) => ({ ...f, nombre: e.target.value }))}
              />
              <CampoTexto
                id="email-usuario"
                etiqueta="Email"
                type="email"
                required
                autoComplete="off"
                value={formulario.email}
                onChange={(e) => setFormulario((f) => ({ ...f, email: e.target.value }))}
              />
              <CampoTexto
                id="password-usuario"
                etiqueta={formulario.id === null ? 'Contraseña' : 'Nueva contraseña (opcional)'}
                type="password"
                required={formulario.id === null}
                minLength={8}
                placeholder={formulario.id === null ? 'Mínimo 8 caracteres' : 'Dejalo vacío para no cambiarla'}
                autoComplete="new-password"
                value={formulario.password}
                onChange={(e) => setFormulario((f) => ({ ...f, password: e.target.value }))}
              />
              <CampoSeleccion
                id="rol-usuario"
                etiqueta="Rol"
                value={formulario.rol}
                onChange={(e) => setFormulario((f) => ({ ...f, rol: e.target.value as Rol }))}
              >
                {ROLES.map((rol) => (
                  <option key={rol} value={rol}>
                    {ETIQUETAS_ROL[rol]}
                  </option>
                ))}
              </CampoSeleccion>

              <div className="flex items-center justify-end gap-3 pt-2">
                <Boton
                  variante="secundario"
                  onClick={() => setFormularioAbierto(false)}
                  disabled={guardando}
                >
                  Cancelar
                </Boton>
                <Boton type="submit" disabled={guardando}>
                  {guardando ? 'Guardando…' : 'Guardar'}
                </Boton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
