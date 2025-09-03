import {atomWithQuery} from "jotai-tanstack-query";
import {buildResourceUri, resourceType} from "@/modules/common/Constants.tsx";
import {PageList} from "@/types";
import apiClient from "@/client/ApiClient.ts";

export const namespaceAtom = atomWithQuery(() => ({
        queryKey: ["namespaces"],
        queryFn: async ({queryKey: [, id]}) => {
            const data = await apiClient.get<PageList>(buildResourceUri(resourceType.clusterResources, resourceType.namespace));
            return data as PageList
        },
    })
);


