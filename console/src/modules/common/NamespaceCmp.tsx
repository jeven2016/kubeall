import {namespaceAtom} from "@/modules/state/namespace.ts";
import {useAtom} from "jotai";
import {Namespace} from "kubernetes-models/v1";
import {Select} from "react-windy-ui";
import {JSX} from "react";

export default function NamespaceCmp({selectedNamespace, onSelect}: {
    selectedNamespace: string,
    onSelect: (selectedNamespace: string) => void
}): JSX.Element {
    const [{data, isError, error}] = useAtom(namespaceAtom)
    const nsList = (data?.items ?? []) as Namespace[]

    return <Select style={{minWidth: "300px"}} value={selectedNamespace} onSelect={onSelect}>
        <Select.Option key="all" id="all">全部命名空间</Select.Option>
        {nsList.map(ns => <Select.Option key={ns.metadata?.name}
                                         id={ns.metadata?.name}>{ns.metadata?.name}</Select.Option>)}
    </Select>
}