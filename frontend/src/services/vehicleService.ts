import api from './api'
import type { Vehicle, VehicleFilter, PaginatedResponse } from '../types'

export const vehicleService = {
  async list(filter?: VehicleFilter): Promise<PaginatedResponse<Vehicle>> {
    const { data } = await api.get('/vehicles', { params: filter })
    return data
  },

  async getById(id: number): Promise<Vehicle> {
    const { data } = await api.get(`/vehicles/${id}`)
    return data
  },

  async create(vehicle: Partial<Vehicle>): Promise<Vehicle> {
    const { data } = await api.post('/vehicles', vehicle)
    return data
  },

  async update(id: number, vehicle: Partial<Vehicle>): Promise<Vehicle> {
    const { data } = await api.put(`/vehicles/${id}`, vehicle)
    return data
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/vehicles/${id}`)
  },

  async getBrands(): Promise<string[]> {
    const { data } = await api.get('/vehicles/brands')
    return data.brands
  },
}
