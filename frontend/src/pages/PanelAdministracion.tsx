import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router'
import { GraficoBarras, type BarraDatos } from '../components/graficos/GraficoBarras'
import { TarjetaMetrica } from '../components/graficos/TarjetaMetrica'
import { EncabezadoPagina } from '../components/ui/EncabezadoPagina'
import { MensajeError } from '../components/ui/MensajeError'
import { ContenidoCargando } from '../components/ui/Spinner'
import { formatearPrecio } from '../utils/formato'
import { api } from '../services/api'
import type { Metricas } from '../types/metricas'

const PERIODOS = [
  { valor: '7', etiqueta: 'Últimos 7 días' },
  { valor: '30', etiqueta: 'Últimos 30 días' },
  { valor: '90', etiqueta: 'Últimos 90 días' },
]

const COLORES_ESTADO: Record<string, string> = {
  disponible: 'bg-emerald-500',
  reservado: 'bg-amber-500',
  vendido: 'bg-acento-500',
  dado_de_baja: 'bg-red-500',
}

const ETIQUETAS_ESTADO: Record<string, string> = {
  disponible: 'Disponibles',
  reservado: 'Reservados',
  vendido: 'Vendidos',
  dado_de_baja: 'Dados de baja',
}

// colorDiasEnStock pinta la barra según cuánto lleva publicado el vehículo:
// cuanto más tiempo, más caliente el color para resaltar la rotación lenta.
function colorDiasEnStock(dias: number): string {
  if (dias >= 90) return 'bg-red-500'
  if (dias >= 60) return 'bg-orange-500'
  if (dias >= 30) return 'bg-amber-500'
  return 'bg-emerald-500'
}

const SECCIONES = [
  {
    a: '/admin/vehiculos',
    titulo: 'Vehículos',
    descripcion: 'Alta, modificación y baja del stock con ficha técnica e imágenes.',
    icono: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M8 7h8M8 12h8M8 17h5M6 3h12a1 1 0 011 1v16a1 1 0 01-1 1H6a1 1 0 01-1-1V4a1 1 0 011-1z"
        />
      </svg>
    ),
  },
  {
    a: '/admin/usuarios',
    titulo: 'Usuarios',
    descripcion: 'Alta, modificación y baja de las cuentas de clientes, vendedores y administradores.',
    icono: (
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-6 w-6">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="M16 11a4 4 0 10-8 0 4 4 0 008 0zM4 20c0-2.5 2.5-4 8-4s8 1.5 8 4"
        />
        <path strokeLinecap="round" d="M19 4v4M17 6h4" />
      </svg>
    ),
  },
]

