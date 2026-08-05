import type { InputHTMLAttributes, ReactNode, SelectHTMLAttributes, TextareaHTMLAttributes } from 'react'

interface PropiedadesCampoBase {
  etiqueta?: string
  error?: string
}

interface PropiedadesCampoTexto extends InputHTMLAttributes<HTMLInputElement>, PropiedadesCampoBase {}

export function CampoTexto({ etiqueta, error, id, className = '', ...resto }: PropiedadesCampoTexto) {
  return (
    <div>
      {etiqueta && (
        <label htmlFor={id} className="etiqueta">
          {etiqueta}
        </label>
      )}
      <input id={id} className={`campo ${error ? 'border-red-500/70' : ''} ${className}`} {...resto} />
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  )
}

interface PropiedadesCampoSeleccion
  extends SelectHTMLAttributes<HTMLSelectElement>,
    PropiedadesCampoBase {
  children: ReactNode
}

export function CampoSeleccion({ etiqueta, error, id, className = '', children, ...resto }: PropiedadesCampoSeleccion) {
  return (
    <div>
      {etiqueta && (
        <label htmlFor={id} className="etiqueta">
          {etiqueta}
        </label>
      )}
      <select id={id} className={`campo ${error ? 'border-red-500/70' : ''} ${className}`} {...resto}>
        {children}
      </select>
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  )
}

interface PropiedadesCampoArea extends TextareaHTMLAttributes<HTMLTextAreaElement>, PropiedadesCampoBase {}

export function CampoArea({ etiqueta, error, id, className = '', ...resto }: PropiedadesCampoArea) {
  return (
    <div>
      {etiqueta && (
        <label htmlFor={id} className="etiqueta">
          {etiqueta}
        </label>
      )}
      <textarea id={id} className={`campo ${error ? 'border-red-500/70' : ''} ${className}`} {...resto} />
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  )
}
