import {Model, ModelData, TypeMeta} from "@kubernetes-models/base";
import {IIoK8sApimachineryPkgApisMetaV1ObjectMeta} from "@kubernetes-models/apimachinery/apis/meta/v1/ObjectMeta";


export interface VmImage extends TypeMeta {
    "apiVersion": "api.kubeall.io/v1";
    "kind": "Image";
    "metadata": IIoK8sApimachineryPkgApisMetaV1ObjectMeta;
    "spec"?: IImageSpec;
    "status"?: IImageStatus;
}


export interface IImageSpec {
    osType?: string;
    osVersion?: string;
    imageType?: string;
    storageClassName?: string;
    imageStorageClassName?: string;
    backend?: string;
    imageFrom?: string;
    sourceStorageClassName?: string;
    storageClassParameters?: {
        [key: string]: string;
    };
}


export interface IImageStatus {
    progress?: number;
    size?: number;
    virtualSize?: number;
    state?: string;
    message?: string;
    lastStateTransitionTime?: string;
}

export interface GlobalSettings extends TypeMeta {
    "apiVersion": "api.kubeall.io/v1";
    "kind": "GlobalSettings";
    "metadata": IIoK8sApimachineryPkgApisMetaV1ObjectMeta;
    "spec"?: GlobalSettingsSpec;
}


export interface GlobalSettingsSpec {
    osTypes?: Array<string>;

}

export interface IPAddressPool {
    apiVersion: "metallb.io/v1beta1" | "metallb.io/v1beta2";
    kind: "IPAddressPool";
    metadata: IIoK8sApimachineryPkgApisMetaV1ObjectMeta;
    spec: {
        addresses?: string[];
        autoAssign?: boolean;
        avoidBuggyIPs?: boolean;
    };
    status?: {
        [key: string]: any; // 状态字段可根据需要扩展
    };
}