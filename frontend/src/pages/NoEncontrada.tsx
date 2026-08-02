import { Link } from 'react-router'

export function NoEncontrada() {
  return (
    <div className="space-y-6 text-center">
      <h1 className="text-3xl font-bold text-gray-900">Página no encontrada</h1>
      <p className="text-gray-700">La página que buscás no existe.</p>
      <Link to="/" className="inline-block text-gray-900 underline hover:text-gray-700">
        Volver al inicio
      </Link>
    </div>
  )
}
