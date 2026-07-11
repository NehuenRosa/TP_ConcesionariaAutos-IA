import { useState, useEffect, useRef, useCallback } from 'react'
import { VehicleCard } from '../components/VehicleCard'
import { FilterPanel } from '../components/FilterPanel'
import { Pagination } from '../components/Pagination'
import { vehicleService } from '../services/vehicleService'
import type { Vehicle, VehicleFilter } from '../types'

function useRevealOnScroll() {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const observerRef = useRef<IntersectionObserver | null>(null)
  const mutationRef = useRef<MutationObserver | null>(null)

  const observeElements = useCallback(() => {
    if (!containerRef.current || !observerRef.current) return
    containerRef.current.querySelectorAll('.scroll-reveal').forEach((el) => {
      observerRef.current?.observe(el)
    })
  }, [])

  const containerCallback = useCallback((node: HTMLDivElement | null) => {
    if (mutationRef.current) mutationRef.current.disconnect()
    if (observerRef.current) observerRef.current.disconnect()
    containerRef.current = node
    if (!node) return

    observerRef.current = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add('revealed')
            observerRef.current?.unobserve(entry.target)
          }
        })
      },
      { threshold: 0.1 }
    )

    observeElements()

    mutationRef.current = new MutationObserver(() => observeElements())
    mutationRef.current.observe(node, { childList: true, subtree: true })
  }, [observeElements])

  return containerCallback
}

export function Catalog() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [brands, setBrands] = useState<string[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [filter, setFilter] = useState<VehicleFilter>({})
  const [loading, setLoading] = useState(true)
  const pageSize = 12
  const gridRef = useRevealOnScroll()

  useEffect(() => {
    vehicleService.getBrands().then(setBrands).catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    vehicleService.list({ ...filter, page, page_size: pageSize }).then((res) => {
      setVehicles(res.data)
      setTotal(res.total)
      setLoading(false)
    }).catch(() => setLoading(false))
  }, [filter, page])

  const handleFilter = (newFilter: VehicleFilter) => {
    setFilter(newFilter)
    setPage(1)
  }

  return (
    <>
      <section className="bg-hero-bg pb-16 pt-12 sm:pt-14">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center max-w-2xl mx-auto">
            <div className="inline-flex items-center gap-2 border border-[#3C3489] rounded-full px-4 py-1.5 mb-5">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
              <span className="text-xs font-medium text-[#AFA9EC]">{total || '...'} vehículos disponibles</span>
            </div>
            <h1 className="text-3xl sm:text-4xl font-semibold text-white mb-3 leading-tight">
              Encontrá tu próximo auto
            </h1>
            <p className="text-base text-text-placeholder max-w-lg mx-auto">
              Explorá nuestro catálogo con los mejores vehículos del mercado
            </p>
          </div>
        </div>
      </section>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 -mt-7 relative z-10" ref={gridRef}>
        <div className="bg-white rounded-[14px] border border-border-subtle shadow-[0_4px_20px_rgba(0,0,0,0.06)] p-4 mb-8">
          <FilterPanel onFilter={handleFilter} brands={brands} />
        </div>

        {loading ? (
          <div className="grid gap-5" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))' }}>
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="bg-white rounded-[14px] border border-border-subtle overflow-hidden animate-pulse">
                <div className="aspect-[4/3] bg-[#DFE1ED]" />
                <div className="p-4 space-y-3">
                  <div className="h-4 bg-surface rounded w-3/4" />
                  <div className="h-3 bg-surface rounded w-1/2" />
                  <div className="h-5 bg-surface rounded w-1/3" />
                  <div className="flex gap-2">
                    <div className="h-6 bg-surface rounded w-16" />
                    <div className="h-6 bg-surface rounded w-14" />
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : vehicles.length === 0 ? (
          <div className="text-center py-20">
            <div className="w-16 h-16 bg-[#DFE1ED] rounded-2xl flex items-center justify-center mx-auto mb-4">
              <svg className="w-8 h-8 text-[#B7BAD0]" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <p className="text-text-secondary font-medium">No se encontraron vehículos</p>
            <p className="text-text-placeholder text-sm mt-1">Probá con otros filtros de búsqueda</p>
          </div>
        ) : (
          <>
            <div className="flex items-center justify-between mb-5">
              <p className="text-sm text-text-secondary">
                <span className="font-medium text-text-primary">{(page - 1) * pageSize + 1}-{Math.min(page * pageSize, total)}</span> de{' '}
                <span className="font-medium text-text-primary">{total}</span> vehículos
              </p>
            </div>
            <div className="grid gap-5" style={{ gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))' }}>
              {vehicles.map((v: Vehicle, i: number) => (
                <div
                  key={v.id}
                  className={i < 3 ? 'animate-slide-up' : 'scroll-reveal'}
                  style={i < 3 ? { animationDelay: `${i * 80}ms` } : { transitionDelay: `${i * 60}ms` }}
                >
                  <VehicleCard vehicle={v} />
                </div>
              ))}
            </div>
            <Pagination page={page} pageSize={pageSize} total={total} onPageChange={setPage} />
          </>
        )}
      </div>
    </>
  )
}
