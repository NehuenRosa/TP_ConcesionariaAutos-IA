import { Link } from 'react-router'
import { Boton } from '../components/ui/Boton'

export function NoEncontrada() {
  return (
    <div className="mx-auto flex max-w-xl flex-col items-center px-4 py-28 text-center sm:px-6">
      <p className="font-display text-[7rem] font-extrabold leading-none tracking-tight text-carbono-600 select-none">
        404
      </p>
      <p className="mb-2 font-display text-xs font-semibold tracking-[0.3em] text-acento-400 uppercase">
        Página no encontrada
      </p>
      <h1 className="font-display text-3xl font-bold text-plata-100">Te perdiste en el camino</h1>
      <p className="mt-3 text-plata-400">La página que buscás no existe o fue movida.</p>
      <div className="mt-8">
        <Boton>
          <Link to="/">Volver al inicio</Link>
        </Boton>
      </div>
    </div>
  )
}
