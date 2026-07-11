import { useState } from 'react'
import type { VehicleFilter } from '../types'

interface Props {
  onFilter: (filter: VehicleFilter) => void
  brands: string[]
}

export function FilterPanel({ onFilter, brands }: Props) {
  const [filters, setFilters] = useState<VehicleFilter>({})
  const [showFilters, setShowFilters] = useState(false)

  const handleChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value || undefined }))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onFilter({ ...filters, page: 1 })
  }

  const handleReset = () => {
    setFilters({})
    onFilter({})
  }

  return (
    <div className="bg-white rounded-lg shadow p-4 mb-6">
      <div className="flex items-center justify-between mb-4">
        <div className="flex-1 max-w-xl">
          <input
            type="text"
            placeholder="Buscar por marca, modelo..."
            value={filters.search || ''}
            onChange={(e) => handleChange('search', e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSubmit(e)}
            className="w-full border rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <button
          onClick={() => setShowFilters(!showFilters)}
          className="ml-4 text-blue-600 hover:text-blue-800"
        >
          {showFilters ? 'Ocultar filtros' : 'Más filtros'}
        </button>
        <button
          onClick={handleSubmit}
          className="ml-2 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
        >
          Buscar
        </button>
      </div>

      {showFilters && (
        <form onSubmit={handleSubmit} className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <select
            value={filters.brand || ''}
            onChange={(e) => handleChange('brand', e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="">Todas las marcas</option>
            {brands.map((b) => (
              <option key={b} value={b}>{b}</option>
            ))}
          </select>
          <select
            value={filters.fuel || ''}
            onChange={(e) => handleChange('fuel', e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="">Todos los combustibles</option>
            <option value="nafta">Nafta</option>
            <option value="diesel">Diesel</option>
            <option value="electrico">Eléctrico</option>
            <option value="hibrido">Híbrido</option>
          </select>
          <select
            value={filters.condition || ''}
            onChange={(e) => handleChange('condition', e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="">Nuevo/Usado</option>
            <option value="nuevo">Nuevo</option>
            <option value="usado">Usado</option>
          </select>
          <select
            value={filters.vehicle_type || ''}
            onChange={(e) => handleChange('vehicle_type', e.target.value)}
            className="border rounded px-3 py-2"
          >
            <option value="">Todos los tipos</option>
            <option value="sedan">Sedán</option>
            <option value="suv">SUV</option>
            <option value="hatchback">Hatchback</option>
            <option value="pickup">Pickup</option>
            <option value="coupe">Coupé</option>
          </select>
          <input
            type="number"
            placeholder="Precio mínimo"
            value={filters.price_from || ''}
            onChange={(e) => handleChange('price_from', e.target.value)}
            className="border rounded px-3 py-2"
          />
          <input
            type="number"
            placeholder="Precio máximo"
            value={filters.price_to || ''}
            onChange={(e) => handleChange('price_to', e.target.value)}
            className="border rounded px-3 py-2"
          />
          <input
            type="number"
            placeholder="Año desde"
            value={filters.year_from || ''}
            onChange={(e) => handleChange('year_from', e.target.value)}
            className="border rounded px-3 py-2"
          />
          <input
            type="number"
            placeholder="Año hasta"
            value={filters.year_to || ''}
            onChange={(e) => handleChange('year_to', e.target.value)}
            className="border rounded px-3 py-2"
          />
          <div className="col-span-2 md:col-span-4 flex gap-2">
            <select
              value={filters.sort_by || ''}
              onChange={(e) => handleChange('sort_by', e.target.value)}
              className="border rounded px-3 py-2"
            >
              <option value="">Ordenar por</option>
              <option value="price">Precio</option>
              <option value="year">Año</option>
            </select>
            <select
              value={filters.sort_order || ''}
              onChange={(e) => handleChange('sort_order', e.target.value)}
              className="border rounded px-3 py-2"
            >
              <option value="desc">Descendente</option>
              <option value="asc">Ascendente</option>
            </select>
            <button type="button" onClick={handleReset} className="text-gray-500 hover:text-gray-700 px-3">
              Limpiar
            </button>
          </div>
        </form>
      )}
    </div>
  )
}