export function PanelAdministracion() {
  const [periodo, setPeriodo] = useState('30')
  const [metricas, setMetricas] = useState<Metricas | null>(null)
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let activo = true
    setCargando(true)
    setError(null)

    api
      .obtenerMetricas(Number(periodo))
      .then((datos) => {
        if (activo) {
          setMetricas(datos)
          setCargando(false)
        }
      })
      .catch((err: unknown) => {
        if (activo) {
          setError(err instanceof Error ? err.message : 'No se pudieron cargar las métricas')
          setCargando(false)
        }
      })

    return () => {
      activo = false
    }
  }, [periodo])

  const datosVehiculos: BarraDatos[] = useMemo(() => {
    if (!metricas) {
      return []
    }
    return Object.entries(ETIQUETAS_ESTADO).map(([estado, etiqueta]) => {
      const encontrado = metricas.vehiculosPorEstado.find((v) => v.estado === estado)
      return {
        etiqueta,
        valor: encontrado?.cantidad ?? 0,
        color: COLORES_ESTADO[estado],
      }
    })
  }, [metricas])

  const datosConsultas: BarraDatos[] = useMemo(
    () =>
      (metricas?.consultasPorPeriodo ?? []).map((dia) => ({
        etiqueta: dia.fecha,
        valor: dia.cantidad,
      })),
    [metricas],
  )

  const datosVentas: BarraDatos[] = useMemo(
    () =>
      (metricas?.ventasPorPeriodo ?? []).map((dia) => ({
        etiqueta: dia.fecha,
        valor: dia.cantidad,
      })),
    [metricas],
  )

  const datosTestDrives: BarraDatos[] = useMemo(
    () =>
      (metricas?.testDrivesPorPeriodo ?? []).map((dia) => ({
        etiqueta: dia.fecha,
        valor: dia.cantidad,
      })),
    [metricas],
  )

  const datosVentasMarca: BarraDatos[] = useMemo(
    () =>
      (metricas?.ventasPorMarca ?? []).map((venta) => ({
        etiqueta: venta.marca,
        valor: venta.cantidad,
      })),
    [metricas],
  )

  const datosDiasEnStock: BarraDatos[] = useMemo(
    () =>
      (metricas?.vehiculosEnStock ?? []).map((vehiculo) => ({
        etiqueta: `${vehiculo.marca} ${vehiculo.modelo} ${vehiculo.anio}`,
        valor: vehiculo.diasEnStock,
        color: colorDiasEnStock(vehiculo.diasEnStock),
      })),
    [metricas],
  )

  const disponibles = useMemo(
    () => metricas?.vehiculosPorEstado.find((v) => v.estado === 'disponible')?.cantidad ?? 0,
    [metricas],
  )

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <EncabezadoPagina
        destacado="Administración"
        titulo="Panel de administración"
        descripcion="Resumen de la operación y accesos a las secciones de gestión."
      />

      {error ? (
        <div className="mt-8">
          <MensajeError>{error}</MensajeError>
        </div>
      ) : cargando || !metricas ? (
        <ContenidoCargando etiqueta="Cargando métricas…" />
      ) : (
        <>
          <div className="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
            <TarjetaMetrica
              etiqueta="Vehículos disponibles"
              valor={disponibles}
              detalle="Unidades en stock para la venta"
              icono={
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-5 w-5 text-emerald-400">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M5 13l1.5 5.5A2 2 0 008.4 20h7.2a2 2 0 001.9-1.5L19 13m-14 0h14m-14 0a7 7 0 1114 0M9 8l1.5 1.5M15 8l-1.5 1.5"
                  />
                </svg>
              }
            />
            <TarjetaMetrica
              etiqueta="Reservas activas"
              valor={metricas.reservasActivas}
              detalle={`${metricas.reservasVendidas} ventas confirmadas`}
              icono={
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-5 w-5 text-amber-400">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              }
            />
            <TarjetaMetrica
              etiqueta="Test drives agendados"
              valor={metricas.testDrivesAgendados}
              detalle={`${metricas.testDrivesCompletados} completados`}
              icono={
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-5 w-5 text-acento-400">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M8 10h.01M16 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0zM8 15h8"
                  />
                </svg>
              }
            />
            <TarjetaMetrica
              etiqueta="Consultas abiertas"
              valor={metricas.consultasAbiertas}
              detalle={`${metricas.totalUsuarios} usuarios registrados`}
              icono={
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" className="h-5 w-5 text-sky-400">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.42-4.03 8-9 8a9.86 9.86 0 01-4.16-.92L3 20l1.05-3.38A7.85 7.85 0 013 12c0-4.42 4.03-8 9-8s9 3.58 9 8z"
                  />
                </svg>
              }
            />
          </div>

          <div className="mt-10 grid grid-cols-1 gap-6 lg:grid-cols-2">
            <section className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
              <h2 className="font-display text-lg font-semibold text-plata-100">Vehículos por estado</h2>
              <div className="mt-5">
                <GraficoBarras datos={datosVehiculos} descripcion="Stock actual por estado." />
              </div>
            </section>

            <section className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
              <h2 className="font-display text-lg font-semibold text-plata-100">Días en stock</h2>
              <div className="mt-5">
                <GraficoBarras
                  datos={datosDiasEnStock}
                  descripcion="Vehículos disponibles con más tiempo publicado. El color señala la rotación lenta."
                />
              </div>
            </section>
          </div>

          <div className="mt-10 flex flex-wrap items-center justify-between gap-3">
            <h2 className="font-display text-lg font-semibold text-plata-100">Evolución por período</h2>
            <div className="flex items-center gap-1 rounded-lg border border-white/10 bg-carbono-900 p-1">
              {PERIODOS.map((opcion) => (
                <button
                  key={opcion.valor}
                  type="button"
                  onClick={() => setPeriodo(opcion.valor)}
                  className={`rounded-md px-3 py-1.5 text-xs font-semibold transition-colors ${
                    periodo === opcion.valor
                      ? 'bg-acento-500/20 text-acento-300'
                      : 'text-plata-400 hover:text-plata-200'
                  }`}
                >
                  {opcion.etiqueta}
                </button>
              ))}
            </div>
          </div>

          <div className="mt-5 grid grid-cols-1 gap-6 lg:grid-cols-2">
            <section className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
              <h2 className="font-display text-lg font-semibold text-plata-100">Consultas por período</h2>
              <div className="mt-5">
                <GraficoBarras
                  datos={datosConsultas}
                  orientacion="vertical"
                  descripcion="Consultas creadas por día en el período seleccionado."
                />
              </div>
            </section>

            <section className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
              <h2 className="font-display text-lg font-semibold text-plata-100">Ventas por período</h2>
              <p className="mt-1 text-sm text-plata-400">
                Ingreso del período:{' '}
                <span className="font-semibold text-emerald-300">{formatearPrecio(metricas.ingresoPorPeriodo)}</span>
              </p>
              <div className="mt-5">
                <GraficoBarras
                  datos={datosVentas}
                  orientacion="vertical"
                  descripcion="Ventas confirmadas por día en el período seleccionado."
                />
              </div>
            </section>

            <section className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
              <h2 className="font-display text-lg font-semibold text-plata-100">Test drives por período</h2>
              <div className="mt-5">
                <GraficoBarras
                  datos={datosTestDrives}
                  orientacion="vertical"
                  descripcion="Turnos agendados por día en el período seleccionado."
                />
              </div>
            </section>

            <section className="rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm">
              <h2 className="font-display text-lg font-semibold text-plata-100">Ventas por marca</h2>
              <div className="mt-5">
                <GraficoBarras datos={datosVentasMarca} descripcion="Unidades vendidas según la marca del vehículo." />
              </div>
            </section>
          </div>

          <div className="mt-10">
            <h2 className="font-display text-lg font-semibold text-plata-100">Gestión</h2>
            <div className="mt-5 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {SECCIONES.map((seccion) => (
                <Link
                  key={seccion.a}
                  to={seccion.a}
                  className="group rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-white/20 hover:shadow-luz"
                >
                  <div className="mb-5 flex h-12 w-12 items-center justify-center rounded-xl border border-acento-500/30 bg-acento-500/10 text-acento-400">
                    {seccion.icono}
                  </div>
                  <h3 className="font-display text-lg font-semibold text-plata-100 transition-colors group-hover:text-white">
                    {seccion.titulo}
                  </h3>
                  <p className="mt-1 text-sm leading-relaxed text-plata-400">{seccion.descripcion}</p>
                  <span className="mt-4 inline-block text-sm font-semibold text-acento-400 opacity-0 transition-opacity group-hover:opacity-100">
                    Ingresar →
                  </span>
                </Link>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  )
}
