import {buildNamespacedUri} from "@/modules/common/Constants.tsx";
import {Message} from "@/modules/common/Message.tsx";
import apiClient from "@/client/ApiClient.ts";


export const ImageApi = {
    baseUri: (namespace: string, resType: string) => {
        return buildNamespacedUri(namespace, resType)
    },
    create: (namespace: string, resType: string, data: any): Promise<any> => {
        const url = ImageApi.baseUri(namespace, resType)
        console.log(url)
        return apiClient.post(url, data);
    },
    deleteResource: (namespace: string, resType: string, name: string): Promise<any> => {
        return apiClient.delete(`${ImageApi.baseUri(namespace, resType)}/${name}`);
    },
    upload: (namespace: string, resType: string, imageName: string, formData: FormData) => {
        const uri = `${ImageApi.baseUri(namespace, resType)}/${imageName}/upload`
        apiClient.upload(uri, formData).then((resp) => {
            Message.ok("上传成功", `${imageName}镜像文件上传成功`)
        }).catch(err => {
            Message.error("上传失败", "镜像文件上传失败: " + err, null)
        })
    }
}