import { Link } from 'react-router'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'

const SECCIONES = [
  {
    a: '/admin/vehiculos',
    titulo: 'Vehículos',
    descripcion: 'Alta, modificación y baja del stock con ficha técnica e imágenes.',
    icono: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M8 7h8M8 12h8M8 17h5M6 3h12a1 1 0 011 1v16a1 1 0 01-1 1H6a1 1 0 01-1-1V4a1 1 0 011-1z"
        />
      </svg>
    ),
  },
  {
    a: '/admin/usuarios',
    titulo: 'Usuarios',
    descripcion: 'Alta, modificación y baja de las cuentas de clientes, vendedores y administradores.',
    icono: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M16 11a4 4 0 10-8 0 4 4 0 008 0zM4 20c0-2.5 2.5-4 8-4s8 1.5 8 4"
        />
        <path strokeLinecap="round" d="M19 4v4M17 6h4" />
      </svg>
    ),
  },
]

export function PanelAdministracion() {
  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <EncabezadoPagina
        destacado="Administración"
        titulo="Panel de administración"
        descripcion="Gestioná el sistema desde las siguientes secciones."
      />

      <div className="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {SECCIONES.map((seccion) => (
          <Link
            key={seccion.a}
            to={seccion.a}
            className="group rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-white/20 hover:shadow-luz"
          >
            <div className="mb-5 flex h-12 w-12 items-center justify-center rounded-xl border border-acento-500/30 bg-acento-500/10 text-acento-400">
              {seccion.icono}
            </div>
            <h2 className="font-display text-lg font-semibold text-plata-100 transition-colors group-hover:text-white">
              {seccion.titulo}
            </h2>
            <p className="mt-1 text-sm leading-relaxed text-plata-400">{seccion.descripcion}</p>
            <span className="mt-4 inline-block text-sm font-semibold text-acento-400 opacity-0 transition-opacity group-hover:opacity-100">
              Ingresar →
            </span>
          </Link>
        ))}
      </div>
    </div>
  )
}
