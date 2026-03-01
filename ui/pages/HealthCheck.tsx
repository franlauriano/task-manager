import { CheckCircle, XCircle, Loader2 } from 'lucide-react'
import { useHealthcheck } from '@/hooks/useHealthcheck'

export default function HealthCheck() {
  const { isSuccess, isError, isPending } = useHealthcheck()

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="flex flex-col items-center gap-4 rounded-2xl bg-white px-12 py-10 shadow-md">
        <h1 className="text-xl font-semibold text-gray-700">API Status</h1>

        {isPending && (
          <div className="flex items-center gap-2 text-gray-400">
            <Loader2 className="animate-spin" size={28} />
            <span className="text-base">Checking…</span>
          </div>
        )}

        {isSuccess && (
          <div className="flex items-center gap-2 text-green-600">
            <CheckCircle size={28} />
            <span className="text-base font-medium">API is running</span>
          </div>
        )}

        {isError && (
          <div className="flex items-center gap-2 text-red-500">
            <XCircle size={28} />
            <span className="text-base font-medium">API is not running</span>
          </div>
        )}
      </div>
    </div>
  )
}
