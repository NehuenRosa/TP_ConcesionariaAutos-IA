import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { vehicleService } from '../../services/vehicleService'

export function VehicleForm() {
  const { id } = useParams()
  const navigate = useNavigate()
  const isEditing = !!id
  const [form, setForm] = useState({
    brand: '', model: '', year: new Date().getFullYear(), price: 0, mileage: 0,
    fuel: 'nafta', transmission: 'manual', condition: 'nuevo', color: '',
    description: '', images: '', vehicle_type: 'sedan',
  })
  const [error, setError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  useEffect(() => {
    if (id) {
      vehicleService.getById(Number(id)).then((v) => {
        setForm({
          brand: v.brand, model: v.model, year: v.year, price: v.price, mileage: v.mileage,
          fuel: v.fuel, transmission: v.transmission, condition: v.condition, color: v.color || '',
          description: v.description || '', images: v.images?.join('\n') || '', vehicle_type: v.vehicle_type,
        })
      }).catch(() => navigate('/admin/vehiculos'))
    }
  }, [id])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setSubmitting(true)
    try {
      const payload: Record<string, unknown> = {
        ...form,
        images: form.images.split('\n').filter(Boolean),
      }
      if (isEditing) {
        await vehicleService.update(Number(id), payload as any)
      } else {
        await vehicleService.create(payload as any)
      }
      navigate('/admin/vehiculos')
    } catch (err: any) {
      setError(err.response?.data?.error || 'Error al guardar el vehículo')
    } finally {
      setSubmitting(false)
    }
  }

  const update = (key: string, value: unknown) => setForm((prev) => ({ ...prev, [key]: value }))

  return (
    <div className="max-w-3xl mx-auto px-4 sm:px-6 py-8 animate-fade-in">
      <Link to="/admin/vehiculos" className="inline-flex items-center gap-1.5 text-sm font-medium text-gray-500 hover:text-brand-600 mb-6 transition-colors">
        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
        Volver a gestión
      </Link>

      <div className="card p-8">
        <div className="flex items-center gap-3 mb-8">
          <div className="w-12 h-12 bg-brand-100 rounded-xl flex items-center justify-center">
            <svg className="w-6 h-6 text-brand-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">{isEditing ? 'Editar' : 'Nuevo'} Vehículo</h1>
            <p className="text-sm text-gray-500">{isEditing ? 'Modificá los datos del vehículo' : 'Completá los datos para agregar un vehículo'}</p>
          </div>
        </div>

        {error && (
          <div className="flex items-center gap-2 bg-red-50 text-red-700 p-3 rounded-xl mb-6 text-sm border border-red-100">
            <svg className="w-4 h-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Marca</label>
              <input placeholder="Ej: Toyota" value={form.brand} onChange={(e) => update('brand', e.target.value)} required className="input-field" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Modelo</label>
              <input placeholder="Ej: Corolla" value={form.model} onChange={(e) => update('model', e.target.value)} required className="input-field" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Año</label>
              <input type="number" placeholder="2024" value={form.year} onChange={(e) => update('year', Number(e.target.value))} required className="input-field" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Precio ($)</label>
              <input type="number" placeholder="25000000" value={form.price} onChange={(e) => update('price', Number(e.target.value))} required className="input-field" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Kilometraje</label>
              <input type="number" placeholder="0" value={form.mileage} onChange={(e) => update('mileage', Number(e.target.value))} className="input-field" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Color</label>
              <input placeholder="Ej: Blanco" value={form.color} onChange={(e) => update('color', e.target.value)} className="input-field" />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Combustible</label>
              <select value={form.fuel} onChange={(e) => update('fuel', e.target.value)} className="input-field">
                <option value="nafta">Nafta</option>
                <option value="diesel">Diesel</option>
                <option value="electrico">Eléctrico</option>
                <option value="hibrido">Híbrido</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Transmisión</label>
              <select value={form.transmission} onChange={(e) => update('transmission', e.target.value)} className="input-field">
                <option value="manual">Manual</option>
                <option value="automatico">Automático</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Condición</label>
              <select value={form.condition} onChange={(e) => update('condition', e.target.value)} className="input-field">
                <option value="nuevo">Nuevo</option>
                <option value="usado">Usado</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">Tipo</label>
              <select value={form.vehicle_type} onChange={(e) => update('vehicle_type', e.target.value)} className="input-field">
                <option value="sedan">Sedán</option>
                <option value="suv">SUV</option>
                <option value="hatchback">Hatchback</option>
                <option value="pickup">Pickup</option>
                <option value="coupe">Coupé</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">Descripción</label>
            <textarea placeholder="Descripción del vehículo..." value={form.description} onChange={(e) => update('description', e.target.value)} rows={3} className="input-field resize-none" />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">URLs de imágenes</label>
            <textarea placeholder="https://ejemplo.com/imagen1.jpg&#10;https://ejemplo.com/imagen2.jpg" value={form.images} onChange={(e) => update('images', e.target.value)} rows={3} className="input-field resize-none font-mono text-xs" />
            <p className="text-xs text-gray-400 mt-1">Una URL por línea</p>
          </div>

          <div className="flex items-center gap-3 pt-2">
            <button type="submit" disabled={submitting} className="btn-primary flex items-center gap-2">
              {submitting ? (
                <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                </svg>
              ) : (
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
                </svg>
              )}
              {isEditing ? 'Actualizar vehículo' : 'Crear vehículo'}
            </button>
            <Link to="/admin/vehiculos" className="btn-secondary text-sm">
              Cancelar
            </Link>
          </div>
        </form>
      </div>
    </div>
  )
}
