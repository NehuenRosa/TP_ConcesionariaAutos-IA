import { Link, Outlet } from 'react-router'
import { useAuth } from '../hooks/useAuth'

export function LayoutBase() {
  const { usuario, cargando, cerrarSesion, esAdministrador } = useAuth()

  return (
    <div className="flex min-h-screen flex-col">
      <header className="border-b border-gray-200 bg-white">
        <nav className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4">
          <Link to="/" className="text-lg font-bold text-gray-900">
            Concesionaria de Autos
          </Link>
          <div className="flex items-center gap-4">
            <Link to="/catalogo" className="text-sm text-gray-700 hover:text-gray-900">
              Catálogo
            </Link>
            {esAdministrador && (
              <Link to="/admin" className="text-sm text-gray-700 hover:text-gray-900">
                Administración
              </Link>
            )}
            {!cargando && usuario && (
              <>
                <span className="text-sm text-gray-700">Hola, {usuario.nombre}</span>
                <button
                  type="button"
                  onClick={cerrarSesion}
                  className="rounded-md border border-gray-300 px-3 py-1.5 text-sm text-gray-700 hover:bg-gray-50"
                >
                  Cerrar sesión
                </button>
              </>
            )}
            {!cargando && !usuario && (
              <>
                <Link to="/registro" className="text-sm text-gray-700 hover:text-gray-900">
                  Registrarse
                </Link>
                <Link
                  to="/login"
                  className="rounded-md bg-gray-900 px-3 py-1.5 text-sm text-white hover:bg-gray-700"
                >
                  Iniciar sesión
                </Link>
              </>
            )}
          </div>
        </nav>
      </header>

      <main className="mx-auto w-full max-w-7xl flex-1 px-4 py-8">
        <Outlet />
      </main>

      <footer className="border-t border-gray-200 bg-white">
        <div className="mx-auto max-w-7xl px-4 py-4 text-center text-sm text-gray-500">
          © {new Date().getFullYear()} Concesionaria de Autos
        </div>
      </footer>
    </div>
  )
}
