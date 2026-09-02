'use client';

import Image from 'next/image';
import { useState, useRef } from 'react';

export function ImageCarousel({ images, alt }: { images: string[]; alt: string }) {
  const [index, setIndex] = useState(0);
  const touchStartX = useRef<number | null>(null);
  const list = images.length > 0 ? images : [];

  function go(delta: number) {
    setIndex((i) => (i + delta + list.length) % list.length);
  }

  function onTouchStart(e: React.TouchEvent) {
    touchStartX.current = e.touches[0].clientX;
  }
  function onTouchEnd(e: React.TouchEvent) {
    if (touchStartX.current === null) return;
    const delta = e.changedTouches[0].clientX - touchStartX.current;
    if (Math.abs(delta) > 40) go(delta > 0 ? -1 : 1);
    touchStartX.current = null;
  }

  if (list.length === 0) {
    return <div className="h-full w-full flex items-center justify-center text-gray-400 text-5xl bg-gray-100">📦</div>;
  }

  return (
    <div
      className="relative h-full w-full bg-gray-100 overflow-hidden select-none"
      onTouchStart={onTouchStart}
      onTouchEnd={onTouchEnd}
    >
      <Image
        key={index}
        src={list[index]}
        alt={alt}
        fill
        sizes="(max-width: 768px) 100vw, 50vw"
        className="object-cover"
        priority
      />

      {list.length > 1 && (
        <>
          <button
            onClick={() => go(-1)}
            aria-label="Imagem anterior"
            className="absolute left-2 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-white/80 flex items-center justify-center text-gray-700 hover:bg-white transition"
          >
            ‹
          </button>
          <button
            onClick={() => go(1)}
            aria-label="Próxima imagem"
            className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 rounded-full bg-white/80 flex items-center justify-center text-gray-700 hover:bg-white transition"
          >
            ›
          </button>
          <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1.5">
            {list.map((_, i) => (
              <button
                key={i}
                onClick={() => setIndex(i)}
                aria-label={`Imagem ${i + 1}`}
                className={`h-1.5 rounded-full transition-all ${i === index ? 'w-5 bg-white' : 'w-1.5 bg-white/60'}`}
              />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
