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
    <div>
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <div className="flex-1 relative">
          <svg className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-text-placeholder" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            type="text"
            placeholder="Buscá por marca, modelo..."
            value={filters.search || ''}
            onChange={(e) => handleChange('search', e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSubmit(e)}
            className="input-field pl-10"
          />
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setShowFilters(!showFilters)}
            className={`btn-secondary text-sm flex items-center gap-1.5 ${showFilters ? 'bg-accent-light text-accent-text border-brand-500' : ''}`}
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
            </svg>
            Filtros
          </button>
          <button onClick={handleSubmit} className="btn-primary text-sm">
            Buscar
          </button>
        </div>
      </div>

      {showFilters && (
        <form onSubmit={handleSubmit} className="grid grid-cols-2 md:grid-cols-4 gap-3 mt-4 pt-4 border-t border-border-subtle animate-slide-up">
          <select
            value={filters.brand || ''}
            onChange={(e) => handleChange('brand', e.target.value)}
            className="input-field"
          >
            <option value="">Todas las marcas</option>
            {brands.map((b) => (
              <option key={b} value={b}>{b}</option>
            ))}
          </select>
          <select
            value={filters.fuel || ''}
            onChange={(e) => handleChange('fuel', e.target.value)}
            className="input-field"
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
            className="input-field"
          >
            <option value="">Nuevo/Usado</option>
            <option value="nuevo">Nuevo</option>
            <option value="usado">Usado</option>
          </select>
          <select
            value={filters.vehicle_type || ''}
            onChange={(e) => handleChange('vehicle_type', e.target.value)}
            className="input-field"
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
            className="input-field"
          />
          <input
            type="number"
            placeholder="Precio máximo"
            value={filters.price_to || ''}
            onChange={(e) => handleChange('price_to', e.target.value)}
            className="input-field"
          />
          <input
            type="number"
            placeholder="Año desde"
            value={filters.year_from || ''}
            onChange={(e) => handleChange('year_from', e.target.value)}
            className="input-field"
          />
          <input
            type="number"
            placeholder="Año hasta"
            value={filters.year_to || ''}
            onChange={(e) => handleChange('year_to', e.target.value)}
            className="input-field"
          />
          <div className="col-span-2 md:col-span-4 flex flex-wrap items-center gap-2">
            <select
              value={filters.sort_by || ''}
              onChange={(e) => handleChange('sort_by', e.target.value)}
              className="input-field w-auto"
            >
              <option value="">Ordenar por</option>
              <option value="price">Precio</option>
              <option value="year">Año</option>
            </select>
            <select
              value={filters.sort_order || ''}
              onChange={(e) => handleChange('sort_order', e.target.value)}
              className="input-field w-auto"
            >
              <option value="desc">Descendente</option>
              <option value="asc">Ascendente</option>
            </select>
            <button type="button" onClick={handleReset} className="text-sm text-text-placeholder hover:text-text-secondary px-3 py-2 hover:bg-surface rounded-lg transition-colors">
              Limpiar filtros
            </button>
          </div>
        </form>
      )}
    </div>
  )
}
