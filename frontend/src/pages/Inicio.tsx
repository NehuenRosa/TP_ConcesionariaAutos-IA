import { Link } from 'react-router'

export function Inicio() {
  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Bienvenido a la Concesionaria</h1>
      <p className="text-gray-700">
        Explorá nuestro catálogo de vehículos disponibles y reservá tu próxima unidad.
      </p>
      <Link
        to="/catalogo"
        className="inline-block rounded-md bg-gray-900 px-4 py-2 text-white hover:bg-gray-700"
      >
        Ver catálogo
      </Link>
    </div>
  )
}
