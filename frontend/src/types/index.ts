export interface User {
  id: number
  name: string
  email: string
  role: 'cliente' | 'vendedor' | 'administrador'
  phone?: string
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Vehicle {
  id: number
  created_at: string
  updated_at: string
  brand: string
  model: string
  year: number
  price: number
  mileage: number
  fuel: 'nafta' | 'diesel' | 'electrico' | 'hibrido'
  transmission: 'manual' | 'automatico'
  condition: 'nuevo' | 'usado'
  color?: string
  description?: string
  images: string[]
  status: 'disponible' | 'reservado' | 'vendido'
  vehicle_type: string
}

export interface VehicleFilter {
  search?: string
  brand?: string
  model?: string
  year_from?: number
  year_to?: number
  price_from?: number
  price_to?: number
  fuel?: string
  condition?: string
  vehicle_type?: string
  sort_by?: string
  sort_order?: string
  page?: number
  page_size?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
}

export interface Consultation {
  id: number
  created_at: string
  updated_at: string
  client_id: number
  client?: User
  vehicle_id: number
  vehicle?: Vehicle
  message: string
  status: 'pendiente' | 'en_conversacion' | 'cerrada'
  assigned_to?: number
  seller?: User
  has_unread_messages?: boolean
  has_unread_for_client?: boolean
  responses?: ConsultationResponse[]
}

export interface ConsultationResponse {
  id: number
  created_at: string
  consultation_id: number
  user_id: number
  user?: User
  message: string
}

export interface TestDrive {
  id: number
  created_at: string
  updated_at: string
  client_id: number
  client?: User
  vehicle_id: number
  vehicle?: Vehicle
  scheduled_at: string
  status: 'pendiente' | 'confirmado' | 'cancelado' | 'completado'
  notes?: string
}

export interface Reservation {
  id: number
  created_at: string
  updated_at: string
  client_id: number
  client?: User
  vehicle_id: number
  vehicle?: Vehicle
  status: 'activa' | 'confirmada' | 'cancelada'
  notes?: string
}

export interface DashboardMetrics {
  vehiculos: {
    total: number
    disponible: number
    reservado: number
    vendido: number
  }
  consultas: {
    total: number
  }
  test_drives_totales: number
  test_drives_pendientes: number
  reservas_activas: number
}
