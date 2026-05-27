import { useQuery, type UseQueryOptions } from "@tanstack/react-query";

type QueryKey = readonly unknown[];

export function useApiQuery<TData>(
  key: QueryKey,
  fetcher: () => Promise<TData>,
  options?: Omit<UseQueryOptions<TData, Error, TData, QueryKey>, "queryKey" | "queryFn">
) {
  return useQuery<TData, Error, TData, QueryKey>({
    queryKey: key,
    queryFn: fetcher,
    ...options,
  });
}
