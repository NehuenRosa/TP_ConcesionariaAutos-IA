import { useState } from 'react'

interface Props {
  images: string[]
  alt: string
}

export function ImageGallery({ images, alt }: Props) {
  const [selectedIndex, setSelectedIndex] = useState(0)
  const imgs = images.length > 0 ? images : ['https://placehold.co/600x400?text=Sin+Imagen']

  return (
    <div>
      <img
        src={imgs[selectedIndex]}
        alt={`${alt} - imagen ${selectedIndex + 1}`}
        className="w-full h-96 object-cover rounded-lg"
      />
      {imgs.length > 1 && (
        <div className="flex gap-2 mt-2">
          {imgs.map((img, i) => (
            <button
              key={i}
              onClick={() => setSelectedIndex(i)}
              className={`border-2 rounded overflow-hidden ${
                i === selectedIndex ? 'border-blue-600' : 'border-transparent'
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
