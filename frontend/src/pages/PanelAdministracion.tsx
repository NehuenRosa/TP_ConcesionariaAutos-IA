import { Link } from 'react-router'

export function PanelAdministracion() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Panel de administración</h1>
      <p className="text-gray-700">Gestioná el sistema desde las siguientes secciones.</p>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <Link
          to="/admin/vehiculos"
          className="rounded-lg border border-gray-200 bg-white p-4 transition hover:shadow-md"
        >
          <h2 className="text-lg font-semibold text-gray-900">Vehículos</h2>
          <p className="mt-1 text-sm text-gray-600">
            Alta, modificación y baja del stock con ficha técnica e imágenes.
          </p>
        </Link>
      </div>
    </div>
  )
}
