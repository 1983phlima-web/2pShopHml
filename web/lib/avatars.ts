// 12 preset avatars — gradient background + emoji, fully self-contained
// (no external avatar service dependency). Selected randomly by the
// backend for every new user, and changeable anytime from the profile
// popup.
export interface AvatarPreset {
  id: number;
  gradient: string;
  emoji: string;
}

export const AVATAR_PRESETS: AvatarPreset[] = [
  { id: 1, gradient: 'from-indigo-400 to-purple-500', emoji: '🦊' },
  { id: 2, gradient: 'from-amber-400 to-orange-500', emoji: '🐯' },
  { id: 3, gradient: 'from-emerald-400 to-teal-500', emoji: '🐸' },
  { id: 4, gradient: 'from-rose-400 to-pink-500', emoji: '🐱' },
  { id: 5, gradient: 'from-sky-400 to-blue-500', emoji: '🐼' },
  { id: 6, gradient: 'from-violet-400 to-fuchsia-500', emoji: '🦄' },
  { id: 7, gradient: 'from-lime-400 to-green-500', emoji: '🐢' },
  { id: 8, gradient: 'from-yellow-400 to-amber-500', emoji: '🦁' },
  { id: 9, gradient: 'from-cyan-400 to-sky-500', emoji: '🐬' },
  { id: 10, gradient: 'from-red-400 to-rose-500', emoji: '🦉' },
  { id: 11, gradient: 'from-teal-400 to-emerald-500', emoji: '🐨' },
  { id: 12, gradient: 'from-purple-400 to-indigo-500', emoji: '🐵' },
];

export function getPreset(id: number): AvatarPreset {
  return AVATAR_PRESETS.find((p) => p.id === id) || AVATAR_PRESETS[0];
}

// Parses the stored `avatar` field: "preset:N" or a base64 data URI.
export function parseAvatar(avatar?: string): { type: 'preset'; preset: AvatarPreset } | { type: 'image'; src: string } {
  if (avatar && avatar.startsWith('preset:')) {
    const id = parseInt(avatar.split(':')[1], 10);
    return { type: 'preset', preset: getPreset(id) };
  }
  if (avatar && avatar.startsWith('data:image')) {
    return { type: 'image', src: avatar };
  }
  return { type: 'preset', preset: AVATAR_PRESETS[0] };
}
