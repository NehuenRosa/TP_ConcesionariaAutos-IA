import { useEffect, useState } from 'react'
import { Link } from 'react-router'
import { useAuth } from '../hooks/useAuth'
import { api, ErrorApi } from '../services/api'
import type { ConsultaResumen } from '../types/consulta'
import type { TurnoTestDrive } from '../types/testDrive'
import { Boton } from '../components/ui/Boton'
import { ContenidoCargando } from '../components/ui/Spinner'
import { formatearFranja } from '../utils/formato'
import {
  EtiquetaEstado,
  estilosEstadoConsulta,
  etiquetasEstadoConsulta,
  estilosEstadoTestDrive,
  etiquetasEstadoTestDrive,
} from '../components/ui/EtiquetaEstado'

const CATEGORIAS = [
  { tipo: 'SUV', valor: 'suv', imagen: 'https://images.unsplash.com/photo-1571607388263-1044f9ea01dd?w=1200&q=80' },
  { tipo: 'Sedán', valor: 'sedán', imagen: 'https://images.unsplash.com/photo-1549317661-bd32c8ce0db2?w=1200&q=80' },
  { tipo: 'Deportivo', valor: 'coupe', imagen: 'https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=1200&q=80' },
]

const VALORES = [
  {
    titulo: 'Stock real',
    descripcion: 'Todas las unidades que ves en el catálogo están disponibles en nuestro concesionario.',
  },
  {
    titulo: 'Test drive',
    descripcion: 'Agendá una prueba de manejo en la fecha y hora que más te convenga.',
  },
  {
    titulo: 'Asistente digital',
    descripcion: 'Consultá por chat cualquier duda sobre los vehículos y el stock disponible.',
  },
]

const SECCIONES_ADMIN = [
  {
    a: '/admin/vehiculos',
    titulo: 'Vehículos',
    descripcion: 'Alta, modificación y baja del stock con ficha técnica e imágenes.',
  },
  {
    a: '/admin/usuarios',
    titulo: 'Usuarios',
    descripcion: 'Alta, modificación y baja de las cuentas del sistema.',
  },
  {
    a: '/admin',
    titulo: 'Panel completo',
    descripcion: 'Ingresá al panel de administración con todas sus secciones.',
  },
]

function HeroInvitado() {
  return (
    <section className="relative isolate flex min-h-[86vh] items-center overflow-hidden">
      <div className="absolute inset-0 -z-10">
        <img
          src="https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=2000&q=80"
          alt="Auto deportivo premium"
          className="h-full w-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-r from-carbono-950 via-carbono-950/75 to-carbono-950/20" />
        <div className="absolute inset-x-0 bottom-0 h-40 bg-gradient-to-t from-carbono-950 to-transparent" />
      </div>

      <div className="mx-auto w-full max-w-7xl px-4 py-24 sm:px-6 lg:px-8">
        <div className="max-w-2xl">
          <p className="mb-4 font-display text-xs font-semibold tracking-[0.35em] text-acento-400 uppercase">
            Concesionaria premium
          </p>
          <h1 className="font-display text-4xl font-extrabold leading-tight tracking-tight text-plata-100 sm:text-6xl">
            El auto de tus sueños,
            <span className="block text-plata-300">a un clic.</span>
          </h1>
          <p className="mt-6 max-w-xl text-base leading-relaxed text-plata-300 sm:text-lg">
            Explorá nuestro catálogo de vehículos nuevos y usados, consultá sin compromiso,
            reservá tu unidad o agendá un test drive con nuestro asistente digital.
          </p>
          <div className="mt-8 flex flex-wrap gap-4">
            <Boton tamano="lg">
              <Link to="/catalogo">Explorar catálogo</Link>
            </Boton>
            <Boton variante="secundario" tamano="lg">
              <Link to="/registro">Crear mi cuenta</Link>
            </Boton>
          </div>
        </div>
      </div>
    </section>
  )
}

function HeroCliente() {
  return (
    <section className="relative isolate flex min-h-[70vh] items-center overflow-hidden">
      <div className="absolute inset-0 -z-10">
        <img
          src="https://images.unsplash.com/photo-1503376780353-7e6692767b70?w=2000&q=80"
          alt="Auto deportivo premium"
          className="h-full w-full object-cover"
        />
        <div className="absolute inset-0 bg-gradient-to-r from-carbono-950 via-carbono-950/75 to-carbono-950/20" />
        <div className="absolute inset-x-0 bottom-0 h-40 bg-gradient-to-t from-carbono-950 to-transparent" />
      </div>

      <div className="mx-auto w-full max-w-7xl px-4 py-24 sm:px-6 lg:px-8">
        <div className="max-w-2xl">
          <p className="mb-4 font-display text-xs font-semibold tracking-[0.35em] text-acento-400 uppercase">
            Bienvenido
          </p>
          <h1 className="font-display text-4xl font-extrabold leading-tight tracking-tight text-plata-100 sm:text-6xl">
            Tu próxima unidad
            <span className="block text-plata-300">te está esperando.</span>
          </h1>
          <p className="mt-6 max-w-xl text-base leading-relaxed text-plata-300 sm:text-lg">
            Seguí explorando el catálogo, revisá tus consultas o agendá un test drive
            en el horario que prefieras.
          </p>
          <div className="mt-8 flex flex-wrap gap-4">
            <Boton tamano="lg">
              <Link to="/catalogo">Explorar catálogo</Link>
            </Boton>
            <Boton variante="secundario" tamano="lg">
              <Link to="/mis-consultas">Mis consultas</Link>
            </Boton>
          </div>
        </div>
      </div>
    </section>
  )
}

