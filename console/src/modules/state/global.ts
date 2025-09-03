import {atom} from "jotai";
import {atomWithQuery} from "jotai-tanstack-query";
import {buildClusterUri, queryKey, resourceType} from "@/modules/common/Constants.tsx";
import apiClient from "@/client/ApiClient.ts";
import {PageList} from "@/types";


export const selectedModule = atom("vmMgr");
export const selectedMenuItem = atom("");

export const globalSettingsAtom = atomWithQuery(() => ({
        queryKey: [queryKey.globalSettings],
        queryFn: async (): Promise<PageList> => {
            return apiClient.get<PageList>(buildClusterUri(resourceType.globalSettings));
        },
    })
);

export const osTypesAtom = atom((get): string[] => {
    const globalSettings = get(globalSettingsAtom);
    const items = globalSettings.data?.items ?? [];
    if (items.length > 0) {
        return items[0].spec.osTypes;
    }
    return []
})