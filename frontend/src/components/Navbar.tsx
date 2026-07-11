import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export function Navbar() {
  const { isAuthenticated, user, logout, isAdmin, isSeller } = useAuth()

  return (
    <nav className="bg-white shadow-md">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between h-16 items-center">
          <div className="flex items-center gap-8">
            <Link to="/" className="text-xl font-bold text-blue-600">
              Concesionaria
            </Link>
            <div className="hidden md:flex gap-4">
              <Link to="/catalogo" className="text-gray-600 hover:text-blue-600">
                Catálogo
              </Link>
            </div>
          </div>

          <div className="flex items-center gap-4">
            {isAdmin && (
              <Link to="/admin/dashboard" className="text-gray-600 hover:text-blue-600">
                Admin
              </Link>
            )}
            {isSeller && (
              <Link to="/seller/consultas" className="text-gray-600 hover:text-blue-600">
                Bandeja
              </Link>
            )}
            {isAuthenticated ? (
              <>
                <span className="text-sm text-gray-500">{user?.name}</span>
                <button onClick={logout} className="text-red-500 hover:text-red-700">
                  Salir
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-gray-600 hover:text-blue-600">
                  Ingresar
                </Link>
                <Link
                  to="/register"
                  className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                >
                  Registrarse
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  )
}
