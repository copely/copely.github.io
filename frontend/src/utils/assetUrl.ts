/**
 * Resolve a public asset path against Vite's configured BASE_URL so assets
 * load correctly both at "/" and under subpath hosts like GitHub Pages.
 */
export function assetUrl(path: string): string {
  if (!path || /^(https?:)?\/\//i.test(path) || path.startsWith('data:')) {
    return path;
  }
  const base = import.meta.env.BASE_URL || '/';
  const normalizedBase = base.endsWith('/') ? base : `${base}/`;
  if (normalizedBase !== '/' && path.startsWith(normalizedBase)) {
    return path;
  }
  return `${normalizedBase}${path.replace(/^\/+/, '')}`;
}
