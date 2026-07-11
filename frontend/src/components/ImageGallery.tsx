import { useState } from 'react'

interface Props {
  images: string[]
  alt: string
}

export function ImageGallery({ images, alt }: Props) {
  const [selectedIndex, setSelectedIndex] = useState(0)
  const imgs = images.length > 0 ? images : ['https://placehold.co/600x400?text=Sin+Imagen']

  return (
    <div className="space-y-3">
      <div className="relative group rounded-2xl overflow-hidden bg-gray-100">
        <img
          src={imgs[selectedIndex]}
          alt={`${alt} - imagen ${selectedIndex + 1}`}
          className="w-full h-96 object-cover transition-transform duration-500 group-hover:scale-105"
        />
        {imgs.length > 1 && (
          <>
            <button
              onClick={() => setSelectedIndex(Math.max(0, selectedIndex - 1))}
              disabled={selectedIndex === 0}
              className="absolute left-3 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/80 backdrop-blur-sm flex items-center justify-center text-gray-700 hover:bg-white disabled:opacity-30 transition-all opacity-0 group-hover:opacity-100"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <button
              onClick={() => setSelectedIndex(Math.min(imgs.length - 1, selectedIndex + 1))}
              disabled={selectedIndex === imgs.length - 1}
              className="absolute right-3 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-white/80 backdrop-blur-sm flex items-center justify-center text-gray-700 hover:bg-white disabled:opacity-30 transition-all opacity-0 group-hover:opacity-100"
            >
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
              </svg>
            </button>
            <div className="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-1.5">
              {imgs.map((_, i) => (
                <button
                  key={i}
                  onClick={() => setSelectedIndex(i)}
                  className={`w-2 h-2 rounded-full transition-all duration-300 ${
                    i === selectedIndex ? 'bg-white w-6' : 'bg-white/50 hover:bg-white/80'
                  }`}
                />
              ))}
            </div>
          </>
        )}
      </div>
      {imgs.length > 1 && (
        <div className="flex gap-2 overflow-x-auto pb-1">
          {imgs.map((img, i) => (
            <button
              key={i}
              onClick={() => setSelectedIndex(i)}
              className={`flex-shrink-0 rounded-xl overflow-hidden transition-all duration-200 ${
                i === selectedIndex
                  ? 'ring-2 ring-brand-500 ring-offset-2 scale-105'
                  : 'opacity-60 hover:opacity-100'
              }`}
            >
              <img src={img} alt="" className="w-20 h-16 object-cover" />
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