function SeccionValores() {
  return (
    <section className="border-y border-white/8 bg-carbono-900/50">
      <div className="mx-auto grid max-w-7xl gap-10 px-4 py-16 sm:px-6 md:grid-cols-3 lg:px-8">
        {VALORES.map((valor) => (
          <div key={valor.titulo}>
            <div className="mb-4 h-px w-10 bg-acento-500" />
            <h3 className="font-display text-lg font-semibold text-plata-100">{valor.titulo}</h3>
            <p className="mt-2 text-sm leading-relaxed text-plata-400">{valor.descripcion}</p>
          </div>
        ))}
      </div>
    </section>
  )
}

function SeccionCategorias() {
  return (
    <section className="mx-auto max-w-7xl px-4 py-20 sm:px-6 lg:px-8">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
            Categorías
          </p>
          <h2 className="font-display text-3xl font-bold text-plata-100">Explorá por tipo</h2>
        </div>
        <Link
          to="/catalogo"
          className="font-display text-sm font-semibold text-plata-300 transition-colors hover:text-plata-100"
        >
          Ver catálogo completo →
        </Link>
      </div>

      <div className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {CATEGORIAS.map((categoria) => (
          <Link
            key={categoria.tipo}
            to={`/catalogo?tipo=${encodeURIComponent(categoria.valor)}`}
            className="group relative aspect-[4/3] overflow-hidden rounded-xl border border-white/8"
          >
            <img
              src={categoria.imagen}
              alt={`Vehículos tipo ${categoria.tipo}`}
              loading="lazy"
              className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-carbono-950/90 via-carbono-950/20 to-transparent" />
            <div className="absolute inset-x-0 bottom-0 flex items-end justify-between p-5">
              <span className="font-display text-xl font-semibold text-plata-100">{categoria.tipo}</span>
              <span className="text-plata-400 transition-colors group-hover:text-plata-100">→</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  )
}

function etiquetaFranja(franja: string): string {
  return formatearFranja(franja)
}

function HomeVendedor() {
  const { usuario } = useAuth()
  const [consultas, setConsultas] = useState<ConsultaResumen[]>([])
  const [testDrives, setTestDrives] = useState<TurnoTestDrive[]>([])
  const [cargando, setCargando] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelado = false
    Promise.all([api.listarBandeja(), api.listarTestDrives()])
      .then(([consultaLista, turnos]) => {
        if (cancelado) return
        setConsultas(consultaLista)
        setTestDrives(turnos)
      })
      .catch((e: unknown) => {
        if (cancelado) return
        setError(e instanceof ErrorApi ? e.message : 'Ocurrió un error inesperado al cargar el panel.')
      })
      .finally(() => {
        if (!cancelado) setCargando(false)
      })

    return () => {
      cancelado = true
    }
  }, [])

  const consultasAbiertas = consultas.filter((c) => c.estado !== 'cerrada')
  const testNuevos = testDrives.filter((t) => t.estado === 'solicitado')

  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
        Panel del vendedor
      </p>
      <h1 className="font-display text-3xl font-bold tracking-tight text-plata-100 sm:text-4xl">
        Hola, {usuario?.nombre}
      </h1>
      <p className="mt-2 max-w-2xl text-sm leading-relaxed text-plata-400">
        Este es tu resumen diario: respondé las consultas pendientes y gestioná los test drives
        nuevos que llegaron.
      </p>

      {cargando && <div className="mt-8"><ContenidoCargando etiqueta="Cargando el panel…" /></div>}

      {!cargando && error && (
        <div className="mt-8">
          <div className="rounded-xl border border-red-500/30 bg-red-500/10 p-4 text-sm text-red-300">{error}</div>
        </div>
      )}

      {!cargando && !error && (
        <div className="mt-10 grid gap-8 lg:grid-cols-2">
          <div className="rounded-2xl border border-white/8 bg-carbono-850/50 p-6 backdrop-blur-sm">
            <div className="flex items-center justify-between">
              <h2 className="font-display text-lg font-semibold text-plata-100">Consultas recibidas</h2>
              <Link to="/vendedor/bandeja" className="text-sm font-semibold text-acento-400 transition-colors hover:text-acento-300">
                Ver bandeja →
              </Link>
            </div>
            <p className="mt-1 text-sm text-plata-400">
              {consultasAbiertas.length} consulta{consultasAbiertas.length === 1 ? '' : 's'} en curso
            </p>

            <div className="mt-5 space-y-3">
              {consultasAbiertas.length === 0 && (
                <p className="rounded-xl border border-dashed border-white/10 px-4 py-6 text-center text-sm text-plata-500">
                  No tenés consultas pendientes. ¡Todo al día!
                </p>
              )}
              {consultasAbiertas.slice(0, 5).map((consulta) => (
                <Link
                  key={consulta.id}
                  to={`/vendedor/bandeja/${consulta.id}`}
                  className="flex items-center justify-between gap-3 rounded-xl border border-white/8 bg-carbono-900/40 px-4 py-3 transition-colors hover:border-white/20"
                >
                  <div className="min-w-0">
                    <p className="truncate text-sm font-semibold text-plata-100">
                      {consulta.vehiculo.marca} {consulta.vehiculo.modelo}
                    </p>
                    <p className="truncate text-xs text-plata-400">
                      {consulta.cliente.nombre}
                      {consulta.mensajesNuevos > 0 && (
                        <span className="ml-2 font-semibold text-acento-400">
                          · {consulta.mensajesNuevos} mensaje{consulta.mensajesNuevos === 1 ? '' : 's'} nuevo{consulta.mensajesNuevos === 1 ? '' : 's'}
                        </span>
                      )}
                    </p>
                  </div>
                  <EtiquetaEstado
                    estado={consulta.estado}
                    estilos={estilosEstadoConsulta}
                    etiqueta={etiquetasEstadoConsulta[consulta.estado]}
                  />
                </Link>
              ))}
            </div>
          </div>

          <div className="rounded-2xl border border-white/8 bg-carbono-850/50 p-6 backdrop-blur-sm">
            <div className="flex items-center justify-between">
              <h2 className="font-display text-lg font-semibold text-plata-100">Test drives nuevos</h2>
              <Link to="/vendedor/test-drives" className="text-sm font-semibold text-acento-400 transition-colors hover:text-acento-300">
                Gestionar →
              </Link>
            </div>
            <p className="mt-1 text-sm text-plata-400">
              {testNuevos.length} turno{testNuevos.length === 1 ? '' : 's'} por confirmar
            </p>

            <div className="mt-5 space-y-3">
              {testNuevos.length === 0 && (
                <p className="rounded-xl border border-dashed border-white/10 px-4 py-6 text-center text-sm text-plata-500">
                  No hay test drives nuevos solicitados.
                </p>
              )}
              {testNuevos.slice(0, 5).map((turno) => (
                <Link
                  key={turno.id}
                  to="/vendedor/test-drives"
                  className="flex items-center justify-between gap-3 rounded-xl border border-white/8 bg-carbono-900/40 px-4 py-3 transition-colors hover:border-white/20"
                >
                  <div className="min-w-0">
                    <p className="truncate text-sm font-semibold text-plata-100">
                      {turno.vehiculo.marca} {turno.vehiculo.modelo}
                    </p>
                    <p className="truncate text-xs text-plata-400">
                      {turno.cliente.nombre} · {turno.fecha} · {etiquetaFranja(turno.franja)}
                    </p>
                  </div>
                  <EtiquetaEstado
                    estado={turno.estado}
                    estilos={estilosEstadoTestDrive}
                    etiqueta={etiquetasEstadoTestDrive[turno.estado]}
                  />
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function HomeAdministrador() {
  return (
    <div className="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
      <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
        Administración
      </p>
      <h1 className="font-display text-3xl font-bold tracking-tight text-plata-100 sm:text-4xl">
        Accesos rápidos
      </h1>
      <p className="mt-2 max-w-2xl text-sm leading-relaxed text-plata-400">
        Ingresá directo a las secciones que más usás del panel de administración.
      </p>

      <div className="mt-10 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {SECCIONES_ADMIN.map((seccion) => (
          <Link
            key={seccion.a}
            to={seccion.a}
            className="group rounded-2xl border border-white/8 bg-carbono-850/60 p-6 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-white/20 hover:shadow-luz"
          >
            <h2 className="font-display text-lg font-semibold text-plata-100 transition-colors group-hover:text-white">
              {seccion.titulo}
            </h2>
            <p className="mt-1 text-sm leading-relaxed text-plata-400">{seccion.descripcion}</p>
            <span className="mt-4 inline-block text-sm font-semibold text-acento-400 opacity-0 transition-opacity group-hover:opacity-100">
              Ingresar →
            </span>
          </Link>
        ))}
      </div>
    </div>
  )
}

export function Inicio() {
  const { usuario, cargando } = useAuth()

  if (cargando) {
    return (
      <div className="flex min-h-[60vh] items-center justify-center">
        <ContenidoCargando etiqueta="Cargando…" />
      </div>
    )
  }

  if (usuario?.rol === 'vendedor') return <HomeVendedor />
  if (usuario?.rol === 'administrador') return <HomeAdministrador />

  return (
    <div>
      {usuario?.rol === 'cliente' ? <HeroCliente /> : <HeroInvitado />}
      <SeccionCategorias />
      <SeccionValores />
    </div>
  )
}
