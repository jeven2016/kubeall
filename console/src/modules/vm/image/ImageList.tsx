import {TFunction} from "i18next";
import {useTranslation} from "react-i18next";
import {FilterQuery, ListParam, PageList, VmImageRow} from "@/types";
import {Box, Button, Dropdown, Input, PopConfirm, Progress, Space, Tooltip} from "react-windy-ui";
import ResourceList from "@/modules/layout/ResourceList.tsx";
import getOsIcon from "@/modules/icon/OsIcon.tsx";
import {TiDelete} from "react-icons/ti";
import {buildNamespacedUri, constants, resourceType} from "@/modules/common/Constants.tsx";
import {QueryClient, useMutation, useQuery, useQueryClient} from "@tanstack/react-query";
import {NavigateFunction, useLocation, useNavigate} from "react-router";
import React, {useState} from "react";
import {VmImage} from "@/crd";
import {ImageApi} from "@/modules/api/ImageApi.ts";
import {Message} from "@/modules/common/Message.tsx";
import {MdGppGood} from "react-icons/md";
import {delay} from "lodash";
import NamespaceCmp from "@/modules/common/NamespaceCmp.tsx";
import apiClient from "@/client/ApiClient.ts";
import {getImageTypeText} from "@/modules/common/VmCommon.tsx";

export default function ImageList() {
    const {t} = useTranslation();
    const [checkedImageKeys, setCheckedImageKeys] = useState<string[]>([]);
    const navigate: NavigateFunction = useNavigate();
    const queryClient: QueryClient = useQueryClient();
    const {pathname} = useLocation();
    const [queryData, setQueryData] = useState<FilterQuery>({});
    const [namespace, setNamespace] = useState<string>(resourceType.allNamespace);

    const {data, isLoading, isSuccess, isError, error} = useQuery({
        refetchInterval: 5000,
        queryKey: ["images"],
        queryFn: async () => {
            let filter: string = "";
            if (Object.keys(queryData).length > 0) {
                filter = `?filter=${JSON.stringify(queryData)}`;
            }
            const baseUrl = buildNamespacedUri(namespace, resourceType.vmImage);
            return apiClient.get<PageList>(`${baseUrl}${filter}`);
        },
    });
    const pageList = data as PageList;

    const deleteImageMutation = useMutation({
        mutationFn: (name: string) => {
            const vmList = pageList.items.filter(item => item.metadata.name === name)
            let vmNamespace: string = "default";
            if (vmList.length > 0) {
                vmNamespace = vmList[0].metadata.namespace;
            }
            return ImageApi.deleteResource(vmNamespace, resourceType.vmImage, name);
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
        title: "镜像",
        description: "虚拟机镜像管理，在此处上传虚拟机基础镜像后，可进行虚拟机的创建。",
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
        onCheck: (img: VmImageRow, checked: boolean) => {
            const name = img.name;
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
        <PopConfirm body="您确定要删除所选镜像吗?" onOk={deleteSelectedImages} disabled={!hasRows}>
            <Button color="red" hasBorder={false}>删除</Button>
        </PopConfirm>
        <Button color="teal" hasBorder={false}
                onClick={() => queryClient.invalidateQueries({queryKey: ["images"]})}>刷新</Button>
    </Space>
}

const getCells = (t: TFunction) => {
    return [
        {head: '名称', paramName: 'name'},
        {head: '类型', paramName: 'type'},
        {head: '操作系统', paramName: 'os', width: '120px'},
        {head: '命名空间', paramName: 'namespace'},
        {head: '状态', paramName: 'status'},
        {head: '镜像来源', paramName: 'imageFrom', width: '120px'},
        {head: '创建时间', paramName: 'createTime'},
        {head: '操作', paramName: 'action'},
    ]
}

const getData = (data: PageList, isSuccess: boolean) => {
    if (isSuccess) {
        return data.items.map(((item: VmImage, index) => {
            const imgStatus = item.status;
            let status = <span>waiting...</span>;
            if (imgStatus) {
                if (imgStatus.state == "ready") {
                    status = <span><MdGppGood color="green" fontSize="2rem"/></span>
                } else {
                    status = <Tooltip active={true} body={imgStatus.state ?? "请等待..."}><Progress
                        percentValue={imgStatus.progress ?? 0}
                        hasContent={true}
                        hasStripe={true}
                        hasAnimation={true}
                    />
                    </Tooltip>
                }
            }

            const annotations = item.metadata?.annotations;
            let displayName = "";
            if (annotations) {
                displayName = annotations[constants.annotations.displayName];
                if (/^\s$/.test(displayName)) {
                    displayName = item.metadata.name!!;
                }
            }
            return {
                key: item.metadata.name,
                name: displayName,
                type: getImageTypeText(item.spec?.imageType),
                os: getOsIcon(item.spec?.osType),
                osVersion: item.spec?.osVersion,
                namespace: item.metadata?.namespace,
                status: status,
                imageFrom: item.spec?.imageFrom,
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
