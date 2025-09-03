import React from "react";


export interface MenuItem {
    order: number;
    name: string;
    displayName: string;
    link: string;
    children?: MenuItem[];
}

export interface PageList {
    items: any[];
    page: number;
    pageSize: number;
    totalItems: number;
    totalPages: number;
}

export interface FilterQuery{
    [key: string]: string;
}


export interface CellDefinition {
    head?: string;
    width?: string;
    paramName?: string
}

export interface QueryData {
    name?: string;
    filters?:{
        [key: string]: string;
    };
}

export interface ListParam {
    title?: string;
    description?: string;
    filterArea?: React.ReactNode;
    cells?: Array<any>;
    data?: Array<any>;
    buttons?: React.ReactNode;
    checkedRows?: Array<any>;
    onCheck?: (rowData?: object, checked?: boolean, e?: MouseEvent) => void;
    onCheckAll?: (checked?: boolean, e?: MouseEvent) => void;
}

export interface UploadParam {
    imageName: string;
    formData: FormData;
}

export interface VmImageRow {
    key: string;
    index: number;
    name: string;
    os?: string;
    osVersion?: string;
    namespace: string;
    status?: string;
    imageFrom?: string;
    createTime?: string;
    action?: any
}

export interface VmFormData {
    name: string;
    namespace: string;
    displayName?: string;
    osType: "windows";
    description?: string;
    arch: "amd64";
    cpu: string;
    memory: string;
    enableDisc: boolean;
    discName: string;
    windowsDriverDisc: string;
    systemDisk: diskData;
    dataDisks: Array<diskData>;
    primaryNetwork: vmNetwork;
    advancedSettings: advancedSettings;
}

export interface diskData {
    image?: string;
    capacity: string;
    type: string; //pvc, image, hostDisk
    hostDiskPath?: string;
    storageClass?: string;
    bus?: string;
    volumeType: "Block" | "Filesystem";
}

export interface vmNetwork {
    networkType: string;
    deviceModel: string;
    interfaceType: string;
}

export interface advancedSettings {
    bootLoader: string;
    cpuModel?: string;
    timeZone?: string;
}
