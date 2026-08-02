import { useParams } from 'react-router'

export function DetalleVehiculo() {
  const { id } = useParams()

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-900">Detalle del vehículo</h1>
      <p className="text-gray-700">
        El detalle del vehículo {id ?? ''} estará disponible al implementar el caso de uso CU-03.
      </p>
    </div>
  )
}
