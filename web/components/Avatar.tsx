'use client';

import Image from 'next/image';
import { parseAvatar } from '@/lib/avatars';

export function Avatar({ avatar, size = 40 }: { avatar?: string; size?: number }) {
  const parsed = parseAvatar(avatar);
  const style = { width: size, height: size };

  if (parsed.type === 'image') {
    return (
      <div className="relative rounded-full overflow-hidden shrink-0" style={style}>
        <Image src={parsed.src} alt="Avatar" fill sizes={`${size}px`} className="object-cover" unoptimized />
      </div>
    );
  }

  return (
    <div
      className={`rounded-full bg-gradient-to-br ${parsed.preset.gradient} flex items-center justify-center shrink-0`}
      style={{ ...style, fontSize: size * 0.55 }}
    >
      {parsed.preset.emoji}
    </div>
  );
}
