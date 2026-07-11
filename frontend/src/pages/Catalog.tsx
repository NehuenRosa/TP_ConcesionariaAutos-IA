import { useState, useEffect } from 'react'
import { VehicleCard } from '../components/VehicleCard'
import { FilterPanel } from '../components/FilterPanel'
import { Pagination } from '../components/Pagination'
import { vehicleService } from '../services/vehicleService'
import type { Vehicle, VehicleFilter } from '../types'

export function Catalog() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [brands, setBrands] = useState<string[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [filter, setFilter] = useState<VehicleFilter>({})
  const [loading, setLoading] = useState(true)
  const pageSize = 12

  useEffect(() => {
    vehicleService.getBrands().then(setBrands).catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    vehicleService.list({ ...filter, page, page_size: pageSize }).then((res) => {
      setVehicles(res.data)
      setTotal(res.total)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [filter, page])

  const handleFilter = (newFilter: VehicleFilter) => {
    setFilter(newFilter)
    setPage(1)
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Catálogo de Vehículos</h1>
      <FilterPanel onFilter={handleFilter} brands={brands} />
      {loading ? (
        <div className="text-center py-12 text-gray-500">Cargando...</div>
      ) : vehicles.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          No se encontraron vehículos con los filtros seleccionados.
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {vehicles.map((v: Vehicle) => (
              <VehicleCard key={v.id} vehicle={v} />
            ))}
          </div>
          <Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} />
        </>
      )}
    </div>
  )
}
