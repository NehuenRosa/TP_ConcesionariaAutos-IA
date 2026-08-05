export function formatearPrecio(precio: number): string {
  return new Intl.NumberFormat('es-AR', {
    style: 'currency',
    currency: 'ARS',
    maximumFractionDigits: 0,
  }).format(precio)
}

export function formatearKilometraje(kilometraje: number): string {
  return `${new Intl.NumberFormat('es-AR').format(kilometraje)} km`
}

export function formatearFecha(fecha: string | Date): string {
  const d = fecha instanceof Date ? fecha : new Date(fecha)
  if (Number.isNaN(d.getTime())) return String(fecha)
  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  }).format(d)
}

export function formatearFechaHora(fecha: string | Date): string {
  const d = fecha instanceof Date ? fecha : new Date(fecha)
  if (Number.isNaN(d.getTime())) return String(fecha)
  return new Intl.DateTimeFormat('es-AR', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(d)
}

export function capitalizar(texto: string): string {
  if (!texto) return texto
  return texto.charAt(0).toUpperCase() + texto.slice(1)
}

export function capitalizarTodo(texto: string): string {
  return texto
    .split(/\s+/)
    .map((palabra) => capitalizar(palabra))
    .join(' ')
}

export function etiquetaCondicion(condicion: string): string {
  return condicion === 'nuevo' ? 'Nuevo' : 'Usado'
}

export function imagenVehiculo(imagen: string | undefined | null): string {
  if (imagen) return imagen
  return 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=1200&q=80'
}
