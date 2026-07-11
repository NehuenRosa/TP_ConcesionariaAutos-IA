import api from './api'
import type { AuthResponse } from '../types'

export const authService = {
  async login(email: string, password: string): Promise<AuthResponse> {
    const { data } = await api.post('/auth/login', { email, password })
    return data
  },

  async register(name: string, email: string, password: string, phone?: string): Promise<AuthResponse> {
    const { data } = await api.post('/auth/register', { name, email, password, phone })
    return data
  },

  async me() {
    const { data } = await api.get('/auth/me')
    return data
  },
}
