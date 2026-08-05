export type Rol = 'cliente' | 'vendedor' | 'administrador'

export interface Usuario {
  id: number
  nombre: string
  email: string
  rol: Rol
}

export interface DatosRegistro {
  nombre: string
  email: string
  password: string
}

export interface DatosLogin {
  email: string
  password: string
}

export interface RespuestaLogin {
  token: string
  usuario: Usuario
}

export interface DatosUsuarioAdmin {
  nombre: string
  email: string
  password?: string
  rol: Rol
}

export const ETIQUETAS_ROL: Record<Rol, string> = {
  cliente: 'Cliente',
  vendedor: 'Vendedor',
  administrador: 'Administrador',
}

export const ROLES: Rol[] = ['cliente', 'vendedor', 'administrador']
