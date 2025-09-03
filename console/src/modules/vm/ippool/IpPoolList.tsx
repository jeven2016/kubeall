import {TFunction} from "i18next";
import {useTranslation} from "react-i18next";
import {FilterQuery, ListParam, PageList} from "@/types";
import {Box, Button, Dropdown, Input, PopConfirm, Space} from "react-windy-ui";
import ResourceList from "@/modules/layout/ResourceList.tsx";
import {TiDelete} from "react-icons/ti";
import {buildNamespacedUri, constants, resourceType} from "@/modules/common/Constants.tsx";
import {QueryClient, useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {NavigateFunction, useLocation, useNavigate} from "react-router";
import React, {useState} from "react";
import {ImageApi} from "@/modules/api/ImageApi.ts";
import {Message} from "@/modules/common/Message.tsx";
import {delay} from "lodash";
import NamespaceCmp from "@/modules/common/NamespaceCmp.tsx";
import apiClient from "@/client/ApiClient.ts";
import {IPAddressPool} from "@/crd";

export default function IpPoolList() {
    const {t} = useTranslation();
    const [checkedImageKeys, setCheckedImageKeys] = useState<string[]>([]);
    const navigate: NavigateFunction = useNavigate();
    const queryClient: QueryClient = useQueryClient();
    const {pathname} = useLocation();
    const [queryData, setQueryData] = useState<FilterQuery>({});
    const [namespace, setNamespace] = useState<string>(resourceType.allNamespace);

    const {data: pageList, isLoading, isSuccess, isError, error} = useQuery<PageList>({
        refetchInterval: 15000,
        queryKey: ["vmList"],
        queryFn: async () => {
            let filter: string = "";
            if (Object.keys(queryData).length > 0) {
                filter = `?filter=${JSON.stringify(queryData)}`;
            }
            const baseUrl = buildNamespacedUri(namespace, resourceType.ipPool);
            return apiClient.get<PageList>(`${baseUrl}${filter}`);
        },
    });

    const deleteImageMutation = useMutation({
        mutationFn: (name: string) => {
            const vmList = pageList!!.items.filter(item => item.name === name)
            let vmNamespace: string = "default";
            if (vmList.length > 0) {
                vmNamespace = vmList[0].metadata.namespace;
            }
            return ImageApi.deleteResource(vmNamespace, resourceType.vm, name);
        },
        onError: (error, variables, context) => {
            Message.error("删除失败", "错误信息：" + error.message, null)
        },
        onSuccess: (data, variables, context) => {
            Message.ok("删除成功", `删除${variables}成功`)
        },
        onSettled: (data, error, variables, context) => {
            delay(() => queryClient.invalidateQueries({queryKey: ["images"]}), 1000)
        },
    })

    const deleteSelectedImages = async () => {
        try {
            await Promise.all(checkedImageKeys.map((name) => deleteImageMutation.mutateAsync(name)));
        } catch (error) {
            console.error('Error during parallel batch deletion:', error);
        }
    }

    const tableData = getData(pageList, isSuccess);
    const hasRows = tableData.length > 0;
    const listParam = {
        title: "IP地址池",
        description: "通过metallb组件提供的二层IP地址池，可以给每个虚拟机分配一个外部可访问的IP。",
        filterArea: <Box
            block
            left={
                <Space>
                    <Input placeholder="输入名称搜索" onChange={(e) => {
                        if (e.target.value) {
                            queryData["name"] = e.target.value;
                            setQueryData(queryData)
                        }
                    }} onBlur={(): Promise<void> => queryClient.invalidateQueries({queryKey: ["images"]})}/>
                    <NamespaceCmp selectedNamespace={namespace}
                                  onSelect={(val: string): void => setNamespace(val)}/>
                </Space>
            }
        />,
        cells: getCells(t),
        data: tableData,
        checkedRows: checkedImageKeys,
        onCheck: (img: IPAddressPool, checked: boolean) => {
            const name = img.metadata.name;
            if (checked) {
                setCheckedImageKeys([...checkedImageKeys, name!!])
            } else {
                const newImages = checkedImageKeys.filter((existingKey) =>
                    existingKey !== name);
                setCheckedImageKeys(newImages);
            }
        },
        onCheckAll: (checked: boolean) => {
            if (!checked) {
                setCheckedImageKeys([]);
            } else {
                const allImageKeys = tableData.map((item) => item.key!!);
                setCheckedImageKeys(allImageKeys)
            }
        },
        buttons: getButtons(queryClient, navigate, pathname, deleteSelectedImages, hasRows)
    } as ListParam;

    console.log(pageList)
    return <ResourceList listParam={listParam} pageList={pageList}/>
}

const getButtons = (queryClient: QueryClient, navigate: NavigateFunction,
                    pathname: string, deleteSelectedImages: () => Promise<void>, hasRows: boolean): React.ReactNode => {
    return <Space>
        <Button color="purple" hasBorder={false} onClick={() => navigate(`${pathname}/creation`)}>
            新建
        </Button>
        <PopConfirm body="您确定要删除所选行吗?" onOk={deleteSelectedImages} disabled={!hasRows}>
            <Button color="red" hasBorder={false}>删除</Button>
        </PopConfirm>
        <Button color="teal" hasBorder={false}
                onClick={() => queryClient.invalidateQueries({queryKey: ["images"]})}>刷新</Button>
    </Space>
}

const getCells = (t: TFunction) => {
    return [
        {head: '名称', paramName: 'name'},
        {head: '命名空间', paramName: 'namespace'},
        {head: '可用IP', paramName: 'addresses'},
        {head: '自动分配IP', paramName: 'autoAssign'},
        {head: '避开问题IP', paramName: 'avoidBuggyIp'},
        {head: '操作', paramName: 'action'},
    ]
}

const getData = (data?: PageList, isSuccess: boolean) => {
    if (isSuccess) {
        return data?.items.map(((item: IPAddressPool, index) => {
            let displayName = item.metadata?.annotations?.[constants.annotations.displayName] ?? "";

            if (/^\s*$/.test(displayName)) {
                displayName = item.metadata?.name ?? "";
            }
            return {
                key: item.metadata.name,
                name: displayName,
                namespace: item.metadata.namespace,
                addresses: item.spec.addresses?.map(addr => <div>{addr}</div>),
                autoAssign: item.spec.autoAssign,
                avoidBuggyIPs: item.spec.avoidBuggyIPs,
                createTime: item.metadata?.creationTimestamp,
                action: <div>
                    <Button inverted color="red" leftIcon={<TiDelete size="1.5rem"/>}>删除</Button>
                    <Dropdown title={<Button inverted color="purple">更多</Button>} activeBy="hover">
                        <Dropdown.Menu type="primary" popupSubMenu>
                            <Dropdown.Item>禁用</Dropdown.Item>
                            <Dropdown.Item>详情</Dropdown.Item>
                        </Dropdown.Menu>
                    </Dropdown>
                </div>
            }
        })) ?? [];
    }
    return []
};
