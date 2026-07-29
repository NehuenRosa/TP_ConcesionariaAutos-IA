import { useState, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'
import api from '../services/api'

export function Navbar() {
  const { isAuthenticated, user, logout, isAdmin, isSeller } = useAuth()
  const location = useLocation()
  const [menuOpen, setMenuOpen] = useState(false)
  const [scrolled, setScrolled] = useState(false)
  const [notifTotal, setNotifTotal] = useState(0)

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 10)
    window.addEventListener('scroll', handleScroll, { passive: true })
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  useEffect(() => {
    if (!isAuthenticated) return
    const fetchCount = () => {
      api.get('/consultations/notifications/count').then(({ data }) => {
        setNotifTotal(data.total ?? 0)
      }).catch(() => {})
    }
    fetchCount()
    const interval = setInterval(fetchCount, 15000)
    return () => clearInterval(interval)
  }, [isAuthenticated])

  const isActive = (path: string) => location.pathname.startsWith(path)

  return (
    <nav className={`bg-white border-b border-border-subtle sticky top-0 z-40 transition-shadow duration-200 ${scrolled ? 'shadow-sm' : ''}`}>
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-14 items-center">
          <div className="flex items-center gap-6">
            <Link to="/" className="flex items-center gap-2">
              <div className="w-8 h-8 bg-brand-500 rounded-lg flex items-center justify-center">
                <svg className="w-4 h-4 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <span className="text-base font-semibold text-text-primary">AutoPrime</span>
            </Link>
            <div className="hidden md:flex items-center gap-1">
              <Link
                to="/catalogo"
                className={isActive('/catalogo') ? 'nav-link-active' : 'nav-link-inactive'}
              >
                Catálogo
              </Link>
              {isAdmin && (
                <Link
                  to="/admin/dashboard"
                  className={isActive('/admin') ? 'nav-link-active' : 'nav-link-inactive'}
                >
                  Admin
                </Link>
              )}
              {isSeller && (
                <Link
                  to="/seller/consultas"
                  className={`relative ${isActive('/seller') ? 'nav-link-active' : 'nav-link-inactive'}`}
                >
                  Bandeja
                  {notifTotal > 0 && (
                    <span className="absolute -top-1.5 -right-3 bg-red-500 text-white text-[10px] font-bold rounded-full w-4 h-4 flex items-center justify-center">
                      {notifTotal}
                    </span>
                  )}
                </Link>
              )}
              {isAuthenticated && !isSeller && !isAdmin && (
                <Link
                  to="/mis-consultas"
                  className={`relative ${isActive('/mis-consultas') ? 'nav-link-active' : 'nav-link-inactive'}`}
                >
                  Mis Consultas
                  {notifTotal > 0 && (
                    <span className="absolute -top-1.5 -right-3 bg-red-500 text-white text-[10px] font-bold rounded-full w-4 h-4 flex items-center justify-center">
                      {notifTotal}
                    </span>
                  )}
                </Link>
              )}
            </div>
          </div>

          <div className="hidden md:flex items-center gap-3">
            {isAuthenticated ? (
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2 px-2.5 py-1.5 bg-surface rounded-lg">
                  <div className="w-6 h-6 bg-brand-500 rounded-md flex items-center justify-center">
                    <span className="text-xs font-semibold text-white">{user?.name?.charAt(0).toUpperCase()}</span>
                  </div>
                  <span className="text-sm font-medium text-text-primary">{user?.name}</span>
                </div>
                <button
                  onClick={logout}
                  className="text-sm font-medium text-text-secondary hover:text-red-600 px-2.5 py-1.5 rounded-lg hover:bg-red-50 transition-colors"
                >
                  Salir
                </button>
              </div>
            ) : (
              <>
                <Link to="/login" className="text-sm font-medium text-text-secondary hover:text-text-primary px-3 py-1.5 transition-colors">
                  Ingresar
                </Link>
                <Link to="/register" className="btn-primary text-sm">
                  Registrarse
                </Link>
              </>
            )}
          </div>

          <button
            onClick={() => setMenuOpen(!menuOpen)}
            className="md:hidden p-2 rounded-lg text-text-secondary hover:bg-surface transition-colors"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              {menuOpen ? (
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              ) : (
                <path strokeLinecap="round" strokeLinejoin="round" d="M4 6h16M4 12h16M4 18h16" />
              )}
            </svg>
          </button>
        </div>

        {menuOpen && (
          <div className="md:hidden pb-3 border-t border-border-subtle mt-2 pt-3 space-y-1 animate-fade-in">
            <Link to="/catalogo" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-sm font-medium text-text-primary hover:bg-surface">
              Catálogo
            </Link>
            {isAdmin && (
              <Link to="/admin/dashboard" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-sm font-medium text-text-primary hover:bg-surface">
                Admin
              </Link>
            )}
            {isSeller && (
              <Link to="/seller/consultas" onClick={() => setMenuOpen(false)} className="relative block px-3 py-2 rounded-lg text-sm font-medium text-text-primary hover:bg-surface">
                Bandeja
                {notifTotal > 0 && (
                  <span className="ml-1.5 bg-red-500 text-white text-[10px] font-bold rounded-full px-1.5 py-0.5 inline-flex items-center justify-center">
                    {notifTotal}
                  </span>
                )}
              </Link>
            )}
            {isAuthenticated && !isSeller && !isAdmin && (
              <Link to="/mis-consultas" onClick={() => setMenuOpen(false)} className="relative block px-3 py-2 rounded-lg text-sm font-medium text-text-primary hover:bg-surface">
                Mis Consultas
                {notifTotal > 0 && (
                  <span className="ml-1.5 bg-red-500 text-white text-[10px] font-bold rounded-full px-1.5 py-0.5 inline-flex items-center justify-center">
                    {notifTotal}
                  </span>
                )}
              </Link>
            )}
            {isAuthenticated ? (
              <div className="border-t border-border-subtle pt-2 mt-2">
                <div className="px-3 py-2 text-sm text-text-secondary">{user?.name}</div>
                <button onClick={() => { logout(); setMenuOpen(false) }} className="w-full text-left px-3 py-2 rounded-lg text-sm font-medium text-red-600 hover:bg-red-50">
                  Salir
                </button>
              </div>
            ) : (
              <div className="border-t border-border-subtle pt-2 mt-2 space-y-1">
                <Link to="/login" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-sm font-medium text-text-primary hover:bg-surface">
                  Ingresar
                </Link>
                <Link to="/register" onClick={() => setMenuOpen(false)} className="block px-3 py-2 rounded-lg text-sm font-medium text-white bg-brand-500 hover:bg-brand-600 text-center rounded-lg">
                  Registrarse
                </Link>
              </div>
            )}
          </div>
        )}
      </div>
    </nav>
  )
}
