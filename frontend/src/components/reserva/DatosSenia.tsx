import type { DatosTransferencia } from '../../types/reserva'
import { formatearPrecio } from '../../utils/formato'

// DatosSenia muestra CBU, alias y monto; si la concesionaria no cargó sus
// datos, avisa que el personal los va a pasar. Se usa en el formulario de
// reserva y en Mis Reservas mientras la reserva esté activa.
export function DatosSenia({ datos }: { datos: DatosTransferencia | null }) {
  return (
    <div className="rounded-xl border border-amber-400/30 bg-amber-400/10 p-4 text-sm">
      <p className="font-display font-semibold text-amber-200">Seña por transferencia bancaria (5 %)</p>
      <p className="mt-1 leading-relaxed text-plata-300">
        Transferí{' '}
        {datos ? (
          <>
            <span className="texto-numerico font-semibold text-amber-200">{formatearPrecio(datos.monto)}</span>{' '}
            a la cuenta de la concesionaria y subí el comprobante para que un vendedor revise tu reserva.
          </>
        ) : (
          <>el 5 % del valor del vehículo a la cuenta de la concesionaria; el personal te va a pasar los datos.</>
        )}
      </p>
      {datos && (
        <dl className="mt-3 space-y-1 text-plata-300">
          {datos.cbu && (
            <div className="flex gap-2">
              <dt className="w-14 shrink-0 text-plata-500">CBU</dt>
              <dd className="texto-numerico break-all font-medium text-plata-100">{datos.cbu}</dd>
            </div>
          )}
          {datos.alias && (
            <div className="flex gap-2">
              <dt className="w-14 shrink-0 text-plata-500">Alias</dt>
              <dd className="break-all font-medium text-plata-100">{datos.alias}</dd>
            </div>
          )}
        </dl>
      )}
    </div>
  )
}
