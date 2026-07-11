import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
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
    }
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">{isEditing ? 'Editar' : 'Nuevo'} Vehículo</h1>
      {error && <div className="bg-red-100 text-red-700 p-3 rounded mb-4">{error}</div>}
      <form onSubmit={handleSubmit} className="bg-white p-6 rounded-lg shadow space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <input placeholder="Marca" value={form.brand} onChange={(e) => setForm({ ...form, brand: e.target.value })} required className="border rounded px-3 py-2" />
          <input placeholder="Modelo" value={form.model} onChange={(e) => setForm({ ...form, model: e.target.value })} required className="border rounded px-3 py-2" />
          <input type="number" placeholder="Año" value={form.year} onChange={(e) => setForm({ ...form, year: Number(e.target.value) })} required className="border rounded px-3 py-2" />
          <input type="number" placeholder="Precio" value={form.price} onChange={(e) => setForm({ ...form, price: Number(e.target.value) })} required className="border rounded px-3 py-2" />
          <input type="number" placeholder="Kilometraje" value={form.mileage} onChange={(e) => setForm({ ...form, mileage: Number(e.target.value) })} className="border rounded px-3 py-2" />
          <input placeholder="Color" value={form.color} onChange={(e) => setForm({ ...form, color: e.target.value })} className="border rounded px-3 py-2" />
          <select value={form.fuel} onChange={(e) => setForm({ ...form, fuel: e.target.value })} className="border rounded px-3 py-2">
            <option value="nafta">Nafta</option><option value="diesel">Diesel</option>
            <option value="electrico">Eléctrico</option><option value="hibrido">Híbrido</option>
          </select>
          <select value={form.transmission} onChange={(e) => setForm({ ...form, transmission: e.target.value })} className="border rounded px-3 py-2">
            <option value="manual">Manual</option><option value="automatico">Automático</option>
          </select>
          <select value={form.condition} onChange={(e) => setForm({ ...form, condition: e.target.value })} className="border rounded px-3 py-2">
            <option value="nuevo">Nuevo</option><option value="usado">Usado</option>
          </select>
          <select value={form.vehicle_type} onChange={(e) => setForm({ ...form, vehicle_type: e.target.value })} className="border rounded px-3 py-2">
            <option value="sedan">Sedán</option><option value="suv">SUV</option>
            <option value="hatchback">Hatchback</option><option value="pickup">Pickup</option>
            <option value="coupe">Coupé</option>
          </select>
        </div>
        <textarea placeholder="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} rows={3} className="w-full border rounded px-3 py-2" />
        <textarea placeholder="URLs de imágenes (una por línea)" value={form.images} onChange={(e) => setForm({ ...form, images: e.target.value })} rows={3} className="w-full border rounded px-3 py-2" />
        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700">
          {isEditing ? 'Actualizar' : 'Crear'} vehículo
        </button>
      </form>
    </div>
  )
}
