import { useEffect, useState } from 'react'
import { Link, NavLink, Outlet } from 'react-router'
import { useAuth } from '../hooks/useAuth'
import { useNotificaciones } from '../hooks/useNotificaciones'
import { Boton } from '../components/ui/Boton'
import { Chatbot } from '../components/Chatbot'

function Marca() {
  return (
    <Link to="/" className="group flex items-center gap-2.5">
      <span className="flex h-8 w-8 items-center justify-center border border-acento-500/60 bg-carbono-900">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.4" className="h-4.5 w-4.5 text-acento-400">
          <path d="M3 17l6-9 4 5.5L17 7l4 10H3z" strokeLinejoin="round" />
        </svg>
      </span>
      <span className="font-display text-sm font-extrabold tracking-[0.3em] text-plata-100 uppercase">
        Aurum<span className="text-acento-400">·</span>Motors
      </span>
    </Link>
  )
}

export function LayoutBase() {
  const { usuario, cargando, cerrarSesion, esAdministrador } = useAuth()
  const esCliente = usuario?.rol === 'cliente'
  const esVendedor = usuario?.rol === 'vendedor'
  const { cantidad, nuevoAviso, descartarAviso } = useNotificaciones(!cargando && !!usuario)
  const [menuAbierto, setMenuAbierto] = useState(false)

  const destinoMensajes = esCliente ? '/mis-consultas' : '/vendedor/bandeja'

  useEffect(() => {
    if (!nuevoAviso) return
    const temporizador = setTimeout(descartarAviso, 8000)
    return () => clearTimeout(temporizador)
  }, [nuevoAviso, descartarAviso])

  const claseEnlace = ({ isActive }: { isActive: boolean }) =>
    `relative px-1 py-1.5 font-display text-[13px] font-medium tracking-wide transition-colors ${
      isActive ? 'text-plata-100' : 'text-plata-400 hover:text-plata-100'
    }`

  const enlaces = [
    { a: '/catalogo', texto: 'Catálogo', mostrar: true },
    { a: '/mis-consultas', texto: 'Mis consultas', mostrar: esCliente, notificacion: esCliente },
    { a: '/mis-test-drives', texto: 'Test drives', mostrar: esCliente },
    { a: '/mis-reservas', texto: 'Reservas', mostrar: esCliente },
    { a: '/vendedor/bandeja', texto: 'Bandeja', mostrar: esVendedor, notificacion: esVendedor },
    { a: '/vendedor/test-drives', texto: 'Test drives', mostrar: esVendedor },
    { a: '/vendedor/reservas', texto: 'Reservas', mostrar: esVendedor },
    { a: '/admin', texto: 'Administración', mostrar: esAdministrador },
  ].filter((enlace) => enlace.mostrar)

  return (
    <div className="flex min-h-screen flex-col">
      <header className="fixed inset-x-0 top-0 z-50 border-b border-white/8 bg-carbono-950/80 backdrop-blur-xl">
        <nav className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
          <Marca />

          <div className="hidden items-center gap-7 md:flex">
            {enlaces.map((enlace) => (
              <NavLink key={enlace.a} to={enlace.a} className={claseEnlace} end={enlace.a === '/admin'}>
                {enlace.texto}
                {enlace.notificacion && cantidad > 0 && (
                  <span className="absolute -top-0.5 -right-2 flex h-2 w-2">
                    <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-acento-400 opacity-60" />
                    <span className="relative inline-flex h-2 w-2 rounded-full bg-acento-400" />
                  </span>
                )}
              </NavLink>
            ))}
          </div>

          <div className="hidden items-center gap-3 md:flex">
            {!cargando && usuario && (
              <>
                <span className="hidden text-sm text-plata-400 lg:block">
                  Hola, <span className="font-medium text-plata-100">{usuario.nombre}</span>
                </span>
                <Boton variante="secundario" tamano="sm" onClick={cerrarSesion}>
                  Cerrar sesión
                </Boton>
              </>
            )}
            {!cargando && !usuario && (
              <>
                <Link
                  to="/registro"
                  className="font-display text-[13px] font-medium text-plata-300 transition-colors hover:text-plata-100"
                >
                  Registrarse
                </Link>
                <Boton tamano="sm">
                  <Link to="/login">Iniciar sesión</Link>
                </Boton>
              </>
            )}
          </div>

          <button
            type="button"
            onClick={() => setMenuAbierto((a) => !a)}
            className="inline-flex h-10 w-10 items-center justify-center rounded-md border border-white/10 text-plata-200 md:hidden"
            aria-label="Abrir menú"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-5 w-5">
              {menuAbierto ? (
                <path strokeLinecap="round" d="M6 18L18 6M6 6l12 12" />
              ) : (
                <path strokeLinecap="round" d="M4 7h16M4 12h16M4 17h16" />
              )}
            </svg>
          </button>
        </nav>

        {menuAbierto && (
          <div className="border-t border-white/8 bg-carbono-900/95 px-4 py-4 backdrop-blur-xl md:hidden">
            <div className="flex flex-col gap-1">
              {enlaces.map((enlace) => (
                <NavLink
                  key={enlace.a}
                  to={enlace.a}
                  onClick={() => setMenuAbierto(false)}
                  className={({ isActive }) =>
                    `rounded-md px-3 py-2.5 font-display text-sm font-medium ${
                      isActive ? 'bg-white/5 text-plata-100' : 'text-plata-400 hover:bg-white/5 hover:text-plata-100'
                    }`
                  }
                  end={enlace.a === '/admin'}
                >
                  {enlace.texto}
                </NavLink>
              ))}
              <div className="mt-3 flex flex-col gap-2 border-t border-white/8 pt-3">
                {!cargando && usuario ? (
                  <>
                    <p className="px-3 text-sm text-plata-400">
                      Hola, <span className="text-plata-100">{usuario.nombre}</span>
                    </p>
                    <Boton variante="secundario" tamano="sm" onClick={cerrarSesion}>
                      Cerrar sesión
                    </Boton>
                  </>
                ) : (
                  !cargando && (
                    <>
                      <Boton variante="secundario" tamano="sm">
                        <Link to="/registro" onClick={() => setMenuAbierto(false)}>
                          Registrarse
                        </Link>
                      </Boton>
                      <Boton tamano="sm">
                        <Link to="/login" onClick={() => setMenuAbierto(false)}>
                          Iniciar sesión
                        </Link>
                      </Boton>
                    </>
                  )
                )}
              </div>
            </div>
          </div>
        )}
      </header>

      <main className="flex-1 pt-16">
        <Outlet />
      </main>

      <footer className="border-t border-white/8 bg-carbono-900/60">
        <div className="mx-auto max-w-7xl px-4 py-14 sm:px-6 lg:px-8">
          <div className="grid gap-10 md:grid-cols-4">
            <div className="md:col-span-2">
              <Marca />
              <p className="mt-4 max-w-sm text-sm leading-relaxed text-plata-400">
                Concesionaria de autos premium. Explorá el catálogo, solicitá una cotización o
                agendá tu test drive con nuestro asistente digital.
              </p>
            </div>
            <div>
              <p className="font-display text-xs font-semibold tracking-[0.25em] text-plata-500 uppercase">
                Explorar
              </p>
              <ul className="mt-4 space-y-2.5 text-sm">
                <li>
                  <Link to="/catalogo" className="text-plata-300 transition-colors hover:text-plata-100">
                    Catálogo
                  </Link>
                </li>
                {usuario?.rol === 'cliente' && (
                  <li>
                    <Link to="/mis-consultas" className="text-plata-300 transition-colors hover:text-plata-100">
                      Mis consultas
                    </Link>
                  </li>
                )}
                {usuario?.rol === 'vendedor' && (
                  <li>
                    <Link to="/vendedor/bandeja" className="text-plata-300 transition-colors hover:text-plata-100">
                      Bandeja de consultas
                    </Link>
                  </li>
                )}
                {!usuario && !cargando && (
                  <li>
                    <Link to="/registro" className="text-plata-300 transition-colors hover:text-plata-100">
                      Crear cuenta
                    </Link>
                  </li>
                )}
                {!usuario && !cargando && (
                  <li>
                    <Link to="/login" className="text-plata-300 transition-colors hover:text-plata-100">
                      Iniciar sesión
                    </Link>
                  </li>
                )}
              </ul>
            </div>
            <div>
              <p className="font-display text-xs font-semibold tracking-[0.25em] text-plata-500 uppercase">
                Contacto
              </p>
              <ul className="mt-4 space-y-2.5 text-sm text-plata-300">
                <li>Av. de los Autos 1234</li>
                <li>+54 11 5555 0199</li>
                <li>hola@aurummotors.com.ar</li>
              </ul>
            </div>
          </div>
          <div className="mt-12 flex flex-col items-center justify-between gap-3 border-t border-white/8 pt-6 sm:flex-row">
            <p className="text-xs text-plata-500">
              © {new Date().getFullYear()} Aurum Motors. Todos los derechos reservados.
            </p>
            <p className="text-xs text-plata-600">Diseñado para apasionados del motor.</p>
          </div>
        </div>
      </footer>

      {nuevoAviso && (esCliente || esVendedor) && (
        <div className="fixed right-4 bottom-4 z-50 flex w-[calc(100%-2rem)] max-w-sm items-center gap-3 rounded-2xl border border-acento-500/40 bg-carbono-850/95 p-4 shadow-luz backdrop-blur-xl sm:right-6 sm:bottom-6">
          <span className="relative flex h-3 w-3 shrink-0">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-acento-400 opacity-60" />
            <span className="relative inline-flex h-3 w-3 rounded-full bg-acento-400" />
          </span>
          <div className="min-w-0 flex-1">
            <p className="font-display text-sm font-semibold text-plata-100">Tenés mensajes nuevos</p>
            <p className="text-xs text-plata-400">Te escribieron en una consulta.</p>
          </div>
          <Link
            to={destinoMensajes}
            onClick={descartarAviso}
            className="shrink-0 text-sm font-semibold text-acento-400 transition-colors hover:text-acento-300"
          >
            Ver
          </Link>
          <button
            type="button"
            onClick={descartarAviso}
            className="shrink-0 text-plata-500 transition-colors hover:text-plata-200"
            aria-label="Cerrar aviso"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-4 w-4">
              <path strokeLinecap="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      )}

      <Chatbot />
    </div>
  )
}
