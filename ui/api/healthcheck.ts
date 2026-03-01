export async function fetchHealthcheck(): Promise<boolean> {
  const res = await fetch('/healthcheck')
  if (!res.ok) throw new Error('API unavailable')
  return true
}
