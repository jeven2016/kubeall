import {atomWithQuery} from "jotai-tanstack-query";
import apiClient from "@/client/ApiClient.ts";
import {buildNamespacedUri, resourceType} from "@/modules/common/Constants.tsx";
import {VmImage} from "@/crd";


export const discImagesAtom = atomWithQuery(() => ({
        queryKey: ["discImages"],
        queryFn: async ({queryKey: [, id]}): Promise<VmImage[]> => {
            const uri = `${buildNamespacedUri(resourceType.allNamespace,
                resourceType.vmImage)}?type=iso`;
            return await apiClient.get<VmImage[]>(uri);
        },
    })
)
