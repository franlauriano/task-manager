import { useQuery } from '@tanstack/react-query'
import { fetchHealthcheck } from '@/api/healthcheck'

export function useHealthcheck() {
  return useQuery({
    queryKey: ['healthcheck'],
    queryFn: fetchHealthcheck,
    retry: false,
    refetchInterval: 10_000,
  })
}
